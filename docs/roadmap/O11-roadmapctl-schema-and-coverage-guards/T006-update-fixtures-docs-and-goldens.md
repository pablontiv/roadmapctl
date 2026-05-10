---
estado: Completed
tipo: task
---
# T006: Actualizar fixtures, docs y goldens

**Outcome**: [Roadmapctl schema and coverage guards](README.md)

[[blocked_by:./T003-wire-schema-compatibility-diagnostics.md]]
[[blocked_by:./T005-use-coverage-setting-in-check-coverage-script.md]]

## Preserva

- Queda al menos un fixture intencional para schema stale.
- Los fixtures normales no requieren estado manual en outcomes.

## Contexto

La investigación encontró múltiples fixtures con .stem legacy y Outcome README con estado manual que pueden ocultar regresiones.

## Alcance

**In**:
1. Actualizar .claude/skills/roadmap/base.stem al schema canónico.
2. Actualizar fixtures normales a estado requerido solo para T* y sin validate estado non_empty.
3. Remover estado manual de Outcome README en fixtures normales.
4. Agregar fixture invalid-stale-outcome-stem.
5. Actualizar goldens y docs con diagnostics y required_code_coverage.
6. Correr validaciones Go y coverage documentado.

**Out**:
1. Modificar repos externos como pinata dentro de esta task.
2. Cambiar semántica de estados de tasks.

## Estado inicial esperado

Muchos fixtures y el skill base.stem aún reflejan el schema legacy que requiere estado en O*.

## Criterios de Aceptación

- go test ./... pasa.
- Goldens reflejan nuevos diagnostics/config.
- Greps de fixtures solo encuentran estado outcome/schema legacy en fixtures intencionales.
- docs explican required_code_coverage y schema compatibility.

## Fuente de verdad

- testdata/fixtures/
- testdata/golden/
- .claude/skills/roadmap/base.stem
- docs/cli-contract.md
- docs/release.md
