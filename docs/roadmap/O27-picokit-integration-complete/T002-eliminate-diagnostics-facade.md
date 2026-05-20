---
estado: Completed
tipo: task
---
# T002: Eliminate internal/diagnostics facade (consolidate via picokit/diag + internal/reports)

**Outcome**: [O27 Complete picokit integration](README.md)
**Contribuye a**: roadmapctl elimina el facade redundante y consume `picokit/diag` directamente; lo específico de dominio (campo `RoadmapRoot`, constantes `RMC_*`) vive en `internal/reports`.

## Preserva

- INV1: Los outputs CLI de los comandos que renderean Report (`check`, `lint`, etc.) permanecen byte-for-byte idénticos.
  - Verificar: golden tests con fixtures actuales pasan sin tocar.
- INV2: Los códigos de exit (`ExitOK`, `ExitValidation`, etc.) producen los mismos integers que hoy.
  - Verificar: `TestReport_DelegatesToDiag` compara contra el comportamiento de `picokit/diag.ExitCode`.

## Contexto

`internal/diagnostics/report.go:1-61` re-exporta vía type alias 9 símbolos de `picokit/diag` (Summary, Diagnostic, Severity, constantes Exit*/Severity*/SummaryStatus*) más define el `Report` struct con el campo `RoadmapRoot` específico de roadmapctl, y `internal/diagnostics/constants.go:1-49` define las constantes de Diagnostic ID (`RMC_LINT_*`, `RMC_TRANSITION_*`, `RMC_BOOTSTRAP_*`, `RMC_MATERIALIZE_*`).

26 callsites del proyecto importan `internal/diagnostics`. La facade fue creada en O25-T004 con la estrategia de type-alias para minimizar churn; ahora se elimina porque los aliases puros no aportan abstracción y duplican el conocimiento.

Plan:
1. Crear `internal/reports/reports.go` con el `Report` struct + `NewReport(kind, root, roadmapRoot string, diagnostics []picokit.Diagnostic)`.
2. Crear `internal/reports/constants.go` con las constantes `RMC_*`.
3. Reescribir los 26 callsites: importar `github.com/pablontiv/picokit/diag` para tipos genéricos, `roadmapctl/internal/reports` para los de dominio.
4. Borrar `internal/diagnostics/` completo.

## Alcance

**In**:
1. `internal/reports/reports.go` con `Report` struct, `NewReport`, helpers `ExitCode/RenderJSON/RenderText` que delegan a `picokit/diag`.
2. `internal/reports/constants.go` con `RMC_*` constants.
3. Migrar `internal/diagnostics/report_test.go` y `text_test.go` → `internal/reports/` preservando golden fixtures.
4. Reescribir 26 callsites (`grep -rln "roadmapctl/internal/diagnostics"`).
5. Borrar `internal/diagnostics/` (4 archivos).

**Out**:
- No tocar `internal/updater/` (T001).
- No cambiar la API pública del CLI ni el shape JSON de los Reports.

## Estado inicial esperado

- `internal/diagnostics/{report,constants,report_test,text_test}.go` existen.
- 26 imports de `roadmapctl/internal/diagnostics` en el árbol.

## Criterios de Aceptación

- `internal/reports/{reports,constants,reports_test,text_test}.go` existen.
- `grep -rln "roadmapctl/internal/diagnostics" /home/shared/roadmapctl --include="*.go"` retorna vacío.
- `grep -rln "roadmapctl/internal/reports" /home/shared/roadmapctl --include="*.go"` retorna 26 imports (1:1 con la migración).
- `go build ./...` pasa.
- `go test ./... -race -count=1` pasa.
- Golden tests de `check`/`lint` byte-for-byte idénticos contra fixtures previos.
- `scripts/check-coverage.sh` pasa con `internal/reports` ≥85%.
- `roadmapctl check --output json` produce JSON con el mismo shape que hoy (validar con `jq` comparando keys).
