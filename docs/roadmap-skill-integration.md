# `/roadmap` skill integration with roadmapctl

`roadmapctl` is mandatory from day 1 for every implemented `/roadmap` command that writes files, mutates roadmap state, executes tasks, or claims that a roadmap is valid.

This document describes how the skill must use `roadmapctl doctor` and `roadmapctl check` as blocking guards while preserving Rootline as the generic filesystem database and constraint engine. The current product decision is to keep `roadmapctl materialize` as the deterministic writer for roadmap plan files; see `docs/decisions/materialize-writer-vs-guard-flow.md`.

## Applicability

The guard applies to implemented `/roadmap` flows that do any of the following:

- create, materialize, rename, move, or delete roadmap files;
- mutate frontmatter or links, including task status changes;
- execute pending tasks or make commits based on roadmap state;
- declare that a roadmap structure, dependency graph, or status set is valid.

Examples that must be guarded:

- `/roadmap plan` when it materializes Outcomes or Tasks;
- `/roadmap loop` before selecting/executing tasks and after mutating task state;
- any future command that updates roadmap files, validates roadmap state, or runs work from roadmap tasks.

## Conceptual-mode exception

Conceptual planning may run without `roadmapctl` only when it does not write files, does not mutate state, does not execute tasks, and does not claim that any materialized roadmap is valid.

Allowed without `roadmapctl`:

- discussing an approach;
- decomposing a possible Outcome/Task tree in chat;
- explaining tradeoffs before materialization.

Not allowed without `roadmapctl`:

- creating `OXX-*` or `TXXX-*` files;
- updating `estado`, links, dependencies, or `.stem` files;
- saying a materialized roadmap passes validation;
- falling back to a freeform markdown summary such as `*-tasks.md`.

## Blocking policy

If `roadmapctl` is unavailable or any required `roadmapctl` command exits non-zero, the skill must stop before the write, mutation, execution, commit, or validity claim that required the guard.

Blocking behavior:

1. preserve the working tree as-is;
2. report the exact command that failed;
3. report exit code and diagnostic IDs when JSON was produced;
4. do not auto-fix, rewrite, materialize alternative files, or continue with freeform markdown;
5. ask for explicit user/supervisor action only if a product decision is needed.

Warnings do not block unless the command is run with `--strict`. The default integration should use `--strict` for pre/post validation in write and execution flows.

## Required commands

Use repo paths resolved by the skill bootstrap. In single-repo mode, `<repo>` is the repository root. In workspace mode, `<repo>` is the selected workspace member. `<roadmap-root>` is the configured roadmap root relative to `<repo>` unless an absolute override is required.

The bootstrap should prefer deterministic context from roadmapctl:

```bash
roadmapctl context --repo <repo> --roadmap-root <roadmap-root> --output json
```

The returned `helpers`, `status_values`, `done_statuses`, `active_statuses`, operational settings, `root`, and `roadmap_root` fields are the source of truth for prompt placeholders. `<roadmap-root>/.roadmapctl.toml` is the lasting config source; `.claude/roadmap.local.md` is migration input only and is deleted by roadmapctl after a successful TOML load or migration.

### Bootstrap exception for missing roadmap roots

Normal writes require `doctor` and `check` first. The only exception is an explicitly requested missing-root bootstrap, because `check` cannot pass before `<roadmap-root>` and its schema/config files exist.

Allowed bootstrap flows are:

1. Preferred setup flow: `roadmapctl bootstrap inspect`, then `roadmapctl bootstrap init --dry-run`, then `roadmapctl bootstrap init --apply` after explicit approval.
2. Plan-materialization bootstrap: `roadmapctl materialize --dry-run` may propose only bootstrap allowlist paths (`.`, `.stem`, `.roadmapctl.toml`) plus canonical roadmap files; `--apply` is allowed only after explicit approval of that dry-run.

Both apply flows must run a postcheck before success is reported. `bootstrap init --apply` runs `roadmapctl check` internally after writing, and the skill must also run an explicit `roadmapctl check --strict` before claiming success or committing. `materialize --apply` runs its own postcheck and must still be followed by the explicit skill postcheck.

This exception is not an auto-fix path: if `<roadmap-root>` already exists and `doctor` or `check` fails, stop and report diagnostics. Do not use bootstrap to repair invalid roadmaps, rewrite `.stem`, or bypass normal preflight.

## Roadmap context compaction extension

The project-local Pi extension lives at `.pi/extensions/roadmap-context/index.ts`, which Pi auto-discovers from `.pi/extensions/*/index.ts` and reloads with `/reload`. It registers the exact tool name `compact_roadmap_context`.

