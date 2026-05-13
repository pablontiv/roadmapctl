---
estado: Completed
tipo: task
---
# T005: Agregar tests cross-platform para bootstrapApplyDiagnostic y lintNameDiagnostic

## Descripción

**RC5 — `smoke/windows-latest` falla con coverage 84.8% < 85.0%**

Con las correcciones de T001-T004, los tests que cubren `bootstrapApplyDiagnostic`
(internal/cli) y `lintNameDiagnostic` (internal/lint) se saltean en Windows:

- `TestBootstrapInitApplyReportsDiagnosticsOnFileError` → skipeado (`runtime.GOOS == "windows"`)
- `TestCheckFilenamePortabilityReportsCaseCollisionAndReservedName` → skipeado (FS case-insensitive)
- `TestCheckFilenamePortabilityDetectsReservedName` → skipeado (`runtime.GOOS == "windows"`)

Resultado: `bootstrapApplyDiagnostic` y `lintNameDiagnostic` quedan al 0% en Windows.
La cobertura total cae a 84.8%, por debajo del umbral de 85.0%.

**Solución:** Agregar tests unitarios simples que llamen directamente a las dos
funciones sin depender del sistema de archivos ni de comportamiento específico
de plataforma. Estos tests corren en las 3 plataformas (Linux/macOS/Windows).

## Criterios de Aceptación

- `TestBootstrapApplyDiagnosticFormat` en `bootstrap_test.go` cubre `bootstrapApplyDiagnostic`
- `TestLintNameDiagnosticFormat` en `schema_portability_test.go` cubre `lintNameDiagnostic`
- `./scripts/check-coverage.sh` pasa en simulación Windows (skip de 3 tests) con total >= 85.0%
- `golangci-lint run ./...` sin errores
