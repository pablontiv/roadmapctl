---
estado: Specified
tipo: task
---
# T005: Clarify materialize schema source and CLI exposure

**Outcome**: [Stabilize roadmap skill routing and materialization UX](README.md)

## Preserva

- roadmapctl materialize continues to reject free-form prose input.
- The canonical schema remains versioned and documented.

## Contexto

The skill references docs/materialize-plan-schema.md without clarifying whether that path is skill-relative, repo-relative, or installed elsewhere.

## Alcance

**In**:
1. Clarify schema source in the skill and integration docs.
2. Document a deterministic lookup path for agents in normal repos and in the roadmapctl repo.
3. Evaluate and, if accepted, specify a future roadmapctl materialize schema --output json command.

**Out**:
1. No schema version change.
2. No expansion of materialize plan fields.

## Estado inicial esperado

Agents have had to search local docs and run materialize --help to avoid inventing the JSON shape.

## Criterios de Aceptación

- Skill docs explicitly identify the canonical schema source and how an agent should access it.
- No instruction implies docs/materialize-plan-schema.md is relative to the installed skill directory unless it actually is.
- The plan records whether CLI schema exposure is required now or deferred.

## Fuente de verdad

- .claude/skills/roadmap/plan-subcommand.md
- docs/materialize-plan-schema.md
- docs/cli-contract.md
- internal/cli/materialize.go
