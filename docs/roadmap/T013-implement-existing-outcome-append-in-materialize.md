---
estado: Specified
tipo: task
---
# T013: Implement existing Outcome append in materialize

## Preserva

- Do not reintroduce duplicate OXX-<slug> Outcome creation when an Outcome slug already exists.
- Materialize remains the canonical writer; the roadmap skill must not edit README task tables manually.
- Direct task numbering and new Outcome numbering remain unchanged.
- No single-file *-tasks.md fallback is created.
- Apply keeps conflict and stale-change detection before writing when possible.

## Contexto

Session 019e17ab-bf58-76d1-b577-a1a392889337 in /home/shared/cartyx is blocked because adding a follow-up task to existing O08-soporte-scip-todos-los-repos cannot currently be represented by roadmapctl materialize. T012 only rejected duplicate Outcome slugs; it did not implement append semantics.

## Alcance

**In**:
1. Detect an existing canonical Outcome directory by slug during materialize path planning.
2. Calculate the next TXXX number inside the existing Outcome.
3. Generate create changes for OXX-slug/TXXX-task.md under the existing Outcome.
4. Generate a safe update change for OXX-slug/README.md that inserts or refreshes the corresponding ## Tasks table row.
5. Support dry-run and apply for the required README update operation.
6. Reject task slug or path collisions inside the existing Outcome.
7. Add planner, dry-run, apply, and CLI JSON tests for existing Outcome append behavior.
8. Update materialize schema documentation and roadmap skill docs to describe append semantics.

**Out**:
1. Do not implement arbitrary patch operations for files outside the parent Outcome README.md.
2. Do not allow generic updates to any roadmap file.
3. Do not change Rootline schema or roadmap status semantics.
4. Do not implement the Cartyx bugfix task itself; only unblock correct materialization.

## Estado inicial esperado

A materialize plan containing an outcome item whose slug matches an existing Outcome now fails with RMC_MATERIALIZE_PLAN_CONFLICT and therefore cannot create the intended child task inside that Outcome.

## Criterios de Aceptación

- A plan with slug soporte-scip-todos-los-repos and an existing O08-soporte-scip-todos-los-repos proposes update O08-soporte-scip-todos-los-repos/README.md plus create O08-soporte-scip-todos-los-repos/T012-*.md.
- The same dry-run does not include O09-soporte-scip-todos-los-repos or any duplicate Outcome directory.
- Apply writes the new child task and updates the parent README.md ## Tasks table through roadmapctl-owned changes.
- Reapplying a stale changeset fails with RMC_MATERIALIZE_PLAN_CONFLICT or an equally specific materialize conflict diagnostic before unsafe overwrite.
- Planner, materialize dry-run/apply, and CLI JSON tests cover existing Outcome append and collision rejection.
- docs/materialize-plan-schema.md and .claude/skills/roadmap/plan-subcommand.md document that v1 Outcome items append tasks when the Outcome slug already exists.
- go test ./internal/roadmap ./internal/materialize ./internal/cli passes.
- go test ./... passes.
- roadmapctl check --repo . --roadmap-root docs/roadmap --output json --strict passes.
- After install, the global roadmapctl can dry-run the Cartyx session plan into the existing O08 Outcome rather than rejecting it or proposing O09.

## Fuente de verdad

- internal/roadmap/numbering.go
- internal/materialize/dryrun.go
- internal/materialize/apply_changes_test.go
- internal/cli/materialize.go
- internal/cli/materialize_test.go
- internal/lint/task_table.go
- docs/materialize-plan-schema.md
- .claude/skills/roadmap/plan-subcommand.md
