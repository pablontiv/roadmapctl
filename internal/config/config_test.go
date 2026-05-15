package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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

func TestLoadPrefersRoadmapctlTOMLAndInfersRoadmapRoot(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `done_statuses = ["Done"]
active_statuses = ["Ready", "Doing"]
leaf_filter = "isIndex == false"
outcome_close_verify = ["go test ./..."]
pr_merge_strategy = "merge"
commit_style = "conventional"
auto_push = false
required_code_coverage = 91.5
loop_max_tasks = 7
parallel = false
autonomy = "manual"
compact_after_task_commit = false
pr_mode = true

[status_values]
in_progress = "Doing"
completed = "Done"
`)

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.RoadmapRoot != filepath.Join(repo, "docs", "roadmap") {
		t.Fatalf("RoadmapRoot = %q", loaded.RoadmapRoot)
	}
	if loaded.ConfigPath != filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml") {
		t.Fatalf("ConfigPath = %q", loaded.ConfigPath)
	}
	if got := loaded.DoneStatuses; len(got) != 1 || got[0] != "Done" {
		t.Fatalf("DoneStatuses = %#v", got)
	}
	if loaded.StatusValues.InProgress != "Doing" || loaded.StatusValues.Completed != "Done" || loaded.StatusValues.Pending != "Pending" {
		t.Fatalf("StatusValues = %#v", loaded.StatusValues)
	}
	if loaded.AutoPush {
		t.Fatal("AutoPush = true, want false")
	}
	if loaded.RequiredCodeCoverage != 91.5 {
		t.Fatalf("RequiredCodeCoverage = %v, want 91.5", loaded.RequiredCodeCoverage)
	}
	if loaded.LoopMaxTasks != 7 || loaded.Parallel || loaded.Autonomy != "manual" || loaded.CompactAfterTaskCommit || !loaded.PRMode {
		t.Fatalf("execution settings = max:%d parallel:%t autonomy:%q compact:%t pr:%t", loaded.LoopMaxTasks, loaded.Parallel, loaded.Autonomy, loaded.CompactAfterTaskCommit, loaded.PRMode)
	}
}

func TestLoadFixtureValidRoadmapctlTOMLDefault(t *testing.T) {
	loaded, err := Load(filepath.Join("..", "..", "testdata", "fixtures", "valid-roadmapctl-toml-default"), Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.RoadmapRootRel != "docs/roadmap" {
		t.Fatalf("RoadmapRootRel = %q", loaded.RoadmapRootRel)
	}
	if filepath.Base(loaded.ConfigPath) != ".roadmapctl.toml" {
		t.Fatalf("ConfigPath = %q", loaded.ConfigPath)
	}
}

func TestLoadUsesRoadmapRootOverrideForTOMLDiscovery(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("custom", "roadmap"), `active_statuses = ["Queued"]
`)

	loaded, err := Load(repo, Options{RoadmapRoot: "custom/roadmap"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.RoadmapRootRel != "custom/roadmap" {
		t.Fatalf("RoadmapRootRel = %q", loaded.RoadmapRootRel)
	}
	if got := loaded.ActiveStatuses; len(got) != 1 || got[0] != "Queued" {
		t.Fatalf("ActiveStatuses = %#v", got)
	}
}

func TestLoadUsesDefaultsWhenTOMLMissingButRoadmapRootExists(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "docs", "roadmap"), 0o755); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.ConfigPath != filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml") {
		t.Fatalf("ConfigPath = %q", loaded.ConfigPath)
	}
	if loaded.RoadmapRootRel != "docs/roadmap" || loaded.StatusValues.Completed != "Completed" {
		t.Fatalf("loaded = %#v", loaded)
	}
	if loaded.RequiredCodeCoverage != 85.0 {
		t.Fatalf("RequiredCodeCoverage = %v, want 85.0", loaded.RequiredCodeCoverage)
	}
}

func TestLoadTOMLParseErrorIsUsageError(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `done_statuses = ["Completed"
`)

	_, err := Load(repo, Options{})
	if err == nil {
		t.Fatal("Load() error = nil, want TOML parse error")
	}
	var cfgErr *Error
	if !errors.As(err, &cfgErr) {
		t.Fatalf("Load() error type = %T, want *Error", err)
	}
	if cfgErr.Code != ErrConfigParse || cfgErr.ExitCode != 2 {
		t.Fatalf("error = %#v", cfgErr)
	}
}

