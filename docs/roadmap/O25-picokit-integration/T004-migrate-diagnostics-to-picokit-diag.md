---
estado: Completed
tipo: task
---
# T004: Replace diagnostics infrastructure with picokit/diag (type-alias strategy)

**Outcome**: [O25 Integrate picokit as a dependency](README.md)
**Contribuye a**: Eliminar la infraestructura duplicada de tipos y funciones de diagnóstico (~150 líneas).

[[blocked_by:./T001-add-picokit-dependency.md]]

## Preserva

- INV1: Los 26 callers de `internal/diagnostics` compilan sin ningún cambio.
  - Verificar: `go build ./...` pasa desde `/home/shared/roadmapctl` sin modificar ningún archivo fuera de `internal/diagnostics/`.
- INV2: El contrato JSON de salida no cambia (mismos campos, mismo orden de serialización).
  - Verificar: `go test ./...` pasa incluyendo tests de integración que verifican JSON output.
- INV3: Coverage general no baja.
  - Verificar: `just coverage-summary` desde `/home/shared/roadmapctl`.

## Contexto

`internal/diagnostics/report.go` duplica la infraestructura de tipos de `picokit/diag`. La estrategia de migración es usar **type aliases** para los tipos hoja genéricos y re-exportar las funciones utilitarias — así los 26 callers no necesitan ningún cambio.

`Report` y `NewReport` **permanecen en roadmapctl** porque son dueños del campo `RoadmapRoot string` que es lógica de dominio de roadmapctl (no pertenece a picokit).

**Tipos que pasan a ser aliases de picokit/diag:**
- `Summary = diag.Summary`
- `Diagnostic = diag.Diagnostic`
- `Severity = diag.Severity`

**Constantes que pasan a ser aliases de picokit/diag:**
- `ExitOK`, `ExitValidation`, `ExitUsage`, `ExitEnvironment`, `ExitInternal`
- `SeverityInfo`, `SeverityWarning`, `SeverityError`
- `SummaryStatusOK`, `SummaryStatusWarning`, `SummaryStatusError`

**Funciones que pasan a ser vars de picokit/diag:**
- `ExitCode = diag.ExitCode`
- `RenderJSON = diag.RenderJSON`
- `RenderText = diag.RenderText`

**Se mantiene en roadmapctl** (dueños de dominio):
- `Report` struct — incluye `RoadmapRoot string`
- `NewReport(kind, root, roadmapRoot string, diagnostics []Diagnostic) Report`
- Todos los `Diagnostic*` string constants (RMC_LINT_*, RMC_TRANSITION_*, etc.)

**Resultado**: `internal/diagnostics/report.go` pasa de ~170 líneas a ~50 líneas. Los `Diagnostic*` constants se mueven a un archivo separado `constants.go` para claridad.

API de picokit/diag después de T002 de picokit:
- `diag.Summary`, `diag.Diagnostic`, `diag.Severity` — tipos genéricos
- `diag.NewReport(kind, root string, diagnostics []Diagnostic) Report` — 3 params (sin RoadmapRoot)
- `diag.ExitCode`, `diag.RenderJSON`, `diag.RenderText` — funciones utilitarias
- Exit code constants: `diag.ExitOK`, etc.
- Severity constants: `diag.SeverityInfo`, etc.

## Alcance

**In**:
1. Reemplazar el contenido de `internal/diagnostics/report.go` con:
   - Type aliases: `Summary = diag.Summary`, `Diagnostic = diag.Diagnostic`, `Severity = diag.Severity`.
   - Const aliases para exit codes, severity y summary status.
   - Var aliases para `ExitCode`, `RenderJSON`, `RenderText`.
   - `Report` struct (mantener, dueño de `RoadmapRoot`).
   - `NewReport` function (mantener, llama a `diag.NewReport` internamente y setea `RoadmapRoot`).
2. Crear `internal/diagnostics/constants.go` moviendo todos los `Diagnostic*` string constants (RMC_*) desde `report.go` si estaban allí, o desde donde estén actualmente.
3. Eliminar `internal/diagnostics/doc.go` si solo contiene la declaración de paquete (sin contenido útil).

**Out**:
- No modificar ningún archivo fuera de `internal/diagnostics/`.
- No cambiar los 26 archivos callers.
- No modificar el contrato JSON (mismos campos `json:"..."` en Report y Diagnostic).

## Estado inicial esperado

- `/home/shared/roadmapctl/internal/diagnostics/report.go` contiene la infraestructura completa duplicada.
- `grep -r '"github.com/pablontiv/roadmapctl/internal/diagnostics"' /home/shared/roadmapctl --include="*.go" | wc -l` retorna 26.
- picokit v0.1.0 está en go.mod (T001 completada).

## Criterios de Aceptación

- `grep -r "RoadmapRoot\|roadmap_root" /home/shared/picokit/diag/` retorna vacío (picokit limpio — prerequisito externo de picokit T002).
- `go build ./...` pasa desde `/home/shared/roadmapctl` sin modificar callers.
- `go test ./...` pasa desde `/home/shared/roadmapctl`.
- `just check` pasa (gofmt + go vet).
- `wc -l /home/shared/roadmapctl/internal/diagnostics/report.go` muestra ≤ 60 líneas.
- La infraestructura duplicada (tipos Summary, Diagnostic, Severity, funciones ExitCode/RenderJSON/RenderText) ya no está definida en roadmapctl sino aliased desde picokit/diag.

## Fuente de verdad

- `/home/shared/roadmapctl/internal/diagnostics/report.go`
- `/home/shared/roadmapctl/internal/diagnostics/doc.go`
- `/home/shared/picokit/diag/report.go` (referencia de API destino)
- `/home/shared/roadmapctl/internal/diagnostics/report_test.go` (si existe, verificar que sigue pasando)
