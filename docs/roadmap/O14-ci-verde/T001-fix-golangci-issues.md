---
estado: Completed
tipo: task
---
# T001: Corregir los 28 issues de golangci-lint

**Outcome**: [O14 CI verde](README.md)
**Contribuye a**: job `ci / Lint` verde en CI

## Preserva

- INV1: `go test ./...` pasa sin errores antes y después del cambio.
  - Verificar: `go test ./...`

## Contexto

El job `ci / Lint` usa golangci-lint v2.10.1 vía crossbeam `go-ci.yml`. Ahora que
el linter corre por primera vez, encontró 28 issues reales en el codebase:

- **gosec G204** (`bootstrap_test.go:119`, `rootlinecli/client.go:280`): subprocess
  lanzado con variable. Son patrones legítimos del CLI — agregar a `excludes` en
  `.golangci.yml`.
- **gosec G302** (`materialize/dryrun.go:151,197`): permisos 0644 en archivos
  Markdown. Legítimo para archivos no-secretos — agregar a `excludes`.
- **ineffassign** (`cli/doctor.go:34-35`): `repoRoot` y `roadmapRoot` asignados
  pero nunca usados.
- **unused**: `plannedTask` (dryrun.go:78), `newUpdateChange` (dryrun.go:620),
  `sortedKeys` (lint/schema_portability.go:178) — código muerto, eliminar.
- **staticcheck SA1012** (`cli/pathplan_test.go:56,100`): `nil` context pasado
  donde se espera `context.Context` — reemplazar con `context.TODO()`.
- **staticcheck SA1019** (`lint/markdown.go:125`): `typed.Text` deprecado —
  usar la propiedad recomendada por el nodo.
- **staticcheck SA9003** (`materialize/pathplan.go:187`): rama `if` vacía.
- **staticcheck QF1001** (`materialize/pathplan.go:204`,
  `roadmap/structure.go:195`): aplicar De Morgan's law.

## Alcance

**In**:
1. Agregar `G204` y `G302` a `linters.settings.gosec.excludes` en `.golangci.yml`
2. Eliminar `plannedTask`, `newUpdateChange`, `sortedKeys`
3. Corregir `ineffassign` en `doctor.go` (eliminar asignaciones muertas)
4. Reemplazar `nil` context con `context.TODO()` en `pathplan_test.go`
5. Corregir `SA1019` en `lint/markdown.go`
6. Corregir `SA9003` en `materialize/pathplan.go`
7. Aplicar De Morgan en `materialize/pathplan.go` y `roadmap/structure.go`

**Out**:
- No cambiar comportamiento observable del CLI
- No agregar tests nuevos (eso es T003)

## Estado inicial esperado

- `.golangci.yml` en versión v2 con syntax `settings:` nested (ya corregido)
- `go test ./...` pasa localmente
- golangci-lint v2.10.1 reporta 28 issues en CI

## Criterios de Aceptación

- `go test ./...` pasa
- `golangci-lint run ./...` reporta 0 issues localmente
- job `ci / Lint` verde en el siguiente push a master

## Fuente de verdad

- `.golangci.yml`
- `internal/cli/doctor.go`
- `internal/cli/pathplan_test.go`
- `internal/lint/markdown.go`
- `internal/lint/schema_portability.go`
- `internal/materialize/dryrun.go`
- `internal/materialize/pathplan.go`
- `internal/roadmap/structure.go`
- `internal/rootlinecli/client.go` (solo lectura para entender G204)
