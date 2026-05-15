---
estado: Specified
tipo: task
---
# T004: Tests para internal/updater

**Outcome**: [O21 Auto-update staged async](README.md)
**Contribuye a**: Cobertura del paquete updater para mantener CI verde (≥85%)

[[blocked_by:./T002-updater-apply-reexec.md]]

## Preserva

- INV1: `go test ./...` pasa en CI (Linux + macOS).
  - Verificar: El workflow de CI no falla en el job de tests.
- INV2: Los tests no hacen llamadas de red reales.
  - Verificar: Tests pasan con red desconectada o usando `httptest.NewServer` local.

## Contexto

El paquete `internal/updater` tiene dos funciones principales:
- `FetchAndStage(currentVersion string)` — chequea API y descarga a staging
- `ApplyStagedIfAvailable()` — detecta staged y re-exec

Usar `httptest.NewServer` para mockear la GitHub API sin red real. El re-exec real no es testeable directamente; usar inyección de dependencia o variable de función (`var execFunc = syscall.Exec`) para mockear en tests.

## Alcance

**In**:
1. Crear `internal/updater/updater_test.go` con tests de `FetchAndStage`
2. Crear `internal/updater/apply_test.go` con tests de `ApplyStagedIfAvailable`
3. Usar `httptest.NewServer` para mockear GitHub API
4. Inyectar función de exec para tests de re-exec (sin exec real en tests)

**Out**:
- No tests de integración que requieran red real
- No tests E2E del binario completo (eso es el smoke test existente en CI)

## Estado inicial esperado

- T001 y T002 completadas: `internal/updater/updater.go` y `apply.go` existen

## Criterios de Aceptación

- `TestFetchAndStage_SkipsDevVersion`: `FetchAndStage("dev")` no llama red, retorna nil
- `TestFetchAndStage_SkipsNoUpdateEnv`: con `ROADMAPCTL_NO_UPDATE=1` no llama red, retorna nil
- `TestFetchAndStage_SkipsIfAlreadyStaged`: si binario staged existe, no re-descarga
- `TestFetchAndStage_VerifiesSHA256`: SHA256 incorrecto retorna error sin escribir staging
- `TestApply_SkipsIfNothingStaged`: `ApplyStagedIfAvailable()` retorna nil si no hay staged
- `TestApply_SkipsIfNotNewer`: versión staged igual o menor no aplica update
- `go test ./internal/updater/... -count=1` pasa en Linux y macOS
- Cobertura del paquete `internal/updater` ≥ 85%

## Fuente de verdad

- `internal/updater/updater_test.go` (crear)
- `internal/updater/apply_test.go` (crear)
- `internal/updater/updater.go` (T001)
- `internal/updater/apply.go` (T002)
