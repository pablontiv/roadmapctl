---
estado: Pending
tipo: outcome
---
# O08: Cutover de skill y gobernanza

## Objetivo

El skill `/roadmap` queda como adapter conversacional/orquestador, mientras la lógica determinística de bootstrap, state, transitions y materialization vive en `roadmapctl`.

## Criterios de Éxito

- CE1: El skill usa `roadmapctl context` para bootstrap/config.
  - Verificar: Pi headless bootstrap.
- CE2: Pending/next/decision y transitions/materialization delegan en comandos roadmapctl.
  - Verificar: grep + Pi headless.
- CE3: No queda lógica determinística duplicada como source of truth en el skill.
  - Verificar: audit final.

## Invariantes

- INV1: Cambios de skill/guard requieren Pi headless verification.
  - Verificar: evidencia guardada.
- INV2: El skill conserva conversación, aprobación y ejecución de código.
  - Verificar: docs del skill.

## Alcance

**In**:
- Cutovers progresivos del skill.
- Automatización de evidencia Pi.
- Audit final de duplicación.

**Out**:
- Cambiar Rootline.
- Implementar comandos faltantes.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-update-skill-bootstrap-to-use-context.md) | Skill bootstrap usa `roadmapctl context` |
| [T002](T002-cutover-pending-next-decision-docs.md) | Skill pending/next/decision usa comandos nuevos |
| [T003](T003-cutover-transition-loop-docs.md) | Skill loop/status usa transition |
| [T004](T004-cutover-plan-materialization-docs.md) | Skill plan usa materialize |
| [T005](T005-automate-pi-headless-verification-evidence.md) | Automatizar evidencia Pi headless |
| [T006](T006-final-thin-skill-adapter-audit.md) | Audit final para eliminar duplicación determinística |
