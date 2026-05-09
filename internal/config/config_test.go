package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigErrorFormatsPathAndUnwrapsCause(t *testing.T) {
	cause := errors.New("cause")
	err := &Error{Code: ErrConfigParse, Message: "bad config", Path: ".claude/roadmap.local.md", Cause: cause}
	if got := err.Error(); got != "RMC_CONFIG_PARSE: .claude/roadmap.local.md: bad config" {
		t.Fatalf("Error() = %q", got)
	}
	if !errors.Is(err, cause) {
		t.Fatal("Unwrap did not expose cause")
	}
}

func TestLoadResolvesValidRoadmapRootInsideRepo(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, "roadmap-root: docs/roadmap\n")

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	wantRoot := filepath.Join(repo, "docs", "roadmap")
	if loaded.RoadmapRoot != wantRoot {
		t.Fatalf("RoadmapRoot = %q, want %q", loaded.RoadmapRoot, wantRoot)
	}
	if loaded.RoadmapRootRel != filepath.ToSlash(filepath.Join("docs", "roadmap")) {
		t.Fatalf("RoadmapRootRel = %q", loaded.RoadmapRootRel)
	}
	if loaded.ConfigPath != filepath.Join(repo, ".claude", "roadmap.local.md") {
		t.Fatalf("ConfigPath = %q", loaded.ConfigPath)
	}
}

func TestLoadRejectsParentEscape(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, "roadmap-root: ../outside\n")

	_, err := Load(repo, Options{})
	if err == nil {
		t.Fatal("Load() error = nil, want path escape error")
	}
	var cfgErr *Error
	if !errors.As(err, &cfgErr) {
		t.Fatalf("Load() error type = %T, want *Error", err)
	}
	if cfgErr.Code != ErrRoadmapRootEscape {
		t.Fatalf("error code = %q, want %q", cfgErr.Code, ErrRoadmapRootEscape)
	}
	if cfgErr.ExitCode != 2 {
		t.Fatalf("exit code = %d, want 2", cfgErr.ExitCode)
	}
}

func TestLoadMissingConfigIsUsageError(t *testing.T) {
	repo := t.TempDir()

	_, err := Load(repo, Options{})
	if err == nil {
		t.Fatal("Load() error = nil, want missing config error")
	}
	var cfgErr *Error
	if !errors.As(err, &cfgErr) {
		t.Fatalf("Load() error type = %T, want *Error", err)
	}
	if cfgErr.Code != ErrConfigMissing {
		t.Fatalf("error code = %q, want %q", cfgErr.Code, ErrConfigMissing)
	}
	if cfgErr.ExitCode != 2 {
		t.Fatalf("exit code = %d, want 2", cfgErr.ExitCode)
	}
}

func TestLoadAcceptsWindowsStyleSeparatorsInRoadmapRoot(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, "roadmap-root: docs\\\\roadmap\n")

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := filepath.Join(repo, "docs", "roadmap")
	if loaded.RoadmapRoot != want {
		t.Fatalf("RoadmapRoot = %q, want %q", loaded.RoadmapRoot, want)
	}
	if loaded.RoadmapRootRel != "docs/roadmap" {
		t.Fatalf("RoadmapRootRel = %q, want docs/roadmap", loaded.RoadmapRootRel)
	}
}

func TestLoadAppliesDocumentedDefaultsAndParsesOverrides(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, `roadmap-root: docs/roadmap
done-statuses: ['Done', 'Archived']
active-statuses: ['Ready', 'Doing']
status-values:
  in-progress: Doing
leaf-filter: 'isIndex == false'
outcome-close-verify: ['go test ./...', 'go build ./cmd/roadmapctl']
pr-merge-strategy: merge
commit-style: conventional
auto-push: false
`)

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := loaded.DoneStatuses; len(got) != 2 || got[0] != "Done" || got[1] != "Archived" {
		t.Fatalf("DoneStatuses = %#v, want overrides", got)
	}
	if got := loaded.ActiveStatuses; len(got) != 2 || got[0] != "Ready" || got[1] != "Doing" {
		t.Fatalf("ActiveStatuses = %#v, want overrides", got)
	}
	if loaded.StatusValues.InProgress != "Doing" {
		t.Fatalf("StatusValues.InProgress = %q, want Doing", loaded.StatusValues.InProgress)
	}
	if loaded.StatusValues.Pending != "Pending" {
		t.Fatalf("StatusValues.Pending = %q, want default", loaded.StatusValues.Pending)
	}
	if loaded.LeafFilter != "isIndex == false" {
		t.Fatalf("LeafFilter = %q", loaded.LeafFilter)
	}
	if got := loaded.OutcomeCloseVerify; len(got) != 2 || got[0] != "go test ./..." || got[1] != "go build ./cmd/roadmapctl" {
		t.Fatalf("OutcomeCloseVerify = %#v", got)
	}
	if loaded.PRMergeStrategy != "merge" {
		t.Fatalf("PRMergeStrategy = %q", loaded.PRMergeStrategy)
	}
	if loaded.CommitStyle != "conventional" {
		t.Fatalf("CommitStyle = %q", loaded.CommitStyle)
	}
	if loaded.AutoPush {
		t.Fatal("AutoPush = true, want false override")
	}
}

func TestLoadRoadmapRootOverride(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, "roadmap-root: docs/roadmap\n")

	loaded, err := Load(repo, Options{RoadmapRoot: "custom\\roadmap"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := filepath.Join(repo, "custom", "roadmap")
	if loaded.RoadmapRoot != want {
		t.Fatalf("RoadmapRoot = %q, want %q", loaded.RoadmapRoot, want)
	}
	if loaded.RoadmapRootRel != "custom/roadmap" {
		t.Fatalf("RoadmapRootRel = %q", loaded.RoadmapRootRel)
	}
}

func writeConfig(t *testing.T, repo string, body string) {
	t.Helper()
	dir := filepath.Join(repo, ".claude")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\n" + body + "---\n"
	if err := os.WriteFile(filepath.Join(dir, "roadmap.local.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
