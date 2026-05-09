package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/materialize"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
	"github.com/spf13/cobra"
)

type materializeReport struct {
	Version     int                      `json:"version"`
	Kind        string                   `json:"kind"`
	Summary     diagnostics.Summary      `json:"summary"`
	Root        string                   `json:"root"`
	RoadmapRoot string                   `json:"roadmap_root"`
	Applied     bool                     `json:"applied"`
	Changes     []materialize.Change     `json:"changes"`
	Diagnostics []diagnostics.Diagnostic `json:"diagnostics"`
}

func newMaterializeCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	var planPath string
	var changesPath string
	var target string
	var dryRun bool
	var apply bool
	cmd := &cobra.Command{Use: "materialize", Short: "Validate and materialize approved structured roadmap plans.", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true, RunE: func(cmd *cobra.Command, args []string) error {
		report := runMaterialize(context.Background(), *options, planPath, changesPath, target, dryRun, apply)
		*exitCode = renderMaterialize(report, options.Output, stdout, stderr)
		return nil
	}}
	cmd.Flags().StringVar(&planPath, "plan", "", "structured materialize plan JSON file")
	cmd.Flags().StringVar(&changesPath, "changes", "", "materialize dry-run JSON report used as a frozen change set")
	cmd.Flags().StringVar(&target, "target", "", "single canonical roadmap file path to apply from --changes")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show proposed materialization without writing")
	cmd.Flags().BoolVar(&apply, "apply", false, "write approved materialization plan")
	return cmd
}

func runMaterialize(ctx context.Context, options Options, planPath string, changesPath string, target string, dryRun bool, apply bool) materializeReport {
	repoRoot := absoluteClean(options.Repo)
	if options.Output != "text" && options.Output != "json" {
		return newMaterializeReport(repoRoot, "", false, nil, []diagnostics.Diagnostic{{ID: "RMC_MATERIALIZE_OUTPUT_INVALID", Severity: diagnostics.SeverityError, Message: "unsupported output format", ExitCode: diagnostics.ExitUsage}})
	}
	if (planPath == "") == (changesPath == "") {
		return newMaterializeReport(repoRoot, "", false, nil, []diagnostics.Diagnostic{{ID: diagnostics.DiagnosticMaterializeInputFieldMissing, Severity: diagnostics.SeverityError, Message: "materialize requires exactly one of --plan or --changes", ExitCode: diagnostics.ExitUsage}})
	}
	if dryRun == apply {
		return newMaterializeReport(repoRoot, "", false, nil, []diagnostics.Diagnostic{{ID: "RMC_MATERIALIZE_MODE_INVALID", Severity: diagnostics.SeverityError, Message: "materialize requires exactly one of --dry-run or --apply", ExitCode: diagnostics.ExitUsage}})
	}
	if changesPath != "" && (!apply || dryRun) {
		return newMaterializeReport(repoRoot, "", false, nil, []diagnostics.Diagnostic{{ID: "RMC_MATERIALIZE_MODE_INVALID", Severity: diagnostics.SeverityError, Message: "--changes requires --apply", ExitCode: diagnostics.ExitUsage}})
	}
	if target != "" && changesPath == "" {
		return newMaterializeReport(repoRoot, "", false, nil, []diagnostics.Diagnostic{{ID: "RMC_MATERIALIZE_MODE_INVALID", Severity: diagnostics.SeverityError, Message: "--target requires --changes", ExitCode: diagnostics.ExitUsage}})
	}
	cfg, err := loadMaterializeConfig(options)
	if err != nil {
		return newMaterializeReport(repoRoot, "", false, nil, []diagnostics.Diagnostic{configDiagnostic(repoRoot, err)})
	}
	var result materialize.Result
	var found []diagnostics.Diagnostic
	if changesPath != "" {
		changes, readFound := readMaterializeChanges(changesPath)
		found = append(found, readFound...)
		if len(found) == 0 {
			result, found, err = materialize.ApplyTarget(cfg.RoadmapRoot, changes, target)
		}
	} else {
		var plan materialize.Plan
		plan, found = readMaterializePlan(planPath)
		if len(found) == 0 {
			if apply {
				result, found, err = materialize.Apply(cfg.RoadmapRoot, plan)
			} else {
				result, found, err = materialize.DryRun(cfg.RoadmapRoot, plan)
			}
		}
	}
	if err != nil {
		found = append(found, diagnostics.Diagnostic{ID: "RMC_MATERIALIZE_DRY_RUN_FAILED", Severity: diagnostics.SeverityError, Message: err.Error(), ExitCode: diagnostics.ExitValidation})
	}
	if apply && len(found) == 0 {
		found = append(found, validateMaterializedFiles(ctx, cfg, options, result.Changes)...)
		if changesPath == "" {
			postOptions := options
			postOptions.Repo = cfg.RepoRoot
			postOptions.RoadmapRoot = cfg.RoadmapRootRel
			postcheck := runCheck(ctx, postOptions)
			found = append(found, postcheck.Diagnostics...)
		}
	}
	return newMaterializeReport(cfg.RepoRoot, cfg.RoadmapRoot, apply && len(found) == 0, result.Changes, found)
}

