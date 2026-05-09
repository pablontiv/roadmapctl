package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/config"
	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
)

type contextReport struct {
	Version                int                      `json:"version"`
	Kind                   string                   `json:"kind"`
	Summary                diagnostics.Summary      `json:"summary"`
	Root                   string                   `json:"root"`
	RoadmapRoot            string                   `json:"roadmap_root"`
	ConfigPath             string                   `json:"config_path"`
	ConfigSource           string                   `json:"config_source"`
	RootlineVersion        string                   `json:"rootline_version"`
	Schema                 contextSchema            `json:"schema"`
	StatusValues           config.StatusValues      `json:"status_values"`
	DoneStatuses           []string                 `json:"done_statuses"`
	ActiveStatuses         []string                 `json:"active_statuses"`
	OutcomeCloseVerify     []string                 `json:"outcome_close_verify"`
	PRMergeStrategy        string                   `json:"pr_merge_strategy"`
	CommitStyle            string                   `json:"commit_style"`
	AutoPush               bool                     `json:"auto_push"`
	LoopMaxTasks           int                      `json:"loop_max_tasks"`
	Parallel               bool                     `json:"parallel"`
	Autonomy               string                   `json:"autonomy"`
	CompactAfterTaskCommit bool                     `json:"compact_after_task_commit"`
	PRMode                 bool                     `json:"pr_mode"`
	Helpers                contextHelpers           `json:"helpers"`
	Repos                  []workspaceRepoContext   `json:"repos,omitempty"`
	Diagnostics            []diagnostics.Diagnostic `json:"diagnostics"`
}

type contextSchema struct {
	Estado []string `json:"estado"`
	Tipo   []string `json:"tipo"`
}

type workspaceRepoContext struct {
	Name                   string         `json:"name"`
	Root                   string         `json:"root"`
	RoadmapRoot            string         `json:"roadmap_root"`
	ConfigPath             string         `json:"config_path"`
	OutcomeCloseVerify     []string       `json:"outcome_close_verify"`
	PRMergeStrategy        string         `json:"pr_merge_strategy"`
	CommitStyle            string         `json:"commit_style"`
	AutoPush               bool           `json:"auto_push"`
	LoopMaxTasks           int            `json:"loop_max_tasks"`
	Parallel               bool           `json:"parallel"`
	Autonomy               string         `json:"autonomy"`
	CompactAfterTaskCommit bool           `json:"compact_after_task_commit"`
	PRMode                 bool           `json:"pr_mode"`
	Helpers                contextHelpers `json:"helpers"`
}

type contextHelpers struct {
	WhereLeaf    string `json:"where_leaf"`
	WhereNotDone string `json:"where_not_done"`
	WhereActive  string `json:"where_active"`
}

func runContext(ctx context.Context, options Options) contextReport {
	if options.Workspace {
		return runWorkspaceContext(options)
	}
	repoRoot := absoluteClean(options.Repo)
	cfg, err := config.Load(options.Repo, config.Options{RoadmapRoot: options.RoadmapRoot})
	if err != nil {
		diagnostic := configDiagnostic(repoRoot, err)
		return newContextReport(repoRoot, "", "", "", "", nil, contextSchema{}, []diagnostics.Diagnostic{diagnostic})
	}

	found := configWarnings(cfg)
	client := rootlinecli.New(rootlinecli.Options{Binary: options.Rootline, Dir: cfg.RepoRoot, Timeout: options.Timeout})
	rootlineVersion := ""
	if version, err := client.Version(ctx); err != nil {
		found = append(found, rootlineDiagnostic(err))
	} else {
		rootlineVersion = strings.TrimSpace(string(version.Stdout))
	}

	schema := contextSchema{}
	if describe, err := client.Describe(ctx, dirPath(cfg.RoadmapRoot)); err != nil {
		found = append(found, rootlineDiagnostic(err))
	} else {
		schema.Estado = contextSchemaValues(describe.Decoded, "estado")
		schema.Tipo = contextSchemaValues(describe.Decoded, "tipo")
	}

	return newContextReport(cfg.RepoRoot, cfg.RoadmapRoot, relToRoot(cfg.RepoRoot, cfg.ConfigPath), configSource(cfg), rootlineVersion, cfg, schema, found)
}

