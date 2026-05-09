package rootlinecli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func TestResolveBinaryPrefersExplicitThenRootlineBinThenPath(t *testing.T) {
	dir := t.TempDir()
	explicit := writeExecutable(t, dir, "explicit-rootline")
	envRootline := writeExecutable(t, dir, "env-rootline")
	pathRootline := writeExecutable(t, dir, "rootline")

	got, err := ResolveBinary(explicit, []string{"ROOTLINE_BIN=" + envRootline, "PATH=" + dir})
	if err != nil {
		t.Fatalf("ResolveBinary explicit error = %v", err)
	}
	if got != explicit {
		t.Fatalf("ResolveBinary explicit = %q, want %q", got, explicit)
	}

	got, err = ResolveBinary("", []string{"ROOTLINE_BIN=" + envRootline, "PATH=" + dir})
	if err != nil {
		t.Fatalf("ResolveBinary ROOTLINE_BIN error = %v", err)
	}
	if got != envRootline {
		t.Fatalf("ResolveBinary ROOTLINE_BIN = %q, want %q", got, envRootline)
	}

	got, err = ResolveBinary("", []string{"PATH=" + dir})
	if err != nil {
		t.Fatalf("ResolveBinary PATH error = %v", err)
	}
	if got != pathRootline {
		t.Fatalf("ResolveBinary PATH = %q, want %q", got, pathRootline)
	}
}

func TestRootlineErrorFormatsPathAndUnwrapsCause(t *testing.T) {
	cause := errors.New("cause")
	err := &Error{Kind: ErrorExecution, Message: "failed", Path: "docs/roadmap", Err: cause}
	if got := err.Error(); !strings.Contains(got, "docs/roadmap: failed") {
		t.Fatalf("Error() = %q", got)
	}
	if !errors.Is(err, cause) {
		t.Fatalf("Unwrap did not expose cause")
	}
}

func TestClientVersionUsesRootlineVersionCommand(t *testing.T) {
	executor := &recordingExecutor{stdout: []byte("rootline version test")}
	client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Executor: executor})

	result, err := client.Version(context.Background())
	if err != nil {
		t.Fatalf("Version error = %v", err)
	}
	if string(result.Stdout) != "rootline version test" {
		t.Fatalf("Stdout = %q", result.Stdout)
	}
	if !reflect.DeepEqual(executor.commands[0].Args, []string{"--version"}) {
		t.Fatalf("Args = %#v", executor.commands[0].Args)
	}
}

func TestMissingRootlineProducesEnvironmentDiagnostic(t *testing.T) {
	_, err := ResolveBinary("", []string{"PATH=" + t.TempDir()})
	if err == nil {
		t.Fatal("ResolveBinary error = nil, want missing binary error")
	}

	var rootlineErr *Error
	if !errors.As(err, &rootlineErr) {
		t.Fatalf("error type = %T, want *Error", err)
	}
	if rootlineErr.Kind != ErrorMissingBinary {
		t.Fatalf("Kind = %q, want %q", rootlineErr.Kind, ErrorMissingBinary)
	}
	if rootlineErr.ExitCode != diagnostics.ExitEnvironment {
		t.Fatalf("ExitCode = %d, want %d", rootlineErr.ExitCode, diagnostics.ExitEnvironment)
	}
	diagnostic := rootlineErr.Diagnostic()
	if diagnostic.ID != diagnostics.DiagnosticRootlineMissing || diagnostic.ExitCode != diagnostics.ExitEnvironment || diagnostic.Severity != diagnostics.SeverityError {
		t.Fatalf("Diagnostic() = %#v", diagnostic)
	}
}

