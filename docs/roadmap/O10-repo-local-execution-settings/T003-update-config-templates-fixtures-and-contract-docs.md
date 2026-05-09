---
estado: Completed
tipo: task
---
# T003: Update config templates, fixtures, and CLI contract docs

**Outcome**: [Repo-local roadmap execution settings](README.md)

[[blocked_by:./T001-implement-config-fields-and-forced-migration.md]]

## Preserva

- Generated TOML remains valid for roadmapctl Load.
- Docs keep TOML as the per-repo authority and do not reintroduce legacy as a lasting source.

## Contexto

The bootstrap and materialize packages both carry default .roadmapctl.toml templates. Docs currently describe legacy fallback and conflict warnings that will be replaced by forced migration/deletion.

## Alcance

**In**:
1. Update defaultRoadmapctlTOML in internal/cli/bootstrap.go and internal/materialize/dryrun.go.
2. Update test fixtures containing .roadmapctl.toml as needed.
3. Update docs/cli-contract.md configuration keys, defaults, precedence, and conflict/migration policy.
4. Update any config-related docs that still say legacy remains a fallback for implemented write/execution flows.

**Out**:
1. No skill loop behavior edits in this task.
2. No Pi extension work in this task.

## Estado inicial esperado

Default .roadmapctl.toml templates include status/workflow keys but not loop execution settings. docs/cli-contract.md still documents legacy fallback and TOML/legacy conflict warnings.

## Especificación Técnica

Keep TOML key order stable across templates and renderTOMLConfig: status lists, leaf filter, outcome_close_verify, PR/commit/push settings, then loop settings, then [status_values] is acceptable if tests expect it. Update docs to say command-line flags do not override behavior settings except --max and --filter in the skill layer.

## Criterios de Aceptación

- go test ./internal/cli ./internal/materialize passes.
- Generated default TOML contains loop_max_tasks, parallel, autonomy, compact_after_task_commit, and pr_mode with approved defaults.
- docs/cli-contract.md states that legacy .claude/roadmap.local.md is migration input only and is deleted after successful config load/migration.
- No docs under docs/ describe --parallel, --worktree, --self-pace, --skip-reviews, --checkpoint-interval, or --pr as active /roadmap loop behavior flags.

## Fuente de verdad

- docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md
- internal/cli/bootstrap.go
- internal/materialize/dryrun.go
- docs/cli-contract.md
- testdata/fixtures/valid-roadmapctl-toml-default/docs/roadmap/.roadmapctl.toml
