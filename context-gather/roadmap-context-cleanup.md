# Code Context

## Files Retrieved

1. `.claude/skills/roadmap/loop-subcommand.md` (lines 1-182) - primary loop procedure, commit boundary, checkpoint interval, self-pace/worktree/parallel prompt behavior.
2. `.claude/skills/roadmap/SKILL.md` (lines 1-205) - skill frontmatter hooks/tools, bootstrap/context/config contract, legacy config template, mandatory bootstrap checkpoint.
3. `.claude/skills/roadmap/pr-workflow.md` (lines 1-90) - PR branch/push/create/merge/post-merge cleanup workflow for `loop --pr`.
4. `docs/roadmap-skill-integration.md` (lines 1-120, 158-224, 232-274) - authoritative skill/roadmapctl boundary, guard/postcheck/commit blocking rules, headless verification.
5. `docs/cli-contract.md` (lines 67-115, 176-188) - `.roadmapctl.toml` config contract and read-only context command shape.
6. `docs/transition-controller.md` (lines 1-66) - transition controller boundaries; explicitly excludes AC/tests/commits/PR from roadmapctl.
7. `README.md` (lines 1-65) - layer responsibilities: roadmapctl deterministic governance, skill orchestration, implementing agent commits.
8. `internal/config/config.go` (lines 43-59, 84-160, 227-267, 269-307, 340-382) - config struct, TOML fields, defaults, legacy YAML mapping, diff/warning behavior.
9. `internal/cli/context.go` (lines 18-53, 55-83, 95-150) - JSON context report fields and workspace repo context; current output omits workflow config fields.
10. `internal/cli/transition.go` (lines 18-147) - transition JSON/report, apply flow using Rootline `set`, validate-one, and postcheck.
11. `internal/config/config_test.go` (lines 22-105, 255-300) - tests proving TOML precedence/defaults and legacy config parsing for workflow fields.
12. `internal/cli/context_test.go` (lines 9-90) - tests for workspace context and helper/schema fields.
13. `testdata/fixtures/valid-roadmapctl-toml-default/docs/roadmap/.roadmapctl.toml` (lines 1-15) - fixture with current default TOML keys.
14. `.claude/roadmap.local.md` (lines 1-17) - current repo legacy config; no `docs/roadmap/.roadmapctl.toml` exists in this working tree.
15. `.githooks/post-merge` (lines 1-5) - only repo hook found; installs/syncs roadmap skill after merge.
16. `scripts/sync-roadmap-skill.sh` (lines 1-123) - sync/check script for canonical skill install; cleanup only temp/backup dirs during install.
17. `scripts/verify-roadmap-skill-headless.sh` (lines 1-88) - headless guard verification; asserts no modifications/commit/push in verification.

## Key Code

### Existing commit/checkpoint/context-management behavior

- Loop options include only prompt-level checkpointing/self-pacing, not context cleanup:

```md
# .claude/skills/roadmap/loop-subcommand.md:9-16
- `--checkpoint-interval N`: checkpoint de calidad cada N tasks (default 5).
- `--worktree`: crear git worktree aislado ... al cerrar se limpia con `ExitWorktree`.
- `--self-pace`: usar `ScheduleWakeup` ... Mantiene el cache caliente ...
- `--parallel`: ejecutar tasks independientes ... via `Agent` tool ...
```

- Per-task commit boundary is after ACs/invariants and `roadmapctl transition complete --apply`; if transition/check fails, do not commit:

```md
# .claude/skills/roadmap/loop-subcommand.md:136-140
roadmapctl transition complete <task.md> --apply ...
Ejecutar este comando solo después de que ACs e invariantes pasaron. ... detenerse antes de declarar completada la iteración o commitear. Si pasa: `git add` específico, commit según `<commit-style>`, push según `<auto-push>` y `--pr`.
```

- Quality checkpoint is a manual/prompt review of accumulated diff, triggered every N tasks/scope change/user stop:

