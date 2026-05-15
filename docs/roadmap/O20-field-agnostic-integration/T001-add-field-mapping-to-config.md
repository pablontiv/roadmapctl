---
estado: Completed
tipo: task
---
# T001: Add field mapping to `config.go`

**Outcome**: [O20 Field-agnostic integration](README.md)
**Contribuye a**: roadmapctl puede configurar el vocabulario de campos sin recompilar

## Preserva

- INV1: Repos existentes sin `[fields]` en `.roadmapctl.toml` funcionan igual que hoy
  - Verificar: `roadmapctl check --repo /home/shared/rootline --roadmap-root docs/roadmap --output json --strict`
- INV2: `go test ./...` verde
  - Verificar: `cd /home/shared/roadmapctl && go test ./...`

## Contexto

`internal/config/config.go` define el struct `Config` que parsea `.roadmapctl.toml`. Actualmente no tiene sección para field names. Los consumers como `model.go`, `status.go`, `next.go` etc. hardcodean strings como `"estado"`, `"tipo"`, `"blocked_by"`.

Agregar una sección `Fields` al struct con defaults que preservan el comportamiento actual:

```toml
[fields]
lifecycle       = "estado"
record_type     = "tipo"
task_value      = "task"
outcome_value   = "outcome"
display_name    = "titulo"
dependency_link = "blocked_by"
```

Los defaults deben aplicarse cuando la sección `[fields]` no existe en el TOML.

## Alcance

**In**:
1. Agregar struct `FieldsConfig` con campos: Lifecycle, RecordType, TaskValue, OutcomeValue, DisplayName, DependencyLink
2. Agregar `Fields FieldsConfig` a `Config` struct
3. Configurar defaults retrocompatibles (lifecycle="estado", record_type="tipo", etc.)
4. Tests unitarios que verifican defaults y override desde TOML

**Out**:
- No cambiar ningún consumer todavía (model.go, status.go, etc.)
- No cambiar comportamiento de ningún comando

## Estado inicial esperado

- `Config` struct en `internal/config/config.go` no tiene sección `Fields`
- `go test ./...` pasa

## Criterios de Aceptación

- `Config.Fields.Lifecycle` retorna `"estado"` cuando no hay `[fields]` en TOML
- `Config.Fields.DependencyLink` retorna `"blocked_by"` por default
- Todos los campos de `FieldsConfig` tienen defaults que preservan comportamiento actual
- `go test ./internal/config/...` verde con tests de defaults y override

## Fuente de verdad

- `internal/config/config.go` — Config struct, defaults
