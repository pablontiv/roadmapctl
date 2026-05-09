package cli

import (
	"context"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	roadmaplint "github.com/pablontiv/roadmapctl/internal/lint"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
)

func runLint(ctx context.Context, options Options) diagnostics.Report {
	repoRoot := absoluteClean(options.Repo)
	cfg, err := config.Load(options.Repo, config.Options{RoadmapRoot: options.RoadmapRoot})
	if err != nil {
		return diagnostics.NewReport("roadmapctl/lint", repoRoot, "", []diagnostics.Diagnostic{configDiagnostic(repoRoot, err)})
	}
	var found []diagnostics.Diagnostic
	for _, check := range []func(string) ([]diagnostics.Diagnostic, error){
		roadmaplint.CheckOutcomeTaskTables,
		roadmaplint.CheckTaskSections,
		roadmaplint.CheckFilenamePortability,
	} {
		checkDiagnostics, err := check(cfg.RoadmapRoot)
		if err != nil {
			found = append(found, diagnostics.Diagnostic{ID: "RMC_LINT_READ_FAILED", Severity: diagnostics.SeverityError, Message: err.Error(), ExitCode: diagnostics.ExitValidation})
			continue
		}
		found = append(found, checkDiagnostics...)
	}
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: cfg.RepoRoot, Timeout: options.Timeout})
	describe, err := client.Describe(ctx, cfg.RoadmapRoot)
	if err != nil {
		found = append(found, rootlineDiagnostic(err))
	} else {
		found = append(found, roadmaplint.CheckSchemaCompatibility(describe.Decoded)...)
	}
	return diagnostics.NewReport("roadmapctl/lint", cfg.RepoRoot, cfg.RoadmapRoot, found)
}
