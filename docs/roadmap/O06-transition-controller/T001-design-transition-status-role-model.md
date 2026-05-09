---
estado: Pending
tipo: task
---
# T001: Diseñar modelo de roles y transiciones

**Outcome**: [O06 Controlador de transiciones](README.md)
**Contribuye a**: CE1, CE2, CE3

[[blocked_by:../O02-post-mvp-foundations/T003-validate-operational-status-roles-separately.md]]
[[blocked_by:../O04-readonly-roadmap-state/T001-normalize-rootline-tree-query-graph-data.md]]

## Preserva

- INV1: Los valores concretos de estado vienen de config/schema, no hardcodes.
  - Verificar: fixtures con labels custom.

## Contexto

La skill hoy cambia estados con `rootline set` en prose. `roadmapctl` debe decidir cuándo una transición está permitida y qué valor concreto usar según roles operacionales.

## Alcance

**In**:
1. Definir acciones `can-start`, `can-complete`, `start`, `complete`, `set-status`.
2. Definir status roles requeridos: pending, specified, in-progress, completed, blocked, obsolete.
3. Definir reglas dependency-satisfied usando `done_statuses`.
4. Definir output JSON de allowed/blocked reasons/changes.

**Out**:
- Ejecutar ACs de una task.
- Commits/PRs.

## Estado inicial esperado

- Read model y roles config validados existen.

## Criterios de Aceptación

- Documento/contrato de transiciones queda claro.
- Se especifica cómo se comporta con `On Hold` y estados no-role válidos.
- Se identifican diagnostics `RMC_TRANSITION_*`.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
- `internal/config/*`
- `internal/roadmap/*`
