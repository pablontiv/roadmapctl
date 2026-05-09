package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestPendingJSONListsOnlyNotDoneTasks(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"pending", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("pending exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind  string `json:"kind"`
		Count int    `json:"count"`
		Tasks []struct {
			Path        string `json:"path"`
			Status      string `json:"status"`
			OutcomePath string `json:"outcome_path"`
		} `json:"tasks"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/pending" || report.Count != 2 || len(report.Tasks) != 2 {
		t.Fatalf("report = %#v", report)
	}
	for _, task := range report.Tasks {
		if task.Status == "Completed" || task.Path == "" || task.OutcomePath != "O01-work" {
			t.Fatalf("task = %#v", task)
		}
	}
}

func TestPendingWorkspaceGroupsByRepo(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"pending", "--workspace", "--repo", doctorFixturePath("valid-workspace"), "--output", "json"}, &stdout, &stderr)
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
