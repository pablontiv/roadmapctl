---
estado: Completed
tipo: task
---
# T003: Implementar `roadmapctl next`

**Outcome**: [O04 Estado read-only](README.md)
**Contribuye a**: CE2

[[blocked_by:./T001-normalize-rootline-tree-query-graph-data.md]]

## Preserva

- INV1: Selección determinística, no juicio AI.
  - Verificar: fixture/golden estable.

## Contexto

El loop del skill calcula readiness y orden topológico en prose. `roadmapctl next` debe exponer qué tasks están listas, cuáles bloqueadas y por qué.

## Alcance

**In**:
1. Agregar comando `roadmapctl next`.
2. Calcular tasks ready si todas sus dependencias están en `done_statuses`.
3. Mostrar blockers para tasks no ready.
4. Soportar `--limit` y filtros aprobados.
5. Orden estable: topo/order/path como fallback documentado.

**Out**:
- Ejecutar tasks.
- Cambiar status a In Progress.

## Estado inicial esperado

- Modelo read-only normalizado existe.

## Criterios de Aceptación

- Fixture con una task blocked explica dependency faltante.
- Fixture con varias ready produce orden determinístico.
- No se usa `Completed` hardcodeado; se usa `done_statuses`.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
- `internal/roadmap/*`
- `testdata/fixtures/*`
