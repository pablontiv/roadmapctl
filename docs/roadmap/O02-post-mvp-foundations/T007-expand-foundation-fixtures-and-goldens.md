---
estado: Completed
tipo: task
---
# T007: Expandir fixtures y goldens fundacionales

**Outcome**: [O02 Fundaciones post-MVP](README.md)
**Contribuye a**: CE1, CE2, CE3

[[blocked_by:./T002-make-stem-authoritative-for-document-schema.md]]
[[blocked_by:./T003-validate-operational-status-roles-separately.md]]
[[blocked_by:./T004-harden-rootline-json-on-nonzero-exit.md]]

## Preserva

- INV1: Goldens reflejan comportamiento intencional, no esconden regresiones.
  - Verificar: seguir `docs/golden-tests.md`.

## Contexto

Los cambios de boundary schema/config y Rootline JSON necesitan fixtures permanentes para evitar regresiones como `On Hold`.

## Alcance

**In**:
1. Crear fixture `valid-status-on-hold`.
2. Crear fixture `invalid-status-bogus`.
3. Crear fixture `invalid-config-role-not-in-schema`.
4. Agregar test fake Rootline para non-zero + JSON válido.
5. Actualizar o agregar goldens normalizados.

**Out**:
- Cambiar goldens sin justificar comportamiento.
- Tests de comandos aún no implementados.

## Estado inicial esperado

- Fixtures MVP existentes cubren estructura, graph, status mismatch básico y rootline missing.

## Criterios de Aceptación

- `go test ./...` pasa.
- Los nuevos fixtures fallan/pasan según el diseño schema-authoritative.
- `docs/golden-tests.md` sigue siendo correcto.

## Fuente de verdad

- `testdata/fixtures/*`
- `testdata/golden/*`
- `internal/cli/golden_test.go`
- `docs/golden-tests.md`
