---
estado: Pending
tipo: task
---
# T002: Hacer `.stem` autoritativo para schema documental

**Outcome**: [O02 Fundaciones post-MVP](README.md)
**Contribuye a**: CE1, INV1

## Preserva

- INV1: Rootline sigue siendo DBMS/constraint engine genérico.
  - Verificar: no hay cambios en Rootline ni imports internos.

## Contexto

Bug actual: `roadmapctl check` rechaza `estado: On Hold` aunque `docs/roadmap/.stem` lo permite. La causa es que `configuredStatuses(cfg)` trata roles operacionales como enum exhaustivo y `statusDiagnostics` calcula `allowed = config ∩ schema`.

El diseño correcto es: `.stem` efectivo leído vía Rootline define valores documentales (`estado`, `tipo`, links). La config solo define roles operacionales.

## Alcance

**In**:
1. Obtener valores permitidos de `estado` desde `rootline describe <roadmap-root>/ --field schema.estado --output json` o shape equivalente.
2. Obtener valores permitidos de `tipo` desde schema cuando esté disponible.
3. Cambiar validación documental para usar schema como autoridad cuando exista.
4. Usar defaults/config solo como fallback degradado si schema no está disponible, con diagnostic claro.
5. Agregar fixture `valid-status-on-hold`.

**Out**:
- Rediseñar config completa.
- Implementar `.roadmapctl.toml`.

## Estado inicial esperado

- `docs/roadmap/.stem` incluye `On Hold`.
- `internal/cli/check.go` pasa `AllowedStatuses` desde roles config.
- `internal/roadmap/status.go` intersecta config con schema.

## Criterios de Aceptación

- `roadmapctl check` no emite `RMC_STATUS_UNKNOWN` para `On Hold` si schema lo permite.
- Un estado no presente en schema sigue fallando.
- `tipo` se valida contra schema cuando Rootline lo expone.
- Tests cubren schema top-level y nested shape si ambos existen.

## Fuente de verdad

- `internal/cli/check.go`
- `internal/roadmap/status.go`
- `internal/roadmap/dependencies.go`
- `docs/roadmap/.stem`
- `testdata/fixtures/valid-status-on-hold`
