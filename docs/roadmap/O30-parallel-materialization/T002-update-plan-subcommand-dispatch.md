---
estado: Completed
tipo: task
---
# T002: Actualizar sección 3.3 de `plan-subcommand.md` — despacho paralelo de Agents

**Outcome**: [O30 Materialización paralela](README.md)
**Contribuye a**: Describir el mecanismo concreto de dispatch paralelo que el skill debe ejecutar en fase 3

## Preserva

- INV1: Wave 0 (Outcomes) debe completarse antes de wave 1+ (Tasks) — los Tasks referencian el README del Outcome
  - Verificar: el texto de sección 3.3 menciona explícitamente que los Outcomes se crean antes de despachar Tasks

## Contexto

La sección 3.3 de `plan-subcommand.md` actualmente dice:

> "Crear directorios padre si aplican, luego escribir con Write tool en paralelo"

"Write tool en paralelo" se refiere a múltiples Write calls del skill mismo, no Agents separados. El cambio reemplaza esto por un modelo de dispatch explícito:

1. Wave 0: el coordinador escribe los Outcome README.md directamente (o despacha 1 Agent si son múltiples Outcomes)
2. Wave 1+: despacha N Agents en paralelo, cada uno con un subset de Task files
3. Cada Agent recibe en su prompt: lista de (path, contenido completo) — el contenido está pre-decidido, no se recalcula
4. Cada Agent hace 1 Write call por archivo
5. El coordinador espera a todos y ejecuta el postcheck global

La sección también debe incluir el template de prompt para los Agents de materialización.

Sección a modificar: `## Fase 3: Materialización > **3.3 Escritura en paralelo**` en `/home/shared/roadmapctl/.claude/skills/roadmap/plan-subcommand.md`.

## Alcance

**In**:
1. Reemplazar "escribir con Write tool en paralelo" por el modelo de waves + Agent dispatch
2. Incluir template de prompt para Agent de materialización (recibe paths + contenido, escribe, confirma)
3. Especificar criterio de partición de subsets (ej. ~3 archivos por Agent, o 1 Agent por Outcome + sus Tasks)

**Out**:
- No modificar secciones 3.1, 3.2, 3.4, 3.5 de `plan-subcommand.md`
- No modificar `task-guide.md` (eso es T003)
- No modificar `SKILL.md` (eso es T001)

## Estado inicial esperado

- `grep "Write tool en paralelo" /home/shared/roadmapctl/.claude/skills/roadmap/plan-subcommand.md` retorna al menos 1 match

## Criterios de Aceptación

- `grep "Write tool en paralelo" /home/shared/roadmapctl/.claude/skills/roadmap/plan-subcommand.md` retorna 0 matches
- `grep -i "agent\|wave\|subset\|dispatch" /home/shared/roadmapctl/.claude/skills/roadmap/plan-subcommand.md` retorna al menos 1 match en la sección 3.3
- La sección 3.3 contiene un template de prompt explícito para el Agent de materialización
- La sección menciona que Wave 0 (Outcomes) precede a Wave 1 (Tasks)

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/plan-subcommand.md` — sección `**3.3 Escritura en paralelo**`
