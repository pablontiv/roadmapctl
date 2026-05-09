package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBootstrapInspectIsReadOnlyAndReportsMissing(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "inspect", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json"}, &stdout, &stderr)
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
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--dry-run", "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("init dry-run exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Kind    string `json:"kind"`
		Changes []struct {
			Path    string `json:"path"`
			Applied bool   `json:"applied"`
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
	}
	if _, err := os.Stat(filepath.Join(repo, "docs", "roadmap", ".stem")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote .stem: %v", err)
	}
}

func TestBootstrapInitApplyWritesAllowedFiles(t *testing.T) {
	repo := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--apply", "--output", "json"}, &stdout, &stderr)
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

func TestBootstrapInitRequiresExplicitMode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", t.TempDir(), "--output", "json"}, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("init without mode exit = %d, want 2", code)
	}
}
