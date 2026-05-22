# Bootstrap + Gitflow Config Redesign — Spec

## Overview

Refactor `roadmapctl`'s gitflow configuration from a mix of top-level fields and free-text style strings into a structured `[gitflow]` TOML section with explicit enums. Make `bootstrap` base command read-only. Update the integrate skill to consume the new enum fields directly instead of classifying free-text.

## TOML Config Contract

All gitflow-related fields move under `[gitflow]`. Top-level `commit_style`, `auto_push`, and `pr_mode` are deprecated.

```toml
[gitflow]
base_branch    = "master"          # required always
branch_mode    = "direct_push"     # enum: direct_push | scope_branch
branch_style   = ""                # required when branch_mode=scope_branch; must contain {scope} or <scope>
pr_create      = "never"           # enum: never | manual | auto
pr_title_style = ""                # required when pr_create=auto
pr_body_style  = ""                # required when pr_create=auto
commit_style   = "conventional"    # moved from top-level
auto_push      = true              # moved from top-level
```

### Validation rules

- `base_branch` is required for all modes.
- `branch_mode` must be `direct_push` or `scope_branch`.
- `branch_style` is required when `branch_mode=scope_branch` and must contain `{scope}` or `<scope>`.
- `pr_create` must be `never`, `manual`, or `auto`.
- `pr_create=auto` requires `branch_mode=scope_branch`.
- `pr_title_style` and `pr_body_style` are required when `pr_create=auto`.

### Migration from deprecated top-level fields

Top-level `commit_style`, `auto_push`, and `pr_mode` are still parsed. If `[gitflow]` fields are absent and deprecated top-level fields are present, they are used as fallback and `RMC_CONFIG_DEPRECATED_TOPLEVEL` is emitted as a warning. If both are present, `[gitflow]` wins.

`pr_mode` migration:
- `pr_mode=false` → `branch_mode=direct_push`, `pr_create=never`
- `pr_mode=true` → `branch_mode=scope_branch`, `pr_create=auto`

## Go Config Layer

### `GitflowConfig` struct

```go
type GitflowConfig struct {
    BaseBranch   string
    BranchMode   string   // "direct_push" | "scope_branch"
    BranchStyle  string
    PrCreate     string   // "never" | "manual" | "auto"
    PrTitleStyle string
    PrBodyStyle  string
    CommitStyle  string
    AutoPush     bool
}
```

`Config.CommitStyle`, `Config.AutoPush`, and `Config.PRMode` are removed from the top-level `Config` struct. All callers update to read from `cfg.Gitflow.*`.

### Validation

`validateConfig` gains:
- Enum checks for `branch_mode` and `pr_create` → `RMC_GITFLOW_BRANCH_MODE_INVALID`, `RMC_GITFLOW_PR_CREATE_INVALID`
- `branch_style` required + pattern check when `branch_mode=scope_branch` → `RMC_GITFLOW_BRANCH_STYLE_MISSING`
- `pr_create=auto` requires `branch_mode=scope_branch` → `RMC_GITFLOW_PR_AUTO_REQUIRES_SCOPE_BRANCH`
- `pr_title_style` / `pr_body_style` required when `pr_create=auto` → `RMC_GITFLOW_PR_STYLE_MISSING`
- `base_branch` required → `RMC_GITFLOW_BASE_BRANCH_MISSING`

### Defaults

```go
Gitflow: GitflowConfig{
    BranchMode:  "direct_push",
    PrCreate:    "never",
    CommitStyle: "conventional",
    AutoPush:    true,
}
```

`base_branch` has no default — it is required and must be explicit.

## Bootstrap Command

### Base command — read-only

`buildBootstrapConfig` loses both write paths:
- Auto-creation of missing bootstrap files (lines 280–295) is removed.
- `repairStemIfNeeded` call and the `--yes` flag are removed.

The base command now only reads config, validates it, and reports state.

### `bootstrap init --apply`

Gains stem repair: when `.stem` exists but has an incompatible schema, `--apply` repairs it in addition to creating missing files. A new `--yes` flag is added to `bootstrap init` (it does not currently exist there) to bypass the interactive repair prompt in non-interactive use.

### Bootstrap JSON output

`bootstrapConfigReport` struct changes:

**Removed flat fields:** `PRMode`, `CommitStyle`, `AutoPush`, `BranchStyle`, `PRTitleStyle`, `PRBodyStyle`.

**Added:**

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

