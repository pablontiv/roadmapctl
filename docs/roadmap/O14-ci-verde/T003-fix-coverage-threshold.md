---
estado: Specified
tipo: task
---
# T003: Llevar cobertura ≥ 85.0%

**Outcome**: [O14 CI verde](README.md)
**Contribuye a**: Coverage gate verde en smoke jobs

[[blocked_by:./T001-fix-golangci-issues.md]]

## Preserva

- INV1: `go test ./...` pasa después del cambio.
  - Verificar: `go test ./...`

## Contexto

La cobertura está en 84.4%, 0.6pp debajo del umbral de 85.0% configurado en
`scripts/check-coverage.sh` y en `ci.yml` (`coverage-threshold: 85`).

El script `./scripts/check-coverage.sh` corre `go test -coverprofile` sobre
todos los packages y calcula el promedio ponderado. Los packages más bajos son:
- `cmd/roadmapctl`: 50.0% — `main.go` y `version.go` apenas cubiertos
- `internal/diff`: 54.2%

Al eliminar código muerto en T001 (`plannedTask`, `newUpdateChange`, `sortedKeys`),
las líneas eliminadas dejan de contar como no-cubiertas, lo que puede subir la
cobertura por sí solo. Si no alcanza, agregar tests en los packages más bajos.

## Alcance

**In**:
1. Verificar cobertura real después de aplicar T001 (código muerto eliminado)
2. Si cobertura < 85.0%: agregar tests puntuales en `cmd/roadmapctl` o
   `internal/diff` hasta alcanzar el umbral
3. Verificar con `./scripts/check-coverage.sh` antes de pushear

**Out**:
- No bajar el umbral de 85%
- No agregar tests de coverage-farming sin valor real

## Estado inicial esperado

- T001 completado (código muerto eliminado, lint verde)
- Cobertura actual: 84.4%

## Criterios de Aceptación

- `./scripts/check-coverage.sh` reporta `coverage X.X% meets required 85.0%`
- Smoke jobs ubuntu/macos/windows verdes en el step "Coverage gate"

## Fuente de verdad

- `scripts/check-coverage.sh`
- `cmd/roadmapctl/main.go`, `cmd/roadmapctl/version.go`
- `internal/diff/`
