---
estado: Completed
tipo: task
---
# T006: Agregar fixtures/goldens de transición

**Outcome**: [O06 Controlador de transiciones](README.md)
**Contribuye a**: CE1, CE2, CE3

[[blocked_by:./T003-implement-can-start-can-complete.md]]
[[blocked_by:./T004-implement-transition-dry-run.md]]
[[blocked_by:./T005-implement-transition-apply-with-postcheck.md]]

## Preserva

- INV1: Tests de apply escriben solo en temp dirs.
  - Verificar: fixtures base no quedan modificados.

## Contexto

Las transiciones combinan graph, roles, Rootline set y postcheck. Necesitan fixtures específicos.

## Alcance

**In**:
1. Fixture can-start ready.
2. Fixture can-start blocked.
3. Fixture custom status labels.
4. Dry-run golden.
5. Apply temp fixture.
6. Postcheck failure scenario.

**Out**:
- Materialization fixtures.

## Estado inicial esperado

- Transition commands implementados.

## Criterios de Aceptación

- `go test ./...` pasa.
- Goldens cubren allowed=false/true y changes.
- CI Windows pasa.

## Fuente de verdad

- `testdata/fixtures/*`
- `testdata/golden/*`
- `internal/cli/*_test.go`
