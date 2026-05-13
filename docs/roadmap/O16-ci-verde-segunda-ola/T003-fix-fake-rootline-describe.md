---
estado: Completed
tipo: task
---
# T003: Fix fake rootline describe — retornar formato rootline/describe completo

## Descripción

**RC2 — lint tests fallan en `ci/Test` con schema errors**

La función `fakeRootline` en `internal/cli/golden_test.go` maneja `case "describe"` 
retornando:
```json
{"type":"enum","values":["Pending","Specified","In Progress","Completed","Blocked","On Hold","Obsolete"]}
```

Esto es un fragmento del formato antiguo. `lint.CheckSchemaCompatibility(describe.Decoded)`
busca `describe["schema"]["estado"]`, `describe["schema"]["tipo"]`, y
`describe["links"]["rules"]["blocked_by"]`. Con el fake actual, los tres están ausentes
→ 3 diagnostics de error en CADA test de lint → todos los lint tests fallan con
`exit = 1, want 0`:
- `lint_valid`, `lint_missing_table_row`, `lint_stale_table_row`,
  `lint_missing_task_sections`

El formato correcto de `rootline describe` (verificado con rootline 0.9.87) tiene:
```json
{
  "version": 1,
  "kind": "rootline/describe",
  "schema": {
    "estado": {"type":"enum","required":true,"values":[...]},
    "tipo":   {"type":"enum","required":true,"values":["outcome","task"]}
  },
  "links": {"rules": {"blocked_by": {"target":"..."}}},
  "validate": []
}
```

Fix: reemplazar la respuesta del `case "describe"` con el envelope completo que
satisface tanto `CheckSchemaCompatibility` como `CheckOutcomeSchemaCompatibility`.

## Criterios de Aceptación

- `PATH="/usr/bin:/bin" go test ./internal/cli/... -run "TestCheckGoldenJSONFixtures/lint_valid"` pasa
- `PATH="/usr/bin:/bin" go test ./internal/cli/... -run "TestCheckGoldenJSONFixtures/lint_"` muestra 0 FAIL por schema errors
- `CheckOutcomeSchemaCompatibility` no genera falsos positivos con el nuevo describe (validate: [] → no `estado.non_empty` rule)
