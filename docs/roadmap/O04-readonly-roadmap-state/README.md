---
estado: Pending
tipo: outcome
---
# O04: Estado read-only del roadmap

## Objetivo

`roadmapctl` expone APIs read-only para pending, next y decision, moviendo cálculos determinísticos de estado/dependencias fuera del skill y dejando a Rootline como fuente genérica de tree/query/graph.

## Criterios de Éxito

- CE1: `roadmapctl pending` lista trabajo no completado agrupado por outcome/repo.
  - Verificar: golden JSON/text.
- CE2: `roadmapctl next` calcula tasks listas y blockers con explicación determinística.
  - Verificar: fixtures ready/blocked.
- CE3: `roadmapctl decision` produce árbol/scoring determinístico sin juicio AI.
  - Verificar: fixtures de reverse dependencies y quick wins.
- CE4: El skill usa estos comandos en vez de recetas Rootline directas.
  - Verificar: Pi headless después de cutover.

## Invariantes

- INV1: Los comandos son read-only.
  - Verificar: tests no modifican fixtures.
- INV2: Los estados done/active vienen de config operacional y schema Rootline.
  - Verificar: fixtures con labels custom.

## Alcance

**In**:
- Modelo read-only sobre tree/query/graph.
- `pending`, `next`, `decision`.
- Fixtures/goldens y cutover skill.

**Out**:
- Mutar estados.
- Crear archivos.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-normalize-rootline-tree-query-graph-data.md) | Normalizar datos Rootline tree/query/graph |
| [T002](T002-implement-pending-command.md) | Implementar `roadmapctl pending` |
| [T003](T003-implement-next-command.md) | Implementar `roadmapctl next` |
| [T004](T004-implement-decision-command.md) | Implementar `roadmapctl decision` |
| [T005](T005-add-read-api-fixtures-and-goldens.md) | Agregar fixtures/goldens read-only |
| [T006](T006-cutover-skill-pending-next-decision.md) | Actualizar skill para usar comandos read-only |