```md
# .claude/skills/roadmap/loop-subcommand.md:156-162
11. **Checkpoint**
Activar si:
- `checkpoint_task_count >= checkpoint_interval`,
- cambia scope,
- usuario decide parar.
Revisar diff acumulado, reportar findings informativos y resetear checkpoint.
```

- Mandatory bootstrap context is rendered from `roadmapctl context`; config preference is TOML then legacy. Current checkpoint is informational output only:

```md
# .claude/skills/roadmap/SKILL.md:96-120, 193-201
Preferir `roadmapctl context` ...
roadmapctl context --repo <repo-path> --roadmap-root <roadmap-root-si-se-conoce> --output json
Usar el JSON devuelto como fuente de verdad para ... helpers/status...
Checkpoint obligatorio: Bootstrap: roadmap-root ... where_* ...
```

- PR mode cleanup is Git branch/base refresh after merge, not agent context cleanup:

```md
# .claude/skills/roadmap/pr-workflow.md:83-90
### Post-merge cleanup
git -C <repo-path> checkout <base_branch>
git -C <repo-path> pull origin <base_branch>
Registrar `{number, url, scope, status}`...
```

### Existing mechanism vs skill-level prompt behavior

- No repo code/docs mention `compact`, `context cleanup`, or session continuation cleanup. Focused search only found unrelated handoff references in roadmap task source-of-truth and `.pi/agents` lean-agent descriptions.
- The only runtime-ish hook in skill frontmatter is a Stop agent critic check, not compaction/cleanup:

```md
# .claude/skills/roadmap/SKILL.md:31-35
hooks:
  Stop:
    - type: agent
      prompt: "Verify that the critic agent was invoked ..."
      timeout: 60
```

- `roadmapctl` owns deterministic state/policy and transition mutation, but not task execution, commits, PRs, or agent/session lifecycle:

```md
# docs/roadmap-skill-integration.md:267-272
| Loop | read task, implement code, run ACs, commit/push according to config | `transition` owns start/complete policy, status mutation and postcheck |
```

```md
# docs/transition-controller.md:54
`can-complete` does not execute acceptance checks, tests, commits, or PR work. The caller remains responsible...
```

### roadmapctl config/context contracts

Current TOML config keys are snake_case and include loop workflow knobs:

```toml
# docs/cli-contract.md:72-88
done_statuses = ["Completed", "Obsolete"]
active_statuses = ["Pending", "Specified", "In Progress"]
leaf_filter = "isIndex == false"
outcome_close_verify = []
pr_merge_strategy = "squash"
commit_style = "conventional"
auto_push = true

[status_values]
...
```

Go config currently supports only these TOML workflow fields:

```go
// internal/config/config.go:55-59, 151-159
OutcomeCloseVerify []string
PRMergeStrategy    string
CommitStyle        string
AutoPush           bool

type tomlConfig struct { ... OutcomeCloseVerify []string `toml:"outcome_close_verify"`; PRMergeStrategy string `toml:"pr_merge_strategy"`; CommitStyle string `toml:"commit_style"`; AutoPush *bool `toml:"auto_push"` }
```

`roadmapctl context` JSON exposes config source, schema, status values, done/active statuses, and helpers only; it does **not** expose `outcome_close_verify`, `pr_merge_strategy`, `commit_style`, or `auto_push` today:

```go
// internal/cli/context.go:18-33
ConfigPath string `json:"config_path"`
ConfigSource string `json:"config_source"`
Schema contextSchema `json:"schema"`
StatusValues config.StatusValues `json:"status_values"`
DoneStatuses []string `json:"done_statuses"`
ActiveStatuses []string `json:"active_statuses"`
Helpers contextHelpers `json:"helpers"`
```

Likely field placement if adding a repo-configurable cleanup/compact option:

