---
estado: Completed
tipo: task
---
# T004: Document and test materialize postcheck recovery

**Outcome**: [Stabilize roadmap skill routing and materialization UX](README.md)

## Preserva

- Agents must not claim successful materialization until explicit roadmapctl check passes.
- Recovery remains explicit and evidence-backed, not automatic silent fixing.

## Contexto

Session evidence showed roadmapctl materialize can create files and then fail postcheck, leaving agents without a clear recovery path.

## Alcance

**In**:
1. Document recovery steps for partial materialization or failed postcheck.
2. Require inspecting applied changes, running rootline validate on affected paths, repairing or reverting deliberately, and rerunning roadmapctl check --strict.
3. Add a focused test or fixture for postcheck-failure behavior if feasible.

**Out**:
1. No automatic rollback unless separately designed and approved.
2. No broad rootline auto-fix behavior.

## Estado inicial esperado

Materialize docs require postcheck but do not clearly describe recovery when apply partially succeeds and postcheck fails.

## Criterios de Aceptación

- Skill or integration docs describe a concrete recovery path for postcheck failure after materialize.
- Agents are instructed to report partial state and not commit or claim success until roadmapctl check passes.
- A test or documented manual validation covers the partial-write/postcheck-failure scenario.

## Fuente de verdad

- .claude/skills/roadmap/plan-subcommand.md
- .claude/skills/roadmap/common-logic.md
- docs/roadmap-skill-integration.md
- internal/cli/materialize.go
- internal/cli/materialize_test.go
