---
estado: Specified
tipo: task
---
# T006: Fix legacy config diagnostic paths

**Outcome**: [Stabilize roadmap skill routing and materialization UX](README.md)

## Preserva

- Legacy migration behavior remains supported and tested.
- Diagnostics stay stable and machine-readable.

## Contexto

The skill no longer treats .claude/roadmap.local.md as a durable source, but some diagnostics and goldens still point at that legacy path.

## Alcance

**In**:
1. Inspect diagnostics that reference .claude/roadmap.local.md.
2. Use the effective config path for TOML-backed diagnostics where appropriate.
3. Update tests and goldens for intentional legacy versus stale references.

**Out**:
1. No removal of legacy migration support.
2. No broad fixture migration unless required by tests.

## Estado inicial esperado

internal/roadmap/status.go and related tests/goldens still include legacy config paths in some diagnostics.

## Criterios de Aceptación

- TOML-backed config diagnostics point to the effective .roadmapctl.toml path when possible.
- Legacy-only migration tests still explicitly cover .claude/roadmap.local.md.
- go test ./... passes with updated goldens.

## Fuente de verdad

- internal/roadmap/status.go
- internal/roadmap/dependencies_test.go
- internal/config/config.go
- testdata/golden
- docs/cli-contract.md
