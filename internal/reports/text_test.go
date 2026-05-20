package reports

import (
	"bytes"
	"strings"
	"testing"

	diag "github.com/pablontiv/picokit/diag"
)

func TestRenderTextHandlesPathlessAndPathedDiagnostics(t *testing.T) {
	report := NewReport("roadmapctl/check", "/repo", "/repo/docs/roadmap", []diag.Diagnostic{
		{ID: "RMC_INFO", Severity: diag.SeverityInfo, Message: "info message"},
		{ID: "RMC_WARN", Severity: diag.SeverityWarning, Message: "warn message", Path: "docs/roadmap/T001-task.md"},
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
