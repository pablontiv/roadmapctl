---
estado: Pending
tipo: task
---
# T007: Align materialized task initial status

[[blocked_by:./T003-implement-granular-materialize-target-apply.md]]

## Preserva

- Status roles remain configured through roadmapctl context/config.
- Outcome README status remains derived, not manually set.

## Contexto

The task guide describes new AI-ready tasks as Specified, while the materializer emits Pending.

## Alcance

**In**:
1. Decide whether generated tasks should start as Pending or Specified.
2. Update renderer/tests/goldens or task guide/docs accordingly.
3. Verify pending/next/decision behavior remains correct after the decision.

**Out**:
1. Do not reintroduce manual estado on outcome README files.
2. Do not change transition semantics without separate approval.

## Estado inicial esperado

Generated tasks use Pending, while task-guide.md templates indicate Specified for AI-ready tasks.

## Criterios de Aceptación

- Generated task frontmatter and task-guide.md agree on initial status.
- Materialize goldens reflect the chosen status.
- go test ./internal/materialize ./internal/roadmap passes.

## Fuente de verdad

- .claude/skills/roadmap/task-guide.md
- internal/materialize/dryrun.go
- testdata/golden/materialize-dry-run-outcome-and-direct.json
- internal/roadmap/transition.go
