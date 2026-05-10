---
estado: Completed
tipo: task
---
# T007: Add headless roadmap skill regression scenarios

**Outcome**: [Stabilize roadmap skill routing and materialization UX](README.md)

[[blocked_by:./T001-route-empty-roadmap-invocation-to-pending.md]]
[[blocked_by:./T002-make-roadmap-plan-materialization-token-light.md]]
[[blocked_by:./T003-align-materialize-batch-and-target-apply-docs.md]]
[[blocked_by:./T004-document-and-test-materialize-postcheck-recovery.md]]
[[blocked_by:./T005-clarify-materialize-schema-source-and-cli-exposure.md]]

## Preserva

- Headless tests do not modify files or commit changes unless the scenario explicitly allows it.
- Evidence is structured enough to avoid fragile phrase matching.

## Contexto

Existing headless verification caught important guard behavior but also failed on fragile text expectations, and it does not cover empty invocation routing.

## Alcance

**In**:
1. Add or document a no-argument /roadmap scenario that proves roadmapctl pending is used.
2. Add or document a materialize dry-run scenario that proves concise output and no file modifications.
3. Use a structured final marker such as HEADLESS_RESULT for assertions.

**Out**:
1. No full end-to-end implementation loop execution.
2. No dependence on exact prose beyond structured markers.

## Estado inicial esperado

scripts/verify-roadmap-skill-headless.sh covers loop and materialize preflight but not the new no-argument pending route or token-light dry-run behavior.

## Criterios de Aceptación

- Headless evidence shows no-argument invocation uses roadmapctl pending.
- Headless evidence shows materialize dry-run reports concise status/paths and no modifications.
- The verification script avoids brittle text-only assertions where practical.

## Fuente de verdad

- scripts/verify-roadmap-skill-headless.sh
- .claude/skills/roadmap/SKILL.md
- .claude/skills/roadmap/plan-subcommand.md
