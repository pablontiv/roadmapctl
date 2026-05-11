package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/materialize"
	"github.com/spf13/cobra"
)

type pathPlanReport struct {
	Version     int                              `json:"version"`
	Kind        string                           `json:"kind"`
	Summary     diagnostics.Summary              `json:"summary"`
	Root        string                           `json:"root"`
	RoadmapRoot string                           `json:"roadmap_root"`
	Result      materialize.PathPlanResult       `json:"result"`
	Diagnostics []diagnostics.Diagnostic         `json:"diagnostics"`
}

func newPlanPathsCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	var inputPath string
	cmd := &cobra.Command{
		Use:           "plan-paths",
		Short:         "Plan canonical paths for outcomes and tasks without writing content.",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			report := runPlanPaths(context.Background(), *options, inputPath)
			*exitCode = renderPlanPaths(report, options.Output, stdout, stderr)
			return nil
		},
	}
	cmd.Flags().StringVar(&inputPath, "input", "", "compact path plan JSON file")
	_ = cmd.MarkFlagRequired("input")
	return cmd
}

func runPlanPaths(ctx context.Context, options Options, inputPath string) pathPlanReport {
	repoRoot := absoluteClean(options.Repo)

	if options.Output != "text" && options.Output != "json" {
		return newPathPlanReport(repoRoot, "", materialize.PathPlanResult{}, []diagnostics.Diagnostic{
			{
				ID:       "RMC_PATHPLAN_OUTPUT_INVALID",
				Severity: diagnostics.SeverityError,
				Message:  "unsupported output format",
				ExitCode: diagnostics.ExitUsage,
			},
		})
	}

	// Read and parse input
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return newPathPlanReport(repoRoot, "", materialize.PathPlanResult{}, []diagnostics.Diagnostic{
			{
				ID:       diagnostics.DiagnosticMaterializeInputFieldMissing,
				Severity: diagnostics.SeverityError,
				Message:  "read path plan input: " + err.Error(),
				Path:     inputPath,
				ExitCode: diagnostics.ExitUsage,
			},
		})
	}

	var input materialize.PathPlanInput
	if err := json.Unmarshal(data, &input); err != nil {
		return newPathPlanReport(repoRoot, "", materialize.PathPlanResult{}, []diagnostics.Diagnostic{
			{
				ID:       diagnostics.DiagnosticMaterializeInputKindInvalid,
				Severity: diagnostics.SeverityError,
				Message:  "parse path plan input JSON: " + err.Error(),
				Path:     inputPath,
				ExitCode: diagnostics.ExitUsage,
			},
		})
	}

	// Load configuration to get roadmap root
	cfg, err := loadMaterializeConfig(options)
	if err != nil {
		return newPathPlanReport(repoRoot, "", materialize.PathPlanResult{}, []diagnostics.Diagnostic{
			configDiagnostic(repoRoot, err),
		})
	}

	// Plan the paths
	result, found, err := materialize.PlanPaths(cfg.RoadmapRoot, input)
	if err != nil {
		found = append(found, diagnostics.Diagnostic{
			ID:       "RMC_PATHPLAN_FAILED",
			Severity: diagnostics.SeverityError,
			Message:  err.Error(),
			ExitCode: diagnostics.ExitValidation,
		})
	}

	return newPathPlanReport(cfg.RepoRoot, cfg.RoadmapRoot, result, found)
}

func newPathPlanReport(root string, roadmapRoot string, result materialize.PathPlanResult, found []diagnostics.Diagnostic) pathPlanReport {
	report := diagnostics.NewReport("roadmapctl/plan-paths", root, roadmapRoot, found)
	return pathPlanReport{
		Version:     report.Version,
		Kind:        report.Kind,
		Summary:     report.Summary,
		Root:        report.Root,
		RoadmapRoot: report.RoadmapRoot,
		Result:      result,
		Diagnostics: report.Diagnostics,
	}
}

func renderPlanPaths(report pathPlanReport, output string, stdout io.Writer, stderr io.Writer) int {
	if output == "json" {
		if err := json.NewEncoder(stdout).Encode(report); err != nil {
			fmt.Fprintf(stderr, "plan-paths: render JSON report: %v\n", err)
			return ExitInternal
		}
		return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), false)
	}

	// Text output
	fmt.Fprintf(stdout, "%s\nstatus: %s\n", report.Kind, report.Summary.Status)
	fmt.Fprintf(stdout, "\npaths (%d):\n", len(report.Result.Paths))
	for _, path := range report.Result.Paths {
		fmt.Fprintf(stdout, "  %s [%s] (%s)\n", path.Path, path.Operation, path.Type)
	}

	if len(report.Result.Collisions) > 0 {
		fmt.Fprintf(stdout, "\ncollisions (%d):\n", len(report.Result.Collisions))
		for _, collision := range report.Result.Collisions {
			fmt.Fprintf(stdout, "  %s: %s\n", collision.Path, collision.Reason)
		}
	}

	if len(report.Diagnostics) > 0 {
		fmt.Fprintf(stdout, "\ndiagnostics:\n")
		for _, diag := range report.Diagnostics {
			if diag.Path == "" {
				fmt.Fprintf(stdout, "  [%s] %s: %s\n", diag.Severity, diag.ID, diag.Message)
			} else {
				fmt.Fprintf(stdout, "  [%s] %s %s: %s\n", diag.Severity, diag.ID, diag.Path, diag.Message)
			}
		}
	}

	return diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), false)
}
