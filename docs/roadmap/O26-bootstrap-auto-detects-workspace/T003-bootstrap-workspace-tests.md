---
estado: Completed
tipo: task
---
# T003: Tests for bootstrap workspace mode (auto-detect, explicit flag, empty, single-repo regression)

**Outcome**: [O26 Bootstrap auto-detects workspace mode](README.md)
**Contribuye a**: Garantizar que el comportamiento nuevo (T002) está cubierto por tests deterministas y que no introduce regresiones en single-repo.

[[blocked_by:./T002-bootstrap-auto-detect-workspace.md]]

## Preserva

- INV1: Coverage del paquete `internal/cli` no cae por debajo del umbral configurado (`required_code_coverage = 85`).
  - Verificar: `go test ./internal/cli/... -coverprofile=/tmp/cov.out && go tool cover -func=/tmp/cov.out | tail -1`.
- INV2: `go test ./...` pasa desde `/home/shared/roadmapctl`.

## Contexto

`internal/cli/bootstrap_test.go` no contiene hoy ningún test de workspace para el comando `bootstrap` (sí los tiene `pending_test.go`). Los tests existentes usan helpers como `initGitRepo(t, dir)` para preparar single-repos. Para workspace necesitamos una fixture o builder que arme: root sin `.git/` + N subdirs con `.git/` + `docs/roadmap/.roadmapctl.toml` mínimo.

El fixture `testdata/fixtures/valid-workspace/{alpha,beta}/` ya provee dos members válidos y puede reutilizarse o copiarse al `t.TempDir()` cuando se necesite mutar.

## Alcance

**In**:
1. Agregar `TestBootstrapAutoDetectsWorkspaceWhenNoGitAtRoot` en `internal/cli/bootstrap_test.go`:
   - Copia el fixture `valid-workspace` a `t.TempDir()`.
   - Ejecuta el comando bootstrap apuntando al root copiado, sin `--workspace`.
   - Assert: `summary.status == "ok"`, `config_source == "workspace"`, `len(repos) == 2`, `repos[].name` contiene `"alpha"` y `"beta"`.
2. `TestBootstrapHonorsExplicitWorkspaceFlag`:
   - Mismo fixture, pero pasa `options.Workspace = true` explícito.
   - Assert: mismo resultado que (1). Confirma equivalencia entre auto-detección y flag explícito.
3. `TestBootstrapWorkspaceEmptyEmitsDiagnostic`:
   - `t.TempDir()` sin `.git/` y sin members con config.
   - Assert: `summary.status == "error"`, `diagnostics[0].id == "RMC_WORKSPACE_EMPTY"`. Verifica que NO emite `RMC_CONFIG_MISSING`.
4. `TestBootstrapSingleRepoUnaffected`:
   - `t.TempDir()` con `.git/` (vía `initGitRepo`) y un `.roadmapctl.toml` válido (reusar fixture o helper existente).
   - Assert: `config_source == "toml"`, `len(repos) == 0` (o campo ausente, según marshalling con `omitempty`). Verifica que la rama workspace no se activa.
5. Si el helper para crear members workspace no existe, agregarlo en `bootstrap_test.go` o en un `testhelper_test.go` siguiendo el patrón de `initGitRepo`.

**Out**:
- No agregar tests para `bootstrap inspect` ni `bootstrap init` en esta task.
- No modificar tests existentes de bootstrap (solo agregar).
- No agregar benchmarks.

## Estado inicial esperado

- T002 mergeada: la lógica de workspace en `buildBootstrapConfig` funciona.
- `grep -n "TestBootstrap.*Workspace" /home/shared/roadmapctl/internal/cli/bootstrap_test.go` retorna vacío (no hay tests previos de workspace en bootstrap).
- `testdata/fixtures/valid-workspace/{alpha,beta}/` existen con sus configs.

## Criterios de Aceptación

- `go test ./internal/cli/... -run "TestBootstrap.*Workspace|TestBootstrapSingleRepoUnaffected" -v` ejecuta los 4 tests nuevos y todos pasan.
- `go test ./internal/cli/... -coverprofile=/tmp/cov.out && go tool cover -func=/tmp/cov.out | grep total` reporta `>= 85.0%`.
- `go test ./...` pasa desde `/home/shared/roadmapctl`.
- Los tests no dependen del estado del filesystem fuera de `t.TempDir()` ni del repo `/home/shared/roadmapctl` mismo (son hermeticos).

## Fuente de verdad

- `/home/shared/roadmapctl/internal/cli/bootstrap_test.go` (modificado: +4 tests, posibles helpers)
- `/home/shared/roadmapctl/testdata/fixtures/valid-workspace/` (referencia, posible copia in-test)
