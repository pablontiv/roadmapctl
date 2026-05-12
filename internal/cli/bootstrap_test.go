package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
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