When `compact_after_task_commit = true`, the roadmap skill can call `compact_roadmap_context` after a task is fully durable (ACs pass, `roadmapctl transition complete --apply` succeeds, commit/push/PR bookkeeping is done). Pass concise strings for `task_path`, `commit_hash`, `validation_summary`, `next_work`, and `config_summary` so compaction preserves the current roadmap goal, completed task, validation results, next task/wave state, unresolved blockers/conflicts, and effective config values.

The tool queues `ctx.compact({ customInstructions, onComplete, onError })` and returns immediately with queued/failed status. Callback failures warn through Pi UI when available; a compaction failure must not invalidate an already completed task commit.

### Preflight before writes, mutations, or execution

Run `doctor` first:

```bash
roadmapctl doctor \
  --repo <repo> \
  --roadmap-root <roadmap-root> \
  --output json \
  --strict
```

Then run `check` before relying on roadmap state:

```bash
roadmapctl check \
  --repo <repo> \
  --roadmap-root <roadmap-root> \
  --output json \
  --strict
```

### Postcheck after materialization or mutation

After creating or changing roadmap files, run:

```bash
roadmapctl check \
  --repo <repo> \
  --roadmap-root <roadmap-root> \
  --output json \
  --strict
```

This postcheck is mandatory before reporting successful materialization, completing a loop iteration, committing roadmap mutations, or claiming that the roadmap is valid. If materialization wrote files and then postcheck fails, the agent must report the partial state, inspect/validate affected paths, deliberately repair or revert, rerun `roadmapctl check --strict`, and only then commit or claim success.

### Optional binary override

If the user or CI provides an explicit Rootline binary, pass it through `roadmapctl` instead of invoking Rootline directly:

```bash
roadmapctl check \
  --repo <repo> \
  --roadmap-root <roadmap-root> \
  --rootline <path-to-rootline> \
  --output json \
  --strict
```

## `/roadmap plan` integration snippet

Apply this behavior to the materialization phase of `/roadmap plan`.

```text
Before creating or modifying any roadmap file:
  1. Require `roadmapctl` on PATH.
  2. If this is an explicit missing-root bootstrap, run:
     roadmapctl bootstrap inspect --repo <repo> --roadmap-root <roadmap-root> --output json
     roadmapctl bootstrap init --repo <repo> --roadmap-root <roadmap-root> --dry-run --output json
     After approval, run:
     roadmapctl bootstrap init --repo <repo> --roadmap-root <roadmap-root> --apply --output json
     roadmapctl check --repo <repo> --roadmap-root <roadmap-root> --output json --strict
     If any command exits non-zero, stop.
  3. Otherwise run:
     roadmapctl doctor --repo <repo> --roadmap-root <roadmap-root> --output json --strict
     roadmapctl check --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  4. If either normal preflight command exits non-zero, stop. Do not write files and do not fall back to a summary markdown file.

After approval, the skill serializes non-bootstrap plans to `roadmapctl/materialize-plan` JSON and delegates deterministic writes. The versioned schema source is `docs/materialize-plan-schema.md` in the `roadmapctl` repository; that path is not assumed to exist in consuming repos or inside the installed skill directory. Until a future `roadmapctl materialize schema --output json` command is explicitly approved, agents should use the skill's embedded minimal schema plus roadmapctl diagnostics outside this repo. Before serialization, it must classify dependencies strictly:
  - `blocked_by` is only for hard blockers: the task cannot execute or validate until the target is completed.
  - For every proposed `blocked_by`, the skill must be able to answer: “What would objectively fail if this task ran first?”
  - Sequencing preference, shared context, provenance, thematic grouping, or “use its output if available” must stay in task context/source-of-truth prose and must not be serialized as `blocked_by`.

  1. Save the approved structured plan in a temp file and run:
     roadmapctl materialize --plan <plan-json> --dry-run --repo <repo> --roadmap-root <roadmap-root> --output json > <dry-run-json>
  2. Save that dry-run JSON as the frozen change-set and verify it proposes only canonical allowlisted paths and no `*-tasks.md` fallback. Normal review shows only `summary`, `diagnostics`, and per-change `path`, `operation`, `applied`, `preconditions`; do not dump `changes[].content` or full diffs unless explicitly requested or troubleshooting.
  3. After explicit human approval of the dry-run, prefer roadmapctl-owned batch apply for the approved change-set:
     roadmapctl materialize --changes <dry-run-json> --apply --repo <repo> --roadmap-root <roadmap-root> --output json
     or equivalently:
     roadmapctl materialize --plan <plan-json> --apply --repo <repo> --roadmap-root <roadmap-root> --output json
  4. Use target apply only for recovery, troubleshooting, or explicit one-file approval:
     roadmapctl materialize --changes <dry-run-json> --target <target.path> --apply --repo <repo> --roadmap-root <roadmap-root> --output json
  5. Never use prompt-side raw writes for dry-run `content`.
  6. Run:
     roadmapctl check --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  7. If any command exits non-zero, report diagnostics and stop before claiming success or committing. If files were already applied, report `changes[].path` where `applied=true`, run `rootline validate` on affected markdown or rerun `roadmapctl check --strict`, then explicitly repair or revert and rerun the postcheck before any commit.
```

