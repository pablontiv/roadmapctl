# `/roadmap` skill integration with roadmapctl

`roadmapctl` is the required deterministic boundary for implemented `/roadmap`
flows that write roadmap files, mutate roadmap state, execute tasks, or claim
that a roadmap is valid.

The current live CLI surface is:

- `roadmapctl doctor`
- `roadmapctl check`
- `roadmapctl bootstrap`
- `roadmapctl bootstrap inspect`
- `roadmapctl bootstrap init`
- `roadmapctl pending`
- `roadmapctl next`
- `roadmapctl decision`
- `roadmapctl lint`
- `roadmapctl transition`

The skill must not call or document removed command surfaces such as
`roadmapctl context`, `roadmapctl materialize`, `roadmapctl plan-paths`, or the
old `--roadmap-root` flag. The roadmap root is fixed at `docs/roadmap/` under
the selected repo and is exposed as `roadmap_root` in JSON reports.

## Responsibility Boundary

| Layer | Owns | Does not own |
|-------|------|--------------|
| Rootline | Generic Markdown database behavior: `.stem`, frontmatter, links, validation, tree/query/graph/describe primitives. | Roadmap-specific status policy, prioritization, task execution, commits, PRs, or AI planning. |
| roadmapctl | Roadmap-specific guards, config, status roles, read-model commands, transition policy, diagnostics, and JSON contracts. | Free-form planning, prose generation, code implementation, commits, PRs, or project-local agent rules. |
| `/roadmap` skill | Intent clarification, decomposition, user approval, orchestration, rendering, and task execution workflow. | Recomputing roadmap policy already exposed by `roadmapctl`, bypassing guards, or carrying repo-specific rules in the shared skill. |

Conceptual planning may run without `roadmapctl` only when it does not write
files, mutate state, execute tasks, or claim validity.

## Bootstrap

Every implemented subcommand starts by resolving context from `bootstrap`:

```bash
roadmapctl bootstrap --repo <repo-path> --output json
```

The JSON is the source of truth for `root`, `roadmap_root`, `config_path`,
`status_values`, `done_statuses`, `active_statuses`, `helpers`,
`required_code_coverage`, `loop_max_tasks`, `parallel`, `autonomy`,
`compact_after_task_commit`, `pr_mode`, `commit_style`, `auto_push`, and
`gitflow` style fields.

If `bootstrap` returns `RMC_GITFLOW_NOT_CONFIGURED`, the skill may run the
gitflow adoption wizard described in the shared skill docs. That wizard updates
only the `[gitflow]` style fields after user confirmation.

For missing roadmap roots, the only setup path is:

```bash
roadmapctl bootstrap inspect --repo <repo-path> --output json
roadmapctl bootstrap init --repo <repo-path> --dry-run --output json
roadmapctl bootstrap init --repo <repo-path> --apply --output json
roadmapctl check --repo <repo-path> --output json --strict
```

`bootstrap init --apply` requires explicit approval after reviewing the dry-run.
It is not a repair path for invalid existing roadmaps.

## Guards

Before writing roadmap files, mutating task state, executing tasks, committing
roadmap mutations, or claiming validity:

```bash
command -v roadmapctl
roadmapctl doctor --repo <repo-path> --output json --strict
roadmapctl check --repo <repo-path> --output json --strict
```

After any roadmap file creation or mutation:

```bash
roadmapctl check --repo <repo-path> --output json --strict
```

If any required command exits non-zero, the skill stops before the guarded
operation. It reports the command, exit code, and diagnostic IDs when JSON was
produced. It must not auto-fix, write fallback summary markdown, execute tasks,
or commit while the guard is failing.

## Plan Flow

`/roadmap plan` is the approved Markdown writer for roadmap planning files:

1. Resolve repo/config with `bootstrap`.
2. Decompose intent into direct tasks or Outcomes with child tasks.
3. Present the exact tree and content for explicit user approval.
4. Run the mandatory `doctor` and `check` preflight.
5. Write canonical files only:
   - `docs/roadmap/TXXX-task.md`
   - `docs/roadmap/OXX-outcome/README.md`
   - `docs/roadmap/OXX-outcome/TXXX-task.md`
6. Run `roadmapctl check --repo <repo-path> --output json --strict`.
7. Stop after creating planning files; implementation belongs to `/roadmap loop`.

Never create a single fallback file such as `docs/roadmap/feature-tasks.md` for
multiple tasks. Outcome README files contain frontmatter, title, and context;
acceptance criteria live in task files.

## Loop Flow

`/roadmap loop` executes one repo at a time:

1. Resolve repo/config with `bootstrap`.
2. Run `doctor` and `check` preflight.
3. Read readiness from:
   ```bash
   roadmapctl next --repo <repo-path> --output json
   ```
4. Read active task counts from:
   ```bash
   roadmapctl pending --repo <repo-path> --output json
   ```
5. Start a task only through:
   ```bash
   roadmapctl transition can-start <task.md> --repo <repo-path> --output json
   roadmapctl transition start <task.md> --apply --repo <repo-path> --output json
   ```
6. Implement only the active task scope and run its acceptance checks.
7. Complete the task only through:
   ```bash
   roadmapctl transition complete <task.md> --apply --repo <repo-path> --output json
   ```
8. Delegate commit/push/PR work to `/integrate` if integration is required.

The skill must not call `rootline set` directly for task start/completion.

## Read-Only State

For pending, next-task, and decision views, use only roadmapctl JSON:

```bash
roadmapctl pending --repo <repo-path> --output json
roadmapctl next --repo <repo-path> --output json
roadmapctl decision --repo <repo-path> --output json
```

Do not call Rootline `tree`, `query`, or `graph` directly to reconstruct
pending queues, blockers, reverse dependencies, quick wins, or scoring.

Workspace summaries use:

```bash
roadmapctl pending --workspace --repo <workspace-root> --output json
```

For `next` and `decision`, run the single-repo command for each selected member
and render grouped results without recalculating policy.

## Local Agent Rules

Shared skills stay portable across repos. Project-specific verification,
toolchain, CI, lint, or shell-hook rules belong in repo-local docs. In this
repo, use `docs/local-agent-workflow.md` for roadmapctl-specific guidance.

Local docs may add checks, but they must not weaken these shared guards:

- explicit approval before roadmap writes;
- `doctor` and `check` before writes/mutations/execution;
- `check` after roadmap mutations;
- no fallback summary files;
- no direct status mutation outside `roadmapctl transition`.

## Verification

For changes to this integration policy or the shared skills, verify the live
contract rather than historical docs:

```bash
./roadmapctl --help
./roadmapctl bootstrap --help
./roadmapctl decision --help
./roadmapctl transition complete --help
./roadmapctl bootstrap --repo . --output json
./roadmapctl check --repo . --output json --strict
./roadmapctl pending --repo . --output json
./roadmapctl next --repo . --output json
./roadmapctl decision --repo . --output json
./roadmapctl lint --repo testdata/fixtures/lint-valid --output json --strict
scripts/sync-roadmap-skill.sh --check --all
```

Negative guard checks should include a known invalid roadmap fixture and a
missing Rootline binary check.
