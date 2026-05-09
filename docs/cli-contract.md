# roadmapctl CLI contract

`roadmapctl` is the roadmap-specific guard CLI for Rootline-governed roadmaps. It validates environment, configuration, structure and dependency invariants while Rootline remains the generic filesystem database and constraint engine.

The historical MVP exposed only `doctor` and `check`. The current CLI also implements deterministic read, lint, transition, bootstrap, and materialization helpers while preserving the same report/version/exit-code contract.

Implemented commands:

- `roadmapctl doctor`
- `roadmapctl check`
- `roadmapctl context`
- `roadmapctl pending`
- `roadmapctl next`
- `roadmapctl decision`
- `roadmapctl lint`
- `roadmapctl transition`
- `roadmapctl materialize`
- `roadmapctl bootstrap`

`roadmapctl` does not decompose roadmap plans with AI, auto-fix invalid roadmap data, or add roadmap-specific subcommands to `rootline`.

## Mandatory use by `/roadmap`

Any implemented `/roadmap` command that writes, mutates, executes tasks, or claims a roadmap is valid must run `roadmapctl` first and block on failing diagnostics.

- Preflight for writes or execution: `roadmapctl doctor`.
- Structural/dependency validation before and after materialization or mutation: `roadmapctl check`.
- Conceptual planning that does not write files may proceed without `roadmapctl`, but must not claim materialization or validity.

See [roadmap skill integration](roadmap-skill-integration.md) for exact preflight/postcheck commands, blocking policy, `/roadmap plan` and `/roadmap loop` snippets, and expected failure handling.

## Command summary

```text
roadmapctl [global flags] <command> [command flags]

Commands:
  doctor       Diagnose repo/workspace, roadmap config, Rootline availability and schema prerequisites.
  check        Validate canonical roadmap structure, metadata, Rootline graph and blocking dependencies.
  context      Show effective roadmapctl context for skill bootstrap.
  pending      List active roadmap tasks that are not done.
  next         Split active tasks into ready and blocked sets.
  decision     Provide deterministic prioritization recommendations.
  lint         Validate deterministic semantic roadmap conventions.
  transition   Evaluate and apply policy-checked status transitions.
  materialize   Validate and write approved structured roadmap plans.
  bootstrap     Inspect or initialize missing bootstrap files.
```

Commands support `--output text` and `--output json`.

## Global flags

| Flag | Values | Default | Description |
|------|--------|---------|-------------|
| `--repo` | path | cwd | Repository root or workspace member to inspect. |
| `--roadmap-root` | path | inferred from `<roadmap-root>/.roadmapctl.toml` or legacy `.claude/roadmap.local.md` | Override configured roadmap root. The resolved path must stay inside the repo. |
| `--workspace` | bool | auto | Treat `--repo`/cwd as a workspace containing multiple repos. |
| `--output` | `text`, `json` | `text` | Select human or machine-readable output. |
| `--strict` | bool | `false` | Treat warnings as failures when calculating exit code. |
| `--rootline` | path | `ROOTLINE_BIN` or PATH | Rootline executable to invoke. |
| `--timeout` | duration | `10s` | Timeout for each Rootline subprocess call. |

Durations use Go duration syntax, for example `500ms`, `10s`, `2m`.

## Configuration contract

Preferred post-MVP config lives at `<roadmap-root>/.roadmapctl.toml` (for example `docs/roadmap/.roadmapctl.toml`). The roadmap root is inferred from the directory containing the TOML file; `roadmap_root` is intentionally not a TOML key. `--roadmap-root` remains a command-line override for explicit inspection and migration workflows.

Example single-repo config:

```toml
done_statuses = ["Completed", "Obsolete"]
active_statuses = ["Pending", "Specified", "In Progress"]
leaf_filter = "isIndex == false"
outcome_close_verify = []
pr_merge_strategy = "squash"
commit_style = "conventional"
auto_push = true

[status_values]
pending = "Pending"
specified = "Specified"
in_progress = "In Progress"
completed = "Completed"
blocked = "Blocked"
obsolete = "Obsolete"
```

