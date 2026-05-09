---
estado: Pending
tipo: task
---
# T001: Normalizar datos Rootline tree/query/graph

**Outcome**: [O04 Estado read-only](README.md)
**Contribuye a**: CE1, CE2, CE3

[[blocked_by:../O02-post-mvp-foundations/T005-add-roadmap-domain-model-and-tree-wrapper.md]]
[[blocked_by:../O03-config-context-workspace/T004-implement-context-command.md]]

## Preserva

- INV1: Rootline provee datos genéricos; roadmapctl solo normaliza para comandos de roadmap.
  - Verificar: wrappers CLI JSON.

## Contexto

`pending`, `next` y `decision` comparten inputs: `tree`, `query`, `graph`, schema y roles. Necesitan un modelo común para evitar duplicación.

## Alcance

**In**:
1. Normalizar tasks directas y tasks en outcomes.
2. Asociar status, tipo, path, outcome, dependencies y reverse dependencies.
3. Calcular done/active según config operacional.
4. Manejar cycles/broken links como diagnostics.

**Out**:
- Scoring específico de decision.
- Mutaciones.

## Estado inicial esperado

- Rootline wrappers y context existen.

## Criterios de Aceptación

- Tests cubren direct tasks, outcome tasks, blocked links y no pending tasks.
- Modelo no hardcodea `Completed` ni `Pending` salvo defaults config.

## Fuente de verdad

- `internal/rootlinecli/*`
- `internal/roadmap/*`
- `.claude/skills/roadmap/pending-subcommand.md`
- `.claude/skills/roadmap/decision-tree-subcommand.md`
