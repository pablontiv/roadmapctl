package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestNextJSONSeparatesReadyAndBlockedTasks(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"next", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("next exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind  string `json:"kind"`
		Ready []struct {
			Path string `json:"path"`
		} `json:"ready"`
		Blocked []struct {
			Path     string   `json:"path"`
			Blockers []string `json:"blockers"`
		} `json:"blocked"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/next" || len(report.Ready) != 1 || report.Ready[0].Path != "O01-work/T001-ready.md" {
		t.Fatalf("ready report = %#v", report)
	}
	if len(report.Blocked) != 1 || report.Blocked[0].Path != "O01-work/T002-blocked.md" || len(report.Blocked[0].Blockers) != 1 {
		t.Fatalf("blocked report = %#v", report)
	}
}

func TestNextLimitRestrictsReadyTasks(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"next", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--limit", "1", "--output", "json"}, &stdout, &stderr)
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
