---
tipo: outcome
---
# Stabilize roadmap skill routing and materialization UX

The roadmap skill routes default invocations correctly, keeps materialization token-light, aligns materialize documentation, and documents safe recovery for postcheck failures while preserving roadmapctl-owned canonical writes.

## Criterios de Aceptación

- /roadmap with no arguments runs the pending flow through roadmapctl pending instead of decision-tree.
- /roadmap plan materialization keeps large plan and dry-run JSON in temporary files and reports only concise status, diagnostics, and planned paths by default.
- Skill docs, integration docs, and materialize schema docs agree on batch apply as the normal path and target apply as recovery or explicit one-file approval.
- Post-materialize or postcheck failures have documented recovery steps and validation coverage.
- Legacy config and schema compatibility issues are either fixed in diagnostics or clearly routed to existing O11 work.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-route-empty-roadmap-invocation-to-pending.md) | Change the roadmap skill so invoking it without arguments runs the pending view instead of the decision tree. |
| [T002](T002-make-roadmap-plan-materialization-token-light.md) | Update roadmap plan instructions so plan and dry-run JSON stay in temporary files and normal output reports concise summaries instead of full content and diffs. |
| [T003](T003-align-materialize-batch-and-target-apply-docs.md) | Remove contradictions between integration docs and skill docs about batch apply versus per-target apply. |
| [T004](T004-document-and-test-materialize-postcheck-recovery.md) | Define what agents should do when roadmapctl materialize writes files but validation or postcheck fails. |
| [T005](T005-clarify-materialize-schema-source-and-cli-exposure.md) | Make it clear where agents should obtain the roadmapctl materialize plan schema and evaluate exposing it through the CLI. |
| [T006](T006-fix-legacy-config-diagnostic-paths.md) | Replace stale hardcoded legacy config paths in diagnostics with the effective config source when appropriate. |
| [T007](T007-add-headless-roadmap-skill-regression-scenarios.md) | Add robust Pi headless checks for the skill bugs fixed in this outcome. |
| [T008](T008-decide-future-of-materialize-writer-versus-guard-flow.md) | Record a decision on whether roadmapctl materialize remains the canonical writer, becomes a lighter guard/path planner, or is deprecated for direct skill writes plus checks. |
