---
estado: Completed
tipo: task
---
# T006: Agregar fixtures/goldens de lint

**Outcome**: [O05 Lint semántico determinístico](README.md)
**Contribuye a**: CE1, CE2, CE3

[[blocked_by:./T003-implement-outcome-task-table-consistency.md]]
[[blocked_by:./T004-implement-task-section-and-ac-lint.md]]
[[blocked_by:./T005-implement-schema-and-cross-platform-name-lints.md]]

## Preserva

- INV1: Fixtures de lint no modifican archivos durante tests.
  - Verificar: tests en copias temporales si hace falta.

## Contexto

Cada lint necesita fixtures claros para que severities y output sean estables.

## Alcance

**In**:
1. Fixture outcome table missing child.
2. Fixture stale table row.
3. Fixture missing task sections.
4. Fixture missing ACs/Fuente de verdad.
5. Fixture case-insensitive collision.
6. Goldens JSON para `roadmapctl lint`.

**Out**:
- Fixtures de materialization apply.

## Estado inicial esperado

- Lints implementados.

## Criterios de Aceptación

- `go test ./...` pasa.
- Goldens cubren warning vs strict si aplica.
- Docs explican cómo actualizar goldens intencionalmente.

## Fuente de verdad

- `testdata/fixtures/*`
- `testdata/golden/*`
- `docs/golden-tests.md`
