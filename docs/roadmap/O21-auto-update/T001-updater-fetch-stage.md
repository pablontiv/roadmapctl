---
estado: Specified
tipo: task
---
# T001: Crear paquete internal/updater — fetch y staging

**Outcome**: [O21 Auto-update staged async](README.md)
**Contribuye a**: Infraestructura de descarga y verificación SHA256 para el update staged async

## Preserva

- INV1: Los errores de red o SHA256 nunca interrumpen el comando en curso.
  - Verificar: `roadmapctl next` retorna normalmente aunque la API de GitHub sea inalcanzable.

## Contexto

El auto-update usa el patrón staged async: cada invocación descarga en background la nueva versión a `~/.cache/roadmapctl/staged/<version>/`. La siguiente invocación detecta el binario staged y lo aplica (T002). Esta task crea solo la lógica de fetch + staging.

Distribución actual: GitHub Releases con SHA256 checksums. `install.sh` ya implementa este flujo — reusar su lógica como referencia. API: `https://api.github.com/repos/pablontiv/roadmapctl/releases/latest`.

La versión actual se inyecta vía ldflags en `cmd/roadmapctl/version.go` como `var version = "dev"`.

## Alcance

**In**:
1. Crear `internal/updater/updater.go` con función exportada `FetchAndStage(currentVersion string) error`
2. Consultar GitHub API para obtener el latest release tag
3. Comparar con `currentVersion`; skip si igual o menor (semver)
4. Descargar el tarball `.tar.gz` (Linux/macOS) o `.zip` (Windows) para el OS/arch actual
5. Verificar SHA256 contra `checksums.txt` del release
6. Extraer binario a `~/.cache/roadmapctl/staged/<version>/roadmapctl`
7. Skip silencioso si `currentVersion == "dev"` o `ROADMAPCTL_NO_UPDATE=1`
8. Skip silencioso si ya existe un binario staged para esa versión

**Out**:
- No aplica el update (eso es T002)
- No modifica el binario actual
- No agrega output a stdout/stderr del usuario

## Estado inicial esperado

- `internal/updater/` no existe aún
- Releases publicados en `https://github.com/pablontiv/roadmapctl/releases`

## Criterios de Aceptación

- `FetchAndStage("dev")` retorna nil sin hacer ninguna llamada de red
- Con `ROADMAPCTL_NO_UPDATE=1`, retorna nil sin llamadas de red
- Si ya existe `~/.cache/roadmapctl/staged/<version>/roadmapctl`, retorna nil sin re-descargar
- SHA256 incorrecto resulta en error retornado y ningún archivo escrito en staging
- `go build ./...` pasa sin errores

## Fuente de verdad

- `internal/updater/updater.go` (crear)
- `install.sh` (referencia para lógica de download/verify)
- `cmd/roadmapctl/version.go` (fuente de la versión actual)
