package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestPendingJSONListsOnlyNotDoneTasks(t *testing.T) {
	requiresRealRootline(t)
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"pending", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("pending exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind        string `json:"kind"`
		Root        string `json:"root"`
		RoadmapRoot string `json:"roadmap_root"`
		Summary     struct {
			Status string `json:"status"`
		} `json:"summary"`
		Count int `json:"count"`
		Tasks []struct {
			Path        string `json:"path"`
			Status      string `json:"status"`
			OutcomePath string `json:"outcome_path"`
		} `json:"tasks"`
		Diagnostics []any `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/pending" || report.Summary.Status != "ok" || report.Root == "" || report.RoadmapRoot == "" || report.Count != 2 || len(report.Tasks) != 2 || len(report.Diagnostics) != 0 {
		t.Fatalf("report = %#v", report)
	}
	want := []struct {
		Path        string
		Status      string
		OutcomePath string
	}{
		{Path: "O01-work/T001-first.md", Status: "Pending", OutcomePath: "O01-work"},
		{Path: "O01-work/T002-second.md", Status: "Pending", OutcomePath: "O01-work"},
	}
	for i, task := range report.Tasks {
		if task.Path != want[i].Path || task.Status != want[i].Status || task.OutcomePath != want[i].OutcomePath {
			t.Fatalf("task[%d] = %#v, want %#v", i, task, want[i])
		}
	}
}

func TestPendingWorkspaceGroupsByRepo(t *testing.T) {
	requiresRealRootline(t)
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"pending", "--workspace", "--repo", doctorFixturePath("valid-workspace"), "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("pending workspace exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Repos []struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		} `json:"repos"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if len(report.Repos) != 2 || report.Repos[0].Name != "alpha" || report.Repos[0].Count != 2 || report.Repos[1].Name != "beta" {
		t.Fatalf("repos = %#v", report.Repos)
	}
}
