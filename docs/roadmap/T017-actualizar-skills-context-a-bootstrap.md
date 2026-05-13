---
tipo: task
estado: Specified
hard_blockers:
  - T016-bootstrap-idempotente-eliminar-context.md
---

# T017 Actualizar skills: context → bootstrap + separar rootline describe

Con `roadmapctl context` eliminado en T016, actualizar los 7 archivos del skill `/roadmap` para usar `roadmapctl bootstrap` y llamar `rootline describe` directamente donde se necesite el schema.

## Criterios de Aceptación

- AC1: `SKILL.md` usa `roadmapctl bootstrap` en lugar de `roadmapctl context` en todos los lugares relevantes
- AC2: `loop-subcommand.md`, `pending-subcommand.md`, `decision-tree-subcommand.md`, `autonomous-mode.md`, `pr-workflow.md` actualizados — sin referencias a `roadmapctl context`
- AC3: Donde el skill necesite schema (valores de `estado`), llama `rootline describe` directamente
- AC4: `./scripts/sync-roadmap-skill.sh --install && ./scripts/sync-roadmap-skill.sh --check` pasan
- AC5: Headless pi test definido en `SKILL.md` pasa

## Archivos a modificar

- `.claude/skills/roadmap/SKILL.md`
- `.claude/skills/roadmap/loop-subcommand.md`
- `.claude/skills/roadmap/pending-subcommand.md`
- `.claude/skills/roadmap/decision-tree-subcommand.md`
- `.claude/skills/roadmap/autonomous-mode.md`
- `.claude/skills/roadmap/pr-workflow.md`

Después de todos los cambios: `sync-roadmap-skill.sh --install`.