func loadMaterializeConfig(options Options) (*config.Config, error) {
	cfg, err := config.Load(options.Repo, config.Options{RoadmapRoot: options.RoadmapRoot})
	if err == nil {
		return cfg, nil
	}
	var cfgErr *config.Error
	if !errors.As(err, &cfgErr) || cfgErr.Code != config.ErrConfigMissing {
		return nil, err
	}
	root, roadmapRoot, found := bootstrapRoots(options)
	if len(found) > 0 {
		return nil, err
	}
	return &config.Config{RepoRoot: root, RoadmapRoot: roadmapRoot, RoadmapRootRel: relToRoot(root, roadmapRoot)}, nil
}

func validateMaterializedFiles(ctx context.Context, cfg *config.Config, options Options, changes []materialize.Change) []diagnostics.Diagnostic {
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: cfg.RepoRoot, Timeout: options.Timeout})
	var found []diagnostics.Diagnostic
	for _, change := range changes {
		if !change.Applied || filepath.Ext(change.Path) != ".md" {
			continue
		}
		if _, err := client.ValidateOne(ctx, filepath.Join(cfg.RoadmapRoot, filepath.FromSlash(change.Path))); err != nil {
			found = append(found, rootlineDiagnostic(err))
		}
	}
	return found
}

func readMaterializeChanges(path string) ([]materialize.Change, []diagnostics.Diagnostic) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, []diagnostics.Diagnostic{{ID: diagnostics.DiagnosticMaterializeInputFieldMissing, Severity: diagnostics.SeverityError, Message: "read materialize changes: " + err.Error(), Path: path, ExitCode: diagnostics.ExitUsage}}
	}
	var report struct {
		Changes []materialize.Change `json:"changes"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, []diagnostics.Diagnostic{{ID: diagnostics.DiagnosticMaterializeInputKindInvalid, Severity: diagnostics.SeverityError, Message: "parse materialize changes JSON: " + err.Error(), Path: path, ExitCode: diagnostics.ExitUsage}}
	}
	return report.Changes, nil
}

func readMaterializePlan(path string) (materialize.Plan, []diagnostics.Diagnostic) {
	data, err := os.ReadFile(path)
	if err != nil {
		return materialize.Plan{}, []diagnostics.Diagnostic{{ID: diagnostics.DiagnosticMaterializeInputFieldMissing, Severity: diagnostics.SeverityError, Message: "read materialize plan: " + err.Error(), Path: path, ExitCode: diagnostics.ExitUsage}}
	}
	var plan materialize.Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return materialize.Plan{}, []diagnostics.Diagnostic{{ID: diagnostics.DiagnosticMaterializeInputKindInvalid, Severity: diagnostics.SeverityError, Message: "parse materialize plan JSON: " + err.Error(), Path: path, ExitCode: diagnostics.ExitUsage}}
	}
	return plan, nil
}

func newMaterializeReport(root string, roadmapRoot string, applied bool, changes []materialize.Change, found []diagnostics.Diagnostic) materializeReport {
	report := diagnostics.NewReport("roadmapctl/materialize", root, roadmapRoot, found)
	return materializeReport{Version: report.Version, Kind: report.Kind, Summary: report.Summary, Root: report.Root, RoadmapRoot: report.RoadmapRoot, Applied: applied, Changes: changes, Diagnostics: report.Diagnostics}
}

func renderMaterialize(report materializeReport, output string, stdout io.Writer, stderr io.Writer) int {
	if output == "json" {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			fmt.Fprintf(stderr, "materialize: render JSON report: %v\n", err)
			return ExitInternal
		}
		return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), false)
	}
	fmt.Fprintf(stdout, "%s\nstatus: %s\napplied: %t\nchanges: %d\n", report.Kind, report.Summary.Status, report.Applied, len(report.Changes))
	for _, change := range report.Changes {
		fmt.Fprintf(stdout, "\n# %s %s\n%s", change.Operation, change.Path, change.Diff)
	}
	for _, diagnostic := range report.Diagnostics {
		if diagnostic.Path == "" {
			fmt.Fprintf(stdout, "[%s] %s: %s\n", diagnostic.Severity, diagnostic.ID, diagnostic.Message)
		} else {
			fmt.Fprintf(stdout, "[%s] %s %s: %s\n", diagnostic.Severity, diagnostic.ID, diagnostic.Path, diagnostic.Message)
		}
	}
	return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), false)
}
