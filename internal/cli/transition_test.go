package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
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

func TestTransitionStartDryRunReadyTaskShowsExactUnappliedChange(t *testing.T) {
	fixture := doctorFixturePath("valid-next-with-blocked")
	trackedFile := filepath.Join(fixture, "docs", "roadmap", "O01-work", "T001-ready.md")
	before, err := os.ReadFile(trackedFile)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"transition", "start", "O01-work/T001-ready.md", "--dry-run", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("transition exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	after, err := os.ReadFile(trackedFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("dry-run modified fixture")
	}
	var report struct {
		Allowed bool `json:"allowed"`
		Changes []struct {
			Path    string `json:"path"`
			Field   string `json:"field"`
			Before  string `json:"before"`
			After   string `json:"after"`
			Applied bool   `json:"applied"`
		} `json:"changes"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if !report.Allowed || len(report.Changes) != 1 || report.Changes[0].Before != "Pending" || report.Changes[0].After != "In Progress" || report.Changes[0].Applied {
		t.Fatalf("report = %#v", report)
	}
}

func TestTransitionStartDryRunBlockedTaskHasNoApplicableChanges(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"transition", "start", "O01-work/T002-blocked.md", "--dry-run", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("transition exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Allowed              bool  `json:"allowed"`
		Changes              []any `json:"changes"`
		BlockingDependencies []any `json:"blocking_dependencies"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Allowed || len(report.Changes) != 0 || len(report.BlockingDependencies) != 1 {
		t.Fatalf("report = %#v", report)
	}
}

func TestTransitionSetStatusDryRunMapsRole(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"transition", "set-status", "O01-work/T001-ready.md", "--status", "completed", "--dry-run", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("transition exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Allowed      bool   `json:"allowed"`
		TargetStatus string `json:"target_status"`
		Changes      []struct {
			After   string `json:"after"`
			Applied bool   `json:"applied"`
		} `json:"changes"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if !report.Allowed || report.TargetStatus != "Completed" || len(report.Changes) != 1 || report.Changes[0].After != "Completed" || report.Changes[0].Applied {
		t.Fatalf("report = %#v", report)
	}
}

func TestTransitionDryRunFalseReturnsApplyUnsupported(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"transition", "start", "O01-work/T001-ready.md", "--dry-run=false", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("transition exit = %d, want 2; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Diagnostics []struct {
			ID string `json:"id"`
		} `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if len(report.Diagnostics) != 1 || report.Diagnostics[0].ID != "RMC_TRANSITION_APPLY_FAILED" {
		t.Fatalf("report = %#v", report)
	}
}

func TestTransitionSetStatusRequiresStatus(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"transition", "set-status", "O01-work/T001-ready.md", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("transition exit = %d, want 2; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Diagnostics []struct {
			ID string `json:"id"`
		} `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if len(report.Diagnostics) != 1 || report.Diagnostics[0].ID != "RMC_TRANSITION_STATUS_UNKNOWN" {
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
