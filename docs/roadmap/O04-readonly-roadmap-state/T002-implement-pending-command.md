---
estado: Pending
tipo: task
---
# T002: Implementar `roadmapctl pending`

**Outcome**: [O04 Estado read-only](README.md)
**Contribuye a**: CE1

[[blocked_by:./T001-normalize-rootline-tree-query-graph-data.md]]

## Preserva

- INV1: Comando read-only sin writes.
  - Verificar: fixtures inmutables.

## Contexto

El skill `pending-subcommand.md` usa `rootline tree` directamente. Ese comportamiento debe pasar a `roadmapctl pending` para centralizar filtros, workspace y diagnostics.

## Alcance

**In**:
1. Agregar comando directo `roadmapctl pending`.
2. Usar config `done_statuses`, `leaf_filter` y context.
3. Output JSON/text agrupado por outcome y repo.
4. Incluir counts y diagnostics.

**Out**:
- Seleccionar next task.
- Mutar estados.

## Estado inicial esperado

- Modelo read-only normalizado existe.

## Criterios de Aceptación

- `roadmapctl pending --output json` produce `kind: roadmapctl/pending`.
- Tasks completadas/obsolete no aparecen.
- Workspace output agrupa por repo cuando workspace existe.

## Fuente de verdad

- `.claude/skills/roadmap/pending-subcommand.md`
- `internal/cli/*`
- `internal/roadmap/*`
