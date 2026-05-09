# `/roadmap` skill integration with roadmapctl

`roadmapctl` is mandatory from day 1 for every implemented `/roadmap` command that writes files, mutates roadmap state, executes tasks, or claims that a roadmap is valid.

This document describes how the skill must use `roadmapctl doctor` and `roadmapctl check` as blocking guards while preserving Rootline as the generic filesystem database and constraint engine.

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

This postcheck is mandatory before reporting successful materialization, completing a loop iteration, committing roadmap mutations, or claiming that the roadmap is valid.

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
  2. Run:
     roadmapctl doctor --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  3. Run:
     roadmapctl check --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  4. If either command exits non-zero, stop. Do not write files and do not fall back to a summary markdown file.

After materializing canonical files:
  1. Run:
     roadmapctl check --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  2. If it exits non-zero, report diagnostics and stop before claiming success or committing.
  3. Continue to existing rootline validation only as additional evidence, not as a replacement for roadmapctl.
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

Never create a single fallback file such as `<roadmap-root>/feature-tasks.md` for multiple tasks.

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
  1. Use Rootline graph/query data only after the preflight check passed.
  2. Do not execute tasks whose dependencies are invalid or unresolved.

After changing task status, links, dependencies, or task files:
  1. Run:
     roadmapctl check --repo <repo> --roadmap-root <roadmap-root> --output json --strict
  2. If it exits non-zero, stop before marking the iteration complete or committing.
```

`/roadmap loop` may still run targeted task acceptance checks and existing Rootline commands, but those checks are additive. They do not replace `roadmapctl doctor` and `roadmapctl check`.

## Read-only state commands

The `/roadmap` skill must delegate deterministic read-only state to `roadmapctl` instead of rebuilding it from Rootline JSON in prompt logic:

```bash
roadmapctl pending --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl next --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl decision --repo <repo> --roadmap-root <roadmap-root> --output json
```

Workspace mode uses `roadmapctl pending --workspace --repo <workspace-root> --output json` for pending summaries. For decision/next, bootstrap resolves repos and the skill runs the single-repo command per repo, then only renders/group labels the returned JSON.

The skill must not call `rootline tree`, `rootline graph`, or `rootline query` directly for pending/next/decision. It must not recalculate done filters, blockers, reverse dependencies, quick wins, or scoring in prompt text.

Headless verification for this cutover must show the read-only commands selected/used in addition to the mandatory doctor/check preflight:

```bash
roadmapctl pending --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl next --repo <repo> --roadmap-root <roadmap-root> --output json
roadmapctl decision --repo <repo> --roadmap-root <roadmap-root> --output json
```

## Mandatory headless Pi verification

Every change to the `/roadmap` skill, this integration policy, or `roadmapctl` guard behavior must be tested with Pi in headless print mode before commit or release. Static grep is not enough: the verification must prove a real Pi agent loads the installed skill and chooses the `roadmapctl` guard before it would write, mutate, execute, or claim validity.

Run both pressure scenarios from the repository root after syncing the skill:

```bash
./scripts/sync-roadmap-skill.sh --install

PI_SKIP_VERSION_CHECK=1 pi \
  --no-extensions \
  --skill .claude/skills/roadmap/SKILL.md \
  --tools read,bash \
  -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: the user asks "loop autonomo" in this repository. Do not modify files and do not run git commit/push. Perform only the bootstrap and the required preflight checks from the skill, then stop. In your final answer, list the exact commands you ran and whether roadmapctl doctor/check were required and passed.'

PI_SKIP_VERSION_CHECK=1 pi \
  --no-extensions \
  --skill .claude/skills/roadmap/SKILL.md \
  --tools read,bash \
  -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: there is an already approved plan to materialize one direct task, and the user says "crea las tareas". Do not create or modify files and do not run git commit/push. Perform only bootstrap and the required preflight checks that must happen before any roadmap write, then stop. In your final answer, list exact commands run and whether roadmapctl doctor/check were required and passed.'
```

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
