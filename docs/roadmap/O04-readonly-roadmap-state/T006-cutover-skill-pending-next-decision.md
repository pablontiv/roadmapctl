---
estado: Pending
tipo: task
---
# T006: Cortar skill pending/next/decision hacia roadmapctl

**Outcome**: [O04 Estado read-only](README.md)
**Contribuye a**: CE4

[[blocked_by:./T005-add-read-api-fixtures-and-goldens.md]]

## Preserva

- INV1: Cambios de skill/guard requieren Pi headless verification.
  - Verificar: comandos documentados en `docs/roadmap-skill-integration.md`.

## Contexto

Cuando `pending`, `next` y `decision` estén estables, el skill debe dejar de llamar Rootline directamente para estas decisiones determinísticas.

## Alcance

**In**:
1. Actualizar `pending-subcommand.md` para usar `roadmapctl pending`.
2. Actualizar `decision-tree-subcommand.md` para usar `roadmapctl decision`/`next`.
3. Actualizar routing/docs si el empty decision tree cambia.
4. Ejecutar sync del skill y Pi headless verification.

**Out**:
- Cambiar loop transitions.
- Cambiar materialización.

## Estado inicial esperado

- Comandos read-only estables y testeados.

## Criterios de Aceptación

- El skill no duplica lógica read-only ya cubierta por roadmapctl.
- Pi headless muestra comandos roadmapctl requeridos/usados.
- `scripts/sync-roadmap-skill.sh --check` pasa.

## Fuente de verdad

- `.claude/skills/roadmap/pending-subcommand.md`
- `.claude/skills/roadmap/decision-tree-subcommand.md`
- `docs/roadmap-skill-integration.md`
