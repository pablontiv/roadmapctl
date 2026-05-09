---
estado: Pending
tipo: task
---
# T005: Implementar apply de transición con postcheck

**Outcome**: [O06 Controlador de transiciones](README.md)
**Contribuye a**: CE3

[[blocked_by:./T004-implement-transition-dry-run.md]]
[[blocked_by:../O05-deterministic-lint/T001-define-lint-taxonomy-severity-json.md]]

## Preserva

- INV1: Apply requiere flag explícito y postcheck.
  - Verificar: tests de apply y failures.

## Contexto

Una transición aplicada debe usar Rootline `set`, validar el archivo y correr `roadmapctl check` o equivalente interno antes de reportar éxito.

## Alcance

**In**:
1. Implementar `--apply` para start/complete/set-status.
2. Mutar vía Rootline `set`.
3. Validar target y correr postcheck.
4. Reportar changes applied=true y diagnostics si postcheck falla.
5. Usar temp fixtures para tests de escritura.

**Out**:
- Commit/push.
- Ejecutar ACs.
- Auto-rollback complejo si no se aprueba; documentar si no existe.

## Estado inicial esperado

- Dry-run y wrappers Rootline existen.

## Criterios de Aceptación

- Apply cambia estado esperado en temp fixture.
- Postcheck failure se reporta y no se declara éxito.
- JSON mode no imprime salida cruda de Rootline.

## Fuente de verdad

- `internal/rootlinecli/*`
- `internal/roadmap/transition.go`
- `internal/cli/transition.go`