func TestClientUsesArgsWithoutShellAndControlsDirEnvAndTimeout(t *testing.T) {
	executor := &recordingExecutor{stdout: []byte(`{"version":1,"kind":"rootline/validate","valid":true}`)}
	binary := writeExecutable(t, t.TempDir(), "rootline")
	client := New(Options{
		Binary:   binary,
		Dir:      "/repo",
		Env:      []string{"PATH=/bin", "ROOTLINE_BIN=" + binary, "EXTRA=value"},
		Timeout:  250 * time.Millisecond,
		Executor: executor,
	})

	_, err := client.Validate(context.Background(), "docs/roadmap/T001-task.md")
	if err != nil {
		t.Fatalf("Validate error = %v", err)
	}

	if len(executor.commands) != 1 {
		t.Fatalf("recorded commands = %d, want 1", len(executor.commands))
	}
	command := executor.commands[0]
	if command.Path != binary {
		t.Fatalf("Path = %q", command.Path)
	}
	wantArgs := []string{"validate", "docs/roadmap/T001-task.md", "--output", "json"}
	if !reflect.DeepEqual(command.Args, wantArgs) {
		t.Fatalf("Args = %#v, want %#v", command.Args, wantArgs)
	}
	if strings.Contains(strings.Join(command.Args, " "), "|") || strings.Contains(strings.Join(command.Args, " "), ">") || strings.Contains(strings.Join(command.Args, " "), "sh -c") {
		t.Fatalf("Args look shell-like: %#v", command.Args)
	}
	if command.Dir != "/repo" {
		t.Fatalf("Dir = %q, want /repo", command.Dir)
	}
	if !reflect.DeepEqual(command.Env, []string{"PATH=/bin", "ROOTLINE_BIN=" + binary, "EXTRA=value"}) {
		t.Fatalf("Env = %#v", command.Env)
	}
	deadline := executor.deadlines[0]
	if deadline.IsZero() {
		t.Fatal("context has no deadline, want timeout deadline")
	}
	if time.Until(deadline) <= 0 || time.Until(deadline) > time.Second {
		t.Fatalf("deadline = %v, want near future", deadline)
	}
}

func TestClientCapturesStdoutStderrSeparately(t *testing.T) {
	executor := &recordingExecutor{
		stdout: []byte(`{"version":1,"kind":"rootline/query","meta":{"count":0},"rows":[]}`),
		stderr: []byte("warning on stderr"),
	}
	client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Executor: executor})

	result, err := client.Query(context.Background(), "docs/roadmap", "tipo == \"task\"")
	if err != nil {
		t.Fatalf("Query error = %v", err)
	}
	if !strings.Contains(string(result.Stdout), "rootline/query") {
		t.Fatalf("Stdout = %s", result.Stdout)
	}
	if string(result.Stderr) != "warning on stderr" {
		t.Fatalf("Stderr = %q", result.Stderr)
	}
	wantArgs := []string{"query", "docs/roadmap", "--where", "tipo == \"task\"", "--output", "json"}
	if !reflect.DeepEqual(executor.commands[0].Args, wantArgs) {
		t.Fatalf("Args = %#v, want %#v", executor.commands[0].Args, wantArgs)
	}
}