Workspace mode uses the same per-repo `<roadmap-root>/.roadmapctl.toml` file. A future workspace index may point at member repositories, but each repo remains authoritative for its own roadmap root and operational roles.

Config keys:

| TOML key | Type | Default | Meaning |
|----------|------|---------|---------|
| `done_statuses` | list(string) | `["Completed", "Obsolete"]` | Status values treated as dependency-satisfying/done. |
| `active_statuses` | list(string) | `["Pending", "Specified", "In Progress"]` | Status values listed by active/pending workflows. |
| `leaf_filter` | string | `isIndex == false` | Rootline expression selecting leaf records. |
| `outcome_close_verify` | list(string) | `[]` | Optional commands for outcome close checks. |
| `pr_merge_strategy` | enum string | `squash` | Preferred PR merge strategy (`squash`, `merge`, `rebase`). |
| `commit_style` | enum string | `conventional` | Commit message style. |
| `auto_push` | bool | `true` | Whether loop workflows push after commits. |
| `[status_values].pending` | string | `Pending` | Operational pending role value. |
| `[status_values].specified` | string | `Specified` | Operational specified role value. |
| `[status_values].in_progress` | string | `In Progress` | Operational in-progress role value. |
| `[status_values].completed` | string | `Completed` | Operational completed role value. |
| `[status_values].blocked` | string | `Blocked` | Operational blocked role value. |
| `[status_values].obsolete` | string | `Obsolete` | Operational obsolete role value. |

Operational config does not define or constrain the full document schema. The effective Rootline `.stem` remains the source of truth for valid `estado`, `tipo`, and link values. `roadmapctl check` validates operational status values against the schema separately and emits `RMC_CONFIG_STATUS_SCHEMA_MISMATCH` when a role points at a status absent from `.stem`.

Precedence:

1. Command-line flags (`--repo`, `--roadmap-root`, `--rootline`, `--timeout`, `--output`, `--strict`).
2. Preferred `<roadmap-root>/.roadmapctl.toml` discovered under the selected repo/roadmap root.
3. Legacy `.claude/roadmap.local.md` fallback while migration is in progress.
4. Built-in defaults above for omitted optional keys.

Conflict policy:

- If both `.roadmapctl.toml` and legacy `.claude/roadmap.local.md` exist and resolve to different roadmap roots or operational values, prefer `.roadmapctl.toml` and emit a compatibility warning diagnostic.
- If only the legacy file exists, load it as fallback and emit an informational diagnostic identifying the legacy source.
- If neither config exists and no explicit `--roadmap-root` is provided, `doctor` emits a config diagnostic and write/mutation flows must block.

Open decision: workspace-level discovery may later add a root `.roadmapctl.toml` or workspace index. Until that is approved, this contract treats config as per-repository.

## `doctor`

`doctor` answers whether `roadmapctl` can operate in the current environment.

It checks:

1. repo/workspace discovery from `--repo` or cwd;
2. `.claude/roadmap.local.md` existence and parseability;
3. `roadmap-root` resolution and containment inside the repo;
4. Rootline executable discovery and basic invocation;
5. roadmap root and `.stem` presence;
6. relevant config/cache paths for troubleshooting.

`doctor` must not create files, install dependencies, modify hooks, or repair configuration.

Example:

```bash
roadmapctl doctor --repo . --output json
```

## `check`

`check` validates roadmap invariants after `doctor` succeeds.

It checks:

1. canonical filesystem shape:
   - direct tasks are `<roadmap-root>/TXXX-*.md`;
   - outcomes are `<roadmap-root>/OXX-*/README.md` plus task files;
   - outcome tasks are `<roadmap-root>/OXX-*/TXXX-*.md`;
   - no extra nesting below outcomes;
2. no fallback summary files such as `*-tasks.md` represent multiple tasks;
3. no duplicate `OXX` in the root or duplicate `TXXX` inside a scope;
4. Rootline validation for `.stem` and markdown frontmatter;
5. Rootline query results for task leaves;
6. Rootline graph results for cycles and broken links;
7. `blocked_by` links use explicit relative paths and resolve to task files;
8. markdown `estado`/`tipo` values are valid according to the effective Rootline schema;
9. operational status roles from config (`status-values`, `done-statuses`, `active-statuses`) refer to statuses present in the effective Rootline schema.

