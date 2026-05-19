package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/roadmap"
	"github.com/spf13/cobra"
)

type decisionReport struct {
	Version          int                      `json:"version"`
	Kind             string                   `json:"kind"`
	Summary          diagnostics.Summary      `json:"summary"`
	Root             string                   `json:"root"`
	RoadmapRoot      string                   `json:"roadmap_root"`
	Recommendations  []decisionItem           `json:"recommendations"`
	QuickWins        []decisionItem           `json:"quick_wins"`
	CriticalBlockers []decisionItem           `json:"critical_blockers"`
	Blocked          []nextTask               `json:"blocked"`
	Diagnostics      []diagnostics.Diagnostic `json:"diagnostics"`
}

type decisionItem struct {
	Path        string   `json:"path"`
	Status      string   `json:"status"`
	OutcomePath string   `json:"outcome_path,omitempty"`
	Score       int      `json:"score"`
	Unblocks    []string `json:"unblocks,omitempty"`
	Reasons     []string `json:"reasons"`
}

func newDecisionCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	return &cobra.Command{Use: "decision", Short: "Explain deterministic roadmap prioritization data.", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true, RunE: func(cmd *cobra.Command, args []string) error {
		report := runDecision(context.Background(), *options)
		if options.Output == "json" {
			if err := json.NewEncoder(stdout).Encode(report); err != nil {
				fmt.Fprintf(stderr, "decision: render JSON report: %v\n", err)
				*exitCode = ExitInternal
				return nil
			}
		} else {
			fmt.Fprintf(stdout, "%s\nstatus: %s\nrecommendations: %d\nquick_wins: %d\nblocked: %d\n", report.Kind, report.Summary.Status, len(report.Recommendations), len(report.QuickWins), len(report.Blocked))
		}
		*exitCode = diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), options.Strict)
		return nil
	}}
}

func runDecision(ctx context.Context, options Options) decisionReport {
	cfg, err := config.Load(options.Repo)
	if err != nil {
		found := []diagnostics.Diagnostic{configDiagnostic(absoluteClean(options.Repo), err)}
		return newDecisionReport(absoluteClean(options.Repo), "", nil, nil, nil, nil, found)
	}
	model, found := readModelForConfig(ctx, cfg, options)
	recommendations, quickWins, criticalBlockers, blocked := decisionForModel(model, cfg)
	return newDecisionReport(cfg.RepoRoot, cfg.RoadmapRoot, recommendations, quickWins, criticalBlockers, blocked, found)
}

func decisionForModel(model roadmap.ReadModel, cfg *config.Config) ([]decisionItem, []decisionItem, []decisionItem, []nextTask) {
	readyByOutcome := map[string][]decisionItem{}
	var recommendations []decisionItem
	var blocked []nextTask
	for _, task := range model.Tasks {
		if task.Done || !task.Active {
			continue
		}
		blockers := incompleteBlockers(model, task)
		if len(blockers) > 0 {
			blocked = append(blocked, nextTask{Path: task.Path, Status: task.Status, Blockers: blockers})
			continue
		}
		item := scoreDecisionTask(task, cfg)
		recommendations = append(recommendations, item)
		readyByOutcome[task.OutcomePath] = append(readyByOutcome[task.OutcomePath], item)
	}
	sortDecisionItems(recommendations)
	sort.Slice(blocked, func(i int, j int) bool { return blocked[i].Path < blocked[j].Path })
	var quickWins []decisionItem
	for _, items := range readyByOutcome {
		if len(items) == 1 {
			quickWins = append(quickWins, items[0])
		}
	}
	sortDecisionItems(quickWins)
	var critical []decisionItem
	for _, item := range recommendations {
		if len(item.Unblocks) > 0 {
			critical = append(critical, item)
		}
	}
	sortDecisionItems(critical)
	return recommendations, quickWins, critical, blocked
}

func scoreDecisionTask(task roadmap.Task, cfg *config.Config) decisionItem {
	unblocks := append([]string(nil), task.Blocks...)
	sort.Strings(unblocks)
	score := 1 + len(unblocks)*10
	reasons := []string{"ready: dependencies satisfied"}
	if len(unblocks) > 0 {
		reasons = append(reasons, fmt.Sprintf("unblocks %d task(s)", len(unblocks)))
	}
	if task.Status == cfg.StatusValues.InProgress {
		score += 5
		reasons = append(reasons, "in progress work should be finished first")
	}
	return decisionItem{Path: task.Path, Status: task.Status, OutcomePath: task.OutcomePath, Score: score, Unblocks: unblocks, Reasons: reasons}
}

func sortDecisionItems(items []decisionItem) {
	sort.Slice(items, func(i int, j int) bool {
		if items[i].Score != items[j].Score {
			return items[i].Score > items[j].Score
		}
		return items[i].Path < items[j].Path
	})
}

func newDecisionReport(root string, roadmapRoot string, recommendations []decisionItem, quickWins []decisionItem, criticalBlockers []decisionItem, blocked []nextTask, found []diagnostics.Diagnostic) decisionReport {
	report := diagnostics.NewReport("roadmapctl/decision", root, roadmapRoot, found)
	return decisionReport{Version: report.Version, Kind: report.Kind, Summary: report.Summary, Root: report.Root, RoadmapRoot: report.RoadmapRoot, Recommendations: recommendations, QuickWins: quickWins, CriticalBlockers: criticalBlockers, Blocked: blocked, Diagnostics: report.Diagnostics}
}
