---
tipo: outcome
---
# Repo-local roadmap execution settings

Move roadmap loop behavior from per-invocation skill flags into the canonical per-repository .roadmapctl.toml contract, with forced legacy migration, config-driven loop behavior, and roadmap-specific context compaction support.

## Criterios de Aceptación

- roadmapctl uses <roadmap-root>/.roadmapctl.toml as the only durable config source and migrates legacy .claude/roadmap.local.md during config load.
- roadmapctl context exposes the operational execution settings required by the roadmap skill.
- The roadmap skill documents /roadmap loop as config-driven, keeping only --filter and --max as user flags.
- A Pi extension/tool exists for roadmap-specific context compaction, with /compact fallback documented in the skill.
- go test ./... passes and roadmap skill headless verification covers the updated loop/materialization bootstrap guards.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-implement-config-fields-and-forced-migration.md) | Extend roadmapctl config loading with repo-local loop execution fields, defaults, validation, and forced migration/deletion of legacy .claude/roadmap.local.md during config load. |
| [T002](T002-expose-execution-settings-in-context.md) | Extend roadmapctl context JSON and workspace repo context so the skill can consume execution settings without parsing TOML directly. |
| [T003](T003-update-config-templates-fixtures-and-contract-docs.md) | Make generated .roadmapctl.toml files, fixtures, and docs include the new execution settings and the forced legacy migration policy. |
| [T004](T004-add-roadmap-context-compaction-extension.md) | Add a project-local Pi extension/tool that triggers roadmap-specific context compaction after durable task completion. |
| [T005](T005-cutover-roadmap-skill-loop-to-config.md) | Update the roadmap skill docs so /roadmap loop reads behavior from roadmapctl context, removes behavior flags, uses opportunistic waves, and applies autonomy-specific continuation/repair rules. |
| [T006](T006-verify-repo-settings-flow-end-to-end.md) | Add/update tests and headless verification evidence covering config migration, context exposure, skill cutover, and roadmap-specific compaction fallback behavior. |
