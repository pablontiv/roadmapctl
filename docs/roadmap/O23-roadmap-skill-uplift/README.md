---
tipo: outcome
---
# Roadmap skill: tool access, observability docs, README outputs

Tres mejoras al skill `/roadmap` y al README del repo motivadas por un incidente del session previo (Backscroll session `9d8a66cd-01eb-4415-8e33-04f3b9cec020.jsonl`) donde el usuario pidió explícitamente "usa el tool monitor para monitorear cualquier proceso" durante `/roadmap loop`, y el modelo nunca pudo invocarlo.

Root cause confirmado: el bloque `allowed-tools:` del frontmatter de `SKILL.md` enumera taxativamente las herramientas permitidas y omite `Monitor` (y `ScheduleWakeup`). La permisología del harness bloquea cualquier invocación a tools fuera de la lista, independientemente de instrucciones narrativas o pedido del usuario.

Tres ejes de mejora cubiertos por las tasks de este Outcome:

1. **Acceso a la tool roster completa.** Eliminar el whitelist en favor de la convención mayoritaria del ecosistema (skills sin `allowed-tools:` heredan todas las tools de la sesión), eliminando el costo de mantenimiento por cada herramienta nueva que aparezca.
2. **Guidance operativa.** Agregar una sección a `loop-subcommand.md` que dice *cuándo* usar `Monitor` y `ScheduleWakeup` — sin esta guía, el modelo cae a `bash sleep`-polling por defecto incluso teniendo acceso.
3. **Documentación con outputs reales.** El README del repo enumera comandos sin mostrar lo que devuelven; los agentes que lo leen para descubrir contratos no pueden anticipar formato. Adicionalmente todos los ejemplos usan `--repo .` redundante (el default del flag global ya es `.`), entrenando ruido.

Plan que generó este Outcome: `~/.claude/plans/logical-twirling-acorn.md`. El trabajo de enforcement vía PreToolUse hook + regla CLAUDE.md (originalmente Steps 3a/3b/3c del plan) queda diferido — explícitamente puesto on-hold por el usuario durante Fase 2 de `/roadmap plan`. Se reabordará en un Outcome futuro si se mantiene el problema de pivot post-aprobación de plan-mode.