func TestClientParsesJSONForValidateDescribeQueryAndGraph(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) (*JSONResult, error)
		want []string
		json string
	}{
		{
			name: "validate",
			call: func(c *Client) (*JSONResult, error) { return c.Validate(context.Background(), "a.md") },
			want: []string{"validate", "a.md", "--output", "json"},
			json: `{"version":1,"kind":"rootline/validate","valid":true}`,
		},
		{
			name: "describe",
			call: func(c *Client) (*JSONResult, error) {
				return c.Describe(context.Background(), "docs/roadmap", "schema.estado")
			},
			want: []string{"describe", "docs/roadmap", "--field", "schema.estado", "--output", "json"},
			json: `{"type":"enum","values":["Pending"]}`,
		},
		{
			name: "query",
			call: func(c *Client) (*JSONResult, error) {
				return c.Query(context.Background(), "docs/roadmap", "isIndex == false")
			},
			want: []string{"query", "docs/roadmap", "--where", "isIndex == false", "--output", "json"},
			json: `{"version":1,"kind":"rootline/query","rows":[]}`,
		},
		{
			name: "graph",
			call: func(c *Client) (*JSONResult, error) {
				return c.Graph(context.Background(), "docs/roadmap", "isIndex == false")
			},
			want: []string{"graph", "docs/roadmap", "--where", "isIndex == false", "--output", "json"},
			json: `{"version":1,"kind":"rootline/graph","nodes":[],"edges":[]}`,
		},
		{
			name: "tree",
			call: func(c *Client) (*JSONResult, error) {
				return c.Tree(context.Background(), "docs/roadmap", "isIndex == false")
			},
			want: []string{"tree", "docs/roadmap", "--where", "isIndex == false", "--output", "json"},
			json: `{"version":1,"kind":"rootline/tree","root":{"children":[]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &recordingExecutor{stdout: []byte(tt.json)}
			client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Executor: executor})
			result, err := tt.call(client)
			if err != nil {
				t.Fatalf("call error = %v", err)
			}
			if len(result.Decoded) == 0 || len(result.Stdout) == 0 {
				t.Fatalf("result = %#v", result)
			}
			if !reflect.DeepEqual(executor.commands[0].Args, tt.want) {
				t.Fatalf("Args = %#v, want %#v", executor.commands[0].Args, tt.want)
			}
		})
	}
}

func TestClientUsesRawArgsForSetNewAndValidateOne(t *testing.T) {
	executor := &recordingExecutor{stdout: []byte(`{"version":1,"kind":"rootline/validate","valid":true}`)}
	client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Executor: executor})

	if _, err := client.Set(context.Background(), "docs/roadmap/T001-task.md", "estado=Completed"); err != nil {
		t.Fatalf("Set error = %v", err)
	}
	if _, err := client.New(context.Background(), "docs/roadmap/T002-new.md"); err != nil {
		t.Fatalf("New error = %v", err)
	}
	if _, err := client.ValidateOne(context.Background(), "docs/roadmap/T001-task.md"); err != nil {
		t.Fatalf("ValidateOne error = %v", err)
	}

	want := [][]string{
		{"set", "docs/roadmap/T001-task.md", "estado=Completed"},
		{"new", "docs/roadmap/T002-new.md"},
		{"validate", "docs/roadmap/T001-task.md", "--output", "json"},
	}
	for i := range want {
		if !reflect.DeepEqual(executor.commands[i].Args, want[i]) {
			t.Fatalf("command %d args = %#v, want %#v", i, executor.commands[i].Args, want[i])
		}
	}
}

func TestSetReturnsRawOutputAndStructuredExecutionError(t *testing.T) {
	executor := &recordingExecutor{stdout: []byte("partial set output"), stderr: []byte("set failed"), exitCode: 1, err: errors.New("exit status 1")}
	client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Executor: executor})

	result, err := client.Set(context.Background(), "docs/roadmap/T001-task.md", "estado=Completed")
	if err == nil {
		t.Fatal("Set error = nil, want execution error")
	}
	if result == nil || string(result.Stdout) != "partial set output" || string(result.Stderr) != "set failed" || result.ExitCode != 1 {
		t.Fatalf("result = %#v", result)
	}
	var rootlineErr *Error
	if !errors.As(err, &rootlineErr) {
		t.Fatalf("error type = %T, want *Error", err)
	}
	if rootlineErr.Kind != ErrorExecution || rootlineErr.Stderr != "set failed" || rootlineErr.ExitCode != 1 {
		t.Fatalf("error = %#v", rootlineErr)
	}
}

func TestExecutionClassifiesIncompatibleCommand(t *testing.T) {
	executor := &recordingExecutor{stderr: []byte("unknown command \"tree\" for \"rootline\""), exitCode: 1, err: errors.New("exit status 1")}
	client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Executor: executor})

	_, err := client.Tree(context.Background(), "docs/roadmap")
	if err == nil {
		t.Fatal("Tree error = nil, want incompatible command error")
	}
	var rootlineErr *Error
	if !errors.As(err, &rootlineErr) {
		t.Fatalf("error type = %T, want *Error", err)
	}
	if rootlineErr.Kind != ErrorIncompatibleCommand {
		t.Fatalf("Kind = %q, want %q", rootlineErr.Kind, ErrorIncompatibleCommand)
	}
}

func TestTimeoutProducesControlledEnvironmentError(t *testing.T) {
	executor := &recordingExecutor{err: context.DeadlineExceeded}
	client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Timeout: time.Nanosecond, Executor: executor})

	_, err := client.Graph(context.Background(), "docs/roadmap")
	if err == nil {
		t.Fatal("Graph error = nil, want timeout")
	}
	var rootlineErr *Error
	if !errors.As(err, &rootlineErr) {
		t.Fatalf("error type = %T, want *Error", err)
	}
	if rootlineErr.Kind != ErrorTimeout || rootlineErr.ExitCode != diagnostics.ExitEnvironment {
		t.Fatalf("error = %#v", rootlineErr)
	}
}

func TestOSExecutorCapturesStdoutStderrAndExitCode(t *testing.T) {
	executor := OSExecutor{}
	result, err := executor.Run(context.Background(), Command{
		Path: os.Args[0],
		Args: []string{"-test.run=TestHelperProcess", "--", "exit7"},
		Env:  append(os.Environ(), "GO_WANT_HELPER_PROCESS=1"),
	})
	if err == nil {
		t.Fatal("Run error = nil, want exit error")
	}
	if result.ExitCode != 7 || string(result.Stdout) != "helper stdout\n" || string(result.Stderr) != "helper stderr\n" {
		t.Fatalf("result = %#v", result)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprintln(os.Stdout, "helper stdout")
	fmt.Fprintln(os.Stderr, "helper stderr")
	os.Exit(7)
}

func TestInvalidJSONProducesControlledError(t *testing.T) {
	executor := &recordingExecutor{stdout: []byte("not json"), stderr: []byte("rootline stderr")}
	client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Executor: executor})

	_, err := client.Describe(context.Background(), "docs/roadmap")
	if err == nil {
		t.Fatal("Describe error = nil, want invalid JSON")
	}
	var rootlineErr *Error
	if !errors.As(err, &rootlineErr) {
		t.Fatalf("error type = %T, want *Error", err)
	}
	if rootlineErr.Kind != ErrorInvalidJSON {
		t.Fatalf("Kind = %q, want %q", rootlineErr.Kind, ErrorInvalidJSON)
	}
	if rootlineErr.Stderr != "rootline stderr" {
		t.Fatalf("Stderr = %q", rootlineErr.Stderr)
	}
}

func TestClientReturnsParsedJSONWithExecutionError(t *testing.T) {
	executor := &recordingExecutor{
		stdout:   []byte(`{"version":1,"kind":"rootline/validate","valid":false,"summary":{"invalid":1}}`),
		stderr:   []byte("validation failed"),
		exitCode: 1,
		err:      errors.New("exit status 1"),
	}
	client := New(Options{Binary: writeExecutable(t, t.TempDir(), "rootline"), Executor: executor})

	result, err := client.Validate(context.Background(), "--all", "docs/roadmap")
	if err == nil {
		t.Fatal("Validate error = nil, want execution error")
	}
	if result == nil || numberFromJSONResult(result, "summary", "invalid") != 1 {
		t.Fatalf("result = %#v, want parsed invalid summary", result)
	}
	var rootlineErr *Error
	if !errors.As(err, &rootlineErr) {
		t.Fatalf("error type = %T, want *Error", err)
	}
	if rootlineErr.Kind != ErrorExecution || rootlineErr.Stderr != "validation failed" || rootlineErr.ExitCode != 1 {
		t.Fatalf("error = %#v", rootlineErr)
	}
}

func writeExecutable(t *testing.T, dir string, name string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

type recordingExecutor struct {
	stdout   []byte
	stderr   []byte
	exitCode int
	err      error

	commands  []Command
	deadlines []time.Time
}

func (e *recordingExecutor) Run(ctx context.Context, command Command) (Result, error) {
	e.commands = append(e.commands, command)
	deadline, _ := ctx.Deadline()
	e.deadlines = append(e.deadlines, deadline)
	return Result{Stdout: e.stdout, Stderr: e.stderr, ExitCode: e.exitCode}, e.err
}

func numberFromJSONResult(result *JSONResult, keys ...string) int {
	var current any = result.Decoded
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return 0
		}
		current = m[key]
	}
	value, _ := current.(float64)
	return int(value)
}
