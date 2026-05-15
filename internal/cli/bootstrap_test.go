package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBootstrapInspectIsReadOnlyAndReportsMissing(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "inspect", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("inspect exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind    string   `json:"kind"`
		Missing []string `json:"missing"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/bootstrap/inspect" || len(report.Missing) != 3 {
		t.Fatalf("report = %#v", report)
	}
	if _, err := os.Stat(filepath.Join(repo, "docs", "roadmap")); !os.IsNotExist(err) {
		t.Fatalf("inspect wrote roadmap root: %v", err)
	}
}

func TestBootstrapInitDryRunDoesNotWriteAndShowsChanges(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--dry-run", "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("init dry-run exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind    string `json:"kind"`
		Changes []struct {
			Path    string `json:"path"`
			Applied bool   `json:"applied"`
			Content string `json:"content"`
		} `json:"changes"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/bootstrap/init" || len(report.Changes) != 3 {
		t.Fatalf("report = %#v", report)
	}
	for _, change := range report.Changes {
		if change.Applied {
			t.Fatalf("dry-run applied change: %#v", change)
		}
		if change.Path == "docs/roadmap/.roadmapctl.toml" {
			for _, want := range []string{"required_code_coverage = 85.0", "loop_max_tasks = 0", "parallel = true", "autonomy = \"until_done\"", "compact_after_task_commit = true", "pr_mode = false"} {
				if !strings.Contains(change.Content, want) {
					t.Fatalf("bootstrap TOML missing %q:\n%s", want, change.Content)
				}
			}
		}
	}
	if _, err := os.Stat(filepath.Join(repo, "docs", "roadmap", ".stem")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote .stem: %v", err)
	}
}

func TestBootstrapInitApplyWritesAllowedFiles(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--apply", "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("init apply exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	for _, path := range []string{
		filepath.Join(repo, "docs", "roadmap"),
		filepath.Join(repo, "docs", "roadmap", ".stem"),
		filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected path %s: %v", path, err)
		}
	}
}

