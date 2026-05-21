---
estado: Completed
tipo: task
---
# T003: Actualizar paso 3 de `task-guide.md` — modelo de materialización por Agent dispatch

**Outcome**: [O30 Materialización paralela](README.md)
**Contribuye a**: Mantener `task-guide.md` consistente con el nuevo modelo de waves descrito en `plan-subcommand.md`

## Contexto

`task-guide.md` tiene en el paso 3 ("Escritura directa tras aprobación"):

> "El skill escribe archivos `.md` directamente usando Write tool."
> "El skill puede escribir archivos en paralelo si los parents (Outcomes) ya existen o fueron creados en el mismo batch."

Estos dos puntos asumen que el skill escribe directamente. Con el nuevo modelo, el skill despacha Agents. El paso 3 debe reflejar eso: el skill coordina, los Agents escriben.

## Alcance

**In**:
1. Reemplazar "El skill escribe archivos `.md` directamente" por "El skill coordina la materialización despachando Agents"
2. Actualizar el punto sobre paralelismo para describir el modelo de waves
3. Mantener los puntos de validación post-escritura (rootline validate, roadmapctl check --strict)

**Out**:
- No modificar el template de Task file (sección `## Template: Task File`)
- No modificar ninguna otra sección de `task-guide.md`

## Estado inicial esperado

- `grep "escribe archivos.*directamente" /home/shared/roadmapctl/.claude/skills/roadmap/task-guide.md` retorna al menos 1 match

## Criterios de Aceptación

- `grep "escribe archivos.*directamente" /home/shared/roadmapctl/.claude/skills/roadmap/task-guide.md` retorna 0 matches
- `grep -i "agent\|wave\|coordin" /home/shared/roadmapctl/.claude/skills/roadmap/task-guide.md` retorna al menos 1 match en el paso 3

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/task-guide.md` — sección `### Paso 3: Escritura directa tras aprobación`
