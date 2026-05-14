package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/config"
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

type contextHelpers struct {
	WhereLeaf    string `json:"where_leaf"`
	WhereNotDone string `json:"where_not_done"`
	WhereActive  string `json:"where_active"`
}

type bootstrapConfigReport struct {
	Version                int                      `json:"version"`
	Kind                   string                   `json:"kind"`
	Summary                diagnostics.Summary      `json:"summary"`
	Root                   string                   `json:"root"`
	RoadmapRoot            string                   `json:"roadmap_root"`
	ConfigPath             string                   `json:"config_path"`
	ConfigSource           string                   `json:"config_source"`
	RootlineVersion        string                   `json:"rootline_version"`
	StatusValues           config.StatusValues      `json:"status_values"`
	DoneStatuses           []string                 `json:"done_statuses"`
	ActiveStatuses         []string                 `json:"active_statuses"`
	OutcomeCloseVerify     []string                 `json:"outcome_close_verify"`
	PRMergeStrategy        string                   `json:"pr_merge_strategy"`
	CommitStyle            string                   `json:"commit_style"`
	AutoPush               bool                     `json:"auto_push"`
	RequiredCodeCoverage   float64                  `json:"required_code_coverage"`
	LoopMaxTasks           int                      `json:"loop_max_tasks"`
	Parallel               bool                     `json:"parallel"`
	Autonomy               string                   `json:"autonomy"`
	CompactAfterTaskCommit bool                     `json:"compact_after_task_commit"`
	PRMode                 bool                     `json:"pr_mode"`
	Helpers                contextHelpers           `json:"helpers"`
	Diagnostics            []diagnostics.Diagnostic `json:"diagnostics"`
}

func newBootstrapCommand(options *Options, stdin io.Reader, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:           "bootstrap",
		Short:         "Inspect or initialize roadmap bootstrap files, or show effective bootstrap context configuration.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			report := buildBootstrapConfig(ctx, *options)
			if hasRepairTriggerDiagnostics(report.Diagnostics) {
				root, roadmapRoot, _ := bootstrapRoots(*options)
				if root != "" && roadmapRoot != "" {
					applied, extraDiags := repairStemIfNeeded(ctx, *options, root, roadmapRoot, yes, stdin, stderr)
					if applied {
						report = buildBootstrapConfig(ctx, *options)
						if len(extraDiags) > 0 {
							report.Diagnostics = append(report.Diagnostics, extraDiags...)
							report.Summary = diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics).Summary
						}
					} else {
						report.Diagnostics = append(report.Diagnostics, extraDiags...)
						report.Summary = diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics).Summary
					}
				}
			}
			*exitCode = renderBootstrapConfig(report, options.Output, stdout, stderr)
			return nil
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply .stem repair without interactive prompt")
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

func buildBootstrapConfig(ctx context.Context, options Options) bootstrapConfigReport {
	root := absoluteClean(options.Repo)
	cfg, err := config.Load(options.Repo, config.Options{RoadmapRoot: options.RoadmapRoot})
	if err != nil {
		diagnostic := configDiagnostic(root, err)
		return newBootstrapConfigReport(root, "", "", "", "", nil, []diagnostics.Diagnostic{diagnostic})
	}

	found := configWarnings(cfg)
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: cfg.RepoRoot, Timeout: options.Timeout})
	rootlineVersion := ""
	if version, err := client.Version(ctx); err != nil {
		found = append(found, rootlineDiagnostic(err))
	} else {
		rootlineVersion = strings.TrimSpace(string(version.Stdout))
	}

	// Check bootstrap files exist or need to be created
	if len(found) == 0 {
		changes := proposedBootstrapChanges(cfg.RepoRoot, cfg.RoadmapRoot, false)
		if len(changes) > 0 {
			// Apply bootstrap changes to ensure config and .stem exist
			applyErrs := applyBootstrapChanges(cfg.RepoRoot, changes)
			found = append(found, applyErrs...)
			if len(applyErrs) == 0 {
				// Reload config after bootstrap apply
				cfg, err = config.Load(options.Repo, config.Options{RoadmapRoot: options.RoadmapRoot})
				if err != nil {
					diagnostic := configDiagnostic(root, err)
					return newBootstrapConfigReport(root, "", "", "", "", nil, []diagnostics.Diagnostic{diagnostic})
				}
			}
		}
	}

	found = append(found, bootstrapSchemaCompatibilityDiagnostics(ctx, options, cfg.RepoRoot, cfg.RoadmapRoot)...)
	return newBootstrapConfigReport(cfg.RepoRoot, cfg.RoadmapRoot, relToRoot(cfg.RepoRoot, cfg.ConfigPath), configSource(cfg), rootlineVersion, cfg, found)
}

