package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/cli"
)

func TestRunDelegatesToCLI(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--help"}, &stdout, &stderr, "dev")
	if code != cli.ExitOK {
		t.Fatalf("run exit = %d, want %d; stderr=%q", code, cli.ExitOK, stderr.String())
	}
	if !strings.Contains(stdout.String(), "roadmapctl") {
		t.Fatalf("stdout missing help: %q", stdout.String())
	}
}
