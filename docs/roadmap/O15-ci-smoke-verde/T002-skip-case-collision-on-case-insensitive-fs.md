---
estado: Completed
tipo: task
---
# T002: Skip tests de colisión de case en FS case-insensitive

**Outcome**: [O15 CI smoke verde](README.md)
**Contribuye a**: job `Smoke / macOS` verde en CI

## Contexto

Dos tests fallan en macOS (APFS, case-insensitive):

1. `TestCheckFilenamePortabilityReportsCaseCollisionAndReservedName` — crea
   `T001-task.md` y `t001-task.md` en `t.TempDir()`. En macOS ambos son el
   mismo archivo; `ReadDir` retorna solo una entrada → colisión no detectada.

2. `TestCheckGoldenJSONFixtures/lint_case_collision` — el fixture tiene ambos
   `t001-first.md` y `T001-first.md` en git, pero macOS git solo checkoutea uno.
   El lint command retorna exit=0 (warning) en vez de exit=1 (error).

La producción es correcta: en FS case-insensitive no puede haber colisión real.
El problema es que los tests asumen un FS case-sensitive.

## Alcance

**In**:
1. En `TestCheckFilenamePortabilityReportsCaseCollisionAndReservedName`
   (`internal/lint/schema_portability_test.go`): detectar si el FS es
   case-insensitive (intentar crear dos archivos con mismo nombre en distinto
   case y verificar si resultan en uno solo) y `t.Skip(...)` si aplica.
2. En `TestCheckGoldenJSONFixtures` (`internal/cli/golden_test.go`): añadir
   condición para saltear el caso `lint_case_collision` cuando el FS es
   case-insensitive.

**Out**:
- No modificar `checkCaseCollisionsInDir` ni la lógica de producción
- En Linux el test sigue ejercitando el path feliz sin skip

## Criterios de Aceptación

- `go test ./...` pasa en macOS sin skip manual
- `go test ./...` sigue pasando en Linux con el test ejecutándose completo
- `TestCheckFilenamePortabilityReportsCaseCollisionAndReservedName` no falla por FS
