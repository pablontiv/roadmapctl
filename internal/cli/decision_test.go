package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDecisionJSONIncludesDeterministicReasons(t *testing.T) {
	requiresRealRootline(t)
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"decision", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("decision exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind    string `json:"kind"`
		Summary struct {
			Status string `json:"status"`
		} `json:"summary"`
		Recommendations []struct {
			Path        string   `json:"path"`
			Status      string   `json:"status"`
			OutcomePath string   `json:"outcome_path"`
			Score       int      `json:"score"`
			Unblocks    []string `json:"unblocks"`
			Reasons     []string `json:"reasons"`
		} `json:"recommendations"`
		QuickWins []struct {
			Path string `json:"path"`
		} `json:"quick_wins"`
		CriticalBlockers []struct {
			Path     string   `json:"path"`
			Unblocks []string `json:"unblocks"`
		} `json:"critical_blockers"`
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
	if report.Kind != "roadmapctl/decision" || report.Summary.Status != "ok" || len(report.Recommendations) != 1 || len(report.Diagnostics) != 0 {
		t.Fatalf("report = %#v", report)
	}
	if report.Recommendations[0].Path != "O01-work/T001-ready.md" || report.Recommendations[0].Status != "Pending" || report.Recommendations[0].OutcomePath != "O01-work" || report.Recommendations[0].Score <= 0 || len(report.Recommendations[0].Reasons) == 0 {
		t.Fatalf("recommendation = %#v", report.Recommendations[0])
	}
	if len(report.QuickWins) != 1 || report.QuickWins[0].Path != "O01-work/T001-ready.md" {
		t.Fatalf("quick wins = %#v", report.QuickWins)
	}
	if len(report.CriticalBlockers) != 1 || report.CriticalBlockers[0].Path != "O01-work/T001-ready.md" || len(report.CriticalBlockers[0].Unblocks) != 1 || report.CriticalBlockers[0].Unblocks[0] != "O01-work/T002-blocked.md" {
		t.Fatalf("critical blockers = %#v", report.CriticalBlockers)
	}
	if len(report.Blocked) != 1 || report.Blocked[0].Status != "Pending" || len(report.Blocked[0].Blockers) != 1 || report.Blocked[0].Blockers[0] != "O01-work/T001-ready.md" {
		t.Fatalf("blocked = %#v", report.Blocked)
	}
}
