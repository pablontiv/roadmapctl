package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestHelpListsImplementedCommands(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"--help"}, &stdout, &stderr, "dev")

	if code != 0 {
		t.Fatalf("Execute(--help) exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"roadmapctl", "doctor", "check", "materialize", "--output", "--repo"} {
		if !strings.Contains(out, want) {
			t.Fatalf("help output missing %q:\n%s", want, out)
		}
	}
	for _, notWant := range []string{"fix"} {
		if strings.Contains(out, notWant) {
			t.Fatalf("help output unexpectedly contains %q:\n%s", notWant, out)
		}
	}
}

func TestCommandHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "doctor", args: []string{"doctor", "--help"}, want: "Diagnose"},
		{name: "check", args: []string{"check", "--help"}, want: "Validate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Execute(tt.args, &stdout, &stderr, "dev")
			if code != 0 {
				t.Fatalf("Execute(%v) exit = %d, want 0; stderr=%q", tt.args, code, stderr.String())
			}
			if !strings.Contains(stdout.String(), tt.want) {
				t.Fatalf("help output missing %q:\n%s", tt.want, stdout.String())
			}
		})
	}
}

func TestUnsupportedOutputIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"doctor", "--output", "yaml"}, &stdout, &stderr, "dev")

	if code != 2 {
		t.Fatalf("Execute unsupported output exit = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unsupported output format") {
		t.Fatalf("stderr missing unsupported output message: %q", stderr.String())
	}
}

func TestUnexpectedArgumentIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"doctor", "extra"}, &stdout, &stderr, "dev")

	if code != 2 {
		t.Fatalf("Execute unexpected arg exit = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") && !strings.Contains(stderr.String(), "accepts 0 arg") {
		t.Fatalf("stderr missing arg error: %q", stderr.String())
	}
}

func TestUnknownCommandIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"plan"}, &stdout, &stderr, "dev")

	if code != 2 {
		t.Fatalf("Execute unknown command exit = %d, want 2", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("stderr missing unknown command message: %q", stderr.String())
	}
}

func TestGlobalFlagsBeforeCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"--repo", doctorFixturePath("valid-outcome-with-tasks"), "--output", "json", "doctor"}, &stdout, &stderr, "dev")

	if code != 0 {
		t.Fatalf("Execute global flags before command exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var report struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout is not parseable JSON report: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/doctor" {
		t.Fatalf("Kind = %q, want roadmapctl/doctor", report.Kind)
	}
}

func TestJSONOutputEmitsOnlyParseableReport(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"doctor", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--output", "json"}, &stdout, &stderr, "dev")

	if code != 0 {
		t.Fatalf("Execute doctor --output json exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var report struct {
		Version     int    `json:"version"`
		Kind        string `json:"kind"`
		Diagnostics []any  `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout is not parseable JSON report: %v\n%s", err, stdout.String())
	}
	if report.Version != 1 || report.Kind != "roadmapctl/doctor" || report.Diagnostics == nil {
		t.Fatalf("unexpected report: %#v", report)
	}
	if strings.Contains(stdout.String(), "not implemented") {
		t.Fatalf("stdout contains extra text: %q", stdout.String())
	}
}