func newBootstrapConfigReport(root string, roadmapRoot string, configPath string, configSource string, rootlineVersion string, cfg *config.Config, found []diagnostics.Diagnostic) bootstrapConfigReport {
	report := diagnostics.NewReport("roadmapctl/bootstrap", root, roadmapRoot, found)
	result := bootstrapConfigReport{
		Version:         report.Version,
		Kind:            report.Kind,
		Summary:         report.Summary,
		Root:            report.Root,
		RoadmapRoot:     report.RoadmapRoot,
		ConfigPath:      configPath,
		ConfigSource:    configSource,
		RootlineVersion: rootlineVersion,
		Diagnostics:     report.Diagnostics,
	}
	if cfg != nil {
		result.StatusValues = cfg.StatusValues
		result.DoneStatuses = append([]string(nil), cfg.DoneStatuses...)
		result.ActiveStatuses = append([]string(nil), cfg.ActiveStatuses...)
		result.OutcomeCloseVerify = append([]string{}, cfg.OutcomeCloseVerify...)
		result.PRMergeStrategy = cfg.PRMergeStrategy
		result.CommitStyle = cfg.CommitStyle
		result.AutoPush = cfg.AutoPush
		result.RequiredCodeCoverage = cfg.RequiredCodeCoverage
		result.LoopMaxTasks = cfg.LoopMaxTasks
		result.Parallel = cfg.Parallel
		result.Autonomy = cfg.Autonomy
		result.CompactAfterTaskCommit = cfg.CompactAfterTaskCommit
		result.PRMode = cfg.PRMode
		result.Helpers = contextHelpers{
			WhereLeaf:    cfg.LeafFilter,
			WhereNotDone: statusWhere("not", cfg.DoneStatuses),
			WhereActive:  statusWhere("", cfg.ActiveStatuses),
		}
	}
	return result
}

func renderBootstrapConfig(report bootstrapConfigReport, output string, stdout io.Writer, stderr io.Writer) int {
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
		fmt.Fprintf(stdout, "%s\nstatus: %s\nconfig: %s (%s)\nwhere_leaf: %s\nwhere_not_done: %s\nwhere_active: %s\n",
			report.Kind, report.Summary.Status, report.ConfigPath, report.ConfigSource, report.Helpers.WhereLeaf, report.Helpers.WhereNotDone, report.Helpers.WhereActive)
	}
	return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), false)
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

func hasRepairTriggerDiagnostics(diags []diagnostics.Diagnostic) bool {
	for _, d := range diags {
		if d.ID == diagnostics.DiagnosticLintSchemaOutcomeEstadoRequired ||
			d.ID == diagnostics.DiagnosticLintSchemaOutcomeEstadoNonEmpty {
			return true
		}
	}
	return false
}