func TestLoadLegacyOnlyMigratesToTOMLAndDeletesLegacy(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, `roadmap-root: docs/roadmap
done-statuses: ['Done']
auto-push: false
`)

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	tomlPath := filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml")
	if loaded.ConfigPath != tomlPath || loaded.RoadmapRootRel != "docs/roadmap" {
		t.Fatalf("loaded = %#v", loaded)
	}
	if _, err := os.Stat(tomlPath); err != nil {
		t.Fatalf("migrated TOML missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(repo, ".claude", "roadmap.local.md")); !os.IsNotExist(err) {
		t.Fatalf("legacy config still exists after migration: %v", err)
	}
	if loaded.DoneStatuses[0] != "Done" || loaded.AutoPush || loaded.LoopMaxTasks != 0 || !loaded.Parallel || loaded.Autonomy != "until_done" || !loaded.CompactAfterTaskCommit || loaded.PRMode {
		t.Fatalf("loaded config = %#v", loaded)
	}
}

func TestLoadExistingTOMLDeletesLegacyWithoutConflictWarning(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, "roadmap-root: docs/roadmap\ndone-statuses: ['Done']\n")
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `done_statuses = ["Completed"]
`)

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(loaded.Warnings) != 0 {
		t.Fatalf("Warnings = %#v", loaded.Warnings)
	}
	if _, err := os.Stat(filepath.Join(repo, ".claude", "roadmap.local.md")); !os.IsNotExist(err) {
		t.Fatalf("legacy config still exists after TOML load: %v", err)
	}
}

func TestLoadInvalidTOMLDoesNotFallbackToLegacy(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, "roadmap-root: docs/roadmap\n")
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `done_statuses = ["Completed"
`)

	_, err := Load(repo, Options{})
	if err == nil {
		t.Fatal("Load() error = nil, want TOML parse error")
	}
	var cfgErr *Error
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrConfigParse {
		t.Fatalf("Load() error = %#v, want RMC_CONFIG_PARSE", err)
	}
	if _, statErr := os.Stat(filepath.Join(repo, ".claude", "roadmap.local.md")); statErr != nil {
		t.Fatalf("legacy config should remain after invalid TOML: %v", statErr)
	}
}

func TestLoadRejectsInvalidExecutionSettings(t *testing.T) {
	t.Run("invalid autonomy", func(t *testing.T) {
		repo := t.TempDir()
		writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `autonomy = "robot"
`)

		_, err := Load(repo, Options{})
		if err == nil {
			t.Fatal("Load() error = nil, want validation error")
		}
		var cfgErr *Error
		if !errors.As(err, &cfgErr) || cfgErr.Code != ErrConfigParse {
			t.Fatalf("Load() error = %#v, want RMC_CONFIG_PARSE", err)
		}
	})

	t.Run("negative loop max tasks", func(t *testing.T) {
		repo := t.TempDir()
		writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `loop_max_tasks = -1
`)

		_, err := Load(repo, Options{})
		if err == nil {
			t.Fatal("Load() error = nil, want validation error")
		}
		var cfgErr *Error
		if !errors.As(err, &cfgErr) || cfgErr.Code != ErrConfigParse {
			t.Fatalf("Load() error = %#v, want RMC_CONFIG_PARSE", err)
		}
	})

	for _, value := range []string{"-0.1", "100.1"} {
		t.Run("required code coverage out of range "+value, func(t *testing.T) {
			repo := t.TempDir()
			writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), "required_code_coverage = "+value+"\n")

			_, err := Load(repo, Options{})
			if err == nil {
				t.Fatal("Load() error = nil, want validation error")
			}
			var cfgErr *Error
			if !errors.As(err, &cfgErr) || cfgErr.Code != ErrConfigParse {
				t.Fatalf("Load() error = %#v, want RMC_CONFIG_PARSE", err)
			}
		})
	}
}

