---
estado: Pending
tipo: task
---
# T009: Refresh CLI contract command summary

[[blocked_by:./T003-implement-granular-materialize-target-apply.md]]
[[blocked_by:./T005-align-transition-json-contract.md]]
[[blocked_by:./T006-tighten-materialize-dependency-path-validation.md]]
[[blocked_by:./T007-align-materialized-task-initial-status.md]]

## Preserva

- CLI contract remains the authoritative public interface reference.
- No behavior changes are made solely for documentation cleanup.

## Contexto

The command list and summary block are partially stale after adding context, pending, next, decision, transition, and materialize.

## Alcance

**In**:
1. Update docs/cli-contract.md command summary and examples.
2. Ensure materialize, transition, pending, next, decision, context, doctor, check, lint, and bootstrap are represented consistently.
3. Cross-check docs against command registration.

**Out**:
1. Do not change CLI behavior as part of this docs-only cleanup unless a mismatch is discovered and separately approved.

## Estado inicial esperado

docs/cli-contract.md has a stale summary block despite commands being implemented.

## Criterios de Aceptación

- docs/cli-contract.md summary includes all implemented roadmapctl commands.
- Examples use implemented flags and JSON field names.
- go test ./... remains green after docs updates.

## Fuente de verdad

- docs/cli-contract.md
- internal/cli/cli.go
- cmd/roadmapctl/main.go
