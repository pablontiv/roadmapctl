package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

type pendingReport struct {
	Version     int                      `json:"version"`
	Kind        string                   `json:"kind"`
	Summary     diagnostics.Summary      `json:"summary"`
	Root        string                   `json:"root"`
	RoadmapRoot string                   `json:"roadmap_root"`
	Count       int                      `json:"count"`
	Tasks       []pendingTask            `json:"tasks,omitempty"`
	Repos       []pendingRepo            `json:"repos,omitempty"`
	Diagnostics []diagnostics.Diagnostic `json:"diagnostics"`
}

type pendingRepo struct {
	Name  string        `json:"name"`
	Root  string        `json:"root"`
	Count int           `json:"count"`
	Tasks []pendingTask `json:"tasks"`
}

type pendingTask struct {
	Path        string `json:"path"`
	OutcomePath string `json:"outcome_path,omitempty"`
	Status      string `json:"status"`
}

func runPending(ctx context.Context, options Options) pendingReport {
	if options.Workspace {
		return runPendingWorkspace(ctx, options)
	}
	cfg, err := config.Load(options.Repo)
	if err != nil {
		diagnostic := configDiagnostic(absoluteClean(options.Repo), err)
		return newPendingReport(absoluteClean(options.Repo), "", nil, nil, []diagnostics.Diagnostic{diagnostic})
	}
	tasks, found := pendingForConfig(ctx, cfg, options)
	return newPendingReport(cfg.RepoRoot, cfg.RoadmapRoot, tasks, nil, found)
}

func runPendingWorkspace(ctx context.Context, options Options) pendingReport {
	workspaceRoot := absoluteClean(options.Repo)
	repos := workspaceRepoRoots(workspaceRoot)
	seen := map[string]string{}
	var pendingRepos []pendingRepo
	var found []diagnostics.Diagnostic
	for _, repoRoot := range repos {
		name := filepath.Base(repoRoot)
		if first, ok := seen[name]; ok {
			found = append(found, diagnostics.Diagnostic{ID: "RMC_WORKSPACE_REPO_AMBIGUOUS", Severity: diagnostics.SeverityError, Message: "multiple workspace repos share the same name", Path: relToRoot(workspaceRoot, repoRoot), Details: map[string]any{"name": name, "first": first}})
			continue
		}
		seen[name] = relToRoot(workspaceRoot, repoRoot)
		cfg, err := config.Load(repoRoot)
		if err != nil {
			found = append(found, configDiagnostic(workspaceRoot, err))
			continue
		}
		tasks, repoDiagnostics := pendingForConfig(ctx, cfg, options)
		found = append(found, repoDiagnostics...)
		pendingRepos = append(pendingRepos, pendingRepo{Name: name, Root: cfg.RepoRoot, Count: len(tasks), Tasks: tasks})
	}
	sort.Slice(pendingRepos, func(i int, j int) bool { return pendingRepos[i].Name < pendingRepos[j].Name })
	return newPendingReport(workspaceRoot, "", nil, pendingRepos, found)
}

func pendingForConfig(ctx context.Context, cfg *config.Config, options Options) ([]pendingTask, []diagnostics.Diagnostic) {
	model, found := readModelForConfig(ctx, cfg, options)
	if len(found) > 0 {
		return nil, found
	}
	var tasks []pendingTask
	for _, task := range model.Tasks {
		if task.Done {
			continue
		}
		tasks = append(tasks, pendingTask{Path: task.Path, OutcomePath: task.OutcomePath, Status: task.Status})
	}
	sort.Slice(tasks, func(i int, j int) bool { return tasks[i].Path < tasks[j].Path })
	return tasks, found
}

func newPendingReport(root string, roadmapRoot string, tasks []pendingTask, repos []pendingRepo, found []diagnostics.Diagnostic) pendingReport {
	report := diagnostics.NewReport("roadmapctl/pending", root, roadmapRoot, found)
	count := len(tasks)
	for _, repo := range repos {
		count += repo.Count
	}
	return pendingReport{Version: report.Version, Kind: report.Kind, Summary: report.Summary, Root: report.Root, RoadmapRoot: report.RoadmapRoot, Count: count, Tasks: tasks, Repos: repos, Diagnostics: report.Diagnostics}
}

func renderPendingJSON(w io.Writer, report pendingReport) error {
	return json.NewEncoder(w).Encode(report)
}

func renderPendingText(w io.Writer, report pendingReport) error {
	_, err := fmt.Fprintf(w, "%s\nstatus: %s\npending: %d\n", report.Kind, report.Summary.Status, report.Count)
	return err
}

func workspaceRepoRoots(workspaceRoot string) []string {
	var repos []string
	_ = filepath.WalkDir(workspaceRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil || !entry.IsDir() || entry.Name() != ".git" {
			return nil
		}
		repos = append(repos, filepath.Dir(path))
		return filepath.SkipDir
	})
	sort.Strings(repos)
	return repos
}
