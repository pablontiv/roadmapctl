---
estado: Specified
tipo: task
---
# T005: Crear .coverage-floors.toml + reemplazar bash con pkcov

**Outcome**: [O29 Adoptar pkcov + reactivar gate](README.md)
**Contribuye a**: roadmapctl usa el tooling compartido en lugar de bash + python3 local

[[blocked_by:./T001-raise-cmd-roadmapctl-coverage.md]]
[[blocked_by:./T002-cover-or-remove-internal-templates.md]]
[[blocked_by:./T003-raise-internal-workspace-coverage.md]]

## Preserva

- INV1 del outcome: threshold uniforme 85
  - Verificar: `pkcov check` exit 0 con default=85 después del swap

## Contexto

`/home/shared/roadmapctl/scripts/check-coverage.sh` (30+ líneas bash + python3) lee el threshold de `.roadmapctl.toml` `required_code_coverage` (con fallback 85). Esa convención es propia de roadmapctl; coverage-spec v1.0 dice que el threshold vive en `.coverage-floors.toml`. Esta task adopta el formato compartido.

roadmapctl ya consume `github.com/pablontiv/picokit/{diag,pathsec,autoupdate}`, así que el bump a la versión que incluye `coverage`+`pkcov` es trivial (un cambio en `go.mod`, `go mod tidy`).

Dependencia cross-repo: picokit `O03-coverage-tooling` debe estar Completed con tag publicado.

## Alcance

**In**:

1. `go.mod`: bump de `github.com/pablontiv/picokit` al tag que incluye `coverage`/`pkcov`. `go mod tidy`.
2. Crear `/home/shared/roadmapctl/.coverage-floors.toml`:
   ```toml
   default = 85
   packages = [
     # listar paquetes vivos del repo. Generar con:
     # go list ./... | sed 's|^github.com/pablontiv/roadmapctl/||'
   ]
   ```
3. Editar/crear `Justfile` recipes:
   - `coverage`: `go test ./... -coverprofile=coverage.out` + `go run github.com/pablontiv/picokit/cmd/pkcov report`
   - `coverage-check`: anterior + `go run github.com/pablontiv/picokit/cmd/pkcov check`
4. Borrar `scripts/check-coverage.sh`.
5. Si `ci.yml` aún invoca `scripts/check-coverage.sh`, actualizar a `just coverage-check` o `pkcov check` directo.
6. Considerar: `.roadmapctl.toml` `required_code_coverage` campo — si está siendo usado por otros tools (e.g. `roadmapctl loop`), preservar el valor 85 para retrocompatibilidad. Documentar en commit que ahora hay dos fuentes (`.coverage-floors.toml` para el gate, `.roadmapctl.toml` campo para retrocompat de tooling roadmapctl).

**Out**:
- No tocar pre-push (T006).
- No cambiar el threshold (sigue 85).
- No borrar `.roadmapctl.toml.required_code_coverage` sin verificar callers.

## Estado inicial esperado

- T001-T003 completadas.
- T004 puede haberse mergeado ya (reactivación de CI threshold).
- `scripts/check-coverage.sh` existe y referencia `.roadmapctl.toml`.

## Criterios de Aceptación

- `go.mod` incluye `github.com/pablontiv/picokit` en versión que incluye `coverage`.
- `.coverage-floors.toml` existe con `default = 85` y la lista completa de paquetes vivos.
- `Justfile` recipes `coverage` y `coverage-check` invocan `pkcov`.
- `scripts/check-coverage.sh` borrado.
- `just coverage-check` exit 0 local.
- CI verde post-push.
- Decisión sobre `.roadmapctl.toml.required_code_coverage` documentada en commit.

## Fuente de verdad

- `/home/shared/roadmapctl/Justfile`
- `/home/shared/roadmapctl/scripts/check-coverage.sh` (a borrar)
- `/home/shared/roadmapctl/.github/workflows/ci.yml`
- `/home/shared/roadmapctl/docs/roadmap/.roadmapctl.toml` — campo `required_code_coverage`
- `/home/shared/rootline/Justfile` — patrón a imitar
- `/home/shared/picokit/cmd/pkcov/` — binario consumido
