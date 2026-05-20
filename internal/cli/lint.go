package cli

import (
	"context"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/reports"
	diag "github.com/pablontiv/picokit/diag"
	roadmaplint "github.com/pablontiv/roadmapctl/internal/lint"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
)

func runLint(ctx context.Context, options Options) reports.Report {
	repoRoot := absoluteClean(options.Repo)
	cfg, err := config.Load(options.Repo)
	if err != nil {
		return reports.NewReport("roadmapctl/lint", repoRoot, "", []diag.Diagnostic{configDiagnostic(repoRoot, err)})
	}
	var found []diag.Diagnostic
	for _, check := range []func(string) ([]diag.Diagnostic, error){
		roadmaplint.CheckOutcomeTaskTables,
		roadmaplint.CheckTaskSections,
		roadmaplint.CheckFilenamePortability,
	} {
		checkDiagnostics, err := check(cfg.RoadmapRoot)
		if err != nil {
			found = append(found, diag.Diagnostic{ID: "RMC_LINT_READ_FAILED", Severity: diag.SeverityError, Message: err.Error(), ExitCode: diag.ExitValidation})
			continue
		}
		found = append(found, checkDiagnostics...)
	}
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: cfg.RepoRoot, Timeout: options.Timeout})
	describe, err := client.Describe(ctx, cfg.RoadmapRoot)
	if err != nil {
		found = append(found, rootlineDiagnostic(err))
	} else {
		found = append(found, roadmaplint.CheckSchemaCompatibility(cfg, describe.Decoded)...)
		found = append(found, roadmaplint.CheckOutcomeSchemaCompatibility(describe.Decoded)...)
	}
	return reports.NewReport("roadmapctl/lint", cfg.RepoRoot, cfg.RoadmapRoot, found)
}
