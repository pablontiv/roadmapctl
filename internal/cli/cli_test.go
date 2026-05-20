package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestHelpListsImplementedCommands(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"--help"}, &stdout, &stderr, "dev")

	if code != 0 {
		t.Fatalf("Execute(--help) exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	out := stdout.String()
	for _, want := range []string{"roadmapctl", "doctor", "check", "--output", "--repo"} {
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

type recordingTransport struct {
	count   int32
	lastURL atomic.Value
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	atomic.AddInt32(&r.count, 1)
	r.lastURL.Store(req.URL.String())
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}, nil
}

func swapDefaultTransport(t *testing.T, rt http.RoundTripper) {
	t.Helper()
	original := http.DefaultTransport
	http.DefaultTransport = rt
	t.Cleanup(func() { http.DefaultTransport = original })
}

func TestExecute_AutoupdateSkipsOnEnv(t *testing.T) {
	t.Setenv("ROADMAPCTL_NO_UPDATE", "1")
	rt := &recordingTransport{}
	swapDefaultTransport(t, rt)

	var stdout, stderr bytes.Buffer
	_, wg := executeWithSync([]string{"--version"}, &stdout, &stderr, "1.0.0")
	wg.Wait()

	if c := atomic.LoadInt32(&rt.count); c != 0 {
		t.Fatalf("expected 0 HTTP requests with ROADMAPCTL_NO_UPDATE=1, got %d", c)
	}
}

func TestExecute_AutoupdateSkipsOnDevVersion(t *testing.T) {
	t.Setenv("ROADMAPCTL_NO_UPDATE", "")
	rt := &recordingTransport{}
	swapDefaultTransport(t, rt)

	var stdout, stderr bytes.Buffer
	_, wg := executeWithSync([]string{"--version"}, &stdout, &stderr, "dev")
	wg.Wait()

	if c := atomic.LoadInt32(&rt.count); c != 0 {
		t.Fatalf("expected 0 HTTP requests for dev version, got %d", c)
	}
}

func TestExecute_AutoupdateConstructorParams(t *testing.T) {
	t.Setenv("ROADMAPCTL_NO_UPDATE", "")
	rt := &recordingTransport{}
	swapDefaultTransport(t, rt)

	var stdout, stderr bytes.Buffer
	_, wg := executeWithSync([]string{"--version"}, &stdout, &stderr, "0.0.0")
	wg.Wait()

	if c := atomic.LoadInt32(&rt.count); c == 0 {
		t.Fatalf("expected at least 1 HTTP request from FetchAndStage, got 0")
	}
	url, _ := rt.lastURL.Load().(string)
	if !strings.Contains(url, "pablontiv/roadmapctl") {
		t.Fatalf("expected URL to target repo 'pablontiv/roadmapctl', got %q", url)
	}
}

type slowTransport struct {
	delay time.Duration
}

func (s *slowTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	time.Sleep(s.delay)
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}, nil
}

func TestExecute_AutoupdateFetchRunsInGoroutine(t *testing.T) {
	t.Setenv("ROADMAPCTL_NO_UPDATE", "")
	swapDefaultTransport(t, &slowTransport{delay: 500 * time.Millisecond})

	var stdout, stderr bytes.Buffer
	start := time.Now()
	_, wg := executeWithSync([]string{"--version"}, &stdout, &stderr, "0.0.0")
	elapsed := time.Since(start)
	wg.Wait()

	if elapsed > 100*time.Millisecond {
		t.Fatalf("Execute did not return promptly (%s); FetchAndStage may have blocked", elapsed)
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
