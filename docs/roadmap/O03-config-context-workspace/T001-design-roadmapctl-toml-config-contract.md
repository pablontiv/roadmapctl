---
estado: Completed
tipo: task
---
# T001: Diseñar contrato `.roadmapctl.toml`

**Outcome**: [O03 Config/context/workspace](README.md)
**Contribuye a**: CE1, INV1

[[blocked_by:../O02-post-mvp-foundations/T003-validate-operational-status-roles-separately.md]]

## Preserva

- INV1: La config operacional no define schema documental.
  - Verificar: no hay enums `estado`/`tipo` como autoridad en TOML.

## Contexto

Decisión aprobada: mover la config local desde `.claude/roadmap.local.md` a `<roadmap-root>/.roadmapctl.toml`, por defecto `docs/roadmap/.roadmapctl.toml`. El root se infiere del directorio del TOML.

## Alcance

**In**:
1. Definir keys TOML: `done_statuses`, `active_statuses`, `leaf_filter`, `outcome_close_verify`, `pr_merge_strategy`, `commit_style`, `auto_push`, `[status_values]`.
2. Decidir si `roadmap_root` existe en TOML o solo se infiere.
3. Definir defaults, precedence y diagnostics.
4. Documentar legacy fallback y conflicto.

**Out**:
- Implementar loader.
- Migrar archivos existentes.

## Estado inicial esperado

- `.claude/roadmap.local.md` es la fuente actual.

## Criterios de Aceptación

- Hay contrato documentado para TOML y ejemplos single-repo/workspace si aplica.
- El contrato diferencia roles operacionales de schema documental.
- Quedan decisiones abiertas marcadas si no se pueden resolver sin aprobación.

## Fuente de verdad

- `.claude/roadmap.local.md`
- `internal/config/config.go`
- `docs/cli-contract.md`
