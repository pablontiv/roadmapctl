# Bootstrap + Gitflow Config Redesign — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the `pr_mode` bool + free-text `[gitflow]` fields with structured enums, move `commit_style`/`auto_push` under `[gitflow]`, make `bootstrap` base command read-only, and update the integrate skill to consume explicit enum fields.

**Architecture:** Config layer changes first (struct → validation → migration), then bootstrap command and JSON output, then the DefaultRoadmapctlTOML template, then skill updates. Each layer compiles and tests green before the next is touched.

**Tech Stack:** Go 1.23, go-toml v2, cobra, picokit/diag, testify-free Go stdlib tests.

**Spec:** `docs/superpowers/specs/2026-05-22-bootstrap-gitflow-config-redesign.md`

---

## File Structure

| File | Change |
|------|--------|
| `internal/config/config.go` | Extend `GitflowConfig`, remove top-level `CommitStyle`/`AutoPush`/`PRMode`, add validation + migration, remove dead `configDiffers` |
| `internal/config/config_test.go` | Update/add tests for new fields, migration, validation |
| `internal/cli/bootstrap.go` | New `bootstrapGitflowReport` struct, remove flat fields from `bootstrapConfigReport`, remove auto-create and stem-repair from base command, add `--yes` to `init`, move stem repair into `init --apply` |
| `internal/cli/bootstrap_test.go` | Update tests: base command is read-only, new JSON shape, template content |
| `internal/cli/bootstrap_stem_repair_test.go` | Update stem repair tests to use `bootstrap init --apply --yes` |
| `internal/templates/bootstrap.go` | Update `DefaultRoadmapctlTOML` to use `[gitflow]` section |
| `internal/templates/bootstrap_test.go` | Update template content assertions |
| `testdata/golden/bootstrap-valid-legacy-config-fallback.json` | Update golden to new JSON shape |
| `.claude/skills/integrate/SKILL.md` | Remove free-text heuristic, add enum-based paths, add `pr_create=manual` |
| `.claude/skills/roadmap/bootstrap-subcommand.md` | Update wizard field list |
| `.claude/skills/roadmap/bootstrap-reference.md` | Update field refs, add new arrays |
| `.claude/skills/roadmap/config-reference.md` | Update TOML schema docs |

---

### Task 1: Extend GitflowConfig with new fields + migration + validation

**Files:**
- Modify: `internal/config/config.go`
- Modify: `internal/config/config_test.go`

- [ ] **Step 1: Write failing tests for new config fields**

Add these tests to `internal/config/config_test.go` (replace `TestGitflowConfigParsesFields` and add new ones):

```go
func TestGitflowConfigParsesNewFields(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `
[gitflow]
base_branch = "main"
branch_mode = "scope_branch"
branch_style = "feat/{scope}"
pr_create = "manual"
pr_title_style = "feat(scope): title"
pr_body_style = "## Summary"
commit_style = "conventional"
auto_push = false
`)
	loaded, err := Load(repo)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	g := loaded.Gitflow
	if g.BaseBranch != "main" {
		t.Fatalf("BaseBranch = %q, want main", g.BaseBranch)
	}
	if g.BranchMode != "scope_branch" {
		t.Fatalf("BranchMode = %q, want scope_branch", g.BranchMode)
	}
	if g.BranchStyle != "feat/{scope}" {
		t.Fatalf("BranchStyle = %q, want feat/{scope}", g.BranchStyle)
	}
	if g.PrCreate != "manual" {
		t.Fatalf("PrCreate = %q, want manual", g.PrCreate)
	}
	if g.PrTitleStyle != "feat(scope): title" {
		t.Fatalf("PrTitleStyle = %q, want feat(scope): title", g.PrTitleStyle)
	}
	if g.PrBodyStyle != "## Summary" {
		t.Fatalf("PrBodyStyle = %q, want ## Summary", g.PrBodyStyle)
	}
	if g.CommitStyle != "conventional" {
		t.Fatalf("CommitStyle = %q, want conventional", g.CommitStyle)
	}
	if g.AutoPush {
		t.Fatal("AutoPush = true, want false")
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
	if loaded.Gitflow.CommitStyle != "conventional" {
		t.Fatalf("Gitflow.CommitStyle default = %q, want conventional", loaded.Gitflow.CommitStyle)
	}
	if !loaded.Gitflow.AutoPush {
		t.Fatal("Gitflow.AutoPush default = false, want true")
	}
	if loaded.Gitflow.BranchMode != "direct_push" {
		t.Fatalf("Gitflow.BranchMode default = %q, want direct_push", loaded.Gitflow.BranchMode)
	}
	if loaded.Gitflow.PrCreate != "never" {
		t.Fatalf("Gitflow.PrCreate default = %q, want never", loaded.Gitflow.PrCreate)
	}
}

func TestGitflowMigrationFromPRModeTrue(t *testing.T) {
	repo := t.TempDir()
	writeRoadmapctlTOML(t, repo, filepath.Join("docs", "roadmap"), `pr_mode = true
commit_style = "conventional"
auto_push = false
`)
	loaded, err := Load(repo)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Gitflow.BranchMode != "scope_branch" {
		t.Fatalf("migrated BranchMode = %q, want scope_branch", loaded.Gitflow.BranchMode)
	}
	if loaded.Gitflow.PrCreate != "auto" {
		t.Fatalf("migrated PrCreate = %q, want auto", loaded.Gitflow.PrCreate)
	}
	if loaded.Gitflow.CommitStyle != "conventional" {
		t.Fatalf("migrated CommitStyle = %q, want conventional", loaded.Gitflow.CommitStyle)
	}
	if loaded.Gitflow.AutoPush {
		t.Fatal("migrated AutoPush = true, want false")
	}
	// migration emits warning
	found := false
	for _, w := range loaded.Warnings {
		if w.Code == "RMC_CONFIG_DEPRECATED_TOPLEVEL" {
			found = true
		}
	}
	if !found {
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
		t.Fatalf("migrated BranchMode = %q, want direct_push", loaded.Gitflow.BranchMode)
	}
	if loaded.Gitflow.PrCreate != "never" {
		t.Fatalf("migrated PrCreate = %q, want never", loaded.Gitflow.PrCreate)
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
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /home/shared/roadmapctl && go test ./internal/config/... -run "TestGitflowConfig|TestGitflowMigration|TestGitflowValidation|TestGitflowDefaults" -v 2>&1 | tail -20
```

Expected: multiple FAIL lines — new fields don't exist yet.

- [ ] **Step 3: Update `GitflowConfig` and `tomlGitflowConfig` in `internal/config/config.go`**

Replace the `GitflowConfig` struct (currently has 3 fields) with:

```go
type GitflowConfig struct {
	BaseBranch   string `json:"base_branch"`
	BranchMode   string `json:"branch_mode"`
	BranchStyle  string `json:"branch_style"`
	PrCreate     string `json:"pr_create"`
	PrTitleStyle string `json:"pr_title_style"`
	PrBodyStyle  string `json:"pr_body_style"`
	CommitStyle  string `json:"commit_style"`
	AutoPush     bool   `json:"auto_push"`
}
```

Replace `tomlGitflowConfig` (currently has 3 fields) with:

```go
type tomlGitflowConfig struct {
	BaseBranch   string `toml:"base_branch"`
	BranchMode   string `toml:"branch_mode"`
	BranchStyle  string `toml:"branch_style"`
	PrCreate     string `toml:"pr_create"`
	PrTitleStyle string `toml:"pr_title_style"`
	PrBodyStyle  string `toml:"pr_body_style"`
	CommitStyle  string `toml:"commit_style"`
	AutoPush     *bool  `toml:"auto_push"`
}
```

- [ ] **Step 4: Update `applyTOMLConfig` in `internal/config/config.go`**

Replace the three gitflow assignment lines at the bottom of `applyTOMLConfig` with:

```go
	// New [gitflow] fields
	if decoded.Gitflow.BaseBranch != "" {
		cfg.Gitflow.BaseBranch = decoded.Gitflow.BaseBranch
	}
	if decoded.Gitflow.BranchMode != "" {
		cfg.Gitflow.BranchMode = decoded.Gitflow.BranchMode
	}
	if decoded.Gitflow.BranchStyle != "" {
		cfg.Gitflow.BranchStyle = decoded.Gitflow.BranchStyle
	}
	if decoded.Gitflow.PrCreate != "" {
		cfg.Gitflow.PrCreate = decoded.Gitflow.PrCreate
	}
	if decoded.Gitflow.PrTitleStyle != "" {
		cfg.Gitflow.PrTitleStyle = decoded.Gitflow.PrTitleStyle
	}
	if decoded.Gitflow.PrBodyStyle != "" {
		cfg.Gitflow.PrBodyStyle = decoded.Gitflow.PrBodyStyle
	}

	// CommitStyle: prefer [gitflow], fall back to deprecated top-level
	if decoded.Gitflow.CommitStyle != "" {
		cfg.Gitflow.CommitStyle = decoded.Gitflow.CommitStyle
	} else if decoded.CommitStyle != "" {
		cfg.Gitflow.CommitStyle = decoded.CommitStyle
		cfg.Warnings = append(cfg.Warnings, Warning{
			Code:    "RMC_CONFIG_DEPRECATED_TOPLEVEL",
			Message: "commit_style at top-level is deprecated; move to [gitflow].commit_style",
			Path:    path,
		})
	}

	// AutoPush: prefer [gitflow], fall back to deprecated top-level
	if decoded.Gitflow.AutoPush != nil {
		cfg.Gitflow.AutoPush = *decoded.Gitflow.AutoPush
	} else if decoded.AutoPush != nil {
		cfg.Gitflow.AutoPush = *decoded.AutoPush
		cfg.Warnings = append(cfg.Warnings, Warning{
			Code:    "RMC_CONFIG_DEPRECATED_TOPLEVEL",
			Message: "auto_push at top-level is deprecated; move to [gitflow].auto_push",
			Path:    path,
		})
	}

	// pr_mode → branch_mode + pr_create (only when new fields are absent)
	if decoded.PRMode != nil && decoded.Gitflow.BranchMode == "" {
		if *decoded.PRMode {
			cfg.Gitflow.BranchMode = "scope_branch"
			cfg.Gitflow.PrCreate = "auto"
		} else {
			cfg.Gitflow.BranchMode = "direct_push"
			cfg.Gitflow.PrCreate = "never"
		}
		cfg.Warnings = append(cfg.Warnings, Warning{
			Code:    "RMC_CONFIG_DEPRECATED_TOPLEVEL",
			Message: "pr_mode is deprecated; use [gitflow].branch_mode and [gitflow].pr_create",
			Path:    path,
		})
	}
```

- [ ] **Step 5: Update `defaultConfig` in `internal/config/config.go`**

Remove `CommitStyle: "conventional"`, `AutoPush: true`, `PRMode: false` from the top-level struct literal. Add them under `Gitflow`:

```go
func defaultConfig(repo string) *Config {
	autoPushDefault := true
	return &Config{
		RepoRoot:       repo,
		DoneStatuses:   []string{"Completed", "Obsolete"},
		ActiveStatuses: []string{"Pending", "Specified", "In Progress"},
		StatusValues: StatusValues{
			Pending:    "Pending",
			Specified:  "Specified",
			InProgress: "In Progress",
			Completed:  "Completed",
			Blocked:    "Blocked",
			Obsolete:   "Obsolete",
		},
		Fields: FieldsConfig{
			Lifecycle:      "estado",
			RecordType:     "tipo",
			TaskValue:      "task",
			OutcomeValue:   "outcome",
			DisplayName:    "titulo",
			DependencyLink: "blocked_by",
		},
		LeafFilter:             "isIndex == false",
		OutcomeCloseVerify:     []string{},
		RequiredCodeCoverage:   85.0,
		LoopMaxTasks:           0,
		Parallel:               true,
		Autonomy:               "until_done",
		CompactAfterTaskCommit: true,
		Gitflow: GitflowConfig{
			BranchMode:  "direct_push",
			PrCreate:    "never",
			CommitStyle: "conventional",
			AutoPush:    autoPushDefault,
		},
	}
}
```

- [ ] **Step 6: Remove top-level `CommitStyle`, `AutoPush`, `PRMode` from the `Config` struct and remove dead `configDiffers`**

In `internal/config/config.go`, remove these three fields from the `Config` struct definition:

```go
// Remove these three lines:
CommitStyle          string
AutoPush             bool
PRMode                 bool
```

Also remove the entire `configDiffers` function (it has zero callers — confirmed by cartyx).

- [ ] **Step 7: Add gitflow validation to `validateConfig` in `internal/config/config.go`**

Add before the final `return nil`:

```go
	// Gitflow enum validation (only when explicitly set)
	if cfg.Gitflow.BranchMode != "" {
		switch cfg.Gitflow.BranchMode {
		case "direct_push", "scope_branch":
		default:
			return &Error{Code: ErrConfigParse, Message: "gitflow.branch_mode must be direct_push or scope_branch", Path: path, ExitCode: 2}
		}
	}
	if cfg.Gitflow.PrCreate != "" {
		switch cfg.Gitflow.PrCreate {
		case "never", "manual", "auto":
		default:
			return &Error{Code: ErrConfigParse, Message: "gitflow.pr_create must be never, manual, or auto", Path: path, ExitCode: 2}
		}
		if cfg.Gitflow.PrCreate == "auto" && cfg.Gitflow.BranchMode != "scope_branch" {
			return &Error{Code: ErrConfigParse, Message: "gitflow.pr_create=auto requires gitflow.branch_mode=scope_branch", Path: path, ExitCode: 2}
		}
	}
```

- [ ] **Step 8: Update `TestLoadPrefersRoadmapctlTOMLAndInfersRoadmapRoot` in `internal/config/config_test.go`**

The test currently references `loaded.AutoPush` and `loaded.PRMode` directly. Update the assertion on line 58 and 64:

```go
	// Replace lines 58-66 with:
	if loaded.Gitflow.AutoPush {
		t.Fatal("Gitflow.AutoPush = true, want false")
	}
	if loaded.LoopMaxTasks != 7 || loaded.Parallel || loaded.Autonomy != "manual" || loaded.CompactAfterTaskCommit {
		t.Fatalf("execution settings = max:%d parallel:%t autonomy:%q compact:%t", loaded.LoopMaxTasks, loaded.Parallel, loaded.Autonomy, loaded.CompactAfterTaskCommit)
	}
	// pr_mode=true in TOML should migrate to scope_branch + auto
	if loaded.Gitflow.BranchMode != "scope_branch" {
		t.Fatalf("migrated BranchMode = %q, want scope_branch", loaded.Gitflow.BranchMode)
	}
	if loaded.Gitflow.PrCreate != "auto" {
		t.Fatalf("migrated PrCreate = %q, want auto", loaded.Gitflow.PrCreate)
	}
```

- [ ] **Step 9: Fix compile errors in `internal/cli/bootstrap.go` caused by removed top-level fields**

In `bootstrap.go`, `newBootstrapConfigReport` currently assigns:
```go
result.CommitStyle = cfg.CommitStyle   // line 378
result.AutoPush = cfg.AutoPush         // line 379
result.PRMode = cfg.PRMode             // line 385
result.BranchStyle = cfg.Gitflow.BranchStyle
result.PRTitleStyle = cfg.Gitflow.PrTitleStyle
result.PRBodyStyle = cfg.Gitflow.PrBodyStyle
```

These will fail to compile. For now, just delete those 6 assignment lines (the bootstrap JSON output will be updated in Task 2). The `bootstrapConfigReport` struct fields `CommitStyle`, `AutoPush`, `PRMode`, `BranchStyle`, `PRTitleStyle`, `PRBodyStyle` also exist — leave them for now (they'll be replaced in Task 2).

- [ ] **Step 10: Run all config + CLI tests to verify they compile and pass**

```bash
cd /home/shared/roadmapctl && go test ./internal/config/... ./internal/cli/... -count=1 2>&1 | tail -30
```

Expected: all tests pass. If any tests reference the removed top-level fields, fix them now.

- [ ] **Step 11: Commit**

```bash
git -C /home/shared/roadmapctl add internal/config/config.go internal/config/config_test.go internal/cli/bootstrap.go
git -C /home/shared/roadmapctl commit -m "$(cat <<'EOF'
feat(config): add structured [gitflow] fields, migrate pr_mode/commit_style/auto_push
EOF
)"
```

---

### Task 2: Bootstrap JSON output — new `gitflow` nested struct

**Files:**
- Modify: `internal/cli/bootstrap.go`
- Modify: `internal/cli/bootstrap_test.go`
- Modify: `testdata/golden/bootstrap-valid-legacy-config-fallback.json`

- [ ] **Step 1: Write a failing test for the new JSON shape**

Add to `internal/cli/bootstrap_test.go`:

```go
func TestBootstrapJSONIncludesGitflowObject(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	// Initialize bootstrap files first
	var stdout, stderr bytes.Buffer
	Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "json"}, &stdout, &stderr, "dev")

	stdout.Reset()
	stderr.Reset()
	code := Execute([]string{"bootstrap", "--repo", repo, "--output", "json"}, &stdout, &stderr, "dev")
	if code != 0 {
		t.Fatalf("bootstrap exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report struct {
		Gitflow struct {
			BranchMode  string `json:"branch_mode"`
			PrCreate    string `json:"pr_create"`
			CommitStyle string `json:"commit_style"`
			AutoPush    bool   `json:"auto_push"`
		} `json:"gitflow"`
		MissingSettings []string `json:"missing_settings"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout invalid JSON: %v\n%s", err, stdout.String())
	}
	if report.Gitflow.BranchMode != "direct_push" {
		t.Fatalf("gitflow.branch_mode = %q, want direct_push", report.Gitflow.BranchMode)
	}
	if report.Gitflow.PrCreate != "never" {
		t.Fatalf("gitflow.pr_create = %q, want never", report.Gitflow.PrCreate)
	}
	if report.Gitflow.CommitStyle != "conventional" {
		t.Fatalf("gitflow.commit_style = %q, want conventional", report.Gitflow.CommitStyle)
	}
	if !report.Gitflow.AutoPush {
		t.Fatal("gitflow.auto_push = false, want true")
	}
	// base_branch is empty in template → should appear in missing_settings
	found := false
	for _, s := range report.MissingSettings {
		if s == "gitflow.base_branch" {
			found = true
		}
	}
	if !found {
		t.Fatalf("missing_settings should include gitflow.base_branch, got %v", report.MissingSettings)
	}
}