`check` must not write, materialize, fix, or normalize roadmap files.

Example:

```bash
roadmapctl check --repo . --output text --strict
```

## Implemented read/context commands

`context`, `pending`, `next`, and `decision` are read-only. They share the top-level report fields and add command-specific arrays/objects:

| Command | `kind` | Purpose | Key fields |
|---------|--------|---------|------------|
| `context` | `roadmapctl/context` | Resolve effective repo, roadmap root, config source, schema, status roles and prompt helpers. | `config_path`, `config_source`, `rootline_version`, `schema`, `status_values`, `done_statuses`, `active_statuses`, `helpers` |
| `pending` | `roadmapctl/pending` | List active non-done tasks without mutating state. | `count`, `tasks[]`, workspace `repos[]` when applicable |
| `next` | `roadmapctl/next` | Separate ready and blocked active tasks. | `ready[]`, `blocked[]` |
| `decision` | `roadmapctl/decision` | Provide deterministic prioritization data. | `recommendations[]`, `quick_wins[]`, `critical_blockers[]`, `blocked[]` |

These commands must not write files or update statuses. They are the supported agent discovery boundary for pending, next-task, and decision flows: agents consume their JSON directly instead of calling Rootline `tree`/`query`/`graph` or postprocessing Rootline JSON with Python snippets. Internally, they rely on Rootline `tree`/`query`/`graph` data and roadmapctl's status-role config.

## `lint` contract

`lint` is the deterministic semantic check layer. It runs after `doctor`/`check` prerequisites are satisfied and remains read-only: it must not materialize, normalize, auto-fix, or judge subjective writing quality.

Boundary:

- `check`: canonical filesystem shape, Rootline/frontmatter/schema validation, dependency graph invariants.
- `lint`: deterministic documentation and portability conventions that are useful for agents and releases.

Initial lint rule groups:

1. Outcome `## Tasks` table consistency with child `TXXX-*.md` files.
2. Required task sections: `Preserva`, `Contexto`, `Alcance`, `Estado inicial esperado`, `Criterios de Aceptación`, `Fuente de verdad`.
3. Presence-only checks for acceptance criteria and source-of-truth entries.
4. Effective schema compatibility for roadmapctl-required `estado`, `tipo`, and `blocked_by`.
5. Cross-platform filename and name portability.

Example future invocation:

```bash
roadmapctl lint --repo . --output json --strict
```

## CI integration

JSON remains the source of truth for CI. Recommended CI usage runs `check` and, when semantic conventions are desired, `lint` with `--output json --strict`:

```bash
roadmapctl check --repo . --roadmap-root docs/roadmap --output json --strict > roadmapctl-check.json
roadmapctl lint --repo . --roadmap-root docs/roadmap --output json --strict > roadmapctl-lint.json
```

Exit codes are the stable machine contract. Additional formats such as SARIF/JUnit are deferred until explicitly approved; GitHub-specific annotations should be generated from JSON by wrapper workflows rather than changing core command semantics.

## Streams

### JSON output

When `--output json` is selected:

- stdout contains exactly one JSON report and no additional text;
- stderr may contain process-level errors that prevented report creation;
- logs, progress messages and debug output are suppressed or sent to stderr;
- diagnostics are represented in the JSON report whenever possible.

### Text output

When `--output text` is selected:

- stdout contains a human-readable summary and diagnostics;
- stderr is reserved for process-level errors and optional debug output;
- formatting is not a stable API.

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success: no errors. Warnings may exist unless `--strict` is set. |
| `1` | Validation failure: roadmap/config/graph diagnostics contain errors. |
| `2` | Usage or configuration error: invalid flags, unreadable config, invalid path override. |
| `3` | Environment/dependency error: missing Rootline, timeout, permission or subprocess failure. |
| `4` | Internal error: unexpected panic, invariant violation or unsupported report version. |

With `--strict`, warnings are promoted when calculating the exit code and may produce `1`.

## JSON report schema

All successful command executions that reach report construction emit this shape in JSON mode:

