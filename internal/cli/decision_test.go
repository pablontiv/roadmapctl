package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDecisionJSONIncludesDeterministicReasons(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"decision", "--repo", doctorFixturePath("valid-next-with-blocked"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("decision exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind            string `json:"kind"`
		Recommendations []struct {
			Path    string   `json:"path"`
			Score   int      `json:"score"`
			Reasons []string `json:"reasons"`
		} `json:"recommendations"`
		QuickWins []struct {
			Path string `json:"path"`
		} `json:"quick_wins"`
		Blocked []struct {
			Path     string   `json:"path"`
			Blockers []string `json:"blockers"`
		} `json:"blocked"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/decision" || len(report.Recommendations) != 1 {
		t.Fatalf("report = %#v", report)
	}
	if report.Recommendations[0].Path != "O01-work/T001-ready.md" || report.Recommendations[0].Score <= 0 || len(report.Recommendations[0].Reasons) == 0 {
		t.Fatalf("recommendation = %#v", report.Recommendations[0])
	}
	if len(report.QuickWins) != 1 || report.QuickWins[0].Path != "O01-work/T001-ready.md" {
		t.Fatalf("quick wins = %#v", report.QuickWins)
	}
	if len(report.Blocked) != 1 || len(report.Blocked[0].Blockers) != 1 {
		t.Fatalf("blocked = %#v", report.Blocked)
	}
}
