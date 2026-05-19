---
estado: Specified
tipo: task
---
# T026: Loop — re-verificar ACs de código modificado por subagente paralelo

**Contribuye a**: El loop no acepta como válido el resumen de un subagente cuando el AC involucra archivos que ese subagente debió modificar.

## Contexto

Cuando `parallel == true`, cada task de una wave es ejecutada por un `Agent` call independiente. El subagente reporta en su resumen que los ACs pasaron, pero sus edits pueden no haber persistido (por error del tool, condición de carrera u otro motivo). El loop actualmente confía en el resumen del agente como evidencia suficiente para invocar `roadmapctl transition complete --apply`.

El resultado es que `transition complete` se ejecuta sobre una task cuyo código no fue modificado, y el error se detecta recién cuando CI falla en el siguiente push.

## Alcance

**In**:
1. En `loop-subcommand.md`, Fase 3 paso 6 "Verificar ACs e invariantes", agregar regla explícita:
   - Si la task fue ejecutada por un subagente paralelo Y el AC involucra verificar estado de un archivo modificado (grep, contenido, compilación, tests), **re-ejecutar el comando de verificación directamente** en el loop principal antes de llamar `roadmapctl transition complete --apply`.
   - Distinguir entre: (a) ACs verificables solo por el agente (e.g. decisiones de diseño), y (b) ACs verificables con un comando directo (grep, go build, go test, etc.) — solo los segundos requieren re-verificación obligatoria.
2. Agregar ejemplo canónico: "si la task modificó `pkg/foo.go` y el AC es `grep -r 'OldSymbol' pkg/foo/` retorna vacío, correr ese grep directamente antes de `complete`."

**Out**:
- No modificar la lógica de dispatch de agents paralelos.
- No modificar el formato de INTEGRATE_RESULT.

## Criterios de Aceptación

- `loop-subcommand.md` paso 6 contiene una regla nombrada o destacada sobre re-verificación post-subagente.
- La regla diferencia explícitamente ACs verificables con comando directo de ACs no verificables externamente.
- La regla indica que la re-verificación directa es requisito previo a `roadmapctl transition complete --apply`, no opcional.
