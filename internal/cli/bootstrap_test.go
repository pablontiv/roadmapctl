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
	code := Execute([]string{"bootstrap", "inspect", "--repo", repo, "--output", "json"}, &stdout, &stderr, "dev")
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
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--dry-run", "--output", "json"}, &stdout, &stderr, "dev")
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
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "json"}, &stdout, &stderr, "dev")
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
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--rootline", missingRootline, "--apply", "--output", "json"}, &stdout, &stderr, "dev")
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
	Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "json"}, &stdout, &stderr, "dev")

	stdout.Reset()
	stderr.Reset()
	code := Execute([]string{"bootstrap", "--repo", repo, "--output", "json", "--field", "roadmap_root"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap --field exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	output := stdout.String()
	expected := filepath.ToSlash(filepath.Join(repo, "docs", "roadmap"))
	if !strings.Contains(filepath.ToSlash(output), expected) {
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
	Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "json"}, &stdout, &stderr, "dev")

	stdout.Reset()
	stderr.Reset()
	code := Execute([]string{"bootstrap", "--repo", repo, "--output", "json", "--field", "helpers.where_leaf"}, &stdout, &stderr, "dev")
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
	code := Execute([]string{"bootstrap", "--repo", repo, "--output", "json", "--field", "nonexistent"}, &stdout, &stderr, "dev")
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
	code := Execute([]string{"bootstrap", "--repo", repo, "--output", "json", "--field", "helpers"}, &stdout, &stderr, "dev")
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
	code := Execute([]string{"bootstrap", "--repo", repo, "--output", "json", "--field", "diagnostics"}, &stdout, &stderr, "dev")
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
	Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "json"}, &stdout, &stderr, "dev")

	stdout.Reset()
	stderr.Reset()
	code := Execute([]string{"bootstrap", "--repo", repo, "--output", "json"}, &stdout, &stderr, "dev")
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

// mkWorkspaceMember creates a fake git repo with a minimal roadmap config
// inside <root>/<name>/. It does not run `git init` — only the `.git`
// directory needs to exist for workspace member detection.
func mkWorkspaceMember(t *testing.T, root string, name string) string {
	t.Helper()
	repo := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "docs", "roadmap"), 0o755); err != nil {
		t.Fatalf("mkdir docs/roadmap: %v", err)
	}
	cfg := `done_statuses = ["Completed", "Obsolete"]` + "\n" +
		`active_statuses = ["Pending", "Specified", "In Progress"]` + "\n" +
		`leaf_filter = "isIndex == false"` + "\n" +
		`autonomy = "until_done"` + "\n" +
		"[status_values]\n" +
		`pending = "Pending"` + "\n" +
		`specified = "Specified"` + "\n" +
		`in_progress = "In Progress"` + "\n" +
		`completed = "Completed"` + "\n" +
		`blocked = "Blocked"` + "\n" +
		`obsolete = "Obsolete"` + "\n"
	if err := os.WriteFile(filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml"), []byte(cfg), 0o644); err != nil {
		t.Fatalf("write .roadmapctl.toml: %v", err)
	}
	return repo
}

