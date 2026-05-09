package diagnostics

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderTextHandlesPathlessAndPathedDiagnostics(t *testing.T) {
	report := NewReport("roadmapctl/check", "/repo", "/repo/docs/roadmap", []Diagnostic{
		{ID: "RMC_INFO", Severity: SeverityInfo, Message: "info message"},
		{ID: "RMC_WARN", Severity: SeverityWarning, Message: "warn message", Path: "docs/roadmap/T001-task.md"},
	})

	var out bytes.Buffer
	if err := RenderText(&out, report); err != nil {
		t.Fatalf("RenderText error = %v", err)
	}
	text := out.String()
	for _, want := range []string{"roadmapctl/check", "status: warning", "[info] RMC_INFO: info message", "[warning] RMC_WARN docs/roadmap/T001-task.md: warn message"} {
		if !strings.Contains(text, want) {
			t.Fatalf("text missing %q:\n%s", want, text)
		}
	}
}
