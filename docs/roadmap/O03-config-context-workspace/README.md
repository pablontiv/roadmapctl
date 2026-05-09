---
tipo: outcome
---
# O03: Config, context, workspace y bootstrap

## Objetivo

La configuración operacional de roadmap vive junto al roadmap en `<roadmap-root>/.roadmapctl.toml`, `roadmapctl context` expone el contexto efectivo, y el bootstrap/workspace discovery dejan de depender de prose del skill.

## Criterios de Éxito

- CE1: `roadmapctl` carga `docs/roadmap/.roadmapctl.toml` e infiere el roadmap root desde su directorio.
  - Verificar: fixture `valid-roadmapctl-toml-default`.
- CE2: `.claude/roadmap.local.md` sigue funcionando como legacy fallback/migración.
  - Verificar: fixture `valid-legacy-config-fallback`.
- CE3: `roadmapctl context --output json` reporta root, config source, schema, roles y helpers.
  - Verificar: golden JSON.
- CE4: workspace mode tiene comportamiento determinístico y diagnosticable.
  - Verificar: fixtures workspace.

## Invariantes

- INV1: `.roadmapctl.toml` no define enums documentales; solo roles/políticas operacionales.
  - Verificar: schema viene desde Rootline describe.
- INV2: Config nueva no rompe repos legacy sin migración explícita.
  - Verificar: fallback tests.

## Alcance

**In**:
- TOML config en `<roadmap-root>/.roadmapctl.toml`.
- Discovery, legacy fallback, conflicto y migración.
- `context`, `bootstrap inspect/init`, workspace discovery.

**Out**:
- Pending/next/decision.
- Materializar tasks.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-design-roadmapctl-toml-config-contract.md) | Diseñar contrato TOML operacional |
| [T002](T002-implement-config-discovery-and-toml-loader.md) | Implementar discovery y loader TOML |
| [T003](T003-implement-legacy-config-migration.md) | Implementar migración/fallback legacy |
| [T004](T004-implement-context-command.md) | Implementar `roadmapctl context` |
| [T005](T005-implement-bootstrap-inspect-init.md) | Implementar bootstrap inspect/init |
| [T006](T006-implement-workspace-discovery.md) | Implementar workspace discovery |
| [T007](T007-add-config-context-workspace-fixtures.md) | Agregar fixtures/goldens de config/context/workspace |
