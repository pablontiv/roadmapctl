---
estado: Completed
tipo: task
---
# T002: Agregar guidance Monitor/ScheduleWakeup a loop-subcommand.md

**Outcome**: [O23 Roadmap skill: tool access, observability docs, README outputs](README.md)
**Contribuye a**: enseñar al modelo cuándo invocar `Monitor` y `ScheduleWakeup` en lugar de `bash sleep`-polling durante `/roadmap loop`, especialmente cuando hay procesos largos (builds, tests, agentes) o estado externo (CI runs, deploys).

## Preserva

- INV1: La estructura de phases del loop subcommand (Discovery, TodoList, Loop, Teardown) y el handoff a `pr-workflow.md` cuando `pr_mode==true` no cambian.
  - Verificar: `diff` de `loop-subcommand.md` muestra solo adiciones de la nueva sección "Observabilidad de procesos largos"; los headers de Phase 1/2/2.5/3/4 y sus cuerpos quedan idénticos.
- INV2: Las invocaciones a `roadmapctl` (`next`, `pending`, `transition can-start/start/complete`) descritas en el loop subcommand se preservan sin modificación.
  - Verificar: `grep -c 'roadmapctl' loop-subcommand.md` antes y después coincide; los comandos exactos no cambian.

## Contexto

Dropping el whitelist (T001) hace que `Monitor` y `ScheduleWakeup` sean alcanzables desde el skill. Pero el modelo no las usa por defecto: cae a `Bash` con `run_in_background: true` y polling con `sleep`. Esto desperdicia cache (cada cycle de poll es un cache miss después de 5 min), pierde el streaming de eventos línea-por-línea, y no maneja bien estado externo que el harness no notifica (CI runs, deploys remotos, queues).

La nueva sección debe ubicarse antes del cuerpo de Phase 3 (Loop) en `/home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md`, porque Phase 3 es donde se ejecutan los procesos largos (tests/builds/agent dispatches por task).

Reglas que la sección debe codificar:

- **Cuándo usar `Monitor`** (en lugar de `Bash` foreground): para procesos en background cuyo stdout queremos surfacear línea-por-línea. Patrón canónico: launch con `Bash` + `run_in_background: true` teeing a `/tmp/roadmap-<task-id>.log`, luego `Monitor` con `grep -E --line-buffered` filtrando milestones (PASS/FAIL/ERROR/heartbeat). Ejemplo concreto debe usar IDs de task reales (e.g., `Monitor` con `description: "T001 test run"`).
- **Cuándo usar `ScheduleWakeup`** (en lugar de polling con sleep): para esperar estado externo que el harness no notifica — GitHub Actions runs (`gh run watch` bloquea; preferir wakeup + `gh run view --json status`), deploys, remote queues.
- **Prohibición**: no usar `Bash sleep` chained loops para esperar — elegir Monitor (stdout streamable) o ScheduleWakeup (poll interval externo) según el caso.
- **Instrucción directa del usuario**: si el usuario dice "monitorea" / "use monitor" / "watch this", invocar `Monitor` inmediatamente en el siguiente paso de proceso largo — no sustituir silenciosamente por `Bash background + poll`.

La sección debe estar limitada a ~30 líneas e incluir al menos un ejemplo concreto del patrón completo `Bash background + tee log + Monitor con grep`.

## Alcance

**In**:
1. Agregar la sección "Observabilidad de procesos largos" en `/home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md`, ubicada antes del cuerpo descriptivo de Phase 3 (Loop).
2. Incluir las cuatro reglas enumeradas en Contexto.
3. Incluir al menos un bloque de ejemplo con el patrón `Bash background + tee + Monitor grep`.
4. Sincronizar con `./scripts/sync-roadmap-skill.sh --install` y verificar parity con `--check`.

**Out**:
- Modificar el comportamiento descrito en otras phases del loop (Discovery, TodoList, PR mode, Teardown).
- Editar `SKILL.md` (cobertura de T001).
- Agregar guidance Monitor/ScheduleWakeup en otros subcommands (`plan`, `pending`, `decision-tree`) — `loop` es el caso de uso primario.

## Estado inicial esperado

- `loop-subcommand.md` no menciona `Monitor`, `ScheduleWakeup`, ni `run_in_background`.
- T001 no es hard blocker para esta task — la guía narrativa funciona aunque el whitelist exista (solo que el modelo no podría invocar las tools efectivamente). Pueden ejecutarse en paralelo.

## Criterios de Aceptación

- `loop-subcommand.md` contiene una sección titulada (e.g., `## Observabilidad de procesos largos` o equivalente) con las cuatro reglas: cuándo `Monitor`, cuándo `ScheduleWakeup`, prohibición de `bash sleep` loops, prioridad de instrucción directa del usuario.
- La sección contiene al menos un ejemplo concreto del patrón `Bash background + tee + Monitor grep`.
- La sección es ≤30 líneas (medible: contar líneas entre el header de la sección y el siguiente header de mismo nivel).
- `./scripts/sync-roadmap-skill.sh --check` exit 0 después del sync.
- `grep -c '^##' loop-subcommand.md` post-cambio = pre-cambio + 1 (única adición top-level).

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md`
- `/home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` (referencia: orden de phases del loop)
