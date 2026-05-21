package main

import (
	"bytes"
	"os"
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

func TestMainCallsRunWithVersion(t *testing.T) {
	// Mock os.Exit to capture the exit code instead of terminating.
	var exitCode int
	mainFn = func(code int) {
		exitCode = code
	}
	defer func() {
		mainFn = os.Exit
	}()

	// Ensure os.Args[1:] is empty (simulating "roadmapctl" with no args).
	// This should invoke the root command's help/usage and return ExitOK.
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()
	os.Args = []string{"roadmapctl"}

	main()

	if exitCode != cli.ExitOK {
		t.Errorf("main exit code = %d, want %d", exitCode, cli.ExitOK)
	}
}

func TestRunWithVersionFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"--version"}, &stdout, &stderr, "1.2.3")
	if code != cli.ExitOK {
		t.Errorf("run --version exit = %d, want %d", code, cli.ExitOK)
	}
	// The version flag should output the version.
	output := stdout.String()
	if !strings.Contains(output, "1.2.3") {
		t.Errorf("version output missing '1.2.3': %q", output)
	}
}

func TestRunWithInvalidCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"nonexistent-command"}, &stdout, &stderr, "dev")
	// Invalid command should return non-zero exit code.
	if code == cli.ExitOK {
		t.Errorf("run nonexistent-command exit = %d, want non-zero", code)
	}
}