The materialized shape must remain canonical:

```text
<roadmap-root>/OXX-slug/README.md
<roadmap-root>/OXX-slug/TXXX-task.md
```

or:

```text
<roadmap-root>/TXXX-task.md
```

Never create a single fallback file such as `<roadmap-root>/feature-tasks.md` for multiple tasks. The skill must not duplicate numbering, `rootline new`, README task-table edits, dependency-link writing, or dry-run content writes once `roadmapctl materialize` is available; those deterministic writes belong to the CLI. Batch apply is the normal path because roadmapctl validates and writes the approved change-set atomically enough for the roadmap workflow. Granular writes use `--changes <dry-run-json> --target <path> --apply` only for recovery, troubleshooting, or explicit one-file approval so the target content still comes from a frozen CLI change-set.

## `/roadmap loop` integration snippet

Apply this behavior before discovery/execution and after each mutation during `/roadmap loop`.

```text
Before querying or executing pending tasks:
  1. Require `roadmapctl` on PATH.
  2. Run:
     roadmapctl doctor --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  3. Run:
     roadmapctl check --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  4. If either command exits non-zero, stop before selecting or executing tasks.

Before each task execution:
  1. Use `roadmapctl next` or `roadmapctl decision` JSON as the canonical readiness/blocker source after the preflight check passed.
  2. Do not execute tasks whose dependencies are reported invalid, unresolved, or blocked by `roadmapctl`.

After changing task status, links, dependencies, or task files:
  1. Run:
     roadmapctl check --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  2. If it exits non-zero, stop before marking the iteration complete or committing.
```

`/roadmap loop` may still run targeted task acceptance checks and existing Rootline commands, but those checks are additive. They do not replace `roadmapctl doctor` and `roadmapctl check`.

## Loop status transitions

The `/roadmap loop` skill must delegate task status transition policy and mutation to `roadmapctl transition`:

```bash
roadmapctl transition can-start <task.md> --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl transition start <task.md> --apply --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl transition complete <task.md> --apply --repo <repo> --roadmap-root <roadmap-root> --output json
```

The skill must not call `rootline set` directly for loop start/completion once `roadmapctl transition` is available. `transition start/complete --apply` owns Rootline `set`, target validation, and roadmap postcheck. The agent remains responsible for reading the task, implementing code, running task-specific ACs, committing, and pushing.

## Read-only state commands

The `/roadmap` skill must delegate deterministic read-only state to `roadmapctl` instead of rebuilding it from Rootline JSON in prompt logic:

```bash
roadmapctl pending --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl next --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl decision --repo <repo> --roadmap-root <roadmap-root> --output json
```

Workspace mode uses `roadmapctl pending --workspace --repo <workspace-root> --output json` for pending summaries. For decision/next, bootstrap resolves repos and the skill runs the single-repo command per repo, then only renders/group labels the returned JSON.

The skill must not call `rootline tree`, `rootline graph`, or `rootline query` directly for pending/next/decision. It must not postprocess Rootline JSON with Python snippets or recalculate done filters, blockers, reverse dependencies, quick wins, or scoring in prompt text. `roadmapctl` is the postprocessing boundary for roadmap workflows: agents consume its JSON fields directly instead of projecting raw Rootline JSON.

Headless verification for this cutover must show the read-only commands selected/used in addition to the mandatory doctor/check preflight:

```bash
roadmapctl pending --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl next --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl decision --repo <repo> --roadmap-root <roadmap-root> --output json
```

## Mandatory headless Pi verification

Every change to the `/roadmap` skill, this integration policy, or `roadmapctl` guard behavior must be tested with Pi in headless print mode before commit or release. Static grep is not enough: the verification must prove a real Pi agent loads the installed skill and chooses the `roadmapctl` guard before it would write, mutate, execute, or claim validity.

Run the reproducible verifier from the repository root after syncing the skill:

```bash
./scripts/sync-roadmap-skill.sh --install
./scripts/verify-roadmap-skill-headless.sh --evidence-dir /tmp/roadmap-headless-evidence
```

The script runs both pressure scenarios, captures logs, runs negative guard checks, and exits non-zero if required evidence is missing. Use `ROADMAP_HEADLESS_EVIDENCE_DIR=<dir>` or `--evidence-dir <dir>` to choose where logs are saved. For release/cutover reviews, attach or archive that evidence directory. A skill or guard change must not be released without this evidence.

Passing evidence must include:

