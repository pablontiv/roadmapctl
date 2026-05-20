---
estado: Specified
tipo: task
---
# T033: Loop — cómputo determinístico de `is_last_in_scope` antes de /integrate

**Contribuye a**: El loop nunca pasa un valor "intuido" de `is_last_in_scope` a /integrate; siempre lo computa recalculando la cola tras la transición.

## Contexto

En la wave 4 de O25 invoqué /integrate para T002, T003 y T004 con `is_last_in_scope=false` para los tres, aunque T004 era de hecho la última task de O25. Con `pr_mode=false` el flag no afecta merge, pero sí gobierna el "outcome close check" de Fase 3 paso 7 — que no corrió cuando debía.

El flag es estructural: cualquier loop debería computarlo determinísticamente, no depender de que el modelo recuerde el conteo de tasks del scope.

## Alcance

**In**:
1. En `.claude/skills/roadmap/loop-subcommand.md`, Fase 3 paso 9 (invocación de /integrate):
   - Tras `roadmapctl transition complete --apply` para una task, ejecutar `roadmapctl next --repo <repo> --output json` y filtrar `ready[]` por scope (Outcome path o `direct-tasks`).
   - Si `ready[]` filtrado por ese scope queda vacío, pasar `is_last_in_scope=true` a /integrate.
   - En cualquier otro caso, `is_last_in_scope=false`.
2. Prohibir explícitamente "intuir" el valor: el skill debe recalcular, no inferir desde el contador interno de la sesión.

**Out**:
- No cambiar el comportamiento de /integrate cuando recibe `is_last_in_scope=true` (sigue siendo `pr_mode`-dependiente).
- No agregar nuevas señales al output de /integrate.
- No tocar Fase 3 paso 7 (Outcome close check) — el cómputo correcto del flag es lo que habilita que ese paso corra cuando corresponde.

## Criterios de Aceptación

- Fase 3 paso 9 documenta el cómputo determinístico de `is_last_in_scope` usando `roadmapctl next` post-complete.
- El skill prohíbe explícitamente pasar un valor sin recálculo.
- La regla define cómo filtrar `ready[]` por scope (Outcome path prefix, o vacío == `direct-tasks`).

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
