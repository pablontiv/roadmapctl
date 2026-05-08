package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpListsOnlyMVPCommands(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"--help"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Execute(--help) exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"roadmapctl", "doctor", "check", "--output", "--repo"} {
		if !strings.Contains(out, want) {
			t.Fatalf("help output missing %q:\n%s", want, out)
		}
	}
	for _, notWant := range []string{"materialize", "fix"} {
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
			code := Execute(tt.args, &stdout, &stderr)
			if code != 0 {
				t.Fatalf("Execute(%v) exit = %d, want 0; stderr=%q", tt.args, code, stderr.String())
			}
			if !strings.Contains(stdout.String(), tt.want) {
				t.Fatalf("help output missing %q:\n%s", tt.want, stdout.String())
			}
		})
	}
}

func TestUnknownCommandIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"plan"}, &stdout, &stderr)

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
