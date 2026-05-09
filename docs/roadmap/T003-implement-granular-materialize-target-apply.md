---
estado: Completed
tipo: task
---
# T003: Implement granular materialize target apply

## Preserva

- Existing batch materialize --apply behavior remains compatible unless explicitly deprecated.
- Dry-run remains deterministic and no-write.
- roadmapctl remains the owner of numbering, generated content, dependency links, validation, and postchecks.

## Contexto

The /roadmap skill now requires per-file materialization, but the CLI currently applies all changes in a plan-level batch.

## Alcance

**In**:
1. Design and implement a frozen dry-run/change-set or equivalent safe target-apply flow.
2. Add CLI tests and materialize package tests for single-target application and invalid targets.
3. Update CLI contract and skill integration docs after behavior exists.

**Out**:
1. Do not implement prompt-side raw markdown writes as the primary solution.
2. Do not break existing dry-run or batch apply contract without explicit approval.

## Estado inicial esperado

roadmapctl materialize supports --plan with --dry-run or --apply only; Apply writes all computed changes.

## Criterios de Aceptación

- A dry-run/change-set can be used to apply exactly one canonical roadmap file target.
- Applying one target does not create or rewrite sibling roadmap files.
- Unknown, empty, duplicate, or non-file targets fail before writing.
- go test ./internal/materialize ./internal/cli and go test ./... pass.

## Fuente de verdad

- /tmp/pi-subagents-uid-1000/chain-runs/07512486/handoff-final.md
- internal/cli/materialize.go
- internal/materialize/dryrun.go
- docs/roadmap-skill-integration.md
- docs/cli-contract.md
