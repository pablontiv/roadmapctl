---
estado: Specified
tipo: task
---
# T010: Actualizar docs/cli-contract.md y CHANGELOG.md

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: documentación pública refleja los nuevos campos y el nuevo diagnostic

[[blocked_by:./T003-integrate-skill-fix-causa-b.md]]
[[blocked_by:./T007-todos-subcomandos-referencian-bootstrap.md]]

## Preserva

- INV1: los contratos existentes documentados no se eliminan ni contradicen
  - Verificar: los campos y diagnostics existentes siguen documentados en cli-contract.md

## Contexto

Con O28 se agregan:
- 3 campos nuevos al TOML / JSON de bootstrap: `branch_style`, `pr_title_style`, `pr_body_style`
- 1 diagnostic nuevo: `RMC_GITFLOW_NOT_CONFIGURED`
- 1 warning de deprecation: `RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED`
- Cambio de comportamiento de `/integrate`: Fase 2 siempre, `gh pr merge --auto` sin flags

Todo esto debe quedar en `docs/cli-contract.md` y el CHANGELOG.

## Alcance

**In**:
1. `docs/cli-contract.md`: documentar los 3 campos nuevos bajo `[gitflow]`, el diagnostic `RMC_GITFLOW_NOT_CONFIGURED` y el warning `RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED`
2. `CHANGELOG.md`: entry para O28 describiendo los cambios (nuevos campos, fix Causa A/B, deprecated pr_merge_strategy)

**Out**:
- No modificar código ni skills

## Estado inicial esperado

- T003 y T007 completadas

## Criterios de Aceptación

- `grep 'branch_style\|pr_title_style\|pr_body_style' docs/cli-contract.md` → al menos 3 matches
- `grep 'RMC_GITFLOW_NOT_CONFIGURED\|RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED' docs/cli-contract.md` → 2 matches
- `CHANGELOG.md` tiene una entrada para O28

## Fuente de verdad

- `docs/cli-contract.md`
- `CHANGELOG.md`
