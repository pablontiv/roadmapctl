---
estado: Completed
tipo: task
---
# T003: Implementar can-start y can-complete

**Outcome**: [O06 Controlador de transiciones](README.md)
**Contribuye a**: CE1

[[blocked_by:./T001-design-transition-status-role-model.md]]

## Preserva

- INV1: Comandos can-* son read-only.
  - Verificar: tests no modifican fixtures.

## Contexto

Antes de mutar estado, roadmapctl debe poder responder si una task puede iniciar o completarse, con razones claras.

## Alcance

**In**:
1. `roadmapctl transition can-start <task-path>`.
2. `roadmapctl transition can-complete <task-path>`.
3. Validar existencia, tipo task, status permitido y dependencias done.
4. Output JSON con `allowed`, `reasons`, `blocking_dependencies`.

**Out**:
- Mutar estados.
- Ejecutar ACs del proyecto.

## Estado inicial esperado

- Modelo de transiciones diseñado.

## Criterios de Aceptación

- Task bloqueada por dependency incomplete devuelve allowed=false y path dependency.
- Task ready devuelve allowed=true.
- Custom status labels funcionan con `done_statuses`.

## Fuente de verdad

- `internal/roadmap/*`
- `.claude/skills/roadmap/loop-subcommand.md`