- bootstrap checkpoint for the repo;
- `command -v roadmapctl`;
- `roadmapctl doctor --repo ... --roadmap-root ... --output json --strict`;
- `roadmapctl check --repo ... --roadmap-root ... --output json --strict`;
- an explicit statement that both commands were required and passed;
- no file modifications, no task execution, no commit, and no push during the headless verification.

Also run the negative guard checks:

```bash
roadmapctl check --repo testdata/fixtures/invalid-single-summary-file --output json --strict
roadmapctl check --repo testdata/fixtures/valid-outcome-with-tasks --rootline /tmp/no-such-rootline-roadmapctl --output json --strict
```

The first command must exit `1` with `RMC_STRUCTURE_SINGLE_FILE_FALLBACK`; the second must exit `3` with `RMC_ENV_ROOTLINE_MISSING`.

## Thin adapter audit

After the cutover, the `/roadmap` skill intentionally keeps only conversational and orchestration logic:

| Area | Remains in skill | Owned by roadmapctl |
|------|------------------|---------------------|
| Bootstrap | choose repo/workspace target, render checkpoint, stop on guard failure | `context` resolves config/schema/helpers |
| Pending/decision | human presentation and routing | `pending`, `next`, `decision` compute state, blockers and scoring |
| Loop | read task, implement code, run ACs, commit/push according to config | `transition` owns start/complete policy, status mutation and postcheck |
| Plan/materialize | decompose conceptually, ask approval, serialize structured plan, review dry-run | `materialize` owns numbering, canonical paths, writes, README tables, dependency links and postcheck |

Rootline commands may remain only as troubleshooting/reference or for loop graph/query discovery where no roadmapctl command owns that read yet. They must not be used as the primary writer/mutator when a roadmapctl command exists.

## Expected failures and messages

### `roadmapctl` missing

If `roadmapctl` is not installed or not on PATH, block before writes/mutations/execution.

Suggested message:

```text
`roadmapctl` no está instalado o no está en PATH. Es requerido para comandos `/roadmap` implementados que escriben, mutan, ejecutan o declaran validez del roadmap.
Instalar o exponer el binario `roadmapctl` y reintentar.
```

No fallback is allowed when `roadmapctl` is missing.

### Rootline missing

If Rootline is missing, `roadmapctl doctor` or `roadmapctl check` should return exit code `3` with diagnostic `RMC_ENV_ROOTLINE_MISSING`.

Suggested message:

```text
`rootline` no está instalado o no fue encontrado por `roadmapctl`.
Diagnostic: RMC_ENV_ROOTLINE_MISSING
Instalar con: curl -fsSL https://raw.githubusercontent.com/pablontiv/rootline/master/install.sh | bash
```

The skill must not call Rootline-specific fallback logic to bypass this failure.

### Invalid roadmap structure

If `roadmapctl check` reports `RMC_STRUCTURE_SINGLE_FILE_FALLBACK`, `RMC_STRUCTURE_MISSING_OUTCOME_README`, `RMC_STRUCTURE_DUPLICATE_ID`, or related structure diagnostics, block and report the diagnostics.

Do not auto-fix structure or create a replacement summary document. The user or an approved implementation task must correct the canonical files.

### Invalid dependencies or graph

If `roadmapctl check` reports `RMC_GRAPH_INVALID_BLOCKED_BY` or `RMC_GRAPH_CYCLE`, block task execution and materialization success claims.

The skill must use resolved graph data after a successful check; it must not grep wikilinks or ignore broken `blocked_by` links.

### Status or schema mismatch

If `roadmapctl check` reports `RMC_STATUS_UNKNOWN`, `RMC_STATUS_TYPE_UNKNOWN`, or `RMC_ROOTLINE_VALIDATE_FAILED`, block and report the offending path/diagnostic.

Do not mutate statuses to force a pass unless the user explicitly approved a repair task.

## No auto-fix and no freeform fallback

`roadmapctl` is a guard, not a materializer or repair tool. The skill integration must not use it as justification to auto-fix roadmaps.

Prohibited fallback behaviors:

- creating one markdown file that contains multiple task descriptions;
- skipping canonical `OXX-*`/`TXXX-*` files because validation failed;
- editing installed user-scope skills as part of ordinary plan/loop execution;
- modifying Rootline or adding roadmap-specific Rootline commands;
- silently repairing broken links, duplicate IDs, or invalid statuses.

When the guard fails, stop and explain the failure with actionable diagnostics.

## Relationship to Rootline

`roadmapctl` may invoke Rootline as an external executable. `/roadmap` should treat `roadmapctl` as the roadmap-specific policy layer and Rootline as the generic DBMS/constraint engine.

The skill may still use Rootline commands for materialization and querying after `roadmapctl` preflight succeeds, but commands that write, mutate, execute, or claim validity must be guarded by `roadmapctl`.
