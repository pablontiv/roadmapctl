---
estado: Completed
tipo: task
---
# T005: Align transition JSON contract

## Preserva

- Existing transition dry-run/apply behavior remains compatible unless an explicit contract migration is approved.
- Skill loop consumes stable field names.

## Contexto

Docs mention blockers/from/to while implementation and skill use blocking_dependencies and before/after.

## Alcance

**In**:
1. Audit transition JSON field names in implementation, tests, docs, and skill references.
2. Update docs or add backward-compatible aliases with tests.
3. Ensure examples match actual output.

**Out**:
1. Do not change transition state machine semantics.
2. Do not broaden transition to outcome README files.

## Estado inicial esperado

transition-controller docs and implementation use different names for dependency and change fields.

## Criterios de Aceptación

- docs/transition-controller.md examples match implementation output or documented aliases.
- CLI contract examples are consistent with tests.
- go test ./internal/cli passes.

## Fuente de verdad

- docs/transition-controller.md
- internal/cli/transition.go
- internal/cli/transition_test.go
- .claude/skills/roadmap/loop-subcommand.md