func TestLegacyMigrationPlanGeneratesTOMLWithoutWriting(t *testing.T) {
	repo := t.TempDir()
	writeConfig(t, repo, `roadmap-root: docs/roadmap
done-statuses: ['Done', 'Archived']
active-statuses: ['Ready']
status-values:
  completed: Done
auto-push: false
`)

	plan, err := LegacyMigrationPlan(repo, Options{})
	if err != nil {
		t.Fatalf("LegacyMigrationPlan() error = %v", err)
	}
	if plan.TargetPath != filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml") {
		t.Fatalf("TargetPath = %q", plan.TargetPath)
	}
	for _, want := range []string{`done_statuses = ['Done', 'Archived']`, `active_statuses = ['Ready']`, `completed = 'Done'`, `auto_push = false`, `required_code_coverage = 85.0`} {
		if !strings.Contains(plan.Content, want) {
			t.Fatalf("migration content missing %q:\n%s", want, plan.Content)
		}
	}
	if _, err := os.Stat(plan.TargetPath); !os.IsNotExist(err) {
		t.Fatalf("migration wrote target unexpectedly: %v", err)
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
	if loaded.ConfigPath != filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml") {
		t.Fatalf("ConfigPath = %q", loaded.ConfigPath)
	}
	if _, err := os.Stat(filepath.Join(repo, ".claude", "roadmap.local.md")); !os.IsNotExist(err) {
		t.Fatalf("legacy config still exists after migration: %v", err)
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

func writeRoadmapctlTOML(t *testing.T, repo string, roadmapRoot string, body string) {
	t.Helper()
	dir := filepath.Join(repo, roadmapRoot)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".roadmapctl.toml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestFieldsConfigDefaults(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "docs", "roadmap"), 0o755); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Fields.Lifecycle != "estado" {
		t.Fatalf("Fields.Lifecycle = %q, want estado", loaded.Fields.Lifecycle)
	}
	if loaded.Fields.RecordType != "tipo" {
		t.Fatalf("Fields.RecordType = %q, want tipo", loaded.Fields.RecordType)
	}
	if loaded.Fields.TaskValue != "task" {
		t.Fatalf("Fields.TaskValue = %q, want task", loaded.Fields.TaskValue)
	}
	if loaded.Fields.OutcomeValue != "outcome" {
		t.Fatalf("Fields.OutcomeValue = %q, want outcome", loaded.Fields.OutcomeValue)
	}
	if loaded.Fields.DisplayName != "titulo" {
		t.Fatalf("Fields.DisplayName = %q, want titulo", loaded.Fields.DisplayName)
	}
	if loaded.Fields.DependencyLink != "blocked_by" {
		t.Fatalf("Fields.DependencyLink = %q, want blocked_by", loaded.Fields.DependencyLink)
	}
}

func TestFieldsConfigOverride(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `[fields]
lifecycle = "status"
record_type = "kind"
task_value = "tarea"
outcome_value = "resultado"
display_name = "nombre"
dependency_link = "depends_on"
`)

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Fields.Lifecycle != "status" {
		t.Fatalf("Fields.Lifecycle = %q, want status", loaded.Fields.Lifecycle)
	}
	if loaded.Fields.RecordType != "kind" {
		t.Fatalf("Fields.RecordType = %q, want kind", loaded.Fields.RecordType)
	}
	if loaded.Fields.TaskValue != "tarea" {
		t.Fatalf("Fields.TaskValue = %q, want tarea", loaded.Fields.TaskValue)
	}
	if loaded.Fields.OutcomeValue != "resultado" {
		t.Fatalf("Fields.OutcomeValue = %q, want resultado", loaded.Fields.OutcomeValue)
	}
	if loaded.Fields.DisplayName != "nombre" {
		t.Fatalf("Fields.DisplayName = %q, want nombre", loaded.Fields.DisplayName)
	}
	if loaded.Fields.DependencyLink != "depends_on" {
		t.Fatalf("Fields.DependencyLink = %q, want depends_on", loaded.Fields.DependencyLink)
	}
}

func TestFieldsConfigPartialOverride(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `[fields]
lifecycle = "custom_status"
`)

	loaded, err := Load(repo, Options{})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Fields.Lifecycle != "custom_status" {
		t.Fatalf("Fields.Lifecycle = %q, want custom_status", loaded.Fields.Lifecycle)
	}
	if loaded.Fields.RecordType != "tipo" {
		t.Fatalf("Fields.RecordType = %q, want default tipo", loaded.Fields.RecordType)
	}
	if loaded.Fields.DependencyLink != "blocked_by" {
		t.Fatalf("Fields.DependencyLink = %q, want default blocked_by", loaded.Fields.DependencyLink)
	}
}
