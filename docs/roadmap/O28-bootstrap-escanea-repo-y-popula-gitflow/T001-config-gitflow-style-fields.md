---
estado: Completed
tipo: task
---
# T001: Agregar campos de estilo gitflow al Config struct

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: config.go parsea y expone los 3 campos de estilo generativos del TOML

## Preserva

- INV1: `roadmapctl bootstrap --output json` sigue siendo válido sin los nuevos campos (repos que no los tienen aún).
  - Verificar: `roadmapctl bootstrap --repo /home/shared/roadmapctl --output json` sale exit 0

## Contexto

El TOML de `.roadmapctl.toml` necesita 3 nuevos campos bajo `[gitflow]`:
- `branch_style` — descripción en lenguaje natural del patrón de naming de branches
- `pr_title_style` — descripción del estilo de título de PR
- `pr_body_style` — descripción del estilo de body de PR

El binario actualmente ignora campos desconocidos. Hay que agregarlos al struct `Config` y a la serialización TOML para que `roadmapctl bootstrap --output json` los devuelva en el JSON y el skill `/integrate` los reciba en su config snapshot.

`pr_merge_strategy` (campo existente) queda deprecated: si está presente en el TOML, el bootstrap debe emitir `RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED` en diagnostics como warning (no error).

## Alcance

**In**:
1. Agregar `BranchStyle string`, `PrTitleStyle string`, `PrBodyStyle string` al struct `Config` en `internal/config/config.go`, parseables desde `[gitflow]` del TOML
2. Agregar lógica de warning `RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED` si `PRMergeStrategy != ""`
3. Tests unitarios en `internal/config/config_test.go`

**Out**:
- No modificar bootstrap.go ni el JSON output (eso es T002)
- No agregar diagnostic `RMC_GITFLOW_NOT_CONFIGURED` (eso es T002)
- No modificar skills

## Estado inicial esperado

- `internal/config/config.go` existe y compila
- `golangci-lint run ./...` sale exit 0
- `go test ./...` sale exit 0

## Criterios de Aceptación

- `grep -n 'BranchStyle\|PrTitleStyle\|PrBodyStyle' internal/config/config.go` retorna las 3 definiciones
- TOML con `[gitflow]\nbranch_style = "test"` parseado correctamente: `config.BranchStyle == "test"`
- TOML con `pr_merge_strategy = "squash"` → `roadmapctl bootstrap` emite `RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED` en diagnostics (severity warning)
- `go test ./internal/config/... -count=1` sale exit 0
- `golangci-lint run ./...` sale exit 0

## Fuente de verdad

- `internal/config/config.go` — struct Config, función Load, applyTOMLConfig
- `internal/config/config_test.go` — tests a agregar
