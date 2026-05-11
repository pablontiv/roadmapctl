---
estado: Completed
tipo: task
---
# T002: Redefinir el contrato del skill /roadmap plan

**Outcome**: [Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill](README.md)

[[blocked_by:./T001-record-responsibility-separation-decision.md]]

## Preserva

- No se escribe ningún archivo antes de aprobación humana explícita.

## Contexto

El skill debe dejar de serializar prose semántica hacia roadmapctl materialize y debe encargarse de generar Markdown final aprobado.

## Alcance

**In**:
1. Actualizar plan-subcommand.md y referencias comunes.
2. Definir templates de Outcome y Task sin ## Tasks persistido ni AC en Outcome.
3. Exigir escritura paralela de archivos aprobados cuando sea seguro.
4. Definir re-pregunta solo ante divergencia de paths/cantidad/destino.

**Out**:
1. No mover lógica de path planning al skill.
2. No reintroducir fallback *-tasks.md.

## Estado inicial esperado

/roadmap plan actualmente delega la escritura completa a roadmapctl materialize y prohíbe write directo.

## Criterios de Aceptación

- El skill muestra una única propuesta de Outcome/Tasks con AC por Task.
- El skill pide aprobación antes de materializar.
- El skill no muestra ni produce JSON semántico para que roadmapctl renderice Markdown.
- El skill escribe archivos aprobados en paralelo cuando los parents ya existen o fueron creados previamente.
- El skill valida con roadmapctl/Rootline después de escribir.

## Fuente de verdad

- .claude/skills/roadmap/SKILL.md
- .claude/skills/roadmap/plan-subcommand.md
- .claude/skills/roadmap/task-guide.md
- .claude/skills/roadmap/outcome-guide.md
