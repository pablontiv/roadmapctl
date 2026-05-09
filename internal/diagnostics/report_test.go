package diagnostics

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRenderJSONWritesOnlyParseableReport(t *testing.T) {
	report := NewReport("roadmapctl/check", "/repo", "/repo/docs/roadmap", []Diagnostic{
		{ID: DiagnosticSingleFileFallback, Severity: SeverityError, Message: "single file fallback", Path: "docs/roadmap/plan-tasks.md", Details: map[string]any{"expected": "TXXX files"}},
	})

	var out bytes.Buffer
	if err := RenderJSON(&out, report); err != nil {
		t.Fatalf("RenderJSON error = %v", err)
	}

	var decoded Report
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("JSON output is not parseable: %v\n%s", err, out.String())
	}
	if decoded.Version != 1 || decoded.Kind != "roadmapctl/check" {
		t.Fatalf("decoded report = %#v", decoded)
	}
	if decoded.Summary.Status != SummaryStatusError || decoded.Summary.Errors != 1 || decoded.Summary.Warnings != 0 || decoded.Summary.Infos != 0 {
		t.Fatalf("summary = %#v", decoded.Summary)
	}
	if len(decoded.Diagnostics) != 1 || decoded.Diagnostics[0].ID != DiagnosticSingleFileFallback {
		t.Fatalf("diagnostics = %#v", decoded.Diagnostics)
	}
	if strings.Contains(out.String(), "not implemented") {
		t.Fatalf("JSON output contains non-report text: %q", out.String())
	}
}

func TestRenderTextIncludesSummaryAndDiagnostics(t *testing.T) {
	report := NewReport("roadmapctl/doctor", "/repo", "", []Diagnostic{
		{ID: DiagnosticRootlineMissing, Severity: SeverityError, Message: "rootline not found"},
		{ID: "RMC_CONFIG_DEFAULT_USED", Severity: SeverityWarning, Message: "using default", Path: ".claude/roadmap.local.md"},
		{ID: "RMC_ENV_PATH", Severity: SeverityInfo, Message: "PATH checked"},
	})

	var out bytes.Buffer
	if err := RenderText(&out, report); err != nil {
		t.Fatalf("RenderText error = %v", err)
	}
	text := out.String()
	for _, want := range []string{"roadmapctl/doctor", "status: error", "errors: 1", "warnings: 1", "infos: 1", "[error] RMC_ENV_ROOTLINE_MISSING", "[warning] RMC_CONFIG_DEFAULT_USED", ".claude/roadmap.local.md"} {
		if !strings.Contains(text, want) {
			t.Fatalf("text output missing %q:\n%s", want, text)
		}
	}
}

func TestExitCodeDerivationCoversContract(t *testing.T) {
	lintWarningReport := NewReport("roadmapctl/lint", "/repo", "/repo/docs/roadmap", []Diagnostic{{ID: DiagnosticLintTaskSectionMissing, Severity: SeverityWarning, Message: "missing section"}})
	tests := []struct {
		name   string
		report Report
		strict bool
		want   int
	}{
		{name: "clean", report: NewReport("roadmapctl/check", "/repo", "/repo/docs/roadmap", nil), want: 0},
		{name: "lint warning non strict", report: lintWarningReport, want: 0},
		{name: "lint warning strict", report: lintWarningReport, strict: true, want: 1},
		{name: "validation", report: NewReport("roadmapctl/check", "/repo", "/repo/docs/roadmap", []Diagnostic{{ID: DiagnosticSingleFileFallback, Severity: SeverityError, Message: "bad"}}), want: 1},
		{name: "usage config", report: NewReport("roadmapctl/doctor", "/repo", "", []Diagnostic{{ID: DiagnosticConfigMissing, Severity: SeverityError, Message: "missing", ExitCode: ExitUsage}}), want: 2},
		{name: "environment", report: NewReport("roadmapctl/doctor", "/repo", "", []Diagnostic{{ID: DiagnosticRootlineMissing, Severity: SeverityError, Message: "missing", ExitCode: ExitEnvironment}}), want: 3},
		{name: "internal", report: NewReport("roadmapctl/check", "/repo", "", []Diagnostic{{ID: "RMC_INTERNAL_UNSUPPORTED_REPORT_VERSION", Severity: SeverityError, Message: "unsupported", ExitCode: ExitInternal}}), want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExitCode(tt.report, tt.strict); got != tt.want {
				t.Fatalf("ExitCode() = %d, want %d", got, tt.want)
			}
		})
	}
	if lintWarningReport.Summary.Status != SummaryStatusWarning || lintWarningReport.Summary.Warnings != 1 || lintWarningReport.Summary.Errors != 0 {
		t.Fatalf("lint warning summary = %#v", lintWarningReport.Summary)
	}
}
