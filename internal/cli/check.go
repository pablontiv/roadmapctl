package cli

import (
	"context"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/roadmap"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
)

func runCheck(ctx context.Context, options Options) diagnostics.Report {
	repoRoot := absoluteClean(options.Repo)
	cfg, err := config.Load(options.Repo, config.Options{RoadmapRoot: options.RoadmapRoot})
	if err != nil {
		found := []diagnostics.Diagnostic{configDiagnostic(repoRoot, err)}
		return diagnostics.NewReport("roadmapctl/check", repoRoot, "", found)
	}

	found, err := roadmap.CheckStructure(cfg.RoadmapRoot)
	if err != nil {
		found = append(found, diagnostics.Diagnostic{
			ID:       "RMC_STRUCTURE_ERROR",
			Severity: diagnostics.SeverityError,
			Message:  err.Error(),
			ExitCode: diagnostics.ExitValidation,
		})
	}

	client := rootlinecli.New(rootlinecli.Options{
		Binary:  options.Rootline,
		Dir:     cfg.RepoRoot,
		Timeout: options.Timeout,
	})
	rootlineDiagnostics, err := roadmap.CheckRootline(ctx, client, roadmap.RootlineCheckOptions{
		RoadmapRoot:     cfg.RoadmapRoot,
		LeafFilter:      cfg.LeafFilter,
		AllowedStatuses: configuredStatuses(cfg),
	})
	if err != nil {
		found = append(found, diagnostics.Diagnostic{
			ID:       "RMC_ROOTLINE_ERROR",
			Severity: diagnostics.SeverityError,
			Message:  err.Error(),
			ExitCode: diagnostics.ExitEnvironment,
		})
	} else {
		found = append(found, rootlineDiagnostics...)
	}

	return diagnostics.NewReport("roadmapctl/check", cfg.RepoRoot, cfg.RoadmapRoot, found)
}

func configuredStatuses(cfg *config.Config) []string {
	values := []string{
		cfg.StatusValues.Pending,
		cfg.StatusValues.Specified,
		cfg.StatusValues.InProgress,
		cfg.StatusValues.Completed,
		cfg.StatusValues.Blocked,
		cfg.StatusValues.Obsolete,
	}
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}
