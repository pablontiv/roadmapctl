package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigErrorFormatsPathAndUnwrapsCause(t *testing.T) {
	cause := errors.New("cause")
	err := &Error{Code: ErrConfigParse, Message: "bad config", Path: "docs/roadmap/.roadmapctl.toml", Cause: cause}
	if got := err.Error(); got != "RMC_CONFIG_PARSE: docs/roadmap/.roadmapctl.toml: bad config" {
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

	loaded, err := Load(repo)
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
	if loaded.Gitflow.AutoPush {
		t.Fatal("Gitflow.AutoPush = true, want false")
	}
	if loaded.RequiredCodeCoverage != 91.5 {
		t.Fatalf("RequiredCodeCoverage = %v, want 91.5", loaded.RequiredCodeCoverage)
	}
	if loaded.LoopMaxTasks != 7 || loaded.Parallel || loaded.Autonomy != "manual" || loaded.CompactAfterTaskCommit {
		t.Fatalf("execution settings = max:%d parallel:%t autonomy:%q compact:%t", loaded.LoopMaxTasks, loaded.Parallel, loaded.Autonomy, loaded.CompactAfterTaskCommit)
	}
	if loaded.Gitflow.BranchMode != "scope_branch" || loaded.Gitflow.PrCreate != "auto" {
		t.Fatalf("gitflow pr migration = branch_mode:%q pr_create:%q", loaded.Gitflow.BranchMode, loaded.Gitflow.PrCreate)
	}
}

func TestLoadFixtureValidRoadmapctlTOMLDefault(t *testing.T) {
	loaded, err := Load(filepath.Join("..", "..", "testdata", "fixtures", "valid-roadmapctl-toml-default"))
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

func TestLoadUsesDefaultsWhenTOMLMissingButRoadmapRootExists(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "docs", "roadmap"), 0o755); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(repo)
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

	_, err := Load(repo)
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

func TestLoadRejectsInvalidExecutionSettings(t *testing.T) {
	t.Run("invalid autonomy", func(t *testing.T) {
		repo := t.TempDir()
		writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `autonomy = "robot"
`)

		_, err := Load(repo)
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

		_, err := Load(repo)
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

			_, err := Load(repo)
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

func TestLoadResolvesValidRoadmapRootInsideRepo(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), "")

	loaded, err := Load(repo)
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
}

func TestLoadMissingConfigIsUsageError(t *testing.T) {
	repo := t.TempDir()

	_, err := Load(repo)
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

	loaded, err := Load(repo)
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

	loaded, err := Load(repo)
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

	loaded, err := Load(repo)
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

func TestGitflowConfigParsesNewFields(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `[gitflow]
base_branch = "main"
branch_mode = "scope_branch"
branch_style = "test-branch"
pr_create = "manual"
pr_title_style = "test-title"
pr_body_style = "test-body"
commit_style = "custom-commit"
auto_push = false
`)

	loaded, err := Load(repo)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Gitflow.BaseBranch != "main" {
		t.Fatalf("Gitflow.BaseBranch = %q, want main", loaded.Gitflow.BaseBranch)
	}
	if loaded.Gitflow.BranchMode != "scope_branch" {
		t.Fatalf("Gitflow.BranchMode = %q, want scope_branch", loaded.Gitflow.BranchMode)
	}
	if loaded.Gitflow.BranchStyle != "test-branch" {
		t.Fatalf("Gitflow.BranchStyle = %q, want test-branch", loaded.Gitflow.BranchStyle)
	}
	if loaded.Gitflow.PrCreate != "manual" {
		t.Fatalf("Gitflow.PrCreate = %q, want manual", loaded.Gitflow.PrCreate)
	}
	if loaded.Gitflow.PrTitleStyle != "test-title" {
		t.Fatalf("Gitflow.PrTitleStyle = %q, want test-title", loaded.Gitflow.PrTitleStyle)
	}
	if loaded.Gitflow.PrBodyStyle != "test-body" {
		t.Fatalf("Gitflow.PrBodyStyle = %q, want test-body", loaded.Gitflow.PrBodyStyle)
	}
	if loaded.Gitflow.CommitStyle != "custom-commit" {
		t.Fatalf("Gitflow.CommitStyle = %q, want custom-commit", loaded.Gitflow.CommitStyle)
	}
	if loaded.Gitflow.AutoPush {
		t.Fatal("Gitflow.AutoPush = true, want false")
	}
}

func TestGitflowConfigDefaultsApplied(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "docs", "roadmap"), 0o755); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(repo)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Gitflow.BranchMode != "direct_push" {
		t.Fatalf("Gitflow.BranchMode = %q, want direct_push", loaded.Gitflow.BranchMode)
	}
	if loaded.Gitflow.PrCreate != "never" {
		t.Fatalf("Gitflow.PrCreate = %q, want never", loaded.Gitflow.PrCreate)
	}
	if loaded.Gitflow.CommitStyle != "conventional" {
		t.Fatalf("Gitflow.CommitStyle = %q, want conventional", loaded.Gitflow.CommitStyle)
	}
	if !loaded.Gitflow.AutoPush {
		t.Fatal("Gitflow.AutoPush = false, want true")
	}
}

func TestGitflowMigrationFromPRModeTrue(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `pr_mode = true
`)

	loaded, err := Load(repo)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Gitflow.BranchMode != "scope_branch" {
		t.Fatalf("Gitflow.BranchMode = %q, want scope_branch", loaded.Gitflow.BranchMode)
	}
	if loaded.Gitflow.PrCreate != "auto" {
		t.Fatalf("Gitflow.PrCreate = %q, want auto", loaded.Gitflow.PrCreate)
	}
	if !hasWarning(loaded.Warnings, "RMC_CONFIG_DEPRECATED_TOPLEVEL") {
		t.Fatal("expected RMC_CONFIG_DEPRECATED_TOPLEVEL warning")
	}
}

func TestGitflowMigrationFromPRModeFalse(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `pr_mode = false
`)

	loaded, err := Load(repo)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Gitflow.BranchMode != "direct_push" {
		t.Fatalf("Gitflow.BranchMode = %q, want direct_push", loaded.Gitflow.BranchMode)
	}
	if loaded.Gitflow.PrCreate != "never" {
		t.Fatalf("Gitflow.PrCreate = %q, want never", loaded.Gitflow.PrCreate)
	}
}

func TestGitflowValidationRejectsInvalidBranchMode(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `[gitflow]
branch_mode = "chaos"
`)

	_, err := Load(repo)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}
	var cfgErr *Error
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrConfigParse {
		t.Fatalf("Load() error = %#v, want RMC_CONFIG_PARSE", err)
	}
}

func TestGitflowValidationRejectsInvalidPrCreate(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `[gitflow]
pr_create = "sometimes"
`)

	_, err := Load(repo)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}
	var cfgErr *Error
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrConfigParse {
		t.Fatalf("Load() error = %#v, want RMC_CONFIG_PARSE", err)
	}
}

func TestGitflowValidationRejectsPrCreateAutoWithDirectPush(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `[gitflow]
branch_mode = "direct_push"
pr_create = "auto"
`)

	_, err := Load(repo)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}
	var cfgErr *Error
	if !errors.As(err, &cfgErr) || cfgErr.Code != ErrConfigParse {
		t.Fatalf("Load() error = %#v, want RMC_CONFIG_PARSE", err)
	}
}

func hasWarning(warnings []Warning, code string) bool {
	for _, warning := range warnings {
		if warning.Code == code {
			return true
		}
	}
	return false
}
