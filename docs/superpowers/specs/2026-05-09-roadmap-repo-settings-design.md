# Roadmap repo-local execution settings design

## Goal

Make `/roadmap` execution behavior deterministic per repository by moving loop behavior options out of skill flags and into the canonical repo config at `<roadmap-root>/.roadmapctl.toml` (normally `docs/roadmap/.roadmapctl.toml`). Legacy `.claude/roadmap.local.md` is read only once as migration input and then removed.

## Decisions

- Canonical config lives in `<roadmap-root>/.roadmapctl.toml`.
- `roadmapctl` owns config loading, legacy migration, defaults, validation, and `context` exposure.
- The roadmap skill consumes `roadmapctl context`; it does not parse or migrate config itself.
- Legacy `.claude/roadmap.local.md` is not a lasting config source.
- `/roadmap loop` keeps only scope/selection flags: `--filter` and `--max`.
- Behavior flags are removed from the skill: `--parallel`, `--worktree`, `--self-pace`, `--skip-reviews`, `--checkpoint-interval`, and `--pr`.
- In this version, quality gates remain enabled, self-pacing is removed, and checkpoint cadence is not user-configurable unless a later design adds explicit config fields.

## Config model

Add these fields to `.roadmapctl.toml`:

```toml
loop_max_tasks = 0
parallel = true
autonomy = "until_done"
compact_after_task_commit = true
pr_mode = false
```

Defaults:

- `loop_max_tasks = 0`: unlimited loop execution.
- `parallel = true`: use opportunistic parallel waves when safe.
- `autonomy = "until_done"`: continue until the queue is exhausted or the max limit is reached.
- `compact_after_task_commit = true`: compact context after each completed task commit/push/PR bookkeeping.
- `pr_mode = false`: branch/PR workflow is disabled unless enabled by config.

Existing fields remain in the same config:

```toml
pr_merge_strategy = "squash"
commit_style = "conventional"
auto_push = true
```

`--max N` remains as a one-run cap overriding `loop_max_tasks`; `--filter` remains a one-run task selection filter. No `--no-*` behavior overrides are added.

## Migration semantics

When `roadmapctl` reads config:

1. If `<roadmap-root>/.roadmapctl.toml` exists:
   - use TOML as the only effective config source;
   - if `.claude/roadmap.local.md` also exists, delete it;
   - if TOML is invalid, fail and do not fall back to legacy.
2. If TOML does not exist but `.claude/roadmap.local.md` exists:
   - read legacy frontmatter only as migration input;
   - generate `<roadmap-root>/.roadmapctl.toml` with preserved values plus new defaults;
   - validate the generated config;
   - delete `.claude/roadmap.local.md` only after successful write and validation;
   - continue with TOML as the effective config.
3. If neither config exists, use the existing missing/default behavior, but new defaults must be included whenever a config is generated.

Migration must be idempotent and safe: never delete legacy if TOML generation or validation fails.

## Loop semantics

### Autonomy

`autonomy` accepts:

- `manual`: run one task/wave at a time and ask whether to continue. If missing dependency information is discovered, suggest the change and stop.
- `supervised`: continue across tasks/waves without asking, but ask before structural roadmap edits such as adding `blocked_by`.
- `until_done`: continue until `loop_max_tasks`/`--max` or no ready tasks remain. It may apply safe structural roadmap repairs, revalidate, and continue.

### Parallel waves

When `parallel = true`, the loop uses opportunistic waves:

- Query `roadmapctl next` for ready/blocked state.
- Treat `blocked_by` / roadmap dependencies as the only canonical source of dependency ordering.
- Execute the largest wave of ready tasks that have no explicit dependency ordering between them.
- If no parallel wave is available, execute sequentially.
- Parallel task execution must use isolated worktrees or an equivalent conflict-controlled integration path. There is no per-invocation `--worktree` flag in this design.
- After each wave, consolidate results, validate, commit each completed task according to the existing commit policy, then recalculate the queue.

No heuristic file/scope dependency inference is required for deciding parallel eligibility. If integration conflicts prove a dependency was missing:

- `manual`: report the recommended `blocked_by` relation and stop.
- `supervised`: ask before applying the recommended `blocked_by` relation.
- `until_done`: apply the recommended `blocked_by`, run strict checks, recalculate, and continue if safe.

### PR mode

`pr_mode` replaces `--pr`. Existing `pr_merge_strategy` and `auto_push` continue to control merge strategy and push behavior.

### Context compaction

When `compact_after_task_commit = true`, compaction happens after the task is fully durable:

1. task implementation and acceptance checks pass;
2. `roadmapctl transition complete --apply` succeeds;
3. a concise iteration summary is prepared for the commit context;
4. conventional commit is created;
5. push/PR bookkeeping completes if enabled;
6. roadmap context is compacted.

Preferred mechanism: a Pi extension/tool from this repo, e.g. `compact_roadmap_context`, calls Pi compaction with roadmap-specific instructions.

Fallback: if the extension/tool is unavailable, the skill queues or invokes `/compact <roadmap-specific instructions>`. Failure to compact should warn clearly; it should not invalidate an already completed task commit.

## Responsibility boundaries

### `roadmapctl`

Owns:

- config loading;
- forced legacy migration while reading config;
- deletion of legacy config after successful TOML migration or when TOML already exists;
- defaults and validation;
- exposing operational settings via `roadmapctl context`;
- config tests and context JSON contract.

Does not own:

- implementing tasks;
- launching agents;
- compacting Pi context;
- committing, pushing, or creating PRs;
- resolving merge conflicts.

### Roadmap skill

Owns:

- invoking `roadmapctl context` as source of truth;
- executing loop behavior according to config;
- removing behavior flags from user-facing loop docs;
- opportunistic parallel wave orchestration;
- applying autonomy-specific decisions;
- invoking compaction tool or `/compact` fallback.

Does not own:

- parsing legacy config;
- performing migration itself;
- using legacy as fallback after TOML exists.

### Pi extension

Lives in this repo as the canonical implementation for roadmap-specific compaction. It registers a tool that triggers Pi compaction with instructions preserving:

- current roadmap goal and repo;
- completed task path and commit hash;
- validation results;
- next task/wave state;
- unresolved blockers or conflicts;
- config values relevant to continuing the loop.

## Validation plan

Add or update tests for:

- legacy-only config migrates to TOML and deletes legacy;
- migration preserves existing legacy values;
- migration adds new defaults;
- TOML existing plus legacy deletes legacy and uses TOML;
- invalid TOML fails without legacy fallback;
- invalid `autonomy` fails;
- negative `loop_max_tasks` fails;
- `roadmapctl context` exposes new operational fields;
- generated default TOML includes new fields;
- skill docs no longer advertise removed behavior flags;
- loop docs describe config-driven behavior and opportunistic waves.

Suggested verification commands:

```bash
go test ./internal/config
go test ./internal/cli
go test ./...
```

If legacy fixtures remain for migration tests, validate their schema as needed before migration assertions.

## Open risks

- Config loading will now have migration side effects; tests must make this explicit and ensure read-only commands remain predictable after migration.
- Parallel wave execution can surface missing `blocked_by` edges; the autonomy rules define how to repair or stop.
- Compaction is Pi-runtime behavior, so the roadmap skill needs a robust fallback when the extension is absent.
