---
estado: Completed
tipo: task
---
# T004: Cortar plan/materialization docs hacia materialize

**Outcome**: [O08 Cutover de skill](README.md)
**Contribuye a**: CE2

[[blocked_by:../O07-materialization-controller/T007-cutover-plan-materialization-in-skill.md]]

## Preserva

- INV1: No fallback `*-tasks.md`.
  - Verificar: docs y headless materialization scenario.

## Contexto

Después de O07, `plan-subcommand.md` debe dejar de describir rootline new/numbering/table edits manuales como fuente primaria.

## Alcance

**In**:
1. Revisar `plan-subcommand.md`.
2. Documentar flujo conceptual plan -> approval -> structured input -> `roadmapctl materialize`.
3. Mantener guard preflight/postcheck.
4. Ejecutar Pi headless materialization preflight.

**Out**:
- Cambiar materialize implementation.
- Ejecutar apply real en tests headless.

## Estado inicial esperado

- O07/T007 completado.

## Criterios de Aceptación

- Skill no duplica numbering/materialization determinístico.
- Headless scenario no modifica archivos y muestra guard correcto.

## Fuente de verdad

- `.claude/skills/roadmap/plan-subcommand.md`
- `.claude/skills/roadmap/common-logic.md`
