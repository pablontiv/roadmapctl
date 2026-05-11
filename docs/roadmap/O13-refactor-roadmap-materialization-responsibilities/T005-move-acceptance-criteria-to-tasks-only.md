---
estado: Completed
tipo: task
---
# T005: Mover criterios de aceptación exclusivamente a Tasks

**Outcome**: [Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill](README.md)

[[blocked_by:./T001-record-responsibility-separation-decision.md]]

## Preserva

- Las Tasks siguen siendo ejecutables por agentes con criterios pass/fail claros.

## Contexto

Los Outcome agrupan contexto; la verificación accionable pertenece a cada Task.

## Alcance

**In**:
1. Actualizar schema/docs/templates para Outcome sin AC obligatorio.
2. Actualizar Task template para AC requeridos.
3. Actualizar validaciones/goldens que esperan outcome.acceptance_criteria.

**Out**:
1. No prohibir prose descriptiva en Outcome si aporta contexto.

## Estado inicial esperado

El materialize-plan actual requiere acceptance_criteria en Outcome y renderOutcome los escribe.

## Criterios de Aceptación

- Outcome README puede contener solo frontmatter, título y descripción/contexto.
- Task template contiene ## Criterios de Aceptación.
- Tests prueban Outcome sin AC y Task con AC.
- Docs dejan claro que AC vive en Tasks.

## Fuente de verdad

- docs/materialize-plan-schema.md
- internal/materialize/dryrun.go
- .claude/skills/roadmap/outcome-guide.md
- .claude/skills/roadmap/task-guide.md