func TestBootstrapJSONDoesNotIncludeLegacyFlatFields(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	var stdout, stderr bytes.Buffer
	Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--output", "json"}, &stdout, &stderr, "dev")

	stdout.Reset()
	stderr.Reset()
	Execute([]string{"bootstrap", "--repo", repo, "--output", "json"}, &stdout, &stderr, "dev")

	var raw map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &raw); err != nil {
		t.Fatalf("stdout invalid JSON: %v", err)
	}
	for _, legacy := range []string{"pr_mode", "commit_style", "auto_push", "branch_style", "pr_title_style", "pr_body_style"} {
		if _, found := raw[legacy]; found {
			t.Fatalf("bootstrap JSON still contains legacy flat field %q", legacy)
		}
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /home/shared/roadmapctl && go test ./internal/cli/... -run "TestBootstrapJSONIncludesGitflowObject|TestBootstrapJSONDoesNotIncludeLegacyFlatFields" -v 2>&1 | tail -20
```

Expected: FAIL — `gitflow` object not in JSON, legacy flat fields still present.

- [ ] **Step 3: Add `bootstrapGitflowReport` struct and update `bootstrapConfigReport` in `bootstrap.go`**

Add after the `bootstrapConfigReport` struct definition:

```go
type bootstrapGitflowReport struct {
	BaseBranch   string `json:"base_branch"`
	BranchMode   string `json:"branch_mode"`
	BranchStyle  string `json:"branch_style"`
	PrCreate     string `json:"pr_create"`
	PrTitleStyle string `json:"pr_title_style"`
	PrBodyStyle  string `json:"pr_body_style"`
	CommitStyle  string `json:"commit_style"`
	AutoPush     bool   `json:"auto_push"`
}
```

In the `bootstrapConfigReport` struct, remove these fields:

```go
// Remove:
CommitStyle            string               `json:"commit_style"`
AutoPush               bool                 `json:"auto_push"`
PRMode                 bool                 `json:"pr_mode"`
BranchStyle            string               `json:"branch_style,omitempty"`
PRTitleStyle           string               `json:"pr_title_style,omitempty"`
PRBodyStyle            string               `json:"pr_body_style,omitempty"`
```

Add these fields to `bootstrapConfigReport`:

```go
Gitflow          bootstrapGitflowReport `json:"gitflow"`
MissingSettings  []string               `json:"missing_settings,omitempty"`
EmptySettings    []string               `json:"empty_settings,omitempty"`
InvalidSettings  []string               `json:"invalid_settings,omitempty"`
```

- [ ] **Step 4: Add `computeBootstrapGitflowStatus` helper in `bootstrap.go`**

Add this function after `newBootstrapConfigReport`:

```go
func computeBootstrapGitflowStatus(cfg *config.Config) (missing, empty, invalid []string) {
	if cfg.Gitflow.BaseBranch == "" {
		missing = append(missing, "gitflow.base_branch")
	}
	if cfg.Gitflow.BranchMode == "" {
		missing = append(missing, "gitflow.branch_mode")
	}
	if cfg.Gitflow.PrCreate == "" {
		missing = append(missing, "gitflow.pr_create")
	}
	if cfg.Gitflow.BranchMode == "scope_branch" && cfg.Gitflow.BranchStyle == "" {
		empty = append(empty, "gitflow.branch_style")
	}
	if cfg.Gitflow.PrCreate == "auto" {
		if cfg.Gitflow.PrTitleStyle == "" {
			empty = append(empty, "gitflow.pr_title_style")
		}
		if cfg.Gitflow.PrBodyStyle == "" {
			empty = append(empty, "gitflow.pr_body_style")
		}
	}
	return
}
```

- [ ] **Step 5: Update `newBootstrapConfigReport` to populate new fields**

Replace the `if cfg != nil { ... }` block (lines 374–406) with:

```go
	if cfg != nil {
		result.StatusValues = cfg.StatusValues
		result.DoneStatuses = append([]string(nil), cfg.DoneStatuses...)
		result.ActiveStatuses = append([]string(nil), cfg.ActiveStatuses...)
		result.OutcomeCloseVerify = append([]string{}, cfg.OutcomeCloseVerify...)
		result.RequiredCodeCoverage = cfg.RequiredCodeCoverage
		result.LoopMaxTasks = cfg.LoopMaxTasks
		result.Parallel = cfg.Parallel
		result.Autonomy = cfg.Autonomy
		result.CompactAfterTaskCommit = cfg.CompactAfterTaskCommit
		result.Gitflow = bootstrapGitflowReport{
			BaseBranch:   cfg.Gitflow.BaseBranch,
			BranchMode:   cfg.Gitflow.BranchMode,
			BranchStyle:  cfg.Gitflow.BranchStyle,
			PrCreate:     cfg.Gitflow.PrCreate,
			PrTitleStyle: cfg.Gitflow.PrTitleStyle,
			PrBodyStyle:  cfg.Gitflow.PrBodyStyle,
			CommitStyle:  cfg.Gitflow.CommitStyle,
			AutoPush:     cfg.Gitflow.AutoPush,
		}
		result.Helpers = contextHelpers{
			WhereLeaf:    cfg.LeafFilter,
			WhereNotDone: statusWhere("not", cfg.DoneStatuses),
			WhereActive:  statusWhere("", cfg.ActiveStatuses),
		}
		result.MissingSettings, result.EmptySettings, result.InvalidSettings = computeBootstrapGitflowStatus(cfg)
	}
```

Also remove the `RMC_GITFLOW_NOT_CONFIGURED` diagnostic block (currently at lines 396–404 in the old code) since it is replaced by `missing_settings`.

- [ ] **Step 6: Run failing tests to verify they now pass**

```bash
cd /home/shared/roadmapctl && go test ./internal/cli/... -run "TestBootstrapJSONIncludesGitflowObject|TestBootstrapJSONDoesNotIncludeLegacyFlatFields" -v 2>&1 | tail -20
```

Expected: both PASS.

- [ ] **Step 7: Update the golden file**

Run bootstrap on the legacy fixture to regenerate expected output:

```bash
cd /home/shared/roadmapctl && go run ./cmd/roadmapctl bootstrap \
  --repo testdata/fixtures/valid-legacy-config-fallback \
  --output json 2>/dev/null | python3 -m json.tool
```

Copy the output into `testdata/golden/bootstrap-valid-legacy-config-fallback.json`. The new file should have `"gitflow": {...}` with `"branch_mode": "direct_push"` (migrated from `pr_mode = false`), `"missing_settings": ["gitflow.base_branch"]`, and no `pr_mode`/`commit_style`/`auto_push` flat fields. It should have `"diagnostics": [{"id": "RMC_CONFIG_DEPRECATED_TOPLEVEL", ...}]` since the fixture uses legacy top-level fields.

Also update the golden normalization in `golden_test.go` if the rootline version token needs it (check existing pattern — it already handles `<rootline-version>`).

- [ ] **Step 8: Run full test suite**

```bash
cd /home/shared/roadmapctl && go test ./... -count=1 2>&1 | tail -30
```

Expected: all pass. Fix any compilation errors or assertion mismatches before continuing.

- [ ] **Step 9: Commit**

```bash
git -C /home/shared/roadmapctl add internal/cli/bootstrap.go internal/cli/bootstrap_test.go testdata/golden/bootstrap-valid-legacy-config-fallback.json
git -C /home/shared/roadmapctl commit -m "$(cat <<'EOF'
feat(bootstrap): replace flat JSON fields with nested gitflow object + settings arrays
EOF
)"
```

---

### Task 3: Bootstrap base command — fully read-only

**Files:**
- Modify: `internal/cli/bootstrap.go`
- Modify: `internal/cli/bootstrap_test.go`
- Modify: `internal/cli/bootstrap_stem_repair_test.go`

- [ ] **Step 1: Write a failing test confirming base command does not auto-create files**

Add to `internal/cli/bootstrap_test.go`:

```go
func TestBootstrapBaseCommandIsReadOnly(t *testing.T) {
	repo := t.TempDir()
	// Do NOT call init — repo has no roadmap files
	var stdout, stderr bytes.Buffer
	// Should return non-zero (config missing) but NOT create any files
	Execute([]string{"bootstrap", "--repo", repo, "--output", "json"}, &stdout, &stderr, "dev")
	if _, err := os.Stat(filepath.Join(repo, "docs", "roadmap")); !os.IsNotExist(err) {
		t.Fatal("base bootstrap command created docs/roadmap — should be read-only")
	}
	if _, err := os.Stat(filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml")); !os.IsNotExist(err) {
		t.Fatal("base bootstrap command created .roadmapctl.toml — should be read-only")
	}
}
```

- [ ] **Step 2: Run test to verify it fails (base command currently auto-creates files)**

```bash
cd /home/shared/roadmapctl && go test ./internal/cli/... -run "TestBootstrapBaseCommandIsReadOnly" -v 2>&1 | tail -10
```

Expected: FAIL — base command writes files.

- [ ] **Step 3: Remove the auto-create block from `buildBootstrapConfig`**

In `bootstrap.go`, `buildBootstrapConfig`, delete lines 280–295 (the block that calls `proposedBootstrapChanges` and `applyBootstrapChanges`):

```go
// DELETE this entire block:
if len(found) == 0 {
    changes := proposedBootstrapChanges(cfg.RepoRoot, cfg.RoadmapRoot, false, cfg.Fields.DependencyLink)
    if len(changes) > 0 {
        applyErrs := applyBootstrapChanges(cfg.RepoRoot, changes)
        found = append(found, applyErrs...)
        if len(applyErrs) == 0 {
            cfg, err = config.Load(options.Repo)
            if err != nil {
                diagnostic := configDiagnostic(root, err)
                return newBootstrapConfigReport(root, "", "", "", "", nil, []diag.Diagnostic{diagnostic})
            }
        }
    }
}
```

- [ ] **Step 4: Remove the stem-repair call and `--yes` flag from the base bootstrap command**

In `newBootstrapCommand`, delete the `--yes` flag registration and the entire `hasRepairTriggerDiagnostics` block from `RunE`:

```go
// DELETE: cmd.Flags().BoolVar(&yes, "yes", false, "apply .stem repair without interactive prompt")
// DELETE the var declaration: var yes bool
// DELETE from RunE:
if hasRepairTriggerDiagnostics(report.Diagnostics) {
    root, roadmapRoot, _ := bootstrapRoots(*options)
    if root != "" && roadmapRoot != "" {
        applied, extraDiags := repairStemIfNeeded(ctx, *options, root, roadmapRoot, yes, stdin, stderr)
        // ... entire block
    }
}
```

Also remove `stdin io.Reader` from `newBootstrapCommand`'s signature if it is no longer used by the base command (check — it was only used for stem repair prompt).

- [ ] **Step 5: Run test to verify it now passes**

```bash
cd /home/shared/roadmapctl && go test ./internal/cli/... -run "TestBootstrapBaseCommandIsReadOnly" -v 2>&1 | tail -10
```

Expected: PASS.

- [ ] **Step 6: Update `bootstrap_stem_repair_test.go` — tests that used `--yes` on the base command**

Find tests in `bootstrap_stem_repair_test.go` that call `bootstrap --yes`. Update them to use `bootstrap init --apply --yes` instead. For example, a test that was:

```go
code := Execute([]string{"bootstrap", "--repo", repo, "--yes", "--output", "json"}, ...)
```

Becomes:

```go
code := Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--yes", "--output", "json"}, ...)
```

(The `--yes` flag on `bootstrap init` will be added in Task 4. For now, just update the test to show the intended interface, even if the flag isn't wired yet — tests will fail but show the right target.)

- [ ] **Step 7: Run full test suite to check for regressions**

```bash
cd /home/shared/roadmapctl && go test ./... -count=1 2>&1 | tail -30
```

Expected: `TestBootstrapBaseCommandIsReadOnly` passes. Stem repair tests may fail because `--yes` on `bootstrap init` is not yet implemented — that is fine. Fix any unexpected failures.

- [ ] **Step 8: Commit**

```bash
git -C /home/shared/roadmapctl add internal/cli/bootstrap.go internal/cli/bootstrap_test.go internal/cli/bootstrap_stem_repair_test.go
git -C /home/shared/roadmapctl commit -m "$(cat <<'EOF'
feat(bootstrap): make base command fully read-only, remove auto-create and stem repair
EOF
)"
```

---

### Task 4: Bootstrap init — stem repair + `--yes` flag + template update

**Files:**
- Modify: `internal/cli/bootstrap.go`
- Modify: `internal/cli/bootstrap_test.go`
- Modify: `internal/cli/bootstrap_stem_repair_test.go`
- Modify: `internal/templates/bootstrap.go`
- Modify: `internal/templates/bootstrap_test.go`

- [ ] **Step 1: Write failing tests for stem repair via `bootstrap init --apply --yes`**

Add to `internal/cli/bootstrap_stem_repair_test.go` (or update existing test that used `bootstrap --yes`):

```go
func TestBootstrapInitApplyRepairsIncompatibleStem(t *testing.T) {
	requiresRealRootline(t)
	staleStem := `version: 2
scope:
  match: "*.md"
schema:
  estado:
    type: enum
    match: ["O*", "T*"]
    values: [Pending, Specified, In Progress, Completed, Blocked, Obsolete]
  tipo:
    type: enum
    required:
      match: ["O*", "T*"]
    match: ["O*", "T*"]
    values: [outcome, task]
`
	repo := setupRepairRepo(t, staleStem)
	// Write a minimal TOML so config loads
	if err := os.WriteFile(
		filepath.Join(repo, "docs", "roadmap", ".roadmapctl.toml"),
		[]byte("autonomy = \"until_done\"\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--apply", "--yes", "--output", "json"}, &stdout, &stderr, "dev")
	assertRepairExitCode(t, code, 0, &stdout, &stderr)

	// Verify .stem was rewritten to canonical content
	stemBytes, err := os.ReadFile(filepath.Join(repo, "docs", "roadmap", ".stem"))
	if err != nil {
		t.Fatalf("read .stem: %v", err)
	}
	canonical := templates.GenerateStemContent("blocked_by")
	if string(stemBytes) != canonical {
		t.Fatalf(".stem content not canonical after repair:\n%s", string(stemBytes))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
cd /home/shared/roadmapctl && go test ./internal/cli/... -run "TestBootstrapInitApplyRepairs" -v 2>&1 | tail -10
```

Expected: FAIL — `--yes` not recognized on `bootstrap init`.

- [ ] **Step 3: Add `--yes` to `bootstrap init` and move stem repair into `buildBootstrapInit`**

In `newBootstrapInitCommand`, add the `--yes` flag and thread it to the builder:

```go
func newBootstrapInitCommand(options *Options, stdin io.Reader, stdout io.Writer, stderr io.Writer, exitCode *int) *cobra.Command {
	var dryRun bool
	var apply bool
	var yes bool
	cmd := &cobra.Command{Use: "init", Short: "Initialize missing bootstrap files with explicit dry-run or apply.", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true, RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun == apply {
			return fmt.Errorf("bootstrap init requires exactly one of --dry-run or --apply")
		}
		report := buildBootstrapInit(context.Background(), *options, apply, yes, stdin, stderr)
		*exitCode = renderBootstrap(report, options.Output, stdout, stderr)
		return nil
	}}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show proposed bootstrap files without writing")
	cmd.Flags().BoolVar(&apply, "apply", false, "write allowed bootstrap files")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply .stem repair without interactive prompt")
	return cmd
}
```

Update `buildBootstrapInit` signature to accept `yes bool, stdin io.Reader, stderr io.Writer` and add repair logic after file creation:

```go
func buildBootstrapInit(ctx context.Context, options Options, apply bool, yes bool, stdin io.Reader, stderr io.Writer) bootstrapReport {
	root, roadmapRoot, diagnosticsFound := bootstrapRoots(options)
	report := bootstrapReport{Version: 1, Kind: "roadmapctl/bootstrap/init", Root: root, RoadmapRoot: roadmapRoot, Diagnostics: diagnosticsFound}
	if len(diagnosticsFound) == 0 {
		report.Diagnostics = append(report.Diagnostics, bootstrapSchemaCompatibilityDiagnostics(ctx, options, root, roadmapRoot)...)
	}
	if len(report.Diagnostics) == 0 {
		depLink := defaultDependencyLink(options)
		report.Changes = proposedBootstrapChanges(root, roadmapRoot, apply, depLink)
		if apply {
			report.Diagnostics = append(report.Diagnostics, applyBootstrapChanges(root, report.Changes)...)
		}
	}
	// Stem repair: run even if there were schema compat diagnostics, when apply=true
	if apply && hasRepairTriggerDiagnostics(report.Diagnostics) {
		repaired, repairDiags := repairStemIfNeeded(ctx, options, root, roadmapRoot, yes, stdin, stderr)
		if repaired {
			// Clear schema compat diagnostics, re-run postcheck
			report.Diagnostics = repairDiags
		} else {
			report.Diagnostics = append(report.Diagnostics, repairDiags...)
		}
	}
	if apply && len(report.Diagnostics) == 0 {
		postOptions := options
		postOptions.Repo = root
		postcheck := runCheck(ctx, postOptions)
		report.Diagnostics = append(report.Diagnostics, postcheck.Diagnostics...)
	}
	report.Summary = reports.NewReport(report.Kind, root, roadmapRoot, report.Diagnostics).Summary
	return report
}
```

- [ ] **Step 4: Run failing test to verify it now passes**

```bash
cd /home/shared/roadmapctl && go test ./internal/cli/... -run "TestBootstrapInitApplyRepairs" -v 2>&1 | tail -10
```

Expected: PASS.

- [ ] **Step 5: Write failing test for updated template content**

Add to `internal/templates/bootstrap_test.go` (or update existing test):

```go
func TestDefaultRoadmapctlTOMLHasGitflowSection(t *testing.T) {
	toml := DefaultRoadmapctlTOML
	for _, want := range []string{
		"[gitflow]",
		`commit_style = "conventional"`,
		"auto_push = true",
		`branch_mode = "direct_push"`,
		`pr_create = "never"`,
	} {
		if !strings.Contains(toml, want) {
			t.Fatalf("DefaultRoadmapctlTOML missing %q:\n%s", want, toml)
		}
	}
	for _, notWant := range []string{"pr_mode", "commit_style = "} {
		// commit_style should only appear inside [gitflow], not top-level
		// We check that the top-level occurrence is absent
	}
	if strings.Contains(toml[:strings.Index(toml, "[gitflow]")], "commit_style") {
		t.Fatal("commit_style appears before [gitflow] section — should be nested")
	}
	if strings.Contains(toml[:strings.Index(toml, "[gitflow]")], "auto_push") {
		t.Fatal("auto_push appears before [gitflow] section — should be nested")
	}
	if strings.Contains(toml, "pr_mode") {
		t.Fatal("DefaultRoadmapctlTOML should not contain deprecated pr_mode")
	}
}
```

- [ ] **Step 6: Run test to verify it fails**

```bash
cd /home/shared/roadmapctl && go test ./internal/templates/... -run "TestDefaultRoadmapctlTOMLHasGitflowSection" -v 2>&1 | tail -10
```

Expected: FAIL.

- [ ] **Step 7: Update `DefaultRoadmapctlTOML` in `internal/templates/bootstrap.go`**

Replace the `const DefaultRoadmapctlTOML` value with:

```go
const DefaultRoadmapctlTOML = `done_statuses = ["Completed", "Obsolete"]
active_statuses = ["Pending", "Specified", "In Progress"]
leaf_filter = "isIndex == false"
outcome_close_verify = []
required_code_coverage = 85.0
loop_max_tasks = 0
parallel = true
autonomy = "until_done"
compact_after_task_commit = true

[status_values]
pending = "Pending"
specified = "Specified"
in_progress = "In Progress"
completed = "Completed"
blocked = "Blocked"
obsolete = "Obsolete"

[gitflow]
base_branch = ""
branch_mode = "direct_push"
branch_style = ""
pr_create = "never"
pr_title_style = ""
pr_body_style = ""
commit_style = "conventional"
auto_push = true
`
```

- [ ] **Step 8: Update `bootstrap_test.go` assertion for template content (line 63)**

The existing test checks for `"pr_mode = false"` in the dry-run template content. Replace with:

```go
for _, want := range []string{
    "required_code_coverage = 85.0",
    "loop_max_tasks = 0",
    "parallel = true",
    `autonomy = "until_done"`,
    "compact_after_task_commit = true",
    "[gitflow]",
    `branch_mode = "direct_push"`,
    `pr_create = "never"`,
} {
    if !strings.Contains(change.Content, want) {
        t.Fatalf("bootstrap TOML missing %q:\n%s", want, change.Content)
    }
}
```

- [ ] **Step 9: Run full test suite**

```bash
cd /home/shared/roadmapctl && go test ./... -count=1 2>&1 | tail -30
```

Expected: all pass. Fix any remaining assertions that reference `pr_mode`.

- [ ] **Step 10: Commit**

```bash
git -C /home/shared/roadmapctl add internal/cli/bootstrap.go internal/cli/bootstrap_test.go internal/cli/bootstrap_stem_repair_test.go internal/templates/bootstrap.go internal/templates/bootstrap_test.go
git -C /home/shared/roadmapctl commit -m "$(cat <<'EOF'
feat(bootstrap): move stem repair to init --apply, update DefaultRoadmapctlTOML to [gitflow]
EOF
)"
```

---

### Task 5: Update integrate skill

**Files:**
- Modify: `.claude/skills/integrate/SKILL.md`

- [ ] **Step 1: Read the current skill to understand what's changing**

The skill currently has:
- Inputs: `pr_mode`, `branch_style` (used as mode classifier), `commit_style`, `auto_push`, `base_branch` (auto-detected)
- Fase 1: classifies `branch_style` text as `scope-branch` / `direct-push` / `ambiguous`
- Fase 5: creates PR when `effective_pr_mode == true`
- No `pr_create=manual` path

The new skill must:
- Inputs: replace `pr_mode` with `branch_mode` + `pr_create`; `commit_style`/`auto_push`/`base_branch` read from `gitflow.*`
- Fase 1: read enum directly, no text classification
- Fase 5: three sub-paths based on `pr_create` value

- [ ] **Step 2: Update the `argument-hint` front-matter line**

Change from:
```
argument-hint: "task_path=<path> scope=<scope> ... pr_mode=<bool> commit_style=<style> auto_push=<bool> branch_style=<style> ..."
```

To:
```
argument-hint: "task_path=<path> scope=<scope> previous_scope=<scope> repo_path=<path> branch_mode=<direct_push|scope_branch> pr_create=<never|manual|auto> commit_style=<style> auto_push=<bool> branch_style=<template> pr_title_style=<style> pr_body_style=<style> base_branch=<branch> autonomy=<mode> [commit_files=<files>] [commit_message=<msg>] [is_last_in_scope=<bool>]"
```

- [ ] **Step 3: Update the Inputs table**

Replace the inputs table with:

```markdown
| Campo | Tipo | Descripción |
|-------|------|-------------|
| `task_path` | string | Path relativo de la task recién completada |
| `scope` | string | Outcome activo o `direct-tasks` |
| `previous_scope` | string | Scope anterior; vacío si es la primera task del loop |
| `repo_path` | string | Path absoluto al repo |
| `branch_mode` | enum | `direct_push` \| `scope_branch` — leído de `gitflow.branch_mode` |
| `pr_create` | enum | `never` \| `manual` \| `auto` — leído de `gitflow.pr_create` |
| `commit_style` | string | Leído de `gitflow.commit_style` |
| `auto_push` | bool | Leído de `gitflow.auto_push` |
| `base_branch` | string | Leído de `gitflow.base_branch` — no se auto-detecta |
| `branch_style` | string | Template para nombre de branch (e.g. `feat/{scope}`). Solo relevante cuando `branch_mode=scope_branch` |
| `pr_title_style` | string | Template para título de PR. Requerido cuando `pr_create=auto` |
| `pr_body_style` | string | Template para body de PR. Requerido cuando `pr_create=auto` |
| `commit_files[]` | string[] | (opcional) Lista de archivos a `git add` |
| `commit_message` | string | (opcional) Override del mensaje de commit |
| `is_last_in_scope` | bool | (opcional) `true` si es la última task del scope |
```

- [ ] **Step 4: Replace Fase 1 entirely**

Delete the current "Fase 1: Gitflow mode y scope change" section and replace with:

```markdown
## Fase 1: Mode determination y scope change

