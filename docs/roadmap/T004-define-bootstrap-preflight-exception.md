---
estado: Completed
tipo: task
---
# T004: Define bootstrap preflight exception

## Preserva

- Existing strict doctor/check guards remain mandatory for normal writes.
- Bootstrap does not become an auto-fix fallback for invalid roadmaps.

## Contexto

Normal preflight requires doctor/check before writes, but missing-root bootstrap cannot satisfy normal check before the root exists.

## Alcance

**In**:
1. Choose and document the canonical bootstrap path: bootstrap init, materialize bootstrap, or both with distinct rules.
2. Update skill and CLI integration docs to describe the bootstrap exception and postcheck requirements.
3. Add or adjust tests/docs if the bootstrap contract changes.

**Out**:
1. Do not weaken strict preflight for existing roadmaps.
2. Do not silently create roadmap roots outside explicit bootstrap flows.

## Estado inicial esperado

Skill preflight and materialize bootstrap docs are ambiguous for missing roadmap roots.

## Criterios de Aceptación

- Docs identify exactly which command/flow is allowed to create a missing roadmap root.
- The allowed bootstrap flow has explicit guard and postcheck steps.
- Normal materialization still stops when doctor/check fail outside bootstrap.

## Fuente de verdad

- /tmp/pi-subagents-uid-1000/chain-runs/07512486/handoff-final.md
- internal/cli/bootstrap.go
- internal/materialize/dryrun.go
- .claude/skills/roadmap/plan-subcommand.md
- docs/roadmap-skill-integration.md
