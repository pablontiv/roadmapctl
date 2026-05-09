---
estado: Pending
tipo: task
---
# T004: Implementar `roadmapctl context`

**Outcome**: [O03 Config/context/workspace](README.md)
**Contribuye a**: CE3

[[blocked_by:./T002-implement-config-discovery-and-toml-loader.md]]
[[blocked_by:../O02-post-mvp-foundations/T005-add-roadmap-domain-model-and-tree-wrapper.md]]

## Preserva

- INV1: JSON mode emite un único objeto parseable.
  - Verificar: golden de `roadmapctl/context`.

## Contexto

El skill hoy calcula bootstrap, helpers y status filters en prose. `roadmapctl context` debe exponer el contexto efectivo para que el skill sea un adapter fino.

## Alcance

**In**:
1. Agregar comando directo `roadmapctl context`.
2. Reportar repo root, roadmap root, config path/source, Rootline version, schema `estado`/`tipo`, roles y helpers.
3. Soportar `--repo`, `--roadmap-root`, `--workspace`, `--output json|text`, `--strict`.
4. Incluir diagnostics sin romper JSON stdout.

**Out**:
- Listar pending tasks.
- Crear config o `.stem`.

## Estado inicial esperado

- `doctor` reporta paths básicos, pero no contexto operacional completo.

## Criterios de Aceptación

- `roadmapctl context --repo . --output json` incluye `kind: roadmapctl/context`.
- Helpers `where_leaf`, `where_not_done`, `where_active` coinciden con config efectiva.
- Si usa legacy config, el output lo indica.

## Fuente de verdad

- `internal/cli/*`
- `internal/config/*`
- `internal/rootlinecli/*`
- `.claude/skills/roadmap/SKILL.md`
