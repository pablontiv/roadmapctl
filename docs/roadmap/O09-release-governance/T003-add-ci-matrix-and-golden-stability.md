---
estado: Completed
tipo: task
---
# T003: Fortalecer CI matrix y estabilidad de goldens

**Outcome**: [O09 Release/governance](README.md)
**Contribuye a**: CE2, INV2

[[blocked_by:../O02-post-mvp-foundations/T007-expand-foundation-fixtures-and-goldens.md]]

## Preserva

- INV1: Goldens no se actualizan para esconder regresiones.
  - Verificar: docs/golden-tests.

## Contexto

La superficie post-MVP aumenta comandos, JSON outputs y fixtures. CI debe proteger estabilidad cross-platform.

## Alcance

**In**:
1. Revisar CI Linux/macOS/Windows.
2. Normalizar paths en goldens.
3. Agregar tests de help/error/stdout-stderr para Cobra.
4. Evaluar matrix Rootline latest/minimum cuando se defina política.
5. Documentar workflow para actualizar goldens.
6. Enforce total Go statement coverage mínimo de 85% en CI y release gates.

**Out**:
- Publicar releases.

## Estado inicial esperado

- CI corre `go test ./...` y `go build ./cmd/roadmapctl`.

## Criterios de Aceptación

- CI pasa en OS matrix.
- Tests detectan ruido en JSON stdout.
- Goldens nuevos son determinísticos.
- `./scripts/check-coverage.sh` pasa y reporta coverage total >= 85%.

## Fuente de verdad

- `.github/workflows/ci.yml`
- `docs/golden-tests.md`
- `internal/cli/golden_test.go`
