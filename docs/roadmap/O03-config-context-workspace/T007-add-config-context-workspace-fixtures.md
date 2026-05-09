---
estado: Pending
tipo: task
---
# T007: Agregar fixtures/goldens de config, context y workspace

**Outcome**: [O03 Config/context/workspace](README.md)
**Contribuye a**: CE1, CE2, CE3, CE4

[[blocked_by:./T002-implement-config-discovery-and-toml-loader.md]]
[[blocked_by:./T003-implement-legacy-config-migration.md]]
[[blocked_by:./T004-implement-context-command.md]]
[[blocked_by:./T006-implement-workspace-discovery.md]]

## Preserva

- INV1: Fixtures no dependen de estado global del usuario.
  - Verificar: tests usan temp dirs/fake Rootline cuando corresponda.

## Contexto

La migración de config y context discovery introduce muchas rutas de compatibilidad. Necesita fixtures dedicados para evitar regresiones.

## Alcance

**In**:
1. Fixture TOML default en `docs/roadmap/.roadmapctl.toml`.
2. Fixture legacy-only.
3. Fixture TOML + legacy conflict.
4. Fixture root escape.
5. Fixtures workspace válido, repo ambiguo, repo missing config.
6. Goldens de `roadmapctl context --output json`.

**Out**:
- Fixtures de pending/next/materialize.

## Estado inicial esperado

- Existen fixtures MVP bajo `testdata/fixtures`.

## Criterios de Aceptación

- `go test ./...` pasa en Linux/macOS/Windows.
- Goldens normalizan paths absolutos.
- `docs/golden-tests.md` se actualiza si cambia workflow.

## Fuente de verdad

- `testdata/fixtures/*`
- `testdata/golden/*`
- `internal/cli/golden_test.go`
- `docs/golden-tests.md`
