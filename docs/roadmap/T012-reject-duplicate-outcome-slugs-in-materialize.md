---
estado: Completed
tipo: task
---
# T012: Reject duplicate outcome slugs in materialize

## Preserva

- Materialize remains create-only for this bugfix.
- No append/update/patch semantics are introduced.
- Numbering for unrelated new outcomes and direct tasks remains unchanged.
- No single-file *-tasks.md fallback is created.

## Contexto

A materialize plan with an outcome slug matching an existing OXX-<slug> currently plans a duplicate new outcome, because PlanMaterializePaths only uses existing outcome directories to compute maxOutcome and does not detect slug reuse.

## Alcance

**In**:
1. Detect existing canonical outcome directories with matching slugs in PlanMaterializePaths.
2. Return RMC_MATERIALIZE_PLAN_CONFLICT for duplicate existing outcome slugs before incrementing maxOutcome.
3. Add roadmap, materialize dry-run, and CLI test coverage proving duplicate outcomes are not proposed.
4. Document the v1 materialize-plan contract for existing outcome slug collisions.

**Out**:
1. Do not implement append of tasks to existing outcomes.
2. Do not introduce update or patch change operations.
3. Do not modify existing outcome README task tables.
4. Do not modify the roadmap skill in this task.

## Estado inicial esperado

roadmapctl materialize dry-run can propose O09-<slug> when O08-<slug> already exists and the input plan uses the same outcome slug.

## Criterios de Aceptación

- A dry-run with an existing OXX-foo outcome and an input outcome slug foo reports RMC_MATERIALIZE_PLAN_CONFLICT.
- The dry-run changes do not include a duplicate OXX-foo README or child task path.
- Unit and CLI tests cover the path planner, materialize dry-run behavior, and JSON CLI diagnostics.
- docs/materialize-plan-schema.md states that v1 outcome items create new outcomes and existing outcome slug collisions are rejected.
- go test ./internal/roadmap ./internal/materialize ./internal/cli passes.
- go test ./... passes.

## Fuente de verdad

- internal/roadmap/numbering.go
- internal/materialize/dryrun.go
- internal/cli/materialize.go
- docs/materialize-plan-schema.md