El modo se lee directamente de los campos enum — no se clasifica texto libre.

1. Derivar `effective_branch_mode` y `effective_pr_create` directamente desde los inputs:
   - `effective_branch_mode = branch_mode` (siempre explícito)
   - `effective_pr_create = pr_create` (siempre explícito)

2. Generar `branch_target`:
   - Si `effective_branch_mode = "direct_push"`: `branch_target = base_branch`
   - Si `effective_branch_mode = "scope_branch"`: LLM genera `branch_target` sustituyendo
     `{scope}` o `<scope>` en `branch_style` con el valor de `scope` (slugificado).
     Si `branch_style` está vacío o no contiene patrón de scope: emitir
     `RMC_INTEGRATE_BRANCH_STYLE_MISSING` y detener.

3. Detectar cambio de scope:
   ```
   scope_changed = (scope != previous_scope && previous_scope != "")
   ```
   Si `scope_changed == true` y `effective_pr_create != "never"`:
   ```bash
   gh pr list --head <previous_branch_target> --state open --json number,url
   ```
   Si existe PR previo abierto: registrar en diagnostics como informativo.
```

- [ ] **Step 5: Update Fase 5 (PR) to handle three `pr_create` values**

Replace the current Fase 5 with:

```markdown
## Fase 5: PR (si `effective_branch_mode = "scope_branch"`)

