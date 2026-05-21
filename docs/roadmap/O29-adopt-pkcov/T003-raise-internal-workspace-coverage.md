---
estado: Specified
tipo: task
---
# T003: Subir internal/workspace a ≥85% coverage

**Outcome**: [O29 Adoptar pkcov + reactivar gate](README.md)
**Contribuye a**: precondición para reactivar el gate (INV2)

## Preserva

- INV2 del outcome: el pre-work no relaja el contrato
  - Verificar: `go test ./internal/workspace/ -cover` ≥85

## Contexto

`internal/workspace` está hoy en 82.6% (medido `2026-05-21`) — apenas debajo del piso. Cerrar la brecha de ~2.5 puntos.

Estrategia: identificar funciones con menor cobertura usando `go tool cover -func` y agregar tests específicos.

## Alcance

**In**:

1. Identificar funciones de `internal/workspace` con cobertura <85.
2. Agregar tests cubriendo las ramas faltantes (foco en error paths y casos límite).
3. Validar que el paquete cruza 85.

**Out**:
- No refactor.
- No tocar otros paquetes.

## Estado inicial esperado

- `go test ./internal/workspace/ -cover` reporta ~82.6%.

## Criterios de Aceptación

- `go test ./internal/workspace/ -cover` ≥85.0%.
- `go test ./... -race` verde.
- Diff incluye sólo archivos `_test.go` en `internal/workspace/`.
- `golangci-lint run` sin issues nuevos.

## Fuente de verdad

- `/home/shared/roadmapctl/internal/workspace/`
