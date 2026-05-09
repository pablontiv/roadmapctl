---
estado: Specified
tipo: task
---
# T001: Implement config fields and forced legacy migration

**Outcome**: [Repo-local roadmap execution settings](README.md)

## Preserva

- TOML remains the canonical config when present.
- Invalid TOML fails without falling back to legacy.
- Legacy config is never deleted unless TOML generation/loading succeeds.

## Contexto

The design requires roadmapctl, not the skill, to own config migration while reading config. Legacy frontmatter is only migration input and must not remain a lasting config source.

## Alcance

**In**:
1. Add Config fields for loop_max_tasks, parallel, autonomy, compact_after_task_commit, and pr_mode.
2. Add TOML parser fields, defaults, validation, render output, and configDiffers coverage for the new fields.
3. Change Load so TOML plus legacy deletes legacy after TOML loads successfully.
4. Change Load so legacy-only repos are migrated to <roadmap-root>/.roadmapctl.toml, validated through the config parser, and then legacy is deleted.
5. Add unit tests for successful migration, existing TOML legacy deletion, invalid TOML no fallback, invalid autonomy, and negative loop_max_tasks.

**Out**:
1. No roadmap task execution changes in this task.
2. No Pi extension work in this task.
3. No skill markdown cutover in this task.

## Estado inicial esperado

internal/config/config.go supports TOML and legacy fallback with conflict warnings, but it does not model loop execution settings and does not force migration/deletion on Load.

## Especificación Técnica

Use TOML keys loop_max_tasks, parallel, autonomy, compact_after_task_commit, and pr_mode. Defaults must be loop_max_tasks=0, parallel=true, autonomy=until_done, compact_after_task_commit=true, pr_mode=false. Represent booleans as *bool in the decoded TOML struct so explicit false overrides defaults. Add a small validation helper called from Load after applying config values. For legacy-only migration, generate TOML with renderTOMLConfig, mkdir the target roadmap root if it already resolves inside the repo, write the TOML, re-load/validate the generated TOML, and only then remove .claude/roadmap.local.md. If removal fails, return a config error rather than silently leaving two sources.

## Criterios de Aceptación

- go test ./internal/config passes.
- A legacy-only temp repo passed to config.Load receives docs/roadmap/.roadmapctl.toml and no longer has .claude/roadmap.local.md after successful load.
- A repo with existing docs/roadmap/.roadmapctl.toml and legacy config loads TOML and deletes legacy without emitting a conflict warning.
- A repo with invalid docs/roadmap/.roadmapctl.toml returns RMC_CONFIG_PARSE even if legacy exists.
- autonomy accepts only manual, supervised, and until_done; negative loop_max_tasks returns a config parse/validation error.

## Fuente de verdad

- docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md
- internal/config/config.go
- internal/config/config_test.go