### pr_create = "never"

No ejecutar ningún comando PR. Continuar directamente a Fase 6 (que no hace nada en este caso).

### pr_create = "manual"

Detectar si ya existe PR abierto para el scope:

```bash
gh pr list --head <branch_target> --state open --json number,url
```

Si no existe, imprimir el comando sugerido (NO ejecutarlo):

```bash
echo "PR sugerido (ejecutar manualmente):"
echo "gh pr create \\"
echo "  --base <base_branch> \\"
echo "  --head <branch_target> \\"
echo "  --title \"<LLM-generated desde pr_title_style>\" \\"
echo "  --body \"<LLM-generated desde pr_body_style>\""
```

Registrar `pr: null` en `INTEGRATE_RESULT` (no se creó PR automáticamente).

### pr_create = "auto"

Detectar si ya existe PR abierto:

```bash
gh pr list --head <branch_target> --state open --json number,url
```

Si no existe, crear:

```bash
gh pr create \
  --base <base_branch> \
  --head <branch_target> \
  --title "<LLM-generated desde pr_title_style>" \
  --body "<LLM-generated desde pr_body_style>"
```

Registrar número de PR en `INTEGRATE_RESULT.pr`.

Si `gh` no está disponible o `gh auth status` falla:
- `manual`: emitir `RMC_INTEGRATE_GH_AUTH` o `RMC_INTEGRATE_NO_GH`, preguntar si continuar sin PR.
- `supervised` / `until_done`: degradar a modo sin PR; advertir; continuar.
```

- [ ] **Step 6: Update Fase 6 (Merge)**

Replace the condition `effective_pr_mode == true && is_last_in_scope == true` with `effective_pr_create == "auto" && is_last_in_scope == true`. No other changes needed.

- [ ] **Step 7: Update the errores comunes table**

Remove `RMC_INTEGRATE_PR_DISABLED_BY_DIRECT_PUSH` (no longer needed — direct_push is explicit).
Add:

```
| `RMC_INTEGRATE_BRANCH_MODE_MISSING` | `branch_mode` no fue provisto | Verificar que bootstrap JSON fue leído y `gitflow.branch_mode` fue pasado al skill |
| `RMC_INTEGRATE_BASE_BRANCH_MISSING` | `base_branch` vacío | Configurar `[gitflow].base_branch` en `.roadmapctl.toml` y re-ejecutar bootstrap |
```

- [ ] **Step 8: Update the headless verification scenarios**

Replace the four existing pi verification scenarios (A–D) with:

```bash
# Scenario A: direct_push
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con branch_mode=direct_push, pr_create=never, base_branch=master, commit_style=conventional, auto_push=true, autonomy=until_done, task=docs/roadmap/T020-x.md, scope=direct-tasks. Listar los comandos que correrías, SIN ejecutar git/gh ni modificar archivos.'

