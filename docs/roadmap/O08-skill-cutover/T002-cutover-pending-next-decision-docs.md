---
estado: Completed
tipo: task
---
# T002: Cortar pending/next/decision docs hacia roadmapctl

**Outcome**: [O08 Cutover de skill](README.md)
**Contribuye a**: CE2

[[blocked_by:../O04-readonly-roadmap-state/T006-cutover-skill-pending-next-decision.md]]

## Preserva

- INV1: El skill presenta resultados pero no recalcula estado determinístico.
  - Verificar: docs actualizados.

## Contexto

Este task consolida el cutover read-only del skill después de que O04 lo implemente y verifique.

## Alcance

**In**:
1. Revisar `pending-subcommand.md` y `decision-tree-subcommand.md`.
2. Eliminar instrucciones Rootline directas ya reemplazadas.
3. Mantener UX humana del skill.
4. Pi headless para pending/empty decision si aplica.

**Out**:
- Transition/materialize cutovers.

## Estado inicial esperado

- O04/T006 completado.

## Criterios de Aceptación

- No quedan dos fuentes de verdad para pending/next/decision.
- Pi verification o tests equivalentes muestran funcionamiento.

## Fuente de verdad

- `.claude/skills/roadmap/pending-subcommand.md`
- `.claude/skills/roadmap/decision-tree-subcommand.md`