// Fields added to bootstrapConfigReport:
Gitflow         bootstrapGitflowReport `json:"gitflow"`
MissingSettings []string               `json:"missing_settings,omitempty"`
EmptySettings   []string               `json:"empty_settings,omitempty"`
InvalidSettings []string               `json:"invalid_settings,omitempty"`
```

`missing_settings` — required fields absent from TOML (e.g. `"gitflow.base_branch"`).
`empty_settings` — fields present but empty when needed (e.g. `"gitflow.branch_style"` when `branch_mode=scope_branch`).
`invalid_settings` — fields with invalid enum values.

`RMC_GITFLOW_NOT_CONFIGURED` is replaced by entries in these arrays.

## Integrate Skill

### Input changes

| Field | Change |
|-------|--------|
| `pr_mode` | Removed |
| `branch_style` | Kept as branch-name template only (e.g. `feat/{scope}`) |
| `branch_mode` | New — read from `gitflow.branch_mode` |
| `pr_create` | New — read from `gitflow.pr_create` |
| `commit_style` | Read from `gitflow.commit_style` |
| `auto_push` | Read from `gitflow.auto_push` |
| `base_branch` | Read from `gitflow.base_branch` (runtime detection removed) |

### Fase 1 — Mode determination

The free-text classification heuristic (`scope-branch` / `direct-push` / `ambiguous`) is removed entirely. Mode is read directly from `branch_mode` and `pr_create`.

### Execution paths

| `branch_mode` | `pr_create` | Behavior |
|---------------|-------------|----------|
| `direct_push` | `never` | Commit on base branch, push if `auto_push=true`, no PR |
| `scope_branch` | `never` | Create/update scope branch, commit, push, no PR |
| `scope_branch` | `manual` | Create/update scope branch, commit, push, print suggested `gh pr create` command, do not execute it |
| `scope_branch` | `auto` | Create/update scope branch, commit, push, create PR, merge governed by `autonomy` + `is_last_in_scope` |

### Verification scenarios (headless)

Updated to cover all 4 paths above. The existing Scenario A–D are replaced with scenarios keyed on `branch_mode` + `pr_create` combinations.

## Roadmap Skill Updates

### `bootstrap-subcommand.md`

Wizard field order (conditional):

1. `base_branch` — always asked
2. `branch_mode` — menu: `direct_push` | `scope_branch`
3. `branch_style` — only when `branch_mode=scope_branch`; inferred from branch history and pre-filled
4. `pr_create` — menu: `never` | `manual` | `auto`; `auto` only available when `branch_mode=scope_branch`
5. `pr_title_style`, `pr_body_style` — only when `pr_create=auto`; inferred from merged PR history

The wizard uses `missing_settings` and `empty_settings` from the bootstrap JSON to decide which fields to ask about. Already-configured fields are shown for confirmation only (not re-asked) unless `/roadmap bootstrap` is invoked explicitly to force re-configuration.

### `bootstrap-reference.md`

- Replace flat field references (`pr_mode`, `commit_style`, `auto_push`) with `gitflow.*` paths.
- Document `missing_settings`, `empty_settings`, `invalid_settings` arrays.
- Replace `RMC_GITFLOW_NOT_CONFIGURED` diagnostic with the new arrays.

### `config-reference.md`

- Update the documented TOML schema to show all fields under `[gitflow]`.
- Mark `pr_mode`, `commit_style`, `auto_push` as deprecated top-level fields.

## Test Plan

### Go tests

- Bootstrap base command is read-only when files are missing (no auto-create).
- Bootstrap base command is read-only when `.stem` schema is incompatible (no auto-repair).
- `bootstrap init --apply` creates missing files and repairs `.stem`.
- New `[gitflow]` fields parse and serialize correctly.
- Invalid enum values produce the correct diagnostic IDs.
- `branch_style` required+pattern validation fires when `branch_mode=scope_branch`.
- `pr_create=auto` rejects when `branch_mode=direct_push`.
- Legacy `pr_mode=true/false` migrates to correct effective gitflow fields.
- `RMC_CONFIG_DEPRECATED_TOPLEVEL` emitted when top-level legacy fields are present.
- Bootstrap JSON includes `gitflow`, `missing_settings`, `empty_settings`, `invalid_settings`.
- Callers of `cfg.CommitStyle`, `cfg.AutoPush`, `cfg.PRMode` updated and still compile.

### Skill / headless tests

- `/roadmap bootstrap` asks fields conditionally and only writes after confirmation.
- `branch_mode=direct_push` never creates scope branches or PRs.
- `branch_mode=scope_branch` + `pr_create=manual` pushes but prints command without executing.
- `branch_mode=scope_branch` + `pr_create=auto` creates PRs as before.

### Run

```bash
go test ./... -count=1
scripts/sync-roadmap-skill.sh --check --skill roadmap
scripts/sync-roadmap-skill.sh --check --skill integrate
```

## Constraints

- `roadmapctl` owns config loading, defaults, validation, and structured reporting.
- Skills own wizard UX, repo convention inference, and user interaction.
- `pr_mode`, `commit_style`, `auto_push` remain parseable from top-level for migration; they are never the primary model going forward.
- `base_branch` is required in TOML — the integrate skill does not fall back to runtime detection.
