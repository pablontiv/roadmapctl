---
estado: Completed
tipo: task
---
# T001: Actualizar bootstrap del skill para usar `roadmapctl context`

**Outcome**: [O08 Cutover de skill](README.md)
**Contribuye a**: CE1

[[blocked_by:../O03-config-context-workspace/T004-implement-context-command.md]]

## Preserva

- INV1: Conceptual planning puede seguir sin writes si faltan guards.
  - Verificar: docs del skill.

## Contexto

El skill hoy calcula mode, roadmap-root, config y helpers en prose. `roadmapctl context` debe ser la fuente determinística.

## Alcance

**In**:
1. Actualizar `SKILL.md` bootstrap para preferir `roadmapctl context`.
2. Documentar `.roadmapctl.toml` como config preferida.
3. Mantener fallback conceptual/no-write.
4. Ejecutar sync + Pi headless.

**Out**:
- Cutover pending/decision.
- Materialization.

## Estado inicial esperado

- `roadmapctl context` implementado y testeado.

## Criterios de Aceptación

- Headless bootstrap reporta comando `roadmapctl context` o razón de fallback.
- No se duplica cálculo de helpers si context está disponible.
- `scripts/sync-roadmap-skill.sh --check` pasa.

## Fuente de verdad

- `.claude/skills/roadmap/SKILL.md`
- `docs/roadmap-skill-integration.md`
