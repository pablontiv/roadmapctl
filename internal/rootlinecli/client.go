package rootlinecli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

const (
	ErrorMissingBinary       ErrorKind = "missing_binary"
	ErrorTimeout             ErrorKind = "timeout"
	ErrorExecution           ErrorKind = "execution"
	ErrorIncompatibleCommand ErrorKind = "incompatible_command"
	ErrorInvalidJSON         ErrorKind = "invalid_json"
	ErrorInvalidShape        ErrorKind = "invalid_shape"
)

type ErrorKind string

type Error struct {
	Kind     ErrorKind
	Message  string
	Path     string
	Stderr   string
	ExitCode int
	Err      error
}

func (e *Error) Error() string {
	if e.Path == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Diagnostic() diagnostics.Diagnostic {
	id := "RMC_ROOTLINE_ERROR"
	if e.Kind == ErrorMissingBinary {
		id = diagnostics.DiagnosticRootlineMissing
	}
	details := map[string]any{"kind": string(e.Kind)}
	if e.Stderr != "" {
		details["stderr"] = e.Stderr
	}
	return diagnostics.Diagnostic{
		ID:       id,
		Severity: diagnostics.SeverityError,
		Message:  e.Message,
		Path:     e.Path,
		Details:  details,
		ExitCode: e.ExitCode,
	}
}

type Options struct {
	Binary   string
	Dir      string
	Env      []string
	Timeout  time.Duration
	Executor Executor
}

type Client struct {
	binary   string
	dir      string
	env      []string
	timeout  time.Duration
	executor Executor
}

type Command struct {
	Path string
	Args []string
	Dir  string
	Env  []string
}

type Result struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
}

type JSONResult struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Decoded  map[string]any
}

type Executor interface {
	Run(ctx context.Context, command Command) (Result, error)
}

func New(options Options) *Client {
	env := options.Env
	if env == nil {
		env = os.Environ()
	}
	executor := options.Executor
	if executor == nil {
		executor = OSExecutor{}
	}
	return &Client{
		binary:   options.Binary,
		dir:      options.Dir,
		env:      append([]string(nil), env...),
		timeout:  options.Timeout,
		executor: executor,
	}
}

func ResolveBinary(explicit string, env []string) (string, error) {
	if env == nil {
		env = os.Environ()
	}
	if explicit != "" {
		return resolveExecutable(explicit, env)
	}
	if rootlineBin := envValue(env, "ROOTLINE_BIN"); rootlineBin != "" {
		return resolveExecutable(rootlineBin, env)
	}
	if path, ok := lookPathInEnv("rootline", env); ok {
		return path, nil
	}
	return "", &Error{
		Kind:     ErrorMissingBinary,
		Message:  "rootline executable not found via --rootline, ROOTLINE_BIN, or PATH",
		ExitCode: diagnostics.ExitEnvironment,
	}
}

func (c *Client) Version(ctx context.Context) (*Result, error) {
	result, err := c.run(ctx, []string{"--version"})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) Validate(ctx context.Context, paths ...string) (*JSONResult, error) {
	args := append([]string{"validate"}, paths...)
	args = append(args, "--output", "json")
	return c.runJSON(ctx, args)
}

func (c *Client) Describe(ctx context.Context, target string, fields ...string) (*JSONResult, error) {
	args := []string{"describe", target}
	for _, field := range fields {
		args = append(args, "--field", field)
	}
	args = append(args, "--output", "json")
	return c.runJSON(ctx, args)
}

func (c *Client) Query(ctx context.Context, root string, wheres ...string) (*JSONResult, error) {
	args := []string{"query", root}
	for _, where := range wheres {
		args = append(args, "--where", where)
	}
	args = append(args, "--output", "json")
	return c.runJSON(ctx, args)
}

func (c *Client) Graph(ctx context.Context, root string, wheres ...string) (*JSONResult, error) {
	args := []string{"graph", root}
	for _, where := range wheres {
		args = append(args, "--where", where)
	}
	args = append(args, "--output", "json")
	return c.runJSON(ctx, args)
}

func (c *Client) Tree(ctx context.Context, root string, wheres ...string) (*JSONResult, error) {
	args := []string{"tree", root}
	for _, where := range wheres {
		args = append(args, "--where", where)
	}
	args = append(args, "--output", "json")
	return c.runJSON(ctx, args)
}

func (c *Client) Set(ctx context.Context, file string, assignments ...string) (*Result, error) {
	args := append([]string{"set", file}, assignments...)
	result, err := c.run(ctx, args)
	return &result, err
}

