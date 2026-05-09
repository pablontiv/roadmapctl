package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestTransitionCanStartReadyTaskAllows(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"transition", "can-start", "O01-work/T001-ready.md", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("transition exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind                 string `json:"kind"`
		Allowed              bool   `json:"allowed"`
		TargetStatus         string `json:"target_status"`
		BlockingDependencies []any  `json:"blocking_dependencies"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/transition" || !report.Allowed || report.TargetStatus != "In Progress" || len(report.BlockingDependencies) != 0 {
		t.Fatalf("report = %#v", report)
	}
}

func TestTransitionCanStartBlockedTaskExplainsDependency(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"transition", "can-start", "O01-work/T002-blocked.md", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("transition exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Allowed              bool `json:"allowed"`
		BlockingDependencies []struct {
			Path   string `json:"path"`
			Status string `json:"status"`
		} `json:"blocking_dependencies"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Allowed || len(report.BlockingDependencies) != 1 || report.BlockingDependencies[0].Path != "O01-work/T001-ready.md" || report.BlockingDependencies[0].Status != "Pending" {
		t.Fatalf("report = %#v", report)
	}
}

func TestTransitionCanCompleteReadyTaskAllows(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"transition", "can-complete", "O01-work/T001-ready.md", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("transition exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Allowed      bool   `json:"allowed"`
		TargetStatus string `json:"target_status"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if !report.Allowed || report.TargetStatus != "Completed" {
		t.Fatalf("report = %#v", report)
	}
}