# Scenario B: scope_branch + pr_create=auto
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con branch_mode=scope_branch, pr_create=auto, base_branch=main, branch_style="feat/{scope}", scope=O22-slug, previous_scope=O21-slug, is_last_in_scope=false. Listar comandos SIN ejecutar ni modificar archivos.'

# Scenario C: scope_branch + pr_create=manual
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con branch_mode=scope_branch, pr_create=manual, base_branch=master, branch_style="feat/{scope}", scope=O22-slug, previous_scope="", is_last_in_scope=false. Listar comandos SIN ejecutar ni modificar archivos. Verificar que NO ejecuta gh pr create sino que imprime el comando sugerido.'

# Scenario D: scope_branch + branch_style vacío
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con branch_mode=scope_branch, pr_create=auto, base_branch=master, branch_style="", scope=O22-slug. Listar comandos SIN ejecutar. Verificar que emite RMC_INTEGRATE_BRANCH_STYLE_MISSING y detiene.'
```

Expected outcomes:
- Scenario A: `git fetch`, `git checkout master`, `git pull --ff-only`, `git add`, `git commit`, `git push -u origin master`; NO `checkout -B`, NO `gh pr create`; `INTEGRATE_RESULT.pr = null`
- Scenario B: `git checkout -B feat/O22-slug`, `gh pr create`; PR registrado
- Scenario C: branch creado, push, imprime `gh pr create ...` sin ejecutar; `INTEGRATE_RESULT.pr = null`
- Scenario D: `RMC_INTEGRATE_BRANCH_STYLE_MISSING`, no branch ni PR

- [ ] **Step 9: Commit**

```bash
git -C /home/shared/roadmapctl add .claude/skills/integrate/SKILL.md
git -C /home/shared/roadmapctl commit -m "$(cat <<'EOF'
feat(integrate): replace free-text heuristic with explicit branch_mode/pr_create enums
EOF
)"
```

---

### Task 6: Update roadmap skills and sync

**Files:**
- Modify: `.claude/skills/roadmap/bootstrap-subcommand.md`
- Modify: `.claude/skills/roadmap/bootstrap-reference.md`
- Modify: `.claude/skills/roadmap/config-reference.md`

- [ ] **Step 1: Update `bootstrap-reference.md` — field references**

In the "Fuente primaria de contexto" section, replace the `status/config operacional` list item:

Old:
```
- status/config operacional = campos `status_values`, `done_statuses`, `active_statuses`, ..., `pr_mode` y cualquier campo adicional expuesto
```

New:
```
- status/config operacional = campos `status_values`, `done_statuses`, `active_statuses`, `outcome_close_verify`, `required_code_coverage`, `loop_max_tasks`, `parallel`, `autonomy`, `compact_after_task_commit` y el objeto `gitflow` con `branch_mode`, `pr_create`, `commit_style`, `auto_push`, `base_branch`, `branch_style`, `pr_title_style`, `pr_body_style`
- `missing_settings`, `empty_settings`, `invalid_settings` — arrays de dot-path keys con campos no configurados, vacíos o inválidos respectivamente (e.g. `"gitflow.base_branch"`)
```

- [ ] **Step 2: Update `bootstrap-reference.md` — diagnósticos gitflow section**

Replace the `RMC_GITFLOW_NOT_CONFIGURED` diagnostic entry:

Old:
```
- `RMC_GITFLOW_NOT_CONFIGURED` (info) — los campos `[gitflow]` (`branch_style`, `pr_title_style`, `pr_body_style`) están vacíos.
```

New:
```
- `RMC_CONFIG_DEPRECATED_TOPLEVEL` (warning) — campos `commit_style`, `auto_push`, o `pr_mode` están definidos en top-level; migrar a `[gitflow]`.
- Los campos faltantes/vacíos/inválidos se reportan en `missing_settings`, `empty_settings`, `invalid_settings` (no como diagnostics individuales).
```

- [ ] **Step 3: Update `bootstrap-reference.md` — Wizard de adopción gitflow**

Replace the wizard description to reflect the new field order and conditionals:

```markdown
## Wizard de adopción gitflow