func (c *Client) NewFile(ctx context.Context, path string) (*Result, error) {
	result, err := c.run(ctx, []string{"new", path})
	return &result, err
}

func (c *Client) New(ctx context.Context, path string) (*Result, error) {
	return c.NewFile(ctx, path)
}

func (c *Client) runJSON(ctx context.Context, args []string) (*JSONResult, error) {
	result, runErr := c.run(ctx, args)

	decoded := map[string]any{}
	if err := json.Unmarshal(result.Stdout, &decoded); err != nil {
		if runErr != nil {
			return nil, runErr
		}
		return nil, &Error{Kind: ErrorInvalidJSON, Message: "rootline returned invalid JSON", Stderr: string(result.Stderr), ExitCode: diagnostics.ExitEnvironment, Err: err}
	}
	return &JSONResult{
		Stdout:   append([]byte(nil), result.Stdout...),
		Stderr:   append([]byte(nil), result.Stderr...),
		ExitCode: result.ExitCode,
		Decoded:  decoded,
	}, runErr
}

func (c *Client) run(ctx context.Context, args []string) (Result, error) {
	binary, err := ResolveBinary(c.binary, c.env)
	if err != nil {
		return Result{}, err
	}

	commandCtx := ctx
	cancel := func() {}
	if c.timeout > 0 {
		commandCtx, cancel = context.WithTimeout(ctx, c.timeout)
	}
	defer cancel()

	result, err := c.executor.Run(commandCtx, Command{
		Path: binary,
		Args: append([]string(nil), args...),
		Dir:  c.dir,
		Env:  append([]string(nil), c.env...),
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(commandCtx.Err(), context.DeadlineExceeded) {
			return result, &Error{Kind: ErrorTimeout, Message: "rootline command timed out", Stderr: string(result.Stderr), ExitCode: diagnostics.ExitEnvironment, Err: err}
		}
		return result, &Error{Kind: classifyExecutionKind(result.Stderr), Message: "rootline command failed", Stderr: string(result.Stderr), ExitCode: executionExitCode(result), Err: err}
	}
	return result, nil
}

func classifyExecutionKind(stderr []byte) ErrorKind {
	text := strings.ToLower(string(stderr))
	if strings.Contains(text, "unknown command") || strings.Contains(text, "unknown flag") || strings.Contains(text, "unknown shorthand flag") {
		return ErrorIncompatibleCommand
	}
	return ErrorExecution
}

func executionExitCode(result Result) int {
	if result.ExitCode != 0 {
		return result.ExitCode
	}
	return diagnostics.ExitEnvironment
}

type OSExecutor struct{}

func (OSExecutor) Run(ctx context.Context, command Command) (Result, error) {
	cmd := exec.CommandContext(ctx, command.Path, command.Args...)
	cmd.Dir = command.Dir
	cmd.Env = command.Env

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := Result{Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}
	if exitErr := new(exec.ExitError); errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
	}
	return result, err
}

func resolveExecutable(path string, env []string) (string, error) {
	if strings.ContainsAny(path, `/\`) || filepath.IsAbs(path) {
		if isExecutable(path) {
			return path, nil
		}
		return "", &Error{Kind: ErrorMissingBinary, Message: "rootline executable not found or not executable", Path: path, ExitCode: diagnostics.ExitEnvironment}
	}
	if resolved, ok := lookPathInEnv(path, env); ok {
		return resolved, nil
	}
	return "", &Error{Kind: ErrorMissingBinary, Message: "rootline executable not found or not executable", Path: path, ExitCode: diagnostics.ExitEnvironment}
}

func lookPathInEnv(name string, env []string) (string, bool) {
	pathValue := envValue(env, "PATH")
	if pathValue == "" {
		return "", false
	}
	for _, dir := range filepath.SplitList(pathValue) {
		if dir == "" {
			continue
		}
		for _, candidate := range executableCandidates(name) {
			path := filepath.Join(dir, candidate)
			if isExecutable(path) {
				return path, true
			}
		}
	}
	return "", false
}

func executableCandidates(name string) []string {
	if runtime.GOOS != "windows" || strings.Contains(filepath.Base(name), ".") {
		return []string{name}
	}
	return []string{name + ".exe", name + ".bat", name + ".cmd", name}
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return info.Mode()&0o111 != 0
}

func envValue(env []string, key string) string {
	prefix := key + "="
	for _, entry := range env {
		if strings.HasPrefix(entry, prefix) {
			return strings.TrimPrefix(entry, prefix)
		}
	}
	return ""
}
