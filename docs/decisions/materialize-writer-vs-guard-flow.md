# Decision: roadmapctl as guard/policy layer; Pi write owns roadmap materialization

Status: Accepted
Date: 2026-05-11

## Decision

`roadmapctl` evolves from a deterministic writer to a guard, validator, policy layer, path planner, and state query engine. The skill (Pi agent) becomes the sole writer of approved roadmap files after human authorization.

### Semantic boundary

**Outcome README**: does not persist `## Tasks` section. The task table is a computed view generated from child `TXXX-*.md` files. Outcomes are purely container/context documents. `roadmapctl lint` may warn if a stale table exists, but `roadmapctl` never writes it.

**Outcome scope**: Outcomes do not require acceptance criteria in the README. Acceptance criteria live exclusively in child `TXXX-*.md` task files as structured sections. An Outcome README carries narrative, goals, constraints, and motivation; Tasks carry work items and their specific ACs.

**Pi write authority**: After human approval of a roadmap plan (dry-run review), the skill serializes only data—not prose. It sends `roadmapctl materialize --dry-run` for validation and preflight diagnostics. On explicit human approval, the skill then directly writes canonical Outcome README and Task files using approved Markdown content. There is no hidden content generation; the skill writes exactly what the human saw. After writing, `roadmapctl check` validates the written files.

**roadmapctl role**: 
- Guard/policy layer: blocking validation before writes/executions (`roadmapctl doctor`, `roadmapctl check`).
- Path planner: deterministic ID numbering, canonical paths, dependency resolution.
- Validator: structural invariants, schema compliance, graph validation, stale content warnings.
- View/query layer: read-only state queries (`roadmapctl pending`, `roadmapctl next`, `roadmapctl decision`).
- Policy engine: status transition rules and operator muscle memory (`roadmapctl transition`).

`roadmapctl materialize` may be deprecated, deferred, or refactored to a schema validator/formatter once the skill stabilizes as the write owner. No obligatory compatibility is maintained with the current `materialize` writer interface.

**Rootline boundary**: Rootline remains the generic Markdown filesystem database and constraint engine. It validates markdown structure, schema compliance, and graph integrity. `roadmapctl` consumes Rootline outputs and adds roadmap-specific policy interpretation.

## Evidence

This decision addresses the experimental session agreement that:
1. No compatibility is obligatory with the current materialize contract.
2. Outcome README is not a persistent task list (tasks are child files).
3. Acceptance criteria belong to Task files, not Outcome README.
4. A human-approved Markdown write is acceptable after validation, not a burden.

The prior decision (keep deterministic writer) was precautionary. Evidence now supports moving writer authority to the skill:
- The skill has proven it can decompose, serialize, and present structured plans for human review.
- Human approval before write is standard in `/roadmap plan` already.
- Postcheck via `roadmapctl check` provides the same safety guarantee regardless of who writes files.
- Removing duplication (skill generating prose + CLI rendering it) simplifies both layers.

## Revisit criteria

This decision is final unless all of the following occur:
1. A measurable regression in safety or consistency is discovered (fallback files, stale links, invalid schema).
2. Human review of dry-runs is no longer practical (e.g., plan sizes exceed review capacity).
3. A new tooling requirement emerges that requires deterministic CLI-side renders (e.g., Rootline native roadmap types).

Until then, optimize skill write clarity and postcheck coverage rather than reintroducing CLI writers.
