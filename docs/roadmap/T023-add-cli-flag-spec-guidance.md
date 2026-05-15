---
estado: Completed
tipo: task
---
# T023: Add CLI flag spec guidance to task-guide.md

**Contribuye a**: specs de tasks que implementan flags CLI de extracción incluyen contratos de comportamiento suficientes para evitar implementaciones ambiguas

## Preserva

- INV1: `roadmapctl check --strict` verde
  - Verificar: `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict`

## Contexto

Durante el loop de O20, agent-t005 implementó `--field` para bootstrap con dos
errores por spec incompleta:

1. `--field` solo funcionaba cuando también se pasaba `--output json` (la spec
   no decía si debía funcionar de forma independiente)
2. `--field` retornaba strings con comillas JSON (e.g. `"docs/roadmap"` en vez de
   `docs/roadmap`) porque la spec no especificaba el formato de salida para strings

Ambos errores requirieron correcciones manuales post-implementación. La spec de T005
tenía el qué pero no el cómo del contrato de output.

## Alcance

**In**:
1. Agregar una nota en `/home/shared/roadmapctl/.claude/skills/roadmap/task-guide.md`
   sobre contratos para flags CLI de extracción, con los tres puntos:
   - Independencia de otros flags
   - Formato exacto por tipo de dato
   - Comportamiento para no-escalares

**Out**:
- No cambiar código
- No modificar otras secciones del task-guide más allá de la nota

## Estado inicial esperado

`task-guide.md` no tiene ninguna guidance sobre cómo especificar el contrato de
flags CLI de extracción.

## Criterios de Aceptación

- `grep -n "flag.*extrac\|extrac.*flag\|raw string\|escalares\|no-escalar\|CLI" /home/shared/roadmapctl/.claude/skills/roadmap/task-guide.md` retorna al menos una línea
- La nota cubre los 3 puntos: independencia de flags, formato por tipo, no-escalares
- `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict` verde

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/task-guide.md` — archivo a modificar
