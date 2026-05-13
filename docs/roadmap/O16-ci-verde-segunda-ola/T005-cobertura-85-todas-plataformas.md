---
estado: Completed
tipo: task
---
# T005: Llevar cobertura a ≥ 85% en Ubuntu y macOS

## Descripción

**RC6 + RC7 — coverage gates fallan en smoke/ubuntu-latest y smoke/macos-latest**

Ubuntu: 84.8% < 85.0% (gap de 0.2 pp). La brecha está en:
- `cli`: 80.1%
- `fsx`: 83.6%

macOS: 83.0% < 85.0% (gap de 2 pp). Causa: los case-collision tests se skipean
en FS case-insensitive → `lint` cae a 75.7% en macOS.

**Approach:**

1. **lint — tests unitarios platform-agnostic** (`internal/lint/schema_portability_test.go`):
   Añadir tests para `CheckSchemaCompatibility`, `CheckOutcomeSchemaCompatibility`,
   `checkEstadoSchemaCompatibility`, `checkEstadoValidateCompatibility`, y
   `reservedWindowsName` que operan sobre datos in-memory (sin syscalls de FS).
   Objetivo: llevar `lint` a ≥ 88% en macOS (compensando el skip del collision test).

2. **fsx — edge cases** (`internal/fsx/`):
   Identificar funciones con < 80% cobertura con `go tool cover -func` y añadir
   tests para las ramas sin cubrir (típicamente error paths).

3. **cli — código nuevo** (`internal/cli/`):
   Los comandos `bootstrap` e `integration` añadidos recientemente tienen paths
   sin cubrir. Identificar con `go tool cover -func` y añadir tests de integración
   para los casos normales de `bootstrap config`.

## Criterios de Aceptación

- `go test ./... -coverprofile=/tmp/cov.out && go tool cover -func=/tmp/cov.out | tail -1` reporta ≥ 85.0% en Linux
- `GOOS=darwin go test ./... -coverprofile=/tmp/cov.out && go tool cover -func=/tmp/cov.out | tail -1` reporta ≥ 85.0% (o el script check-coverage.sh pasa localmente simulando el comportamiento macOS)
- Todos los nuevos tests pasan en las 3 plataformas
