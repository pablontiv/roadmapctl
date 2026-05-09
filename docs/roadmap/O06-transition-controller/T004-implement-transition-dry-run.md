---
estado: Completed
tipo: task
---
# T004: Implementar dry-run de transición

**Outcome**: [O06 Controlador de transiciones](README.md)
**Contribuye a**: CE2

[[blocked_by:./T002-add-rootline-set-and-validateone-wrappers.md]]
[[blocked_by:./T003-implement-can-start-can-complete.md]]

## Preserva

- INV1: Dry-run no escribe archivos.
  - Verificar: fixtures y git diff limpio tras test.

## Contexto

Antes de permitir apply, transition debe mostrar exactamente qué cambio de frontmatter haría.

## Alcance

**In**:
1. `transition start --dry-run` planea `estado` a role `in-progress`.
2. `transition complete --dry-run` planea `estado` a role `completed`.
3. `transition set-status --dry-run` planea estado explícito/role.
4. Output `changes[]` con path, field, before, after, applied=false.

**Out**:
- Escribir archivos.
- Unified diff si no es necesario para frontmatter simple.

## Estado inicial esperado

- can-start/can-complete existen.

## Criterios de Aceptación

- Dry-run blocked no propone cambios aplicables.
- Dry-run ready muestra cambio exacto.
- JSON stable para goldens.

## Fuente de verdad

- `internal/cli/transition.go` nuevo
- `internal/roadmap/transition.go` nuevo
