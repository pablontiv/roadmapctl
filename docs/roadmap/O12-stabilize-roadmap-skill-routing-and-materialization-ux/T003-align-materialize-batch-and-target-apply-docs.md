---
estado: Completed
tipo: task
---
# T003: Align materialize batch and target apply docs

**Outcome**: [Stabilize roadmap skill routing and materialization UX](README.md)

## Preserva

- Batch apply remains roadmapctl-owned and guarded by postcheck.
- Target apply remains available for recovery, troubleshooting, or explicit one-file approval.

## Contexto

docs/roadmap-skill-integration.md still contains older instructions that conflict with the current batch apply policy in the skill and materialize schema docs.

## Alcance

**In**:
1. Update docs/roadmap-skill-integration.md to match plan-subcommand.md and common-logic.md.
2. Clarify when --changes --target is appropriate.
3. Ensure docs consistently prohibit manual multi-file writes by the skill.

**Out**:
1. No new apply mode implementation.
2. No changes to dependency semantics.

## Estado inicial esperado

Integration docs contain stale per-target-only language while skill docs permit roadmapctl-owned batch apply.

## Criterios de Aceptación

- No integration doc text says the skill must never use batch apply for multi-file plans.
- All materialize docs agree that roadmapctl-owned batch apply is the normal path.
- All materialize docs agree that target apply is recovery or explicit one-file approval only.

## Fuente de verdad

- docs/roadmap-skill-integration.md
- .claude/skills/roadmap/plan-subcommand.md
- .claude/skills/roadmap/common-logic.md
- docs/materialize-plan-schema.md
