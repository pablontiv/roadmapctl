---
estado: Specified
tipo: task
---
# T004: Sync skill y validar instalación

**Outcome**: [O30 Materialización paralela](README.md)
**Contribuye a**: Que los cambios de T001–T003 estén activos en `~/.claude/skills/roadmap/` y no haya drift

[[blocked_by:./T001-update-skill-md-invariant.md]]
[[blocked_by:./T002-update-plan-subcommand-dispatch.md]]
[[blocked_by:./T003-update-task-guide-materialization.md]]

## Contexto

El source-of-truth de los skills roadmap es `/home/shared/roadmapctl/.claude/skills/roadmap/`. Los skills se distribuyen a `~/.claude/skills/roadmap/` via `scripts/sync-roadmap-skill.sh --install`. Sin el sync, los cambios de T001–T003 existen en el repo pero no están activos en Claude Code.

## Alcance

**In**:
1. Ejecutar `scripts/sync-roadmap-skill.sh --install` desde `/home/shared/roadmapctl`
2. Verificar con `scripts/sync-roadmap-skill.sh --check` que no hay drift

**Out**:
- No modificar archivos del skill (eso fue T001–T003)
- No hacer commit en esta task (el commit de los cambios de skill va en T001–T003 vía integrate)

## Estado inicial esperado

- T001, T002 y T003 completadas y commitadas
- `~/.claude/skills/roadmap/SKILL.md` aún refleja el texto anterior ("único writer")

## Criterios de Aceptación

- `scripts/sync-roadmap-skill.sh --check` retorna 0 (sin drift entre repo y ~/.claude/skills/)
- `grep "coordinador" ~/.claude/skills/roadmap/SKILL.md` retorna al menos 1 match
- `grep "único writer" ~/.claude/skills/roadmap/SKILL.md` retorna 0 matches

## Fuente de verdad

- `scripts/sync-roadmap-skill.sh`
- `~/.claude/skills/roadmap/SKILL.md`
- `~/.claude/skills/roadmap/plan-subcommand.md`
- `~/.claude/skills/roadmap/task-guide.md`