func TestBootstrapInitApplyRunsPostcheck(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	var stdout, stderr bytes.Buffer
	missingRootline := filepath.Join(t.TempDir(), "missing-rootline")
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--rootline", missingRootline, "--apply", "--output", "json"}, &stdout, &stderr, "dev")
	if code == 0 {
		t.Fatalf("init apply with failing postcheck exit = 0, want non-zero; stderr=%q stdout=%q", stderr.String(), stdout.String())
	}
	var report struct {
		Diagnostics []struct {
			ID string `json:"id"`
		} `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	for _, diagnostic := range report.Diagnostics {
		if diagnostic.ID == "RMC_ROOTLINE_ERROR" || diagnostic.ID == "RMC_ENV_ROOTLINE_MISSING" {
			return
		}
	}
	t.Fatalf("missing postcheck diagnostic in %#v", report.Diagnostics)
}

func initGitRepo(t *testing.T, repo string) {
	t.Helper()
	cmd := exec.Command("git", "init", repo)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %v\n%s", err, output)
	}
}

func TestBootstrapInitRequiresExplicitMode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", t.TempDir(), "--output", "json"}, &stdout, &stderr, "dev")
	if code != 2 {
		t.Fatalf("init without mode exit = %d, want 2", code)
	}
}

func TestBootstrapInspectOutputFormatValidation(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "inspect", "--repo", repo, "--output", "xml"}, &stdout, &stderr, "dev")
	if code != 2 {
		t.Fatalf("bootstrap inspect with invalid output exit = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unsupported output format") {
		t.Fatalf("stderr missing format error: %s", stderr.String())
	}
}

func TestBootstrapInitTextOutput(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "text"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap init exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "status:") {
		t.Fatalf("text output missing status: %s", output)
	}
	if !strings.Contains(output, "missing:") {
		t.Fatalf("text output missing missing count: %s", output)
	}
}

func TestBootstrapInitTextOutputInvalidFormat(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--dry-run", "--output", "yaml"}, &stdout, &stderr, "dev")
	if code != 2 {
		t.Fatalf("bootstrap init with invalid output exit = %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "unsupported output format") {
		t.Fatalf("stderr missing format error: %s", stderr.String())
	}
}

func TestBootstrapInspectTextOutput(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "inspect", "--repo", repo, "--output", "text"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap inspect exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "status:") {
		t.Fatalf("text output missing status: %s", output)
	}
	if !strings.Contains(output, "changes:") {
		t.Fatalf("text output missing changes count: %s", output)
	}
}

func TestBootstrapInitApplyReportsDiagnosticsOnFileError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod 0o555 does not prevent directory creation on Windows")
	}
	// Create a repo with a directory where a file should be written
	repo := t.TempDir()
	initGitRepo(t, repo)

	// Create a file where we try to write a directory (will cause mkdir to fail)
	// Actually, this is hard to trigger. Let's just verify the init behavior with read-only permissions
	roadmapPath := filepath.Join(repo, "docs")
	if err := os.MkdirAll(roadmapPath, 0o755); err != nil {
		t.Fatal(err)
	}

	// Make it read-only so we can't create subdirs
	if err := os.Chmod(roadmapPath, 0o555); err != nil { //nolint:gosec
		t.Fatal(err)
	}
	defer func() { _ = os.Chmod(roadmapPath, 0o755) }() //nolint:gosec // restore for cleanup

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "json"}, &stdout, &stderr, "dev")
	// Should report diagnostics about the permission error
	if code == 0 {
		t.Fatalf("bootstrap init with permission error should fail, got exit = 0")
	}

	var report struct {
		Diagnostics []struct {
			ID string `json:"id"`
		} `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	found := false
	for _, diag := range report.Diagnostics {
		if diag.ID == "RMC_BOOTSTRAP_APPLY_FAILED" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected RMC_BOOTSTRAP_APPLY_FAILED diagnostic, got %#v", report.Diagnostics)
	}
}

func TestBootstrapApplyDiagnosticFormat(t *testing.T) {
	diag := bootstrapApplyDiagnostic("/some/path", errors.New("permission denied"))
	if diag.ID != "RMC_BOOTSTRAP_APPLY_FAILED" {
		t.Fatalf("bootstrapApplyDiagnostic ID = %q, want RMC_BOOTSTRAP_APPLY_FAILED", diag.ID)
	}
	if diag.Path != "/some/path" {
		t.Fatalf("bootstrapApplyDiagnostic Path = %q, want /some/path", diag.Path)
	}
}

func TestBootstrapFieldExtractionScalarValue(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	// Initialize bootstrap files so we have a valid config
	var stdout, stderr bytes.Buffer
	Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--apply", "--output", "json"}, &stdout, &stderr, "dev")

	stdout.Reset()
	stderr.Reset()
	code := Execute([]string{"bootstrap", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json", "--field", "roadmap_root"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap --field exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	output := stdout.String()
	expected := filepath.Join(repo, "docs", "roadmap")
	if !strings.Contains(output, expected) {
		t.Fatalf("output should contain %q, got: %s", expected, output)
	}
	// Verify raw string output (no JSON quotes)
	if strings.HasPrefix(strings.TrimSpace(output), `"`) {
		t.Fatalf("output should be raw string (no JSON quotes), got: %s", output)
	}
}

func TestBootstrapFieldExtractionNestedValue(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	// Initialize bootstrap files so we have a valid config
	var stdout, stderr bytes.Buffer
	Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--apply", "--output", "json"}, &stdout, &stderr, "dev")

	stdout.Reset()
	stderr.Reset()
	code := Execute([]string{"bootstrap", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json", "--field", "helpers.where_leaf"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap --field nested exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	output := stdout.String()
	// where_leaf should be a string like "isIndex == false"
	if !strings.Contains(output, "isIndex") {
		t.Fatalf("output should contain 'isIndex', got: %s", output)
	}
}

func TestBootstrapFieldExtractionNonexistentField(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json", "--field", "nonexistent"}, &stdout, &stderr, "dev")
	if code == 0 {
		t.Fatalf("bootstrap --field nonexistent exit = 0, want non-zero")
	}
	errMsg := stderr.String()
	if !strings.Contains(errMsg, "key") && !strings.Contains(errMsg, "not found") {
		t.Fatalf("stderr should indicate missing field, got: %s", errMsg)
	}
}

func TestBootstrapFieldExtractionObjectNotAllowed(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json", "--field", "helpers"}, &stdout, &stderr, "dev")
	if code == 0 {
		t.Fatalf("bootstrap --field object exit = 0, want non-zero")
	}
	errMsg := stderr.String()
	if !strings.Contains(errMsg, "object") && !strings.Contains(errMsg, "scalar") {
		t.Fatalf("stderr should indicate object extraction not allowed, got: %s", errMsg)
	}
}

func TestBootstrapFieldExtractionArrayNotAllowed(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json", "--field", "diagnostics"}, &stdout, &stderr, "dev")
	if code == 0 {
		t.Fatalf("bootstrap --field array exit = 0, want non-zero")
	}
	errMsg := stderr.String()
	if !strings.Contains(errMsg, "array") && !strings.Contains(errMsg, "scalar") {
		t.Fatalf("stderr should indicate array extraction not allowed, got: %s", errMsg)
	}
}

func TestBootstrapWithoutFieldStillReturnsFullJSON(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	// Initialize bootstrap files so we have a valid config
	var stdout, stderr bytes.Buffer
	Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--apply", "--output", "json"}, &stdout, &stderr, "dev")

	stdout.Reset()
	stderr.Reset()
	code := Execute([]string{"bootstrap", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap without --field exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	var report map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	// Verify full object is returned
	if _, ok := report["roadmap_root"]; !ok {
		t.Fatalf("full JSON should contain roadmap_root field")
	}
	if _, ok := report["helpers"]; !ok {
		t.Fatalf("full JSON should contain helpers field")
	}
}
