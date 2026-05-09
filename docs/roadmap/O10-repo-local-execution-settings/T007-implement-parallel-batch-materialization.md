---
estado: Specified
tipo: task
---
# T007: Implement parallel batch materialization

**Outcome**: [Repo-local roadmap execution settings](README.md)

[[blocked_by:./T001-implement-config-fields-and-forced-migration.md]]

## Preserva

- The skill does not write roadmap markdown directly.
- roadmapctl remains the deterministic owner of materialization writes.
- Materialization still rejects non-canonical targets such as `*-tasks.md`.

## Contexto

The initial repo-settings design kept the existing granular `--changes --target` apply flow in the skill. The clarified direction is different: if roadmapctl receives a structured plan or frozen change-set, roadmapctl itself should own safe multi-file apply and may use internal parallelism when `parallel = true`. The skill should not orchestrate one process invocation per target.

## Alcance

**In**:
1. Define and implement a roadmapctl-owned batch materialization apply path that can create multiple canonical roadmap files from a plan or frozen change-set in one command.
2. Preserve deterministic dry-run output and per-file diagnostics while allowing internal safe parallelism.
3. Ensure parent/container files such as `OXX-*/README.md` are created before child task files when required.
4. Update materialization docs and skill instructions to allow batch apply through roadmapctl while still prohibiting manual multi-file shell/write operations in the skill.
5. Add tests for successful batch apply, partial failure diagnostics, stale change-set conflicts, and ordering around Outcome README plus task files.

**Out**:
1. No direct markdown writes from the skill.
2. No removal of dry-run review before apply.
3. No agent/task implementation parallelism changes beyond materialization apply.

## Estado inicial esperado

roadmapctl materialize supports `--plan --dry-run`, full `--plan --apply`, and granular frozen change-set target apply. The current skill docs require applying each dry-run target with separate `roadmapctl materialize --changes <dry-run-json> --target <target.path> --apply` invocations.

## Especificación Técnica

Use `parallel` from the effective `.roadmapctl.toml` as the default policy for materialization internals. The external command surface may remain `roadmapctl materialize --plan <plan-json> --apply` and/or add `roadmapctl materialize --changes <dry-run-json> --apply` without `--target` for applying the whole frozen change-set. The implementation must keep dry-run deterministic, report every created path, reject stale or conflicting paths before writing when possible, and run the standard postcheck before success is reported. If internal parallelism is used, serialize dependent writes where necessary: create roadmap root/bootstrap artifacts before roadmap records, create `OXX-*/README.md` before `OXX-*/TXXX-*.md`, and allow independent task files to be written concurrently only after their parent directories/files are safe.

## Criterios de Aceptación

- go test ./internal/materialize ./internal/cli passes.
- A single roadmapctl materialize apply command can create the approved set of canonical files without the skill invoking one command per target.
- Dry-run JSON remains deterministic and can still be saved as a frozen change-set.
- Batch apply reports per-file changes and diagnostics; conflicts identify the concrete path that blocked apply.
- Skill docs explicitly allow roadmapctl-owned batch/parallel materialization and still forbid shell loops, multiple heredocs, and direct multi-file writes by the skill.

## Fuente de verdad

- docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md
- docs/materialize-plan-schema.md
- internal/materialize/dryrun.go
- internal/cli/materialize_test.go
- .claude/skills/roadmap/plan-subcommand.md
- .claude/skills/roadmap/common-logic.md
