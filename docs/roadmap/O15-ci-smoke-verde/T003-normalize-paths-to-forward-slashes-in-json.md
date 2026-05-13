---
estado: Specified
tipo: task
---
# T003: Normalizar Root y RoadmapRoot a forward slashes en JSON

**Outcome**: [O15 CI smoke verde](README.md)
**Contribuye a**: job `Smoke / Windows` verde en CI

## Contexto

En Windows, `root` y `roadmap_root` en el JSON output tienen backslashes:

```json
"root": "D:\\a\\roadmapctl\\roadmapctl\\testdata\\fixtures\\...",
"roadmap_root": "D:\\a\\roadmapctl\\roadmapctl\\testdata\\fixtures\\...\\docs\\roadmap"
```

`AssertNoBackslashes` rechaza cualquier backslash en el JSON report, haciendo
fallar todos los golden tests en Windows smoke.

`diagnostics.NewReport` almacena las strings crudas recibidas (absolute Windows
paths). La normalización debe hacerse al construir el report.

## Alcance

**In**:
1. En `internal/diagnostics/report.go`, dentro de `NewReport`, aplicar
   `filepath.ToSlash(root)` y `filepath.ToSlash(roadmapRoot)` antes de
   asignarlos a `Report.Root` y `Report.RoadmapRoot`
2. Verificar que ningún golden JSON existente se rompa (los goldens de Linux
   ya usan forward slashes)

**Out**:
- No modificar otros campos de `Diagnostic` (paths relativos ya usan `relToRoot`
  que aplica `filepath.ToSlash`)
- No tocar tests o goldens manualmente; si cambian, actualizar con `-update`

## Criterios de Aceptación

- `go test ./...` pasa en Windows (AssertNoBackslashes ok en todos los goldens)
- `go test ./...` sigue pasando en Linux sin regresión
- `root` y `roadmap_root` en todo JSON output usan `/` como separador
