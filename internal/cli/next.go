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

type nextReport struct {
	Version     int                      `json:"version"`
	Kind        string                   `json:"kind"`
	Summary     diagnostics.Summary      `json:"summary"`
	Root        string                   `json:"root"`
	RoadmapRoot string                   `json:"roadmap_root"`
	Ready       []nextTask               `json:"ready"`
	Blocked     []nextTask               `json:"blocked"`
	Diagnostics []diagnostics.Diagnostic `json:"diagnostics"`
}

type nextTask struct {
	Path     string   `json:"path"`
	Status   string   `json:"status"`
	Title    string   `json:"title,omitempty"`
	Blockers []string `json:"blockers,omitempty"`
}

func newNextCommand(options *Options, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	var limit int
	cmd := &cobra.Command{Use: "next", Short: "List ready and blocked roadmap tasks without mutating state.", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true, RunE: func(cmd *cobra.Command, args []string) error {
		report := runNext(context.Background(), *options, limit)
		if options.Output == "json" {
			if err := json.NewEncoder(stdout).Encode(report); err != nil {
				fmt.Fprintf(stderr, "next: render JSON report: %v\n", err)
				*exitCode = ExitInternal
				return nil
			}
		} else {
			fmt.Fprintf(stdout, "%s\nstatus: %s\nready: %d\nblocked: %d\n", report.Kind, report.Summary.Status, len(report.Ready), len(report.Blocked))
		}
		*exitCode = diagnostics.ExitCode(diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics), options.Strict)
		return nil
	}}
	cmd.Flags().IntVar(&limit, "limit", 0, "maximum number of ready tasks to return")
	return cmd
}

func runNext(ctx context.Context, options Options, limit int) nextReport {
	cfg, err := config.Load(options.Repo)
	if err != nil {
		found := []diagnostics.Diagnostic{configDiagnostic(absoluteClean(options.Repo), err)}
		return newNextReport(absoluteClean(options.Repo), "", nil, nil, found)
	}
	model, found := readModelForConfig(ctx, cfg, options)
	var ready []nextTask
	var blocked []nextTask
	for _, task := range model.Tasks {
		if task.Done || !task.Active {
			continue
		}
		blockers := incompleteBlockers(model, task)
		if len(blockers) == 0 {
			ready = append(ready, nextTask{Path: task.Path, Status: task.Status, Title: task.Title})
		} else {
			blocked = append(blocked, nextTask{Path: task.Path, Status: task.Status, Title: task.Title, Blockers: blockers})
		}
	}
	sort.Slice(ready, func(i int, j int) bool { return ready[i].Path < ready[j].Path })
	sort.Slice(blocked, func(i int, j int) bool { return blocked[i].Path < blocked[j].Path })
	if limit > 0 && len(ready) > limit {
		ready = ready[:limit]
	}
	return newNextReport(cfg.RepoRoot, cfg.RoadmapRoot, ready, blocked, found)
}

func incompleteBlockers(model roadmap.ReadModel, task roadmap.Task) []string {
	var blockers []string
	for _, dep := range task.Dependencies {
		dependency := model.TaskByPath[dep]
		if dependency == nil || !dependency.Done {
			blockers = append(blockers, dep)
		}
	}
	sort.Strings(blockers)
	return blockers
}

func newNextReport(root string, roadmapRoot string, ready []nextTask, blocked []nextTask, found []diagnostics.Diagnostic) nextReport {
	report := diagnostics.NewReport("roadmapctl/next", root, roadmapRoot, found)
	return nextReport{Version: report.Version, Kind: report.Kind, Summary: report.Summary, Root: report.Root, RoadmapRoot: report.RoadmapRoot, Ready: ready, Blocked: blocked, Diagnostics: report.Diagnostics}
}
