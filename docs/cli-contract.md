# roadmapctl CLI contract

`roadmapctl` is the roadmap-specific guard CLI for Rootline-governed
roadmaps. Rootline remains the generic Markdown filesystem database and
constraint engine; roadmapctl adds roadmap policy, guards, diagnostics, and
stable JSON outputs.

## Implemented Commands

```text
roadmapctl [global flags] <command> [command flags]

Commands:
  doctor       Diagnose repo/workspace, config, Rootline availability, and schema prerequisites.
  check        Validate canonical roadmap structure, metadata, Rootline graph, and blocking dependencies.
  bootstrap    Inspect or initialize bootstrap files, or show effective agent config.
  pending      List active roadmap tasks that are not done.
  next         Split active tasks into ready and blocked sets.
  decision     Provide deterministic prioritization recommendations.
  lint         Validate deterministic semantic roadmap conventions.
  transition   Evaluate and apply policy-checked task status transitions.
```

Removed or unsupported surfaces are not part of the live contract: the old
context command, materialization command, plan-path helper, and roadmap-root
override flag.

## Global Flags

| Flag | Values | Default | Description |
|------|--------|---------|-------------|
| `--repo` | path | `.` | Repository root or workspace member to inspect. |
| `--workspace` | bool | `false` | Treat `--repo` as a workspace containing multiple repos. |
| `--output` | `text`, `json` | `text` | Select human or machine-readable output. |
| `--strict` | bool | `false` | Treat warnings as failures when calculating exit code. |
| `--rootline` | path | `ROOTLINE_BIN` or `PATH` | Rootline executable to invoke. |
| `--timeout` | duration | `10s` | Timeout for each Rootline subprocess call. |

Durations use Go syntax, for example `500ms`, `10s`, or `2m`.

## Configuration

The roadmap root is fixed at `docs/roadmap/` under the selected repo. It is
not configurable by CLI flag or TOML field. Config lives at
`docs/roadmap/.roadmapctl.toml`.

Example config:

```toml
done_statuses = ["Completed", "Obsolete"]
active_statuses = ["Pending", "Specified", "In Progress"]
leaf_filter = "isIndex == false"
outcome_close_verify = []
commit_style = "conventional"
auto_push = true
required_code_coverage = 85.0
loop_max_tasks = 0
parallel = true
autonomy = "until_done"
compact_after_task_commit = true
pr_mode = false

[status_values]
pending = "Pending"
specified = "Specified"
in_progress = "In Progress"
completed = "Completed"
blocked = "Blocked"
obsolete = "Obsolete"

[gitflow]
branch_style = ""
pr_title_style = ""
pr_body_style = ""
```

Config keys:

| TOML key | Type | Default | Meaning |
|----------|------|---------|---------|
| `done_statuses` | list(string) | `["Completed", "Obsolete"]` | Status values treated as dependency-satisfying/done. |
| `active_statuses` | list(string) | `["Pending", "Specified", "In Progress"]` | Status values listed by active/pending workflows. |
| `leaf_filter` | string | `isIndex == false` | Rootline expression selecting leaf records. |
| `outcome_close_verify` | list(string) | `[]` | Optional commands for outcome close checks. |
| `commit_style` | string | `conventional` | Commit message style consumed by skills. |
| `auto_push` | bool | `true` | Whether integration workflows push after commits. |
| `required_code_coverage` | float | `85.0` | Repo coverage target consumed by local tooling. |
| `loop_max_tasks` | integer | `0` | Repo-default loop cap; `0` means unlimited. |
| `parallel` | bool | `true` | Whether `/roadmap loop` may execute safe independent waves. |
| `autonomy` | enum | `until_done` | Loop policy: `manual`, `supervised`, or `until_done`. |
| `compact_after_task_commit` | bool | `true` | Whether the skill may compact context after a durable task commit. |
| `pr_mode` | bool | `false` | Whether integration workflows use PR mode by default. |
| `[gitflow].branch_style` | string | `""` | Natural-language branch convention for skills. |
| `[gitflow].pr_title_style` | string | `""` | Natural-language PR title convention for skills. |
| `[gitflow].pr_body_style` | string | `""` | Natural-language PR body convention for skills. |
| `[status_values].pending` | string | `Pending` | Operational pending role value. |
| `[status_values].specified` | string | `Specified` | Operational specified role value. |
| `[status_values].in_progress` | string | `In Progress` | Operational in-progress role value. |
| `[status_values].completed` | string | `Completed` | Operational completed role value. |
| `[status_values].blocked` | string | `Blocked` | Operational blocked role value. |
| `[status_values].obsolete` | string | `Obsolete` | Operational obsolete role value. |

## Bootstrap

`bootstrap` is the agent context API:

```bash
roadmapctl bootstrap --repo <path> --output json
```

It returns effective repo, roadmap root, config source, Rootline version,
status roles, operational settings, gitflow style fields, helpers, and
diagnostics. Use `--field <dot-path>` to extract a scalar value from the
bootstrap JSON.

Missing bootstrap files are handled explicitly:

```bash
roadmapctl bootstrap inspect --repo <path> --output json
roadmapctl bootstrap init --repo <path> --dry-run --output json
roadmapctl bootstrap init --repo <path> --apply --output json
```

`bootstrap` may also repair known legacy `.stem` compatibility issues when run
with `--yes`; unsupported custom stems produce
`RMC_BOOTSTRAP_REPAIR_UNSUPPORTED_STEM`.

## Guards

Any `/roadmap` flow that writes, mutates, executes tasks, commits roadmap
state, or claims validity must run:

