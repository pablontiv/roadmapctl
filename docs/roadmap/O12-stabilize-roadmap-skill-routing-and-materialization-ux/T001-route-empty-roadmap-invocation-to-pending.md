---
estado: Completed
tipo: task
---
# T001: Route empty roadmap invocation to pending

**Outcome**: [Stabilize roadmap skill routing and materialization UX](README.md)

## Preserva

- Explicit pending, decision, next, loop, and plan routes remain available.
- The bootstrap through roadmapctl context remains mandatory before dispatch.

## Contexto

The current skill maps empty arguments to decision-tree, but the desired default behavior is to show pending work.

## Alcance

**In**:
1. Update the routing table and dispatch rules in SKILL.md.
2. Document pending as the default no-argument view.
3. Keep decision-tree behavior available for explicit prioritization requests.

**Out**:
1. No Go CLI behavior changes unless a validation test requires it.
2. No changes to roadmap task execution semantics.

## Estado inicial esperado

SKILL.md routes empty arguments to decision-tree-subcommand.md.

## Criterios de Aceptación

- SKILL.md table and dispatch rules route empty arguments to pending-subcommand.md.
- decision or next wording still routes to decision-tree-subcommand.md.
- A headless Pi check demonstrates that the no-argument scenario runs roadmapctl pending, not roadmapctl decision.

## Fuente de verdad

- .claude/skills/roadmap/SKILL.md
- .claude/skills/roadmap/pending-subcommand.md
- .claude/skills/roadmap/decision-tree-subcommand.md
