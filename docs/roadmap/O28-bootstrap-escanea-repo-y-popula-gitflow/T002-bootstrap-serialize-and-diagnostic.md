---
estado: Specified
tipo: task
---
# T002: Bootstrap serializa style fields y emite RMC_GITFLOW_NOT_CONFIGURED

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: el JSON de bootstrap incluye los campos de estilo; skill puede detectar cuando faltan

[[blocked_by:./T001-config-gitflow-style-fields.md]]

## Preserva

- INV1: el JSON output de bootstrap mantiene todos los campos existentes intactos
  - Verificar: campos existentes (`commit_style`, `auto_push`, `pr_mode`, etc.) siguen presentes en el JSON

## Contexto

Con T001 los campos existen en el struct. Ahora hay que:
1. Serializarlos en el JSON output de `roadmapctl bootstrap --output json`
2. Emitir el diagnostic `RMC_GITFLOW_NOT_CONFIGURED` (severity info) cuando los 3 campos están vacíos — esta es la señal que el skill `/roadmap` detecta en Fase 1 para ejecutar el wizard de adopción

El diagnostic no debe ser un error (no rompe el flujo si el usuario aún no adoptó).

## Alcance

**In**:
1. En `internal/cli/bootstrap.go`: incluir `branch_style`, `pr_title_style`, `pr_body_style` en el JSON output
2. Si los 3 campos están vacíos: agregar `RMC_GITFLOW_NOT_CONFIGURED` a `diagnostics[]` con severity `"info"`
3. Tests en `internal/cli/bootstrap_test.go`

**Out**:
- No modificar config.go (ya hecho en T001)
- No modificar skills

## Estado inicial esperado

- T001 completada: `BranchStyle`, `PrTitleStyle`, `PrBodyStyle` existen en Config
- `golangci-lint run ./...` sale exit 0

## Criterios de Aceptación

- `roadmapctl bootstrap --repo /home/shared/roadmapctl --output json` incluye campos `branch_style`, `pr_title_style`, `pr_body_style` en el JSON
- En repo con TOML sin los campos de estilo: el JSON incluye `RMC_GITFLOW_NOT_CONFIGURED` en `diagnostics[]`
- En repo con TOML que tiene los 3 campos no-vacíos: `RMC_GITFLOW_NOT_CONFIGURED` no aparece
- `go test ./internal/cli/... -count=1` sale exit 0
- `golangci-lint run ./...` sale exit 0

## Fuente de verdad

- `internal/cli/bootstrap.go`
- `internal/cli/bootstrap_test.go`
