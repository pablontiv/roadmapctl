---
estado: Completed
tipo: task
---
# T002: Implementar numbering y path planning cross-platform

**Outcome**: [O07 Controlador de materialización](README.md)
**Contribuye a**: CE2, CE3

[[blocked_by:./T001-define-structured-materialize-plan-schema.md]]
[[blocked_by:../O02-post-mvp-foundations/T005-add-roadmap-domain-model-and-tree-wrapper.md]]

## Preserva

- INV1: No usar snippets POSIX/GNU como `find -printf`.
  - Verificar: tests Windows.

## Contexto

El skill hoy describe auto-numbering con shell. `roadmapctl` debe planear IDs y paths en Go, respetando secuencias separadas OXX y TXXX.

## Alcance

**In**:
1. Calcular próximo OXX en raíz.
2. Calcular próximo TXXX directo en raíz.
3. Calcular próximo TXXX dentro de Outcome.
4. Detectar collisions, slugs inválidos, existing files y path escapes.
5. Producir plan de paths antes de escribir.

**Out**:
- Crear archivos.
- Commit/push.

## Estado inicial esperado

- Structure checks y regex OXX/TXXX existen.

## Criterios de Aceptación

- Tests cubren repos vacíos, mixed direct/outcome, collisions y Windows separators.
- Nunca planea `*-tasks.md`.
- Todos los paths resuelven dentro de roadmap root.

## Fuente de verdad

- `internal/roadmap/numbering.go`
- `internal/roadmap/structure.go`
- `.claude/skills/roadmap/common-logic.md`
