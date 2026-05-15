---
estado: Completed
tipo: task
---
# T003: Add `title` to `next` command output

**Outcome**: [O20 Field-agnostic integration](README.md)
**Contribuye a**: `roadmapctl next` retorna el título de cada task; los consumers (CLI, skills) no necesitan hacer un segundo lookup

[[blocked_by:./T001-add-field-mapping-to-config.md]]

## Preserva

- INV1: El output JSON de `next` sigue siendo válido y compatible hacia atrás (campo nuevo, no rompe parsers existentes)
  - Verificar: `roadmapctl next --repo /home/shared/rootline --roadmap-root docs/roadmap --output json | jq '.ready[0]'`
- INV2: `go test ./...` verde
  - Verificar: `cd /home/shared/roadmapctl && go test ./...`

## Contexto

`internal/cli/next.go` define `nextTask` como:
```go
type nextTask struct {
    Path     string   `json:"path"`
    Status   string   `json:"status"`
    Blockers []string `json:"blockers,omitempty"`
}
```

Con T001, `cfg.Fields.DisplayName` = `"titulo"`. Con rootline T001 implementado, `titulo` llega en `frontmatter` del output de rootline tree. Agregar `Title string json:"title,omitempty"` y popularlo con `stringField(frontmatter, cfg.Fields.DisplayName)`.

## Alcance

**In**:
1. Agregar `Title string json:"title,omitempty"` a `nextTask` struct
2. Popularlo con `stringField(frontmatter, cfg.Fields.DisplayName)` al construir cada nextTask

**Out**:
- No cambiar la lógica de qué tasks se consideran "ready" o "blocked"
- No cambiar otros campos del output

## Estado inicial esperado

- T001 completada (cfg.Fields.DisplayName disponible)
- rootline retorna `titulo` en frontmatter (rootline T001 completada)
- `nextTask` no tiene campo `Title`

## Criterios de Aceptación

- `nextTask` struct tiene `Title string json:"title,omitempty"`
- `roadmapctl next --repo /home/shared/rootline --roadmap-root docs/roadmap --output json | jq '.ready[].title'` retorna el H1 de cada task
- Si `titulo` no está en frontmatter, `Title` es omitido (omitempty)
- `go test ./internal/cli/...` verde

## Fuente de verdad

- `internal/cli/next.go` — nextTask struct
- `internal/config/config.go` — FieldsConfig.DisplayName
