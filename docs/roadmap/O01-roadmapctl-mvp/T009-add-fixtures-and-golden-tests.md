---
estado: Completed
tipo: task
---
# T009: Añadir fixtures y golden tests cross-platform

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE1 y CE2

[[blocked_by:./T006-implement-doctor-command.md]]
[[blocked_by:./T007-implement-structure-checks.md]]
[[blocked_by:./T008-implement-rootline-backed-checks.md]]

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: fixtures ejercitan `roadmapctl`, no cambios a Rootline.
- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: tests esperan diagnostics, no modificaciones.

## Contexto

La confianza en `roadmapctl` depende de fixtures que reproduzcan estructuras válidas e inválidas, incluyendo el bug de un único archivo resumen. Los tests deben ser estables en Linux, macOS y Windows.

## Alcance

**In**:
1. Crear fixtures válidos: direct tasks y outcome con tasks.
2. Crear fixtures inválidos: single summary file, missing `.stem`, missing README, bare `blocked_by`, broken `blocked_by`, cycle, duplicate IDs, status mismatch, root escape.
3. Añadir golden JSON tests con paths normalizados.
4. Añadir integración opcional con `ROOTLINE_BIN`.
5. Documentar cómo actualizar goldens.

**Out**:
- Tests end-to-end de package managers.
- Tests de materialización.
- Fixtures que dependan de paths absolutos locales.

## Estado inicial esperado

- `doctor` y `check` implementados o suficientemente mockeables.

## Criterios de Aceptación

- `go test ./...` pasa.
- Fixture `invalid-single-summary-file` falla con exit `1`.
- Fixture `valid-outcome-with-tasks` pasa con exit `0`.
- Tests no dependen de separador `/` del sistema para comparar paths.
- Integración real con Rootline se puede activar con `ROOTLINE_BIN`.

## Fuente de verdad

- `testdata/fixtures/`
- `internal/testutil/`
- `internal/roadmap/`
