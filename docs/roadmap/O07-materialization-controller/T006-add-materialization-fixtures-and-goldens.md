---
estado: Pending
tipo: task
---
# T006: Agregar fixtures/goldens de materialización

**Outcome**: [O07 Controlador de materialización](README.md)
**Contribuye a**: CE1, CE2, CE3

[[blocked_by:./T003-implement-materialize-dry-run-and-diff.md]]
[[blocked_by:./T004-implement-materialize-apply-with-postcheck.md]]
[[blocked_by:./T005-support-bootstrap-materialization.md]]

## Preserva

- INV1: Apply tests usan temp dirs y no modifican fixtures fuente.
  - Verificar: test helpers.

## Contexto

Materialization tiene más riesgo que read-only commands. Necesita matriz amplia de fixtures.

## Alcance

**In**:
1. Fixture direct tasks.
2. Fixture outcome+tasks.
3. Fixture dependencies same/cross outcome.
4. Fixture path escape.
5. Fixture existing file collision.
6. Fixture stale dry-run/drift.
7. Anti-regression no `*-tasks.md`.

**Out**:
- Release packaging.

## Estado inicial esperado

- Dry-run/apply implementados.

## Criterios de Aceptación

- `go test ./...` pasa.
- Goldens muestran changes/diffs estables.
- Invalid input no escribe archivos.

## Fuente de verdad

- `testdata/fixtures/*`
- `testdata/plans/*`
- `testdata/golden/*`
