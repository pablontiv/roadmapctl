---
estado: Completed
tipo: task
---
# T027: Loop — forzar bypass de caché en verificación de tests post-edit

**Contribuye a**: La verificación de ACs de cobertura/tests no produce falsos positivos cuando el test runner retorna un resultado cacheado sobre código recién modificado.

## Contexto

Go test (y otros runners como Jest, pytest) cachean resultados. Si un archivo fue modificado pero el binario del test no se recompiló aún, el runner puede retornar `(cached)` con el resultado anterior. En el loop, si el AC verifica `go test -cover ./pkg/...` inmediatamente después de editar `pkg/foo.go`, el resultado cacheado puede mostrar cobertura ≥ umbral aunque el código nuevo no esté cubierto.

Ejemplo real: T002 del loop picokit — el subagente reportó `coverage: 95.7% (cached)` sobre código que en realidad no había sido modificado, y el loop lo tomó como AC pasado.

## Alcance

**In**:
1. En `loop-subcommand.md`, Fase 3 paso 6 "Verificar ACs e invariantes", agregar nota:
   - Si el AC usa un test runner con caché (go test, jest, pytest, cargo test, etc.) y la task modificó archivos bajo test, **forzar ejecución sin caché**: `-count=1` en Go, `--no-cache` en Jest/cargo, `--cache-clear` en pytest, o equivalente del runner.
   - Un resultado marcado `(cached)` o `from cache` sobre un archivo recién editado es un falso positivo — no cuenta como AC verificado.
2. La nota puede vivir como sub-regla dentro de la regla de re-verificación de T026, o como punto separado en el mismo paso.

**Out**:
- No modificar el comportamiento de dispatch ni de integración.
- No prescribir un runner específico — la regla aplica a cualquier lenguaje.

## Criterios de Aceptación

- `loop-subcommand.md` paso 6 menciona explícitamente el riesgo de resultados cacheados sobre código editado.
- Se indica al menos un ejemplo de flag de bypass (`-count=1` para Go) con nota de que otros runners tienen equivalentes.
- La regla establece que `(cached)` sobre archivo recién modificado no es evidencia válida de AC pasado.