```bash
roadmapctl doctor --repo <path> --output json --strict
roadmapctl check --repo <path> --output json --strict
```

After creating or mutating roadmap files, run:

```bash
roadmapctl check --repo <path> --output json --strict
```

`doctor` and `check` are read-only. They do not create files, install
dependencies, repair configuration, or normalize roadmap content.

## Read Commands

`pending`, `next`, and `decision` are read-only and are the supported agent
boundary for queue, readiness, blockers, and prioritization:

```bash
roadmapctl pending --repo <path> --output json
roadmapctl next --repo <path> --output json
roadmapctl decision --repo <path> --output json
```

Workspace summaries use:

```bash
roadmapctl pending --workspace --repo <workspace-root> --output json
```

Agents must consume these JSON reports directly instead of rebuilding policy
from raw Rootline `tree`, `query`, or `graph` output.

## Lint

`lint` is a deterministic semantic convention layer. It remains read-only.
Warnings exit `0` unless `--strict` is set.

Current lint groups include:

- Outcome task-table consistency with child files.
- Required task sections: `Preserva`, `Contexto`, `Alcance`, `Estado inicial esperado`, `Criterios de Aceptación`, `Fuente de verdad`.
- Presence checks for acceptance criteria and source-of-truth entries.
- Effective schema compatibility for required roadmap fields and links.
- Cross-platform filename and name portability.

## Transition

`transition` owns status policy and status mutation:

```bash
roadmapctl transition can-start <task-path> --repo <path> --output json
roadmapctl transition can-complete <task-path> --repo <path> --output json
roadmapctl transition start <task-path> --repo <path> --apply --output json
roadmapctl transition complete <task-path> --repo <path> --apply --output json
roadmapctl transition set-status <task-path> <status> --repo <path> --apply --output json
```

Mutating actions require `--apply`; otherwise they operate as dry-run plans.
Apply actions validate after mutation before reporting success.

## JSON Output

When `--output json` is selected:

- stdout contains exactly one JSON report when report construction succeeds;
- stderr is reserved for process-level failures and explicit prompts;
- diagnostics are represented in the JSON report whenever possible.

Common top-level fields:

| Field | Type | Description |
|-------|------|-------------|
| `version` | integer | Report schema version. |
| `kind` | string | Command report kind, such as `roadmapctl/check`. |
| `summary` | object | Aggregated status and diagnostic counts. |
| `root` | string | Absolute selected repo or workspace root. |
| `roadmap_root` | string | Absolute fixed roadmap root when known. |
| `diagnostics` | array | Ordered diagnostics. |

`summary.status` is derived from diagnostic severities. `--strict` affects the
process exit code for warnings; it does not rewrite warning severities.

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success: no errors. Warnings may exist unless `--strict` is set. |
| `1` | Validation failure: roadmap/config/graph diagnostics contain errors, or warnings under `--strict`. |
| `2` | Usage or configuration error. |
| `3` | Environment/dependency error such as missing Rootline, timeout, permission, or subprocess failure. |
| `4` | Internal error or unsupported report version. |

## Diagnostics

Diagnostic IDs are stable strings with prefix `RMC_`.

Important IDs include:

| ID | Severity | Meaning |
|----|----------|---------|
| `RMC_CONFIG_MISSING` | error | Roadmap root/config could not be found. |
| `RMC_CONFIG_PARSE` | error | `.roadmapctl.toml` could not be parsed or validated. |
| `RMC_ENV_ROOTLINE_MISSING` | error | Rootline executable could not be found. |
| `RMC_STRUCTURE_SINGLE_FILE_FALLBACK` | error | A summary file appears to represent multiple tasks instead of canonical task files. |
| `RMC_STRUCTURE_MISSING_OUTCOME_README` | error | Outcome directory has no `README.md`. |
| `RMC_STRUCTURE_DUPLICATE_ID` | error | Duplicate `OXX` or `TXXX` IDs in a scope. |
| `RMC_GRAPH_INVALID_BLOCKED_BY` | error | A `blocked_by` link is invalid, unresolved, or not a task file. |
| `RMC_GRAPH_CYCLE` | error | Dependency graph contains a cycle. |
| `RMC_STATUS_UNKNOWN` | error | Roadmap status is not valid for the effective schema. |
| `RMC_CONFIG_STATUS_SCHEMA_MISMATCH` | error | Config status roles point at values absent from `.stem`. |
| `RMC_LINT_TASK_SECTION_MISSING` | warning | Task is missing a required section heading. |
| `RMC_LINT_FILENAME_CASE_COLLISION` | error | Roadmap entries collide on case-insensitive filesystems. |
| `RMC_LINT_FILENAME_RESERVED` | error | Roadmap entry name is reserved or problematic on supported platforms. |
| `RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED` | error | `.stem` requires `estado` for Outcome README records. |
| `RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY` | error | `.stem` has a global non-empty status rule incompatible with Outcome README files. |
| `RMC_GITFLOW_NOT_CONFIGURED` | info | `[gitflow]` style fields are absent or empty. |
| `RMC_BOOTSTRAP_REPAIR_UNSUPPORTED_STEM` | error | Automatic `.stem` repair is not safe for this schema. |

Rootline operation diagnostics use `details.kind` when known, with values such
as `missing_binary`, `timeout`, `execution`, `incompatible_command`,
`invalid_json`, or `invalid_shape`.

## Rootline Boundary

`roadmapctl` invokes Rootline as an external executable with explicit arguments
and a timeout. It must not import Rootline internals or run shell strings.
Roadmap-specific interpretation belongs in roadmapctl, not Rootline.
