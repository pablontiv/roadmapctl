---
estado: Specified
tipo: task
---
# T001: Extract workspace detection into shared internal/workspace package

**Outcome**: [O26 Bootstrap auto-detects workspace mode](README.md)
**Contribuye a**: Reusable detection primitives so both `pending` and `bootstrap` (and future commands) share one canonical implementation.

## Preserva

- INV1: El comportamiento de `roadmapctl pending --workspace` no cambia.
  - Verificar: `go test ./internal/cli/...` pasa sin tocar los tests existentes de pending workspace (`TestPendingWorkspaceGroupsByRepo` en `pending_test.go:52`).
- INV2: `go build ./...` pasa desde `/home/shared/roadmapctl`.

## Contexto

`internal/cli/pending.go:114-125` define `workspaceRepoRoots(workspaceRoot string) []string` — un `filepath.WalkDir` que detecta directorios `.git` y devuelve sus padres como member roots. Esta lógica es genérica y la necesitará también `bootstrap` (T002) y potencialmente otros subcomandos.

La firma actual es lo suficientemente simple para mover como está; lo único que falta es la heurística complementaria que pregunta "¿este root parece un workspace?" — usada solo para auto-detección.

## Alcance

**In**:
1. Crear paquete nuevo `internal/workspace/detect.go` con dos funciones exportadas:
   - `MemberRoots(root string) []string` — copiada literalmente de `workspaceRepoRoots`. Mismo comportamiento (WalkDir buscando `.git`, ordenado lexicográficamente, omite el propio root si tiene `.git`).
   - `IsWorkspaceRoot(root string) bool` — retorna `true` solo si `root` **no** contiene `.git/` Y `MemberRoots(root)` devuelve ≥1 entrada Y al menos una de esas entradas contiene `docs/roadmap/.roadmapctl.toml`. La verificación del config del member usa `os.Stat`, no `config.Load` (evita ciclos de importación y trabajo innecesario).
2. Crear `internal/workspace/detect_test.go` con tests unitarios:
   - Root con `.git/` → `IsWorkspaceRoot` false; `MemberRoots` omite el root.
   - Root sin `.git/`, sin members → `IsWorkspaceRoot` false; `MemberRoots` empty.
   - Root sin `.git/`, members con `.git/` pero sin config → `IsWorkspaceRoot` false; `MemberRoots` los lista.
   - Root sin `.git/`, ≥1 member con `.git/` + `docs/roadmap/.roadmapctl.toml` → `IsWorkspaceRoot` true.
3. Reemplazar en `internal/cli/pending.go`:
   - Eliminar la función local `workspaceRepoRoots`.
   - Importar `roadmapctl/internal/workspace` (path real del módulo) y usar `workspace.MemberRoots` en `runPendingWorkspace` (línea 56).

**Out**:
- No agregar lógica nueva a `pending` ni a `bootstrap` en esta task (solo refactor).
- No mover otras funciones de `pending.go`.
- No cambiar la firma ni el orden de salida de `MemberRoots`.

## Estado inicial esperado

- `internal/workspace/` no existe.
- `grep -n "workspaceRepoRoots" /home/shared/roadmapctl/internal/cli/pending.go` retorna la definición (líneas ~114-125) y la llamada (~línea 56).

## Criterios de Aceptación

- `ls /home/shared/roadmapctl/internal/workspace/detect.go /home/shared/roadmapctl/internal/workspace/detect_test.go` ambos existen.
- `grep -rn "workspaceRepoRoots" /home/shared/roadmapctl --include="*.go"` retorna vacío.
- `grep -n "workspace.MemberRoots" /home/shared/roadmapctl/internal/cli/pending.go` muestra al menos una llamada.
- `go build ./...` pasa desde `/home/shared/roadmapctl`.
- `go test ./internal/workspace/... ./internal/cli/...` pasa.
- `IsWorkspaceRoot` no tiene callers en producción todavía (será consumido por T002), pero está cubierto por su propio test.

## Fuente de verdad

- `/home/shared/roadmapctl/internal/workspace/detect.go` (nuevo)
- `/home/shared/roadmapctl/internal/workspace/detect_test.go` (nuevo)
- `/home/shared/roadmapctl/internal/cli/pending.go` (modificado: eliminar función local, importar paquete nuevo)
