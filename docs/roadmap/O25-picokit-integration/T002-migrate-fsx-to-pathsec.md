---
estado: Specified
tipo: task
---
# T002: Replace internal/fsx with picokit/pathsec

**Outcome**: [O25 Integrate picokit as a dependency](README.md)
**Contribuye a**: Eliminar código duplicado de resolución de paths con seguridad de symlinks.

[[blocked_by:./T001-add-picokit-dependency.md]]

## Preserva

- INV1: El comportamiento de resolución de paths es idéntico al actual.
  - Verificar: `go test ./...` pasa desde `/home/shared/roadmapctl`.
- INV2: No se introducen nuevas dependencias externas.
  - Verificar: `go mod tidy` no agrega módulos distintos a picokit.

## Contexto

`internal/fsx/path.go` es una copia casi idéntica de `picokit/pathsec`. La única diferencia es el nombre del paquete (`fsx` vs `pathsec`) y un comentario de docstring.

API de picokit/pathsec:
- `pathsec.ResolveInside(root, candidate string) (abs, rel string, err error)` — misma firma que `fsx.ResolveInside`.
- `pathsec.ErrPathEscape` — misma semántica que `fsx.ErrPathEscape`.
- `pathsec.ErrAbsolutePath` — misma semántica que `fsx.ErrAbsolutePath`.

Callers actuales (2 archivos):
- `/home/shared/roadmapctl/internal/config/config.go`
- `/home/shared/roadmapctl/internal/cli/bootstrap.go`

## Alcance

**In**:
1. En `config/config.go`: reemplazar import `internal/fsx` por `github.com/pablontiv/picokit/pathsec`. Renombrar `fsx.ResolveInside` → `pathsec.ResolveInside`, `fsx.ErrPathEscape` → `pathsec.ErrPathEscape`, `fsx.ErrAbsolutePath` → `pathsec.ErrAbsolutePath`.
2. En `cli/bootstrap.go`: mismo reemplazo.
3. Eliminar `/home/shared/roadmapctl/internal/fsx/` (archivos `path.go` y `path_test.go` si existe).

**Out**:
- No modificar la lógica de negocio en config.go ni bootstrap.go.
- No modificar otros paquetes que no importen fsx.

## Estado inicial esperado

- `/home/shared/roadmapctl/internal/fsx/path.go` existe.
- `grep -r "internal/fsx" /home/shared/roadmapctl --include="*.go"` muestra 2 archivos.
- picokit v0.1.0 está en go.mod (T001 completada).

## Criterios de Aceptación

- `grep -r "internal/fsx" /home/shared/roadmapctl --include="*.go"` retorna vacío.
- `ls /home/shared/roadmapctl/internal/fsx/` falla (directorio eliminado).
- `go build ./...` pasa desde `/home/shared/roadmapctl`.
- `go test ./...` pasa desde `/home/shared/roadmapctl`.
- `just check` pasa (gofmt + go vet).

## Fuente de verdad

- `/home/shared/roadmapctl/internal/fsx/path.go` (a eliminar)
- `/home/shared/roadmapctl/internal/config/config.go`
- `/home/shared/roadmapctl/internal/cli/bootstrap.go`
- `/home/shared/picokit/pathsec/pathsec.go` (referencia de API destino)
