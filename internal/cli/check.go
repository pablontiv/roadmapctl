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
	cfg, err := config.Load(options.Repo)
	if err != nil {
		found := []diagnostics.Diagnostic{configDiagnostic(repoRoot, err)}
		return diagnostics.NewReport("roadmapctl/check", repoRoot, "", found)
	}

	found := configWarnings(cfg)
	structureDiagnostics, err := roadmap.CheckStructure(cfg, cfg.RoadmapRoot)
	if err != nil {
		found = append(found, diagnostics.Diagnostic{
			ID:       "RMC_STRUCTURE_ERROR",
			Severity: diagnostics.SeverityError,
			Message:  err.Error(),
			ExitCode: diagnostics.ExitValidation,
		})
	}
	found = append(found, structureDiagnostics...)

	client := rootlinecli.New(rootlinecli.Options{
		Binary:  options.Rootline,
		Dir:     cfg.RepoRoot,
		Timeout: options.Timeout,
	})
	rootlineDiagnostics, err := roadmap.CheckRootline(ctx, cfg, client, roadmap.RootlineCheckOptions{
		RoadmapRoot:         cfg.RoadmapRoot,
		LeafFilter:          cfg.LeafFilter,
		AllowedStatuses:     configuredStatuses(cfg),
		OperationalStatuses: operationalStatuses(cfg),
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
	statuses := operationalStatuses(cfg)
	seen := map[string]bool{}
	result := make([]string, 0, len(statuses))
	for _, status := range statuses {
		if status.Value == "" || seen[status.Value] {
			continue
		}
		seen[status.Value] = true
		result = append(result, status.Value)
	}
	return result
}

func operationalStatuses(cfg *config.Config) []roadmap.OperationalStatus {
	configPath := relToRoot(cfg.RepoRoot, cfg.ConfigPath)
	values := []roadmap.OperationalStatus{
		{Source: "status-values.pending", Value: cfg.StatusValues.Pending, Path: configPath},
		{Source: "status-values.specified", Value: cfg.StatusValues.Specified, Path: configPath},
		{Source: "status-values.in-progress", Value: cfg.StatusValues.InProgress, Path: configPath},
		{Source: "status-values.completed", Value: cfg.StatusValues.Completed, Path: configPath},
		{Source: "status-values.blocked", Value: cfg.StatusValues.Blocked, Path: configPath},
		{Source: "status-values.obsolete", Value: cfg.StatusValues.Obsolete, Path: configPath},
	}
	for _, value := range cfg.DoneStatuses {
		values = append(values, roadmap.OperationalStatus{Source: "done-statuses", Value: value, Path: configPath})
	}
	for _, value := range cfg.ActiveStatuses {
		values = append(values, roadmap.OperationalStatus{Source: "active-statuses", Value: value, Path: configPath})
	}
	return values
}
