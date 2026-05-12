package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestContextWorkspaceJSONIncludesRepos(t *testing.T) {
	workspace := doctorFixturePath("valid-workspace")

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"context", "--workspace", "--repo", workspace, "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("context workspace exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind  string `json:"kind"`
		Repos []struct {
			Name                   string   `json:"name"`
			RoadmapRoot            string   `json:"roadmap_root"`
			LoopMaxTasks           int      `json:"loop_max_tasks"`
			Parallel               bool     `json:"parallel"`
			Autonomy               string   `json:"autonomy"`
			CompactAfterTaskCommit bool     `json:"compact_after_task_commit"`
			PRMode                 bool     `json:"pr_mode"`
			RequiredCodeCoverage   float64  `json:"required_code_coverage"`
			AutoPush               bool     `json:"auto_push"`
			CommitStyle            string   `json:"commit_style"`
			PRMergeStrategy        string   `json:"pr_merge_strategy"`
			OutcomeCloseVerify     []string `json:"outcome_close_verify"`
		} `json:"repos"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/context" || len(report.Repos) != 2 || report.Repos[0].Name != "alpha" || report.Repos[1].Name != "beta" {
		t.Fatalf("report = %#v", report)
	}
	for _, repo := range report.Repos {
		if repo.LoopMaxTasks != 0 || !repo.Parallel || repo.Autonomy != "until_done" || !repo.CompactAfterTaskCommit || repo.PRMode || repo.RequiredCodeCoverage != 85.0 {
			t.Fatalf("repo execution settings = %#v", repo)
		}
		if !repo.AutoPush || repo.CommitStyle != "conventional" || repo.PRMergeStrategy != "squash" || repo.OutcomeCloseVerify == nil {
			t.Fatalf("repo operational settings = %#v", repo)
		}
	}
}

func TestContextWorkspaceAmbiguousRepoFails(t *testing.T) {
	workspace := doctorFixturePath("invalid-workspace-ambiguous-repo")

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"context", "--workspace", "--repo", workspace, "--output", "json"}, &stdout, &stderr, "dev")
	if code != 1 {
		t.Fatalf("context workspace ambiguous exit = %d, want 1; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Diagnostics []struct {
			ID string `json:"id"`
		} `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if len(report.Diagnostics) == 0 || report.Diagnostics[0].ID != "RMC_WORKSPACE_REPO_AMBIGUOUS" {
		t.Fatalf("diagnostics = %#v", report.Diagnostics)
	}
}

func TestContextJSONIncludesEffectiveHelpersAndOperationalSettings(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"context", "--repo", doctorFixturePath("valid-roadmapctl-toml-default"), "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("context exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var report struct {
		Kind                   string   `json:"kind"`
		ConfigSource           string   `json:"config_source"`
		LoopMaxTasks           int      `json:"loop_max_tasks"`
		Parallel               bool     `json:"parallel"`
		Autonomy               string   `json:"autonomy"`
		CompactAfterTaskCommit bool     `json:"compact_after_task_commit"`
		PRMode                 bool     `json:"pr_mode"`
		RequiredCodeCoverage   float64  `json:"required_code_coverage"`
		AutoPush               bool     `json:"auto_push"`
		CommitStyle            string   `json:"commit_style"`
		PRMergeStrategy        string   `json:"pr_merge_strategy"`
		OutcomeCloseVerify     []string `json:"outcome_close_verify"`
		Helpers                struct {
			WhereLeaf    string `json:"where_leaf"`
			WhereNotDone string `json:"where_not_done"`
			WhereActive  string `json:"where_active"`
		} `json:"helpers"`
		Schema struct {
			Estado []string `json:"estado"`
			Tipo   []string `json:"tipo"`
		} `json:"schema"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout is not parseable JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/context" {
		t.Fatalf("Kind = %q", report.Kind)
	}
	if report.ConfigSource != "toml" {
		t.Fatalf("ConfigSource = %q", report.ConfigSource)
	}
	if report.LoopMaxTasks != 0 || !report.Parallel || report.Autonomy != "until_done" || !report.CompactAfterTaskCommit || report.PRMode || report.RequiredCodeCoverage != 85.0 {
		t.Fatalf("execution settings = %#v", report)
	}
	if !report.AutoPush || report.CommitStyle != "conventional" || report.PRMergeStrategy != "squash" || report.OutcomeCloseVerify == nil {
		t.Fatalf("operational settings = %#v", report)
	}
	if report.Helpers.WhereLeaf != "isIndex == false" || report.Helpers.WhereNotDone != `not (estado in ["Completed", "Obsolete"])` || report.Helpers.WhereActive != `estado in ["Pending", "Specified", "In Progress"]` {
		t.Fatalf("helpers = %#v", report.Helpers)
	}
	if len(report.Schema.Estado) == 0 || len(report.Schema.Tipo) == 0 {
		t.Fatalf("schema = %#v", report.Schema)
	}
}
