---
estado: Completed
tipo: task
---
# T007: Migrar docs, fixtures y goldens al nuevo boundary

**Outcome**: [Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill](README.md)

[[blocked_by:./T002-redefine-roadmap-plan-skill-contract.md]]
[[blocked_by:./T003-design-roadmapctl-path-planning-guard.md]]
[[blocked_by:./T004-remove-persisted-outcome-task-table.md]]
[[blocked_by:./T005-move-acceptance-criteria-to-tasks-only.md]]

## Preserva

- Los tests siguen capturando estructura canónica y diagnostics estables.

## Contexto

Romper compat está aprobado, pero los contratos nuevos deben quedar cubiertos por docs y pruebas.

## Alcance

**In**:
1. Actualizar docs/cli-contract.md y docs/roadmap-skill-integration.md.
2. Deprecar, eliminar o convertir docs/materialize-plan-schema.md a histórico.
3. Actualizar fixtures/goldens sin ## Tasks ni outcome AC obligatorio.
4. Actualizar skill docs instalables y base.stem si corresponde.

**Out**:
1. No introducir compat legacy salvo que sea necesario para migrar el repo actual.

## Estado inicial esperado

Docs y tests actuales codifican materialize writer, Outcome AC y task table persistida.

## Criterios de Aceptación

- Docs reflejan la separación de responsabilidades acordada.
- Goldens ya no esperan ## Tasks persistido ni outcome.acceptance_criteria.
- Fixtures representan Outcomes como contenedores/contexto y Tasks como unidades verificables.
- La suite documenta los comandos de validación requeridos.

## Fuente de verdad

- docs/cli-contract.md
- docs/roadmap-skill-integration.md
- docs/materialize-plan-schema.md
- testdata/fixtures
- testdata/golden