- TOML: near existing loop workflow config (`outcome_close_verify`, `pr_merge_strategy`, `commit_style`, `auto_push`) in `<roadmap-root>/.roadmapctl.toml`.
- Naming convention: snake_case boolean, e.g. `clean_context_after_task_commit = false` or `compact_context_after_task_commit = false`.
- Legacy YAML analogue, if still supported: kebab-case, e.g. `clean-context-after-task-commit: false`.
- Context JSON: if the skill must decide from `roadmapctl context`, add a snake_case JSON field or nested workflow/config object; otherwise the skill cannot discover the setting from current context output.

## Architecture

- User invokes `/roadmap`; the skill bootstraps with `roadmapctl context`, runs preflight `doctor`/`check`, then uses `next`/`pending` for deterministic queue discovery.
- For each task, skill calls `transition can-start`, applies `transition start --apply`, reads/implements task, runs ACs/invariants, then applies `transition complete --apply`. Only after successful completion transition/postcheck does skill perform `git add`, commit, and optional push/PR.
- `roadmapctl transition --apply` mutates only roadmap status via Rootline `set`, `ValidateOne`, and `runCheck`; it does not inspect or manage LLM/agent context.
- Config loading prefers `<roadmap-root>/.roadmapctl.toml`; legacy `.claude/roadmap.local.md` is still supported. Current repo itself only has `.claude/roadmap.local.md`; fixture TOMLs show intended shape.
- Hooks/scripts are about installing/verifying the skill, not loop runtime cleanup. `sync-roadmap-skill.sh` temp-dir cleanup is shell housekeeping only.

## Existing behavior summary

- Commit boundary: documented in loop and integration docs as after ACs/invariants plus `transition complete --apply` and postcheck; failures block commit.
- Checkpointing: prompt-only quality checkpoint every `--checkpoint-interval` tasks (default 5), on scope change, or user stop; also bootstrap checkpoint prints resolved roadmap helpers.
- Session/context management: only `--self-pace` keeps cache warm with `ScheduleWakeup`; no compaction/cleanup mechanism was found.
- Agent handoffs/parallelism: `--parallel` can use `Agent` only for provably independent tasks; no persistent handoff or context-cleaning contract.
- PR cleanup: branch checkout/pull after merge only.

## Risks/constraints for skill vs roadmapctl

- Skill-side implementation matches current ownership: commit/push and agent/session lifecycle are skill/agent responsibilities, not roadmapctl.
- roadmapctl-side config exposure may still be needed because the skill is instructed to use `roadmapctl context` as source of truth; current context JSON omits workflow knobs, so adding only a TOML parser field may be invisible to the skill.
- A cleanup/compact action likely depends on host/agent capabilities, not Go CLI capabilities; roadmapctl has no API to invoke LLM context compaction and should remain deterministic/read-write roadmap state tooling.
- Running cleanup immediately after commit must not happen before task status completion/postcheck/commit hash summary; otherwise loop evidence/UI could be lost.
- `--parallel`, worktrees, and PR mode create multiple scopes/agents/branches; cleanup timing must not discard unresolved child-agent results or PR bookkeeping.
- Existing headless verification forbids modifications/commit/push during guard tests; any skill/guard behavior change must update/run `scripts/verify-roadmap-skill-headless.sh` evidence expectations.

## Remaining clarification questions

1. Should the option mean “after every successful per-task commit” only, or also after roadmap materialization commits from `/roadmap plan`?
2. Is the desired action named/implemented by the host as “compact” with a callable tool/command, or is this only a prompt instruction the skill can emit?
3. Should default be off for compatibility, or on for long autonomous loops?
4. Should cleanup run before or after `git push`/PR operations when `<auto-push>` or `--pr` is enabled?
5. Must legacy `.claude/roadmap.local.md` support the setting, or only preferred `.roadmapctl.toml`?

## Start Here

Start with `.claude/skills/roadmap/loop-subcommand.md` lines 136-165: it is the exact per-task post-success boundary where commit, UI summary, self-pace, checkpoint, and continuation are currently ordered.

## Supervisor coordination

No supervisor contact needed; repository inspection only, no code changes made except this requested context artifact.
