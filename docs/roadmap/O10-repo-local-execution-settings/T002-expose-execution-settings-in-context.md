---
estado: Specified
tipo: task
---
# T002: Expose execution settings in roadmapctl context

**Outcome**: [Repo-local roadmap execution settings](README.md)

[[blocked_by:./T001-implement-config-fields-and-forced-migration.md]]

## Preserva

- context remains the skill source of truth for effective config.
- Existing status/helper fields remain stable except for intentional golden updates.

## Contexto

The skill must not parse TOML or legacy config. It needs the new config values from roadmapctl context before deciding loop behavior.

## Alcance

**In**:
1. Add execution settings to contextReport JSON.
2. Add execution settings to workspaceRepoContext JSON for repo selection flows.
3. Update context tests and golden snapshots.

**Out**:
1. No changes to rootline schema discovery.
2. No changes to pending/next/decision algorithms.

## Estado inicial esperado

internal/cli/context.go exposes status values, done/active statuses, and helpers, but omits auto_push, commit_style, pr_merge_strategy, outcome_close_verify, and the new loop execution settings.

## Especificación Técnica

Prefer a nested JSON object named execution or operational_config only if all callers/tests are updated consistently; otherwise add explicit snake_case fields at the top level to minimize skill parsing. Whichever shape is chosen, use the same field names in workspace repo entries. Preserve parseable JSON stdout and empty stderr on success.

## Criterios de Aceptación

- go test ./internal/cli -run 'TestContext|TestCheckGoldenJSONFixtures' passes after golden updates.
- roadmapctl context --output json includes loop_max_tasks, parallel, autonomy, compact_after_task_commit, pr_mode, auto_push, commit_style, pr_merge_strategy, and outcome_close_verify for a single repo.
- roadmapctl context --workspace --output json includes the same operational settings for each repo entry.

## Fuente de verdad

- docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md
- internal/cli/context.go
- internal/cli/context_test.go
- testdata/golden/context-valid-legacy-config-fallback.json
- testdata/golden/context-valid-workspace.json
