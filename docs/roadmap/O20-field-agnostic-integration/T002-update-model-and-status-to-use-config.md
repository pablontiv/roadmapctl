---
estado: Specified
tipo: task
---
# T002: Update `model.go` and `status.go` to use config field names

**Outcome**: [O20 Field-agnostic integration](README.md)
**Contribuye a**: model y status no hardcodean "estado"/"tipo"; leen desde cfg.Fields

[[blocked_by:./T001-add-field-mapping-to-config.md]]

## Preserva

- INV1: Comportamiento funcional idéntico al actual con defaults
  - Verificar: `roadmapctl next --repo /home/shared/rootline --roadmap-root docs/roadmap --output json`
- INV2: `go test ./...` verde
  - Verificar: `cd /home/shared/roadmapctl && go test ./...`

## Contexto

`internal/roadmap/model.go` construye el modelo de roadmap desde el output de rootline. Hardcodea:
```go
statusByPath[path] = stringField(frontmatter, "estado")
typeByPath[path]   = stringField(frontmatter, "tipo")
if stringField(frontmatter, "tipo") != "task" { ... }
if stringField(edge, "type") != "blocked_by" { ... }
```

`internal/roadmap/status.go` tiene:
```go
stringField(frontmatter, "estado")
stringField(frontmatter, "tipo")
extractStatusValues(schema, "estado")
```

Todos deben usar `cfg.Fields.*` en vez de strings literales.

## Alcance

**In**:
1. model.go: reemplazar `"estado"` → `cfg.Fields.Lifecycle`, `"tipo"` → `cfg.Fields.RecordType`, `"task"` → `cfg.Fields.TaskValue`, `"blocked_by"` → `cfg.Fields.DependencyLink`
2. status.go: mismos reemplazos, más `extractStatusValues(schema, cfg.Fields.Lifecycle)`
3. Pasar `cfg` a las funciones que lo necesiten (o usar un receiver si ya existe)

**Out**:
- No cambiar la lógica de model ni status, solo los nombres de campo
- No tocar next.go, dependencies.go, structure.go todavía

## Estado inicial esperado

- T001 completada (Config.Fields disponible)
- `go test ./...` pasa

## Criterios de Aceptación

- Ningún string literal `"estado"`, `"tipo"`, `"task"`, `"outcome"`, `"blocked_by"` en model.go o status.go (excepto en tests con config explícita)
- `roadmapctl next --repo /home/shared/rootline --roadmap-root docs/roadmap --output json` retorna resultado idéntico al actual
- `go test ./internal/roadmap/...` verde

## Fuente de verdad

- `internal/roadmap/model.go`
- `internal/roadmap/status.go`
- `internal/config/config.go` — FieldsConfig
