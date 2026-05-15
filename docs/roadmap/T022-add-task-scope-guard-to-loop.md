---
estado: Specified
tipo: task
---
# T022: Add task scope guard to loop-subcommand.md

**Contribuye a**: el loop no implementa trabajo fuera del spec de la task actual, aunque sea relacionado o conveniente

## Preserva

- INV1: `roadmapctl check --strict` verde
  - Verificar: `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict`

## Contexto

Durante la ejecución de la task "agregar skill /retrospective" (commit 24bef02),
el agente también añadió el bloque de binary staleness a loop-subcommand.md —
contenido de T020 que ni existía como task en ese momento. Esta es una forma de
scope creep: implementar trabajo conveniente o relacionado fuera del spec de la task
activa.

La instrucción existente en Fase 3 paso 5 dice "Ejecutar exactamente el alcance de
la task", pero no prohíbe explícitamente el trabajo adicional ni indica qué hacer
cuando se detecta.

## Alcance

**In**:
1. Reforzar el paso 5 "Implementar" en `## Fase 3: Loop` de
   `/home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md`
   con la prohibición explícita y el comportamiento alternativo (anotar, no implementar).

**Out**:
- No cambiar otras secciones del loop skill
- No cambiar código

## Estado inicial esperado

Fase 3 paso 5 dice solo "Ejecutar exactamente el alcance de la task. Si hay una
sección `## Especificación Técnica`, seguirla." Sin prohibición de trabajo adicional.

## Criterios de Aceptación

- `grep -n "fuera del spec\|Prohibido añadir\|anotar.*contexto\|contexto.*anotar" /home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md` retorna al menos una línea
- La adición menciona qué hacer cuando se detecta trabajo útil fuera del spec
- `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict` verde

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md` — archivo a modificar