func repairStemIfNeeded(ctx context.Context, options Options, root string, roadmapRoot string, yes bool, stdin io.Reader, stderr io.Writer) (applied bool, extraDiags []diagnostics.Diagnostic) {
	stemPath := filepath.Join(roadmapRoot, ".stem")
	content, err := os.ReadFile(stemPath)
	if err != nil {
		return false, nil
	}

	if !isStemRecognizedLegacy(string(content)) {
		return false, []diagnostics.Diagnostic{{
			ID:       diagnostics.DiagnosticBootstrapRepairUnsupportedStem,
			Severity: diagnostics.SeverityError,
			Message:  ".stem has unrecognized custom fields; automatic repair is not supported",
			Path:     ".stem",
			ExitCode: diagnostics.ExitValidation,
		}}
	}

	fmt.Fprintf(stderr, "\nBootstrap: .stem schema is incompatible (estado required for outcomes).\n\n")
	fmt.Fprintf(stderr, "--- current .stem\n%s\n+++ canonical .stem\n%s\n", string(content), templates.BaseStemContent)

	if !yes {
		fmt.Fprintf(stderr, "Update .stem to canonical schema? [y/N]: ")
		reader := bufio.NewReader(stdin)
		answer, _ := reader.ReadString('\n')
		if !strings.EqualFold(strings.TrimSpace(answer), "y") {
			return false, nil
		}
	}

	if err := os.WriteFile(stemPath, []byte(templates.BaseStemContent), 0o644); err != nil {
		return false, []diagnostics.Diagnostic{bootstrapApplyDiagnostic(relToRoot(root, stemPath), err)}
	}

	postOptions := options
	postOptions.Repo = root
	postOptions.RoadmapRoot = relToRoot(root, roadmapRoot)
	postcheck := runCheck(ctx, postOptions)
	if len(postcheck.Diagnostics) > 0 {
		return true, postcheck.Diagnostics
	}
	return true, nil
}

// isStemRecognizedLegacy returns true if the stem content only contains the known
// top-level keys and known schema fields, so automatic repair is safe.
func isStemRecognizedLegacy(content string) bool {
	knownTopLevel := map[string]bool{"version": true, "scope": true, "schema": true, "links": true, "validate": true}
	knownSchemaFields := map[string]bool{"estado": true, "tipo": true, "id": true}
	inSchema := false

	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimRight(rawLine, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		// Top-level key (no leading whitespace)
		if len(line) > 0 && line[0] != ' ' && line[0] != '\t' {
			inSchema = false
			parts := strings.SplitN(trimmed, ":", 2)
			key := parts[0]
			if !knownTopLevel[key] {
				return false
			}
			if key == "schema" {
				inSchema = true
			}
			continue
		}
		// Inside schema: check for field names at exactly one level of indentation (2 spaces)
		if inSchema && strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "   ") {
			parts := strings.SplitN(trimmed, ":", 2)
			key := parts[0]
			if key != "" && !knownSchemaFields[key] {
				return false
			}
		}
	}
	return true
}

// Helper functions (previously in context.go)

func statusWhere(prefix string, values []string) string {
	encoded := make([]string, len(values))
	for i, value := range values {
		encoded[i] = fmt.Sprintf("%q", value)
	}
	inner := "estado in [" + strings.Join(encoded, ", ") + "]"
	if prefix == "not" {
		return "not (" + inner + ")"
	}
	return inner
}

func configSource(cfg *config.Config) string {
	if filepath.Base(cfg.ConfigPath) == "roadmap.local.md" {
		return "legacy"
	}
	if _, err := os.Stat(cfg.ConfigPath); err == nil {
		return "toml"
	}
	return "defaults"
}

func contextSchemaValues(decoded map[string]any, field string) []string {
	if values := contextStringsFromArray(decoded["values"]); len(values) > 0 && field == "estado" {
		return values
	}
	schema, _ := decoded["schema"].(map[string]any)
	fieldSchema, _ := schema[field].(map[string]any)
	return contextStringsFromArray(fieldSchema["values"])
}

func contextStringsFromArray(value any) []string {
	items, _ := value.([]any)
	result := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}
