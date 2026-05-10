package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/fsx"
	roadmaplint "github.com/pablontiv/roadmapctl/internal/lint"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
	"github.com/pablontiv/roadmapctl/internal/templates"
	"github.com/spf13/cobra"
)

type bootstrapReport struct {
	Version     int                      `json:"version"`
	Kind        string                   `json:"kind"`
	Summary     diagnostics.Summary      `json:"summary"`
	Root        string                   `json:"root"`
	RoadmapRoot string                   `json:"roadmap_root"`
	Missing     []string                 `json:"missing,omitempty"`
	Changes     []bootstrapChange        `json:"changes,omitempty"`
	Diagnostics []diagnostics.Diagnostic `json:"diagnostics"`
}

type bootstrapChange struct {
	Path      string `json:"path"`
	Operation string `json:"operation"`
	Applied   bool   `json:"applied"`
	Content   string `json:"content,omitempty"`
}

func newBootstrapCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	cmd := &cobra.Command{Use: "bootstrap", Short: "Inspect or initialize roadmap bootstrap files.", SilenceUsage: true, SilenceErrors: true}
	cmd.AddCommand(newBootstrapInspectCommand(options, stdout, stderr, exitCode))
	cmd.AddCommand(newBootstrapInitCommand(options, stdout, stderr, exitCode))
	return cmd
}

func newBootstrapInspectCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	return &cobra.Command{Use: "inspect", Short: "Inspect missing bootstrap files without writing.", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true, RunE: func(cmd *cobra.Command, args []string) error {
		report := buildBootstrapInspect(context.Background(), *options)
		*exitCode = renderBootstrap(report, options.Output, stdout, stderr)
		return nil
	}}
}

func newBootstrapInitCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	var dryRun bool
	var apply bool
	cmd := &cobra.Command{Use: "init", Short: "Initialize missing bootstrap files with explicit dry-run or apply.", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true, RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun == apply {
			return fmt.Errorf("bootstrap init requires exactly one of --dry-run or --apply")
		}
		report := buildBootstrapInit(context.Background(), *options, apply)
		*exitCode = renderBootstrap(report, options.Output, stdout, stderr)
		return nil
	}}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show proposed bootstrap files without writing")
	cmd.Flags().BoolVar(&apply, "apply", false, "write allowed bootstrap files")
	return cmd
}

func buildBootstrapInspect(ctx context.Context, options Options) bootstrapReport {
	root, roadmapRoot, diagnosticsFound := bootstrapRoots(options)
	report := bootstrapReport{Version: 1, Kind: "roadmapctl/bootstrap/inspect", Root: root, RoadmapRoot: roadmapRoot, Diagnostics: diagnosticsFound}
	if len(diagnosticsFound) == 0 {
		report.Missing = missingBootstrapPaths(root, roadmapRoot)
		report.Diagnostics = append(report.Diagnostics, bootstrapSchemaCompatibilityDiagnostics(ctx, options, root, roadmapRoot)...)
	}
	report.Summary = diagnostics.NewReport(report.Kind, root, roadmapRoot, report.Diagnostics).Summary
	return report
}

func buildBootstrapInit(ctx context.Context, options Options, apply bool) bootstrapReport {
	root, roadmapRoot, diagnosticsFound := bootstrapRoots(options)
	report := bootstrapReport{Version: 1, Kind: "roadmapctl/bootstrap/init", Root: root, RoadmapRoot: roadmapRoot, Diagnostics: diagnosticsFound}
	if len(diagnosticsFound) == 0 {
		report.Diagnostics = append(report.Diagnostics, bootstrapSchemaCompatibilityDiagnostics(ctx, options, root, roadmapRoot)...)
	}
	if len(report.Diagnostics) == 0 {
		report.Changes = proposedBootstrapChanges(root, roadmapRoot, apply)
		if apply {
			report.Diagnostics = append(report.Diagnostics, applyBootstrapChanges(root, report.Changes)...)
			if len(report.Diagnostics) == 0 {
				postOptions := options
				postOptions.Repo = root
				postOptions.RoadmapRoot = relToRoot(root, roadmapRoot)
				postcheck := runCheck(ctx, postOptions)
				report.Diagnostics = append(report.Diagnostics, postcheck.Diagnostics...)
			}
		}
	}
	report.Summary = diagnostics.NewReport(report.Kind, root, roadmapRoot, report.Diagnostics).Summary
	return report
}