```json
{
  "version": 1,
  "kind": "roadmapctl/doctor",
  "summary": {
    "status": "ok",
    "errors": 0,
    "warnings": 0,
    "infos": 1
  },
  "root": "/abs/path/to/repo",
  "roadmap_root": "/abs/path/to/repo/docs/roadmap",
  "diagnostics": [
    {
      "id": "RMC_EXAMPLE",
      "severity": "info",
      "message": "human-readable message",
      "path": "docs/roadmap",
      "details": {
        "key": "value"
      }
    }
  ]
}
```

### Top-level fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `version` | integer | yes | Report schema version. MVP uses `1`. |
| `kind` | string | yes | `roadmapctl/doctor`, `roadmapctl/check`, or planned `roadmapctl/lint`. |
| `summary` | object | yes | Aggregated status and diagnostic counts. |
| `root` | string | yes | Absolute repo/workspace member root. |
| `roadmap_root` | string | yes | Absolute resolved roadmap root, when known. Empty if unavailable. |
| `diagnostics` | array | yes | Ordered diagnostics. Empty on clean success. |

### `summary`

| Field | Type | Values |
|-------|------|--------|
| `status` | string | `ok`, `warning`, `error` |
| `errors` | integer | count of error diagnostics |
| `warnings` | integer | count of warning diagnostics |
| `infos` | integer | count of info diagnostics |

`summary.status` is derived only from emitted diagnostic severities. `--strict` affects process exit code for warning diagnostics; it does not rewrite warning severities or change a warning-only summary to `error`.

### `diagnostics[]`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | Stable diagnostic identifier. |
| `severity` | string | yes | `info`, `warning`, or `error`. |
| `message` | string | yes | Human-readable actionable explanation. |
| `path` | string | no | Repo-relative path when applicable. |
| `details` | object | no | Machine-readable diagnostic context. |

## Diagnostic ID convention

Diagnostic IDs are stable strings with the prefix `RMC_`:

```text
RMC_<AREA>_<CONDITION>
```

Areas in the MVP:

- `RMC_CONFIG_*`
- `RMC_ENV_*`
- `RMC_ROOTLINE_*`
- `RMC_STRUCTURE_*`
- `RMC_GRAPH_*`
- `RMC_STATUS_*`
- `RMC_LINT_*` (planned `lint` command)

IDs are part of the machine-readable contract. Messages may change; IDs should not change without a report version bump or compatibility note.

## Required MVP diagnostics

| ID | Severity | Command | Meaning |
|----|----------|---------|---------|
| `RMC_STRUCTURE_SINGLE_FILE_FALLBACK` | error | `check` | A file such as `*-tasks.md` appears to represent multiple tasks instead of canonical `TXXX-*.md` files. |
| `RMC_ENV_ROOTLINE_MISSING` | error | `doctor`, `check` | Rootline executable could not be found via `--rootline`, `ROOTLINE_BIN`, or PATH. |
| `RMC_GRAPH_INVALID_BLOCKED_BY` | error | `check` | A `blocked_by` link is broken, not explicit relative syntax, or does not resolve to a task file. |

Additional MVP diagnostics should reuse the same convention, for example:

- `RMC_CONFIG_MISSING`
- `RMC_CONFIG_ROADMAP_ROOT_ESCAPE`
- `RMC_CONFIG_STATUS_SCHEMA_MISMATCH`
- `RMC_STRUCTURE_MISSING_OUTCOME_README`
- `RMC_STRUCTURE_DUPLICATE_ID`
- `RMC_ROOTLINE_VALIDATE_FAILED`
- `RMC_ROOTLINE_DESCRIBE_FAILED`
- `RMC_ROOTLINE_QUERY_FAILED`
- `RMC_ROOTLINE_GRAPH_FAILED`
- future rootline operation IDs such as `RMC_ROOTLINE_TREE_FAILED`, `RMC_ROOTLINE_SET_FAILED`, `RMC_ROOTLINE_NEW_FAILED`
- `RMC_GRAPH_CYCLE`
- `RMC_STATUS_UNKNOWN`

Rootline operation diagnostics use `details.kind` to distinguish `missing_binary`, `timeout`, `execution`, `incompatible_command`, `invalid_json`, and `invalid_shape` when known.

