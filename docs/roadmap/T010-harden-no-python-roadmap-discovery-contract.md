---
estado: Specified
tipo: task
---
# T010: Harden no-Python roadmap discovery contract

[[blocked_by:./T002-cutover-loop-discovery-to-roadmapctl.md]]

## Preserva

- The roadmap skill remains a thin adapter over roadmapctl for pending, next, and decision flows.
  - Verificar: skill docs do not recalculate blockers or readiness in prompt text.
- Rootline remains the underlying data engine, not the roadmap policy engine.
  - Verificar: direct Rootline query/graph/tree usage is documented only as troubleshooting or low-level inspection.

## Contexto

A live loop session used `rootline graph` and `rootline query` followed by Python snippets to extract cycles, broken links, edges, task titles, statuses, and links. The repository now has `roadmapctl pending`, `roadmapctl next`, and `roadmapctl decision`; these commands should be the canonical no-postprocessing discovery boundary for agents.

## Alcance

**In**:
1. Audit `.claude/skills/roadmap/*.md` for remaining main-flow `rootline graph`, `rootline query`, or `rootline tree` usage in pending, next, decision, and loop discovery.
2. Replace any main-flow Rootline discovery with `roadmapctl pending`, `roadmapctl next`, or `roadmapctl decision` as appropriate.
3. Remove or rewrite Python postprocessing snippets for Rootline JSON from skill docs.
4. Keep Rootline commands only in troubleshooting/reference sections.

**Out**:
- Do not change roadmapctl Go behavior unless the audit exposes a blocking mismatch requiring separate approval.
- Do not modify Rootline itself.

## Estado inicial esperado

- `T002-cutover-loop-discovery-to-roadmapctl.md` has cut over the loop discovery path.
- Some docs may still mention Rootline low-level commands for historical or troubleshooting contexts.

## Criterios de Aceptación

- Pending, next, decision, and loop docs use roadmapctl commands as their primary flow.
- No skill main flow contains Python snippets to postprocess Rootline JSON.
- References to `rootline graph`, `rootline query`, and `rootline tree` are clearly marked as troubleshooting/legacy when present.
- `go test ./...` passes.
- `go run ./cmd/roadmapctl check --output json --strict` passes.

## Fuente de verdad

- `.claude/skills/roadmap/pending-subcommand.md`
- `.claude/skills/roadmap/decision-tree-subcommand.md`
- `.claude/skills/roadmap/loop-subcommand.md`
- `.claude/skills/roadmap/common-logic.md`
- `docs/roadmap-skill-integration.md`
- `docs/cli-contract.md`
