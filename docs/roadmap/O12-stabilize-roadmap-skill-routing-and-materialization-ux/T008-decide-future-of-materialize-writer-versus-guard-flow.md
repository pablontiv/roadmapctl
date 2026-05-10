---
estado: Completed
tipo: task
---
# T008: Decide future of materialize writer versus guard flow

**Outcome**: [Stabilize roadmap skill routing and materialization UX](README.md)

## Preserva

- No implementation behavior changes are bundled into the decision task.
- The decision accounts for historical single-file fallback failures and current token overhead.

## Contexto

Materialize may be over-engineered relative to its original purpose of preventing invalid LLM-created roadmap files, but removing it immediately would re-open known consistency risks.

## Alcance

**In**:
1. Compare keep-writer, guard-only, and direct-write-plus-check alternatives.
2. Use Backscroll evidence, local code constraints, and token usage evidence.
3. Record explicit criteria for revisiting or deprecating materialize.

**Out**:
1. No removal of roadmapctl materialize.
2. No direct skill write reintroduction.

## Estado inicial esperado

Current recommendation is to keep materialize and reduce token overhead, but the product direction is still unsettled.

## Criterios de Aceptación

- A short decision record documents the chosen direction and tradeoffs.
- The decision identifies conditions under which materialize should be simplified or deprecated later.
- No code behavior changes are made solely by this decision task.

## Fuente de verdad

- docs/materialize-plan-schema.md
- docs/roadmap-skill-integration.md
- .claude/skills/roadmap/plan-subcommand.md
- Backscroll session evidence for materialize origin and failures
