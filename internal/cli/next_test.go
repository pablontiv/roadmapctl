package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNextJSONSeparatesReadyAndBlockedTasks(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"next", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("next exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind        string `json:"kind"`
		RoadmapRoot string `json:"roadmap_root"`
		Summary     struct {
			Status string `json:"status"`
		} `json:"summary"`
		Ready []struct {
			Path   string `json:"path"`
			Status string `json:"status"`
		} `json:"ready"`
		Blocked []struct {
			Path     string   `json:"path"`
			Status   string   `json:"status"`
			Blockers []string `json:"blockers"`
		} `json:"blocked"`
		Diagnostics []any `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/next" || report.Summary.Status != "ok" || report.RoadmapRoot == "" || len(report.Diagnostics) != 0 {
		t.Fatalf("report = %#v", report)
	}
	if len(report.Ready) != 1 || report.Ready[0].Path != "O01-work/T001-ready.md" || report.Ready[0].Status != "Pending" {
		t.Fatalf("ready report = %#v", report)
	}
	if len(report.Blocked) != 1 || report.Blocked[0].Path != "O01-work/T002-blocked.md" || report.Blocked[0].Status != "Pending" || len(report.Blocked[0].Blockers) != 1 || report.Blocked[0].Blockers[0] != "O01-work/T001-ready.md" {
		t.Fatalf("blocked report = %#v", report)
	}
}

func TestNextLimitRestrictsReadyTasks(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"next", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--limit", "1", "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("next exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Ready []struct {
			Path string `json:"path"`
		} `json:"ready"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if len(report.Ready) != 1 {
		t.Fatalf("ready = %#v", report.Ready)
	}
}
