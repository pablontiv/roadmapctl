---
estado: Completed
tipo: task
---
# T006: Regenerar golden files desactualizados

## Descripción

**RC8 — golden mismatches residuales**

Después de los fixes de T003 (fake describe) y T004 (skip helper), algunos golden
files pueden quedar desactualizados si el formato de output cambió en los commits
recientes (`feat(cli): replace context with bootstrap`, `fix(cli): require --apply`).

Los golden files afectados incluyen outputs de `pending`, `next`, `decision`, `lint`
y `transition`. Algunos mismatches en `ci/Test` (Ubuntu) eran por fake rootline
vacío, pero otros pueden ser cambios reales de formato.

Proceso:
1. Ejecutar `go test ./internal/cli/... -run TestCheckGoldenJSONFixtures` con rootline real
2. Para cada subtest con `golden mismatch`, comparar want vs got manualmente
3. Si la diferencia es un cambio de formato legítimo (no un bug), actualizar el golden file
4. Verificar que no hay regresiones en ningún test tras la actualización
5. Confirmar que los golden files están en LF (post `.gitattributes` de T001)

## Criterios de Aceptación

- `go test ./internal/cli/... 2>&1 | grep FAIL` → sin FAIL con rootline real
- `go test ./internal/cli/... 2>&1 | grep FAIL` con `PATH="/usr/bin:/bin"` → sin FAIL (solo SKIP para los tests de T004)
- `git diff --check testdata/golden/` → sin CRLF warnings
- `grep -r $'\r' testdata/golden/` → sin resultados

## Bloqueadores

- T003 (fake rootline describe correcto antes de regenerar)
- T004 (skip helper para saber cuáles tests son válidos sin rootline)
