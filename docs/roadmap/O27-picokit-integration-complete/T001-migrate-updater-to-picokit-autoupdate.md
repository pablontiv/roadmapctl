---
estado: Specified
tipo: task
---
# T001: Migrate internal/updater to picokit/autoupdate

**Outcome**: [O27 Complete picokit integration](README.md)
**Contribuye a**: roadmapctl deja de mantener su copia local del updater y consume `picokit/autoupdate` upstream.

## Preserva

- INV1: `ROADMAPCTL_NO_UPDATE=1 ./roadmapctl --version` no hace requests HTTP.
  - Verificar: `t.Setenv("ROADMAPCTL_NO_UPDATE", "1")` + `http.DefaultTransport` instrumentado con `failTransport` que falla si se invoca.
- INV2: `version == "dev"` no dispara update.
  - Verificar: `TestExecute_AutoupdateSkipsOnDevVersion` pasa.
- INV3: `FetchAndStage` corre en goroutine y no bloquea el exit del CLI.
  - Verificar: `TestExecute_AutoupdateFetchRunsInGoroutine` con `httptest.NewServer` mock.

## Contexto

`internal/updater/` (apply.go, exec_unix.go, exec_windows.go, updater.go, ac_check_test.go, apply_test.go, updater_test.go, ~1554 LOC) fue el origen del paquete `picokit/autoupdate`. Hoy roadmapctl mantiene su copia local; el spec original de picokit (`O01-picokit-module/T003-autoupdate.md`) lo registró como "componente principal" y previó la migración.

Wiring actual en `internal/cli/cli.go:33-35`:
```go
updater.CurrentVersion = version
_ = updater.ApplyStagedIfAvailable()
go updater.FetchAndStage(version) //nolint:errcheck
```

Reemplazo:
```go
u := autoupdate.New("pablontiv/roadmapctl", "roadmapctl", "ROADMAPCTL_NO_UPDATE")
_ = u.ApplyStagedIfAvailable()
go u.FetchAndStage(version) //nolint:errcheck
```

Goreleaser ya genera `roadmapctl_{version}_{os}_{arch}.tar.gz` + `checksums.txt` (verificado en `.goreleaser.yml:31-46`), compatible con expectativas de `picokit/autoupdate`.

## Alcance

**In**:
1. Borrar `internal/updater/{updater,apply,exec_unix,exec_windows,ac_check_test,apply_test,updater_test}.go`.
2. Sustituir el wiring en `internal/cli/cli.go:32-36` por la llamada a `picokit/autoupdate.New(...)`.
3. Agregar 4 tests en `internal/cli/cli_test.go`:
   - `TestExecute_AutoupdateConstructorParams`
   - `TestExecute_AutoupdateSkipsOnEnv`
   - `TestExecute_AutoupdateSkipsOnDevVersion`
   - `TestExecute_AutoupdateFetchRunsInGoroutine`
4. Smoke script `scripts/test-autoupdate-smoke.sh` que verifica `ROADMAPCTL_NO_UPDATE=1 ./roadmapctl --version` termina en <100ms.

**Out**:
- No tocar `internal/diagnostics/` (eso es T002).
- No cambiar `.goreleaser.yml` (ya compatible).
- No introducir cambios de API pública del CLI.

## Estado inicial esperado

- `internal/updater/` existe con la copia local del updater.
- `picokit v0.1.1` ya está en `go.mod` (cumplido por O25-T001).
- `internal/cli/cli.go:33-35` referencia `updater.CurrentVersion`, `updater.ApplyStagedIfAvailable`, `updater.FetchAndStage`.

## Criterios de Aceptación

- `grep -rn "roadmapctl/internal/updater" /home/shared/roadmapctl --include="*.go"` retorna vacío.
- `grep -n "picokit/autoupdate" /home/shared/roadmapctl/internal/cli/cli.go` muestra el import.
- `go build ./...` pasa.
- `go test ./... -race -count=1` pasa.
- `golangci-lint run` pasa.
- `scripts/check-coverage.sh` pasa con threshold ≥85%.
- `ROADMAPCTL_NO_UPDATE=1 ./roadmapctl --version` termina en <100ms y no hace requests HTTP (verificable con `strace -e network`).
- Los 4 tests nuevos en `internal/cli/cli_test.go` pasan.