func TestBootstrapAutoDetectsWorkspaceWhenNoGitAtRoot(t *testing.T) {
	root := t.TempDir()
	mkWorkspaceMember(t, root, "alpha")
	mkWorkspaceMember(t, root, "beta")

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "--repo", root, "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		ConfigSource string `json:"config_source"`
		Repos        []struct {
			Name string `json:"name"`
		} `json:"repos"`
		Summary struct {
			Status string `json:"status"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Summary.Status != "ok" {
		t.Fatalf("status = %q, want ok", report.Summary.Status)
	}
	if report.ConfigSource != "workspace" {
		t.Fatalf("config_source = %q, want workspace", report.ConfigSource)
	}
	if len(report.Repos) != 2 {
		t.Fatalf("repos len = %d, want 2; got %#v", len(report.Repos), report.Repos)
	}
	if report.Repos[0].Name != "alpha" || report.Repos[1].Name != "beta" {
		t.Fatalf("repo names = %v, want [alpha beta]", report.Repos)
	}
}

func TestBootstrapHonorsExplicitWorkspaceFlag(t *testing.T) {
	root := t.TempDir()
	mkWorkspaceMember(t, root, "alpha")
	mkWorkspaceMember(t, root, "beta")

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "--workspace", "--repo", root, "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap --workspace exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		ConfigSource string `json:"config_source"`
		Repos        []struct {
			Name string `json:"name"`
		} `json:"repos"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.ConfigSource != "workspace" {
		t.Fatalf("config_source = %q, want workspace", report.ConfigSource)
	}
	if len(report.Repos) != 2 {
		t.Fatalf("repos len = %d, want 2", len(report.Repos))
	}
}

func TestBootstrapWorkspaceEmptyEmitsDiagnostic(t *testing.T) {
	root := t.TempDir()

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "--repo", root, "--output", "json"}, &stdout, &stderr, "dev")
	if code == 0 {
		t.Fatalf("bootstrap on empty dir exit = 0, want non-zero; stdout=%q", stdout.String())
	}
	var report struct {
		ConfigSource string `json:"config_source"`
		Diagnostics  []struct {
			ID string `json:"id"`
		} `json:"diagnostics"`
		Summary struct {
			Status string `json:"status"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Summary.Status != "error" {
		t.Fatalf("status = %q, want error", report.Summary.Status)
	}
	if report.ConfigSource != "workspace" {
		t.Fatalf("config_source = %q, want workspace", report.ConfigSource)
	}
	if len(report.Diagnostics) == 0 || report.Diagnostics[0].ID != "RMC_WORKSPACE_EMPTY" {
		t.Fatalf("expected RMC_WORKSPACE_EMPTY as first diagnostic; got %#v", report.Diagnostics)
	}
	for _, d := range report.Diagnostics {
		if d.ID == "RMC_CONFIG_MISSING" {
			t.Fatalf("must not emit RMC_CONFIG_MISSING in empty workspace case; got %#v", report.Diagnostics)
		}
	}
}

func TestBootstrapWorkspaceMemberWithInvalidConfigIsSkipped(t *testing.T) {
	root := t.TempDir()
	good := mkWorkspaceMember(t, root, "alpha")
	_ = good
	// Create a second member with a broken TOML so config.Load fails for it.
	bad := filepath.Join(root, "beta")
	if err := os.MkdirAll(filepath.Join(bad, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(bad, "docs", "roadmap"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bad, "docs", "roadmap", ".roadmapctl.toml"), []byte("autonomy = \"invalid_value\"\n"), 0o644); err != nil {
		t.Fatalf("write bad toml: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "--repo", root, "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap exit = %d, want 0 (good member loads); stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Repos []struct {
			Name string `json:"name"`
		} `json:"repos"`
		Diagnostics []struct {
			ID       string `json:"id"`
			Severity string `json:"severity"`
		} `json:"diagnostics"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, stdout.String())
	}
	if len(report.Repos) != 1 || report.Repos[0].Name != "alpha" {
		t.Fatalf("expected only alpha to load; got %#v", report.Repos)
	}
	foundSkipped := false
	for _, d := range report.Diagnostics {
		if d.ID == "RMC_WORKSPACE_MEMBER_SKIPPED" {
			foundSkipped = true
			break
		}
	}
	if !foundSkipped {
		t.Fatalf("expected RMC_WORKSPACE_MEMBER_SKIPPED diagnostic; got %#v", report.Diagnostics)
	}
}

func TestBootstrapSingleRepoUnaffected(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	var stdout, stderr bytes.Buffer
	if code := Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "json"}, &stdout, &stderr, "dev"); code != 0 {
		t.Fatalf("init --apply exit = %d; stderr=%q", code, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()
	code := Execute([]string{"bootstrap", "--repo", repo, "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("single-repo bootstrap exit = %d, want 0; stderr=%q", code, stderr.String())
	}
	var report struct {
		ConfigSource string `json:"config_source"`
		Repos        []any  `json:"repos"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.ConfigSource != "toml" {
		t.Fatalf("config_source = %q, want toml (single-repo path must not be reclassified)", report.ConfigSource)
	}
	if len(report.Repos) != 0 {
		t.Fatalf("repos should be empty for single-repo; got %#v", report.Repos)
	}
}
