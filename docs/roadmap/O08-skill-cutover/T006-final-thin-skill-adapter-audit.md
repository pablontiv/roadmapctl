---
estado: Pending
tipo: task
---
# T006: Audit final de skill como adapter fino

**Outcome**: [O08 Cutover de skill](README.md)
**Contribuye a**: CE3

[[blocked_by:./T001-update-skill-bootstrap-to-use-context.md]]
[[blocked_by:./T002-cutover-pending-next-decision-docs.md]]
[[blocked_by:./T003-cutover-transition-loop-docs.md]]
[[blocked_by:./T004-cutover-plan-materialization-docs.md]]
[[blocked_by:./T005-automate-pi-headless-verification-evidence.md]]

## Preserva

- INV1: El skill conserva su rol conversacional y de aprobación.
  - Verificar: autonomous-mode y routing siguen claros.

## Contexto

Después de los cutovers, hay que auditar que la lógica determinística no quede duplicada en Markdown y CLI.

## Alcance

**In**:
1. Revisar todos los archivos `.claude/skills/roadmap/*.md`.
2. Identificar y remover/degradar duplicación determinística reemplazada por roadmapctl.
3. Mantener referencias a comandos roadmapctl y reglas de escalación.
4. Ejecutar sync/check y Pi headless.

**Out**:
- Cambiar comandos roadmapctl.
- Quitar conceptual decomposition.

## Estado inicial esperado

- Cutovers previos completados.

## Criterios de Aceptación

- Audit documenta qué lógica queda en skill y por qué.
- No hay instrucciones primarias conflictivas con roadmapctl.
- Pi verification pasa.

## Fuente de verdad

- `.claude/skills/roadmap/*.md`
- `docs/roadmap-skill-integration.md`
