---
estado: Specified
tipo: task
---
# T011: Verify roadmapctl as postprocessing boundary

[[blocked_by:./T010-harden-no-python-roadmap-discovery-contract.md]]

## Preserva

- `roadmapctl pending`, `roadmapctl next`, and `roadmapctl decision` remain deterministic JSON interfaces for agents.
  - Verificar: golden or focused tests assert stable fields used by the skill.
- The skill does not reconstruct dependency maps or topological order from Rootline JSON.
  - Verificar: docs and tests cover the returned roadmapctl shape directly.

## Contexto

The no-Python contract is only durable if roadmapctl exposes every field the skill needs: task path, title, status, outcome grouping, readiness, blocking dependencies, diagnostics, and deterministic ordering. This task verifies that command outputs cover the former Python-derived data and documents any intentional gaps.

## Alcance

**In**:
1. Compare the observed Python-derived data from the loop session against `roadmapctl pending --output json`, `roadmapctl next --output json`, and `roadmapctl decision --output json`.
2. Add tests or golden assertions for fields required by the skill contract.
3. Update docs to state that agents should consume roadmapctl JSON directly instead of postprocessing Rootline JSON.
4. Record any Rootline projection improvements as upstream/non-blocking references, not roadmapctl requirements.

**Out**:
- Do not implement new Rootline projection features in this repo.
- Do not add broad new roadmap policy unless a missing field blocks the skill and is explicitly scoped.

## Estado inicial esperado

- T010 has removed or quarantined main-flow Python/Rootline postprocessing from skill docs.
- Existing commands already return most or all required discovery data.

## Criterios de Aceptación

- Tests or golden checks prove `pending`, `next`, and `decision` expose the fields consumed by skill docs.
- Documentation names roadmapctl as the postprocessing boundary for roadmap workflows.
- No Python is required to select the next task in documented flows.
- `go test ./...` passes.
- `go run ./cmd/roadmapctl check --output json --strict` passes.

## Fuente de verdad

- `internal/cli/read_model.go`
- `internal/roadmap/model.go`
- `internal/cli/*pending*`
- `internal/cli/*next*`
- `internal/cli/*decision*`
- `docs/cli-contract.md`
- `docs/roadmap-skill-integration.md`
