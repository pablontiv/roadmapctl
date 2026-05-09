---
estado: Completed
tipo: task
---
# T004: Implementar `roadmapctl decision`

**Outcome**: [O04 Estado read-only](README.md)
**Contribuye a**: CE3

[[blocked_by:./T003-implement-next-command.md]]

## Preserva

- INV1: Scoring explicable y determinístico.
  - Verificar: golden JSON incluye razones.

## Contexto

`decision-tree-subcommand.md` calcula quick wins, blockers y reverse dependencies en prompt. Eso debe convertirse en API determinística para que el skill presente opciones sin recalcular lógica.

## Alcance

**In**:
1. Agregar `roadmapctl decision`.
2. Calcular reverse dependencies/unblocks.
3. Identificar quick wins y critical blockers con reglas documentadas.
4. Output JSON/text para UI humana y automation.

**Out**:
- Priorización AI subjetiva.
- Mutaciones o ejecución.

## Estado inicial esperado

- `roadmapctl next` produce readiness y blockers.

## Criterios de Aceptación

- Decision output es estable para el mismo fixture.
- Cada recomendación incluye razones y datos fuente.
- No usa git log u otras fuentes no aprobadas sin documentarlo.

## Fuente de verdad

- `.claude/skills/roadmap/decision-tree-subcommand.md`
- `internal/roadmap/*`
