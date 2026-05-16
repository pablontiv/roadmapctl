---
estado: Completed
tipo: task
---
# T003: Integrar auto-update en internal/cli/cli.go

**Outcome**: [O21 Auto-update staged async](README.md)
**Contribuye a**: Activación del auto-update en cada invocación del CLI

[[blocked_by:./T002-updater-apply-reexec.md]]

## Preserva

- INV1: El contrato de salida de todos los comandos existentes no cambia.
  - Verificar: `roadmapctl next --output json` produce JSON válido sin campos extra.
- INV2: La latencia de arranque no aumenta perceptiblemente cuando no hay nada staged.
  - Verificar: `time roadmapctl --version` no agrega overhead observable.

## Contexto

El auto-update se activa al inicio de `Execute()` en `internal/cli/cli.go`. El patrón es:

```go
updater.ApplyStagedIfAvailable()       // sync: detecta y aplica staged; si hay update, re-exec y no continúa
go updater.FetchAndStage(version)      // goroutine: descarga en background para el próximo run
```

La variable `version` ya existe en `cmd/roadmapctl/version.go` y se pasa a `cli.Execute(version)`.

## Alcance

**In**:
1. En `internal/cli/cli.go`, al inicio de `Execute(version string)`, agregar las dos líneas anteriores
2. Importar `internal/updater`

**Out**:
- No modificar ningún subcomando existente
- No agregar flags ni config en `.roadmapctl.toml`
- No cambiar la firma de `Execute()`

## Estado inicial esperado

- T001 y T002 completadas: `internal/updater` exporta `FetchAndStage` y `ApplyStagedIfAvailable`
- `internal/cli/cli.go` existe con función `Execute(version string)`

## Criterios de Aceptación

- `go build ./...` pasa sin errores
- `roadmapctl --version` reporta la versión actual compilada
- En la siguiente invocación después de un update staged, el binario nuevo está en uso (verificable con `roadmapctl --version`)
- `roadmapctl next --output json` produce JSON válido (contrato no roto)

## Fuente de verdad

- `internal/cli/cli.go` (modificar — agregar 2 líneas al inicio de Execute)
- `internal/updater/updater.go` (T001)
- `internal/updater/apply.go` (T002)
- `cmd/roadmapctl/version.go` (fuente del version string)