func runWorkspaceContext(options Options) contextReport {
	root := absoluteClean(options.Repo)
	repos, found := discoverWorkspaceRepos(root)
	report := newContextReport(root, "", "", "workspace", "", nil, contextSchema{}, found)
	report.Repos = repos
	report.Summary = diagnostics.NewReport(report.Kind, report.Root, report.RoadmapRoot, report.Diagnostics).Summary
	return report
}

func discoverWorkspaceRepos(workspaceRoot string) ([]workspaceRepoContext, []diagnostics.Diagnostic) {
	var repos []workspaceRepoContext
	var found []diagnostics.Diagnostic
	seen := map[string]string{}
	_ = filepath.WalkDir(workspaceRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil || !entry.IsDir() {
			return nil
		}
		if entry.Name() != ".git" {
			return nil
		}
		repoRoot := filepath.Dir(path)
		name := filepath.Base(repoRoot)
		if first, ok := seen[name]; ok {
			found = append(found, diagnostics.Diagnostic{ID: "RMC_WORKSPACE_REPO_AMBIGUOUS", Severity: diagnostics.SeverityError, Message: "multiple workspace repos share the same name", Path: relToRoot(workspaceRoot, repoRoot), Details: map[string]any{"name": name, "first": first}})
			return filepath.SkipDir
		}
		seen[name] = relToRoot(workspaceRoot, repoRoot)
		cfg, err := config.Load(repoRoot, config.Options{})
		if err != nil {
			found = append(found, configDiagnostic(workspaceRoot, err))
			return filepath.SkipDir
		}
		repos = append(repos, workspaceRepoContext{
			Name:                   name,
			Root:                   cfg.RepoRoot,
			RoadmapRoot:            cfg.RoadmapRoot,
			ConfigPath:             relToRoot(cfg.RepoRoot, cfg.ConfigPath),
			OutcomeCloseVerify:     append([]string{}, cfg.OutcomeCloseVerify...),
			PRMergeStrategy:        cfg.PRMergeStrategy,
			CommitStyle:            cfg.CommitStyle,
			AutoPush:               cfg.AutoPush,
			LoopMaxTasks:           cfg.LoopMaxTasks,
			Parallel:               cfg.Parallel,
			Autonomy:               cfg.Autonomy,
			CompactAfterTaskCommit: cfg.CompactAfterTaskCommit,
			PRMode:                 cfg.PRMode,
			Helpers:                contextHelpers{WhereLeaf: cfg.LeafFilter, WhereNotDone: statusWhere("not", cfg.DoneStatuses), WhereActive: statusWhere("", cfg.ActiveStatuses)},
		})
		return filepath.SkipDir
	})
	sort.Slice(repos, func(i int, j int) bool { return repos[i].Name < repos[j].Name })
	return repos, found
}

func newContextReport(root string, roadmapRoot string, configPath string, configSource string, rootlineVersion string, cfg *config.Config, schema contextSchema, found []diagnostics.Diagnostic) contextReport {
	report := diagnostics.NewReport("roadmapctl/context", root, roadmapRoot, found)
	result := contextReport{
		Version:         report.Version,
		Kind:            report.Kind,
		Summary:         report.Summary,
		Root:            report.Root,
		RoadmapRoot:     report.RoadmapRoot,
		ConfigPath:      configPath,
		ConfigSource:    configSource,
		RootlineVersion: rootlineVersion,
		Schema:          schema,
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

func renderContextText(w io.Writer, report contextReport) error {
	_, err := fmt.Fprintf(w, "%s\nstatus: %s\nroot: %s\nroadmap_root: %s\nconfig: %s (%s)\nwhere_leaf: %s\nwhere_not_done: %s\nwhere_active: %s\n", report.Kind, report.Summary.Status, report.Root, report.RoadmapRoot, report.ConfigPath, report.ConfigSource, report.Helpers.WhereLeaf, report.Helpers.WhereNotDone, report.Helpers.WhereActive)
	return err
}

func renderContextJSON(w io.Writer, report contextReport) error {
	return json.NewEncoder(w).Encode(report)
}

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

func dirPath(path string) string {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, `\`) {
		return path
	}
	return path + string(os.PathSeparator)
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
