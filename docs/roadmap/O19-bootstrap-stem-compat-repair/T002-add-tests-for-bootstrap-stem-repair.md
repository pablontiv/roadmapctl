---
estado: Completed
tipo: task
---
# T002: Add tests for bootstrap stem repair

**Outcome**: [O19 Bootstrap stem compat repair](README.md)
**Contribuye a**: cobertura de los paths feliz/triste del repair interactivo

[[blocked_by:./T001-implement-bootstrap-stem-interactive-repair.md]]

## Preserva

- INV1: La suite de tests existente (go test ./...) sigue pasando.
  - Verificar: `go test ./...` sin errores antes y después.
- INV2: Coverage no baja del threshold configurado (85%).
  - Verificar: script de coverage o CI.

## Contexto

T001 agrega lógica nueva a bootstrap. Esta task agrega tests que cubren todos los paths del repair: detección, confirmación interactiva, flag --yes, stem no reconocido, y verificación de postcheck.

Los tests deben usar fixtures de `.stem` legacy (con `required.match: ["O*", "T*"]` y/o `validate estado non_empty`) y fixtures del `.stem` canónico para verificar el resultado.

## Alcance

**In**:
1. Fixture de `.stem` legacy que dispara ambos diagnostics.
2. Test: bootstrap con stem legacy → detecta, reporta diagnostics, muestra diff, prompt.
3. Test: bootstrap con stem legacy + stdin "y" → aplica repair; check --strict pasa después.
4. Test: bootstrap con stem legacy + `--yes` → aplica repair sin prompt interactivo.
5. Test: bootstrap con stem legacy + stdin "N" → no modifica nada; bloqueo persiste.
6. Test: bootstrap con stem custom no reconocido → `RMC_BOOTSTRAP_REPAIR_UNSUPPORTED_STEM`; no modifica nada.
7. Test: bootstrap con stem ya canónico → no activa repair path; bootstrap pasa normalmente.

**Out**:
- No agregar tests de integración que requieran repos externos reales.
- No testear comportamiento de doctor/check standalone.

## Estado inicial esperado

- T001 completado: bootstrap tiene el repair implementado.
- Existen fixtures de `.stem` en el directorio de tests (o se crean en esta task).

## Criterios de Aceptación

- `go test ./...` pasa con los nuevos tests incluidos.
- Los 7 escenarios del alcance tienen cobertura en tests unitarios o de integración.
- Coverage no baja del 85%.
- `golangci-lint run ./...` reporta 0 issues.

## Fuente de verdad

- `internal/cli/bootstrap_test.go` (o equivalente)
- Fixtures en `testdata/` o directorio de tests existente
