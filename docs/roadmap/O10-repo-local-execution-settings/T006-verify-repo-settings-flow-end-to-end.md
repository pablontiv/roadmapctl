---
estado: Specified
tipo: task
---
# T006: Verify repo settings flow end to end

**Outcome**: [Repo-local roadmap execution settings](README.md)

[[blocked_by:./T001-implement-config-fields-and-forced-migration.md]]
[[blocked_by:./T002-expose-execution-settings-in-context.md]]
[[blocked_by:./T003-update-config-templates-fixtures-and-contract-docs.md]]
[[blocked_by:./T004-add-roadmap-context-compaction-extension.md]]
[[blocked_by:./T005-cutover-roadmap-skill-loop-to-config.md]]

## Preserva

- Full Go test suite remains green.
- Skill guard verification still proves roadmapctl doctor/check run before loop/materialization.

## Contexto

The design changes config loading side effects and skill behavior. A final verification pass must prove the new contract works as an integrated workflow.

## Alcance

**In**:
1. Run focused config and CLI tests, then go test ./... .
2. Run the required Pi headless roadmap skill verification commands from .claude/skills/roadmap/SKILL.md after skill changes.
3. Add any missing tests or fixture updates revealed by those commands.
4. Record verification evidence in the final task summary or an existing docs location if the repo uses one.

**Out**:
1. No new feature behavior beyond fixes needed to satisfy the prior tasks' acceptance criteria.

## Estado inicial esperado

The repo has focused Go tests and a documented headless verification script for roadmap skill guard behavior.

## Especificación Técnica

Before running headless verification, use git status --short and preserve user changes. If headless verification would mutate files because config Load now migrates legacy during context, run it in a disposable copy or fixture repo and report that choice. The final evidence must include exact commands and exit statuses.

## Criterios de Aceptación

- go test ./internal/config passes.
- go test ./internal/cli ./internal/materialize passes.
- go test ./... passes.
- The two Pi headless verification commands from .claude/skills/roadmap/SKILL.md pass or produce documented evidence that roadmapctl doctor/check were required before loop/materialization without modifying files.
- git status shows only intentional implementation, test, fixture, doc, roadmap, and extension changes.

## Fuente de verdad

- docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md
- .claude/skills/roadmap/SKILL.md
- scripts/verify-roadmap-skill-headless.sh
