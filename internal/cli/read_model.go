package cli

import (
	"context"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/roadmap"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
)

func readModelForConfig(ctx context.Context, cfg *config.Config, options Options) (roadmap.ReadModel, []diagnostics.Diagnostic) {
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: cfg.RepoRoot, Timeout: options.Timeout})
	var found []diagnostics.Diagnostic
	tree, err := client.Tree(ctx, cfg.RoadmapRoot, cfg.LeafFilter)
	if err != nil {
		found = append(found, rootlineDiagnostic(err))
	}
	query, err := client.Query(ctx, cfg.RoadmapRoot, cfg.LeafFilter, `tipo == "task"`)
	if err != nil {
		found = append(found, rootlineDiagnostic(err))
	}
	graph, err := client.Graph(ctx, cfg.RoadmapRoot, cfg.LeafFilter)
	if err != nil {
		found = append(found, rootlineDiagnostic(err))
	}
	if len(found) > 0 {
		return roadmap.ReadModel{}, found
	}
	model, modelDiagnostics := roadmap.ReadModelFromRootline(tree.Decoded, query.Decoded, graph.Decoded, roadmap.StatusRoleConfig{Done: cfg.DoneStatuses, Active: cfg.ActiveStatuses})
	found = append(found, modelDiagnostics...)
	return model, found
}
