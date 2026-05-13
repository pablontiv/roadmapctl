---
estado: Completed
tipo: task
---

# T016 Bootstrap idempotente: eliminar context

`roadmapctl context` viola la separación de responsabilidades del ecosistema: llama `rootline describe` para obtener el schema del `.stem`, responsabilidad que pertenece a rootline. Además no inicializa repos nuevos.

Solución: convertir `roadmapctl bootstrap` (sin subcomando) en el comando primario idempotente que reemplaza `context`, y eliminar `context` completamente.

## Criterios de Aceptación

- AC1: `roadmapctl bootstrap --output json` retorna config + helpers (mismo JSON que `context` excepto el campo `schema`, que se elimina)
- AC2: Si falta `.roadmapctl.toml`, lo crea con defaults y continúa
- AC3: Si falta `.stem`, lo crea con el template base y continúa
- AC4: `context.go` y `context_test.go` eliminados completamente — sin aliases, sin stubs
- AC5: `cli.go` no registra el comando `context`
- AC6: `go test ./...` pasa; golden tests actualizados para bootstrap

## Especificación Técnica

### Archivos modificados
- `internal/cli/bootstrap.go`: agregar handler para invocación sin subcomando (`bootstrap` directo); ejecutar `proposedBootstrapChanges` + `applyBootstrapChanges` si hay cambios; retornar JSON de config (sin schema de rootline)
- `internal/cli/cli.go`: remover registro de `context`; `bootstrap` sin subcomando retorna config JSON
- `internal/cli/context.go`: **eliminar**
- `internal/cli/context_test.go`: **eliminar** (cobertura migrada a bootstrap_test.go)
- `internal/cli/golden_test.go`: actualizar golden cases de `context` a `bootstrap`

### Separación de responsabilidades
El campo `schema` (valores de `estado`/`tipo` leídos de `.stem` via `rootline describe`) se elimina del output de bootstrap. Quien necesite el schema llama `rootline describe` directamente — roadmapctl no envuelve responsabilidades de rootline.
