---
estado: Completed
tipo: task
---
# T006: Tighten materialize dependency path validation

[[blocked_by:./T003-implement-granular-materialize-target-apply.md]]

## Preserva

- blocked_by links remain explicit relative paths.
- Bare basename dependencies remain invalid.

## Contexto

The schema says dependency.path must point to an existing or concurrently planned task, but implementation currently validates shape more than target existence.

## Alcance

**In**:
1. Decide whether dependency.path is strict at dry-run time or deferred to postcheck.
2. Implement strict validation or update docs/tests to reflect actual behavior.
3. Add tests for unresolved explicit dependency paths.

**Out**:
1. Do not allow bare target dependencies.
2. Do not change ref resolution semantics beyond the approved contract.

## Estado inicial esperado

dependency.path validation rejects invalid shape but may not verify existing or concurrently planned targets.

## Criterios de Aceptación

- The schema, implementation, and tests agree on dependency.path behavior.
- Invalid unresolved paths produce deterministic diagnostics if strict validation is chosen.
- go test ./internal/materialize passes.

## Fuente de verdad

- docs/materialize-plan-schema.md
- internal/materialize/dryrun.go
- internal/materialize/validation_test.go
