---
estado: Completed
tipo: task
---
# T003: Implementar materialize dry-run y diff

**Outcome**: [O07 Controlador de materialización](README.md)
**Contribuye a**: CE1

[[blocked_by:./T002-implement-cross-platform-numbering-path-planning.md]]
[[blocked_by:../O02-post-mvp-foundations/T001-adopt-cobra-and-community-packages.md]]

## Preserva

- INV1: Dry-run no escribe archivos.
  - Verificar: git diff limpio tras dry-run.

## Contexto

Antes de permitir apply, materialize debe mostrar los archivos/cambios que generaría. Para output humano, usar unified diff estable cuando sea útil.

## Alcance

**In**:
1. `roadmapctl materialize --plan FILE --dry-run`.
2. Output JSON con `changes[]`, paths, operations, applied=false.
3. Output text/diff con `go-udiff` o equivalente aprobado.
4. Drift/precondition info para apply futuro.

**Out**:
- Escribir archivos.
- Actualizar README real.

## Estado inicial esperado

- Plan schema y path planning existen.

## Criterios de Aceptación

- Dry-run para outcome+tasks muestra todos los archivos.
- Dry-run para direct tasks muestra paths correctos.
- No se crean archivos durante dry-run.

## Fuente de verdad

- `internal/materialize` nuevo
- `internal/diff` nuevo
- `testdata/plans/*` nuevo
