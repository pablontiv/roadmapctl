---
estado: Specified
tipo: task
---
# T002: Cubrir o borrar internal/templates

**Outcome**: [O29 Adoptar pkcov + reactivar gate](README.md)
**Contribuye a**: precondición para reactivar el gate; aplica política de dead code si corresponde (INV2)

## Preserva

- INV2 del outcome: no se relaja — o el paquete se cubre o se borra (no se ignora)
  - Verificar: o `go test ./internal/templates/ -cover` ≥85, o el paquete no existe

## Contexto

`internal/templates/bootstrap.go:6` contiene `GenerateStemContent` con 0% coverage. Hay dos caminos válidos según coverage-spec v1.0 sección 6:

1. **Cubrir**: si la función está en uso, agregar tests hasta ≥85% del paquete.
2. **Borrar**: si está deprecada o nadie la importa, eliminarla siguiendo política de dead code.

## Alcance

**In**:

1. Verificar uso: `grep -r "GenerateStemContent\|internal/templates" --include="*.go" .` para identificar callers.
2. **Si tiene callers vivos**: agregar tests directos para `GenerateStemContent` hasta ≥85 del paquete.
3. **Si no tiene callers vivos**: borrar `internal/templates/` (o sólo la función + el archivo si hay otras funciones vivas), verificar `go build` y `go test` siguen verdes.
4. Documentar la decisión en commit message.

**Out**:
- No refactor de la lógica.
- No tocar otros paquetes.

## Estado inicial esperado

- `go test ./internal/templates/ -cover` reporta 0% (GenerateStemContent sin tests).

## Criterios de Aceptación

- Una de las dos condiciones se cumple:
  - `go test ./internal/templates/ -cover` ≥85.0%, o
  - `GenerateStemContent` (y el paquete si quedó vacío) borrado.
- `go build ./...` verde.
- `go test ./... -race` verde.
- Commit message documenta cuál camino se tomó y por qué.

## Fuente de verdad

- `/home/shared/roadmapctl/internal/templates/bootstrap.go`
- `/home/shared/roadmapctl/` — buscar callers