## Planned lint diagnostics

| ID | Severity | Command | Meaning |
|----|----------|---------|---------|
| `RMC_LINT_TASK_TABLE_MISSING` | warning | `lint` | Outcome README has no parseable `## Tasks` table. |
| `RMC_LINT_TASK_TABLE_MISSING_ROW` | warning | `lint` | Child task file is absent from the outcome task table. |
| `RMC_LINT_TASK_TABLE_STALE_ROW` | warning | `lint` | Task table row links to no current child task file. |
| `RMC_LINT_TASK_TABLE_INVALID_LINK` | warning | `lint` | Task table link is not an explicit relative child-task link. |
| `RMC_LINT_TASK_SECTION_MISSING` | warning | `lint` | Task is missing a required roadmap section heading. |
| `RMC_LINT_ACCEPTANCE_CRITERIA_MISSING` | warning | `lint` | Task has no acceptance criteria entries detectable by structure. |
| `RMC_LINT_SOURCE_OF_TRUTH_EMPTY` | warning | `lint` | Task `Fuente de verdad` section has no entries. |
| `RMC_LINT_FILENAME_CASE_COLLISION` | error | `lint` | Roadmap entries collide under case-insensitive filesystem rules. |
| `RMC_LINT_FILENAME_RESERVED` | error | `lint` | Roadmap entry name is reserved or problematic on supported platforms. |
| `RMC_LINT_SCHEMA_FIELD_MISSING` | error | `lint` | Effective schema lacks a required field such as `estado` or `tipo`. |
| `RMC_LINT_SCHEMA_LINK_MISSING` | error | `lint` | Effective schema lacks a required link relation such as `blocked_by`. |

Severity policy:

- `warning`: deterministic semantic or documentation consistency issue; exits `0` unless `--strict` is set.
- `error`: deterministic portability or schema problem that can break roadmapctl operation or supported filesystems.
- `lint` must not reclassify MVP `RMC_STRUCTURE_*`, `RMC_GRAPH_*`, `RMC_ROOTLINE_*`, or `RMC_STATUS_*` diagnostics without compatibility notes.

## Materialize plan input contract

`roadmapctl materialize` accepts only structured, approved JSON input. It does not parse free-form chat or call an LLM. The versioned input schema, examples, validation rules, dependency representation, and `RMC_MATERIALIZE_*` diagnostics are specified in [materialize-plan-schema.md](materialize-plan-schema.md).

Implemented modes:

```bash
roadmapctl materialize --plan plan.json --dry-run --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl materialize --plan plan.json --apply --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl materialize --changes dry-run.json --target O02-new-outcome/README.md --apply --repo <repo> --roadmap-root <roadmap-root> --output json
```

JSON adds `applied` and `changes[]`. Each change has `path`, `operation`, `applied`, optional `content`, optional `diff`, and `preconditions[]`. Dry-run must not write. Batch `--plan ... --apply` remains backward-compatible: it writes only allowlisted bootstrap/config/schema files and canonical roadmap markdown, validates created markdown through Rootline, and runs a postcheck before reporting success.

Granular target apply uses a previously saved dry-run JSON report as a frozen change-set. `--changes <dry-run-json> --target <path> --apply` applies exactly one matching canonical roadmap markdown `create` change and validates the created markdown through Rootline. Empty, unknown, duplicate, non-create, or non-canonical targets fail before writing. This mode does not recompute numbering from the plan and does not create sibling roadmap files.

## Transition controller contract

`roadmapctl transition` is specified in [transition-controller.md](transition-controller.md). It defines actions `can-start`, `can-complete`, `start`, `complete`, and `set-status`; required status roles; dependency satisfaction via `done_statuses`; behavior for schema-valid non-role statuses such as `On Hold`; and `RMC_TRANSITION_*` diagnostics.

## Rootline integration boundary

`roadmapctl` may invoke Rootline only as an external executable using explicit arguments and a timeout. It must not import Rootline internals or execute shell strings.

Expected Rootline operations for the MVP include generic commands such as `validate`, `describe`, `query` and `graph`. Roadmap-specific interpretation belongs in `roadmapctl`, not in Rootline.
