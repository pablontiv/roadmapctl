---
estado: Specified
tipo: task
---
# T002: Implementar apply + re-exec en internal/updater

**Outcome**: [O21 Auto-update staged async](README.md)
**Contribuye a**: Activación del binario staged en la siguiente invocación con re-exec transparente

[[blocked_by:./T001-updater-fetch-stage.md]]

## Preserva

- INV1: Un fallo de permisos al reemplazar el binario nunca interrumpe el comando en curso.
  - Verificar: Correr roadmapctl con binario instalado en `/usr/local/bin` sin sudo; el comando funciona normal y el error de permisos es silencioso.

## Contexto

Complemento de T001. Una vez que `FetchAndStage` dejó un binario en staging, la próxima invocación necesita detectarlo, reemplazar el binario actual de forma atómica, y re-ejecutar el proceso con el nuevo binario antes de ejecutar el comando pedido.

Re-exec en Unix: `syscall.Exec(newBinaryPath, os.Args, os.Environ())` — reemplaza el proceso actual en el mismo PID, invisible para el usuario.
Re-exec en Windows: el binario en uso no puede reemplazarse en-place; usar `exec.Command` + `os.Exit(0)` para lanzar el nuevo proceso y terminar el actual.

Atomic rename: `os.Rename(stagedPath, currentBinaryPath)` — atómico en el mismo filesystem (aplica para instalaciones en `$HOME/.local/bin`).

## Alcance

**In**:
1. Crear `internal/updater/apply.go` con función exportada `ApplyStagedIfAvailable() error`
2. Detectar si existe algún binario staged más nuevo en `~/.cache/roadmapctl/staged/`
3. Comparar versión staged con versión actual; skip si no es mayor
4. Atomic rename del binario staged sobre el binario actual (`os.Executable()`)
5. Re-exec con `syscall.Exec` (Unix) o `exec.Command` + `os.Exit` (Windows)
6. Errores de permisos u `os.Rename` silenciosos — retornar nil, no interrumpir

**Out**:
- No descarga nada (eso es T001/FetchAndStage)
- No agrega output visible al usuario

## Estado inicial esperado

- T001 completada: `internal/updater/updater.go` existe y exporta la constante o función que retorna el path de staging

## Criterios de Aceptación

- `ApplyStagedIfAvailable()` retorna nil si no hay nada staged
- Con binario staged más nuevo presente: reemplaza binario actual y re-exec (testeable con mock de exec)
- Con binario staged de versión igual o menor: skip, retorna nil
- Error de permisos en `os.Rename`: retorna nil (silencioso)
- `go build ./...` pasa sin errores

## Fuente de verdad

- `internal/updater/apply.go` (crear)
- `internal/updater/updater.go` (T001, para path de staging)