Cuando `missing_settings` o `empty_settings` contienen campos de `gitflow`, el skill ejecuta el wizard:

1. **Escaneo abierto del repo**: leer CLAUDE.md, git log, gh pr list, .github/workflows/, etc.
2. **Preguntar campos en orden** (uno por vez, solo si ausentes o vacíos en `missing_settings`/`empty_settings`):
   - `base_branch` — siempre requerido
   - `branch_mode` — menú: `direct_push` | `scope_branch`
   - `branch_style` — solo si `branch_mode=scope_branch`; inferir de historial de branches y pre-llenar
   - `pr_create` — menú: `never` | `manual` | `auto`; `auto` solo disponible con `branch_mode=scope_branch`
   - `pr_title_style`, `pr_body_style` — solo si `pr_create=auto`; inferir de PRs mergeados
3. **Dry-run TOML**: presentar el bloque `[gitflow]` propuesto para confirmación.
4. **Escribir solo tras confirmación**: actualizar `<roadmap-root>/.roadmapctl.toml`.
5. **Re-ejecutar bootstrap**: invocar `roadmapctl bootstrap --repo <repo> --output json` y continuar desde el JSON actualizado.

Campos ya configurados (no en `missing_settings` ni `empty_settings`) se muestran para confirmación, no se re-preguntan, salvo que el usuario invoque `/roadmap bootstrap` explícitamente para forzar re-configuración.
```

- [ ] **Step 4: Update `bootstrap-subcommand.md`**

Replace the entire flujo section to reflect the new wizard and read-only base command:

```markdown
## Flujo

1. Ejecutar `roadmapctl bootstrap --repo <repo> --output json` para detectar el estado actual.
2. Si `missing_settings` o `empty_settings` contienen campos de `gitflow`, ejecutar el wizard de adopción gitflow (ver [bootstrap-reference.md](bootstrap-reference.md) sección "Wizard de adopción gitflow").
3. El escaneo es abierto — el LLM decide qué leer según lo que encuentra en el repo.
4. Presentar dry-run del bloque `[gitflow]` propuesto. Solo escribir tras confirmación explícita del usuario.
5. Tras confirmación, escribir los campos bajo `[gitflow]` en `.roadmapctl.toml`.
6. Re-ejecutar `roadmapctl bootstrap --repo <repo> --output json` para confirmar los valores escritos.
7. **Terminar**: no ejecutar plan/loop/pending/decision tras el bootstrap subcommand.

## Nota

Este subcomando gestiona la configuración gitflow. `roadmapctl bootstrap` es ahora read-only;
las escrituras (crear archivos, reparar `.stem`) se hacen con `roadmapctl bootstrap init --apply`.
```

- [ ] **Step 5: Update `config-reference.md` TOML schema (if this file contains a TOML schema example)**

If `config-reference.md` contains a TOML example with `pr_mode`, `commit_style`, or `auto_push` at top-level, update it to show those fields under `[gitflow]`. Mark the deprecated fields:

```toml
# DEPRECATED (still parsed for migration):
# commit_style = "conventional"
# auto_push = true
# pr_mode = false

[gitflow]
base_branch    = "master"
branch_mode    = "direct_push"    # direct_push | scope_branch
branch_style   = ""               # required when branch_mode=scope_branch
pr_create      = "never"          # never | manual | auto
pr_title_style = ""               # required when pr_create=auto
pr_body_style  = ""               # required when pr_create=auto
commit_style   = "conventional"
auto_push      = true
```

- [ ] **Step 6: Sync skills to global install**

```bash
cd /home/shared/roadmapctl && ./scripts/sync-roadmap-skill.sh --install --skill integrate
./scripts/sync-roadmap-skill.sh --install --skill roadmap
```

Expected: both complete without error.

- [ ] **Step 7: Verify skill sync**

```bash
cd /home/shared/roadmapctl && ./scripts/sync-roadmap-skill.sh --check --skill integrate
./scripts/sync-roadmap-skill.sh --check --skill roadmap
```

Expected: both `--check` commands exit 0.

- [ ] **Step 8: Run full test suite one final time**

```bash
cd /home/shared/roadmapctl && go test ./... -count=1 2>&1 | tail -30
```

Expected: all pass.

- [ ] **Step 9: Commit**

```bash
git -C /home/shared/roadmapctl add .claude/skills/roadmap/bootstrap-subcommand.md .claude/skills/roadmap/bootstrap-reference.md .claude/skills/roadmap/config-reference.md
git -C /home/shared/roadmapctl commit -m "$(cat <<'EOF'
docs(skills): update bootstrap and integrate skills for new [gitflow] schema
EOF
)"
```

---

## Self-Review Notes

- **Spec coverage:** All 5 spec sections covered — config layer (Tasks 1), bootstrap JSON (Task 2), bootstrap read-only + init repair (Tasks 3–4), integrate skill (Task 5), roadmap skills (Task 6). ✓
- **Placeholders:** None. Every step contains actual code or commands. ✓
- **Type consistency:** `GitflowConfig` fields match across all tasks. `bootstrapGitflowReport` mirrors `GitflowConfig` field names. `computeBootstrapGitflowStatus` uses the exact field names from Task 1. ✓
- **configDiffers removal:** Task 1 removes it explicitly. No other tasks reference it. ✓
- **`config_test.go` migration test:** `TestLoadPrefersRoadmapctlTOMLAndInfersRoadmapRoot` uses `pr_mode = true` in its fixture — Task 1 Step 8 updates this to assert `BranchMode=scope_branch`. ✓
- **Golden file:** Task 2 Step 7 regenerates from the running binary. The legacy fixture has `pr_mode = false` → expect `branch_mode: direct_push`, `RMC_CONFIG_DEPRECATED_TOPLEVEL` diagnostic. ✓
