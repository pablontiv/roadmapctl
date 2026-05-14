---
estado: Specified
tipo: task
---
# T004: Update remaining roadmapctl consumers to use config field names

**Outcome**: [O20 Field-agnostic integration](README.md)
**Contribuye a**: ningún archivo de roadmapctl hardcodea "blocked_by" u otros nombres de campo

[[blocked_by:./T001-add-field-mapping-to-config.md]]

## Preserva

- INV1: Comportamiento funcional idéntico al actual con defaults
  - Verificar: `roadmapctl check --repo /home/shared/rootline --roadmap-root docs/roadmap --output json --strict`
- INV2: `go test ./...` verde
  - Verificar: `cd /home/shared/roadmapctl && go test ./...`

## Contexto

Quedan cuatro consumers que hardcodean nombres de campo:

**dependencies.go** (`internal/roadmap/dependencies.go`):
```go
stringField(link, "type") != "blocked_by"
```
→ usar `cfg.Fields.DependencyLink`

**structure.go** (`internal/roadmap/structure.go`):
```go
// package-level var
blockedByLinkPattern = regexp.MustCompile(`\[\[blocked_by:([^\]]+)\]\]`)
```
→ construir dinámicamente:
```go
func blockedByPattern(linkName string) *regexp.Regexp {
    return regexp.MustCompile(fmt.Sprintf(`\[\[%s:([^\]]+)\]\]`, regexp.QuoteMeta(linkName)))
}
```
Las funciones de validación reciben el patrón desde config.

**schema_portability.go** (`internal/lint/schema_portability.go`):
```go
rules["blocked_by"]
```
→ usar `rules[cfg.Fields.DependencyLink]`

**bootstrap.go** (`internal/templates/bootstrap.go`):
El template de .stem hardcodea `"blocked_by"` en el link. Usar `cfg.Fields.DependencyLink`.

## Alcance

**In**:
1. dependencies.go: `!= "blocked_by"` → `!= cfg.Fields.DependencyLink`
2. structure.go: eliminar package-level var `blockedByLinkPattern`; agregar helper `blockedByPattern(linkName string) *regexp.Regexp`; pasar el patrón desde config a las funciones de validación
3. schema_portability.go: `rules["blocked_by"]` → `rules[cfg.Fields.DependencyLink]`
4. bootstrap.go: template usa `cfg.Fields.DependencyLink` en el link del .stem generado

**Out**:
- No cambiar la lógica de validación de dependencies/structure, solo el nombre del campo
- No tocar model.go, status.go ni next.go (ya cubiertos en T002/T003)

## Estado inicial esperado

- T001 completada (Config.Fields disponible)
- `go test ./...` pasa

## Criterios de Aceptación

- Ningún string literal `"blocked_by"` en dependencies.go, structure.go, schema_portability.go, bootstrap.go
- `roadmapctl check --repo /home/shared/rootline --roadmap-root docs/roadmap --output json --strict` exit 0
- `go test ./internal/roadmap/... ./internal/lint/... ./internal/templates/...` verde

## Fuente de verdad

- `internal/roadmap/dependencies.go`
- `internal/roadmap/structure.go`
- `internal/lint/schema_portability.go`
- `internal/templates/bootstrap.go`
- `internal/config/config.go` — FieldsConfig.DependencyLink
