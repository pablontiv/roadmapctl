---
estado: Specified
tipo: task
---
# T003: Delete dead internal/diff package

**Outcome**: [O25 Integrate picokit as a dependency](README.md)
**Contribuye a**: Eliminar código duplicado sin callers activos del codebase.

[[blocked_by:./T001-add-picokit-dependency.md]]

## Preserva

- INV1: `go build ./...` sigue pasando.
  - Verificar: `go build ./...` desde `/home/shared/roadmapctl`.

## Contexto

`internal/diff/unified.go` es una copia exacta de `picokit/diff/diff.go`. No tiene ningún caller dentro del codebase de roadmapctl — es dead code. picokit ya provee la implementación canónica.

Verificación previa:
```bash
grep -r '"github.com/pablontiv/roadmapctl/internal/diff"' /home/shared/roadmapctl --include="*.go"
# retorna vacío — confirma que nadie lo importa
```

## Alcance

**In**:
1. Eliminar `/home/shared/roadmapctl/internal/diff/unified.go`.
2. Eliminar `/home/shared/roadmapctl/internal/diff/unified_test.go`.
3. Eliminar el directorio `/home/shared/roadmapctl/internal/diff/` si queda vacío.

**Out**:
- No modificar ningún otro archivo.
- No añadir imports de `picokit/diff` en ningún lugar (no hay callers que los necesiten).

## Estado inicial esperado

- `/home/shared/roadmapctl/internal/diff/` existe con `unified.go` y `unified_test.go`.
- `grep -r '"github.com/pablontiv/roadmapctl/internal/diff"' /home/shared/roadmapctl --include="*.go"` retorna vacío.
- picokit v0.1.0 está en go.mod (T001 completada).

## Criterios de Aceptación

- `ls /home/shared/roadmapctl/internal/diff/` falla (directorio eliminado).
- `go build ./...` pasa desde `/home/shared/roadmapctl`.
- `go test ./...` pasa desde `/home/shared/roadmapctl`.

## Fuente de verdad

- `/home/shared/roadmapctl/internal/diff/` (a eliminar)
