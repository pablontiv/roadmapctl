---
estado: Specified
tipo: task
---
# T004: Añadir helper requiresRealRootline y skipear tests que lo necesitan

## Descripción

**RC3 — tests que esperan fallos de rootline pasan con fake (`ci/Test`)**

El fake rootline siempre devuelve éxito en `validate` y graph vacío en `graph`. 
Los tests que esperan que rootline *encuentre errores* (ciclos, broken_by, status
inválido, decision scoring) reciben `Status:"ok"` en lugar de los errores esperados.

Tests afectados (todos en `internal/cli/`):
- `check_test.go`: `TestCheckInvalidCycleExitsValidation`, `TestCheckBrokenBlockedByExitsValidation`, `TestCheckStatusMismatchExitsValidation`
- `golden_test.go`: `invalid_status_bogus`, `invalid_stale_outcome_stem`, `bare_blocked_by`
- `golden_test.go`: `pending_direct_tasks`, `pending_outcome_tasks`, `next_ready_blocked`, `decision_reverse_dependencies` (fake `query`/`tree`/`graph` devuelve vacío → no hay tasks en la respuesta)
- `golden_test.go`: `can_start_ready`, `can_start_blocked` (transition depende de graph)
- `TestDecisionJSONIncludesDeterministicReasons`, `TestReadOnlyTextGoldens`

Fix: añadir helper en `golden_test.go`:
```go
func requiresRealRootline(t *testing.T) {
    t.Helper()
    if os.Getenv("ROADMAPCTL_FAKE_ROOTLINE") == "1" {
        t.Skip("skipping: requires real rootline")
    }
}
```
Y llamarlo al inicio de cada test/subtest afectado. Estos tests siguen corriendo
en `smoke` (que instala rootline real).

## Criterios de Aceptación

- `PATH="/usr/bin:/bin" go test ./internal/cli/... 2>&1 | grep FAIL` → sin FAIL (solo SKIP)
- `PATH="/usr/bin:/bin" go test ./internal/cli/... 2>&1 | grep -c skip` → ≥ 10 skips
- Con rootline real: `go test ./internal/cli/... 2>&1 | grep FAIL` → sin FAIL

## Bloqueadores

- T003 (fix fake describe primero para minimizar el alcance de los skips)