func bootstrapRoots(options Options) (string, string, []diagnostics.Diagnostic) {
	root := absoluteClean(options.Repo)
	roadmapRootValue := options.RoadmapRoot
	if roadmapRootValue == "" {
		roadmapRootValue = filepath.ToSlash(filepath.Join("docs", "roadmap"))
	}
	roadmapRoot, _, err := fsx.ResolveInside(root, roadmapRootValue)
	if err != nil {
		return root, "", []diagnostics.Diagnostic{{ID: "RMC_CONFIG_ROADMAP_ROOT_ESCAPE", Severity: diagnostics.SeverityError, Message: "roadmap-root must resolve inside repo", ExitCode: diagnostics.ExitUsage}}
	}
	return root, roadmapRoot, nil
}

func missingBootstrapPaths(root string, roadmapRoot string) []string {
	candidates := []string{roadmapRoot, filepath.Join(roadmapRoot, ".stem"), filepath.Join(roadmapRoot, ".roadmapctl.toml")}
	missing := []string{}
	for _, path := range candidates {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			missing = append(missing, relToRoot(root, path))
		}
	}
	return missing
}

func proposedBootstrapChanges(root string, roadmapRoot string, apply bool) []bootstrapChange {
	changes := []bootstrapChange{}
	if _, err := os.Stat(roadmapRoot); os.IsNotExist(err) {
		changes = append(changes, bootstrapChange{Path: relToRoot(root, roadmapRoot), Operation: "mkdir", Applied: apply})
	}
	stemPath := filepath.Join(roadmapRoot, ".stem")
	if _, err := os.Stat(stemPath); os.IsNotExist(err) {
		changes = append(changes, bootstrapChange{Path: relToRoot(root, stemPath), Operation: "create", Applied: apply, Content: templates.BaseStemContent})
	}
	configPath := filepath.Join(roadmapRoot, ".roadmapctl.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		changes = append(changes, bootstrapChange{Path: relToRoot(root, configPath), Operation: "create", Applied: apply, Content: templates.DefaultRoadmapctlTOML})
	}
	return changes
}

func applyBootstrapChanges(root string, changes []bootstrapChange) []diagnostics.Diagnostic {
	var found []diagnostics.Diagnostic
	for _, change := range changes {
		abs := filepath.Join(root, filepath.FromSlash(change.Path))
		switch change.Operation {
		case "mkdir":
			if err := os.MkdirAll(abs, 0o755); err != nil {
				found = append(found, bootstrapApplyDiagnostic(change.Path, err))
			}
		case "create":
			if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
				found = append(found, bootstrapApplyDiagnostic(change.Path, err))
				continue
			}
			if err := os.WriteFile(abs, []byte(change.Content), 0o644); err != nil {
				found = append(found, bootstrapApplyDiagnostic(change.Path, err))
			}
		}
	}
	return found
}

func bootstrapApplyDiagnostic(path string, err error) diagnostics.Diagnostic {
	return diagnostics.Diagnostic{ID: "RMC_BOOTSTRAP_APPLY_FAILED", Severity: diagnostics.SeverityError, Message: err.Error(), Path: path, ExitCode: diagnostics.ExitValidation}
}

func bootstrapSchemaCompatibilityDiagnostics(ctx context.Context, options Options, root string, roadmapRoot string) []diagnostics.Diagnostic {
	stemPath := filepath.Join(roadmapRoot, ".stem")
	if _, err := os.Stat(stemPath); err != nil {
		return nil
	}
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: root, Timeout: options.Timeout})
	describe, err := client.Describe(ctx, ensureRootlineDirPath(roadmapRoot))
	if err != nil {
		return []diagnostics.Diagnostic{rootlineDiagnostic(err)}
	}
	return roadmaplint.CheckOutcomeSchemaCompatibility(describe.Decoded)
}

func renderBootstrap(report bootstrapReport, output string, stdout io.Writer, stderr io.Writer) int {
	if output != "text" && output != "json" {
		fmt.Fprintf(stderr, "bootstrap: unsupported output format %q\n", output)
		return ExitUsage
	}
	if output == "json" {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			fmt.Fprintf(stderr, "bootstrap: render JSON report: %v\n", err)
			return ExitInternal
		}
	} else {
		fmt.Fprintf(stdout, "%s\nstatus: %s\nmissing: %d\nchanges: %d\n", report.Kind, report.Summary.Status, len(report.Missing), len(report.Changes))
	}
	return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), false)
}
