---
estado: Pending
tipo: task
---
# T004: Implementar materialize apply con postcheck

**Outcome**: [O07 Controlador de materialización](README.md)
**Contribuye a**: CE2, CE3

[[blocked_by:./T003-implement-materialize-dry-run-and-diff.md]]
[[blocked_by:../O05-deterministic-lint/T003-implement-outcome-task-table-consistency.md]]

## Preserva

- INV1: Apply requiere `--apply`, preflight y postcheck.
  - Verificar: tests de failure antes/después de write.

## Contexto

Apply crea archivos canónicos y actualiza tablas. Debe validar con Rootline y `roadmapctl check` antes de reportar éxito.

## Alcance

**In**:
1. Crear `OXX/README.md` y `TXXX-*.md` según dry-run.
2. Actualizar `## Tasks` en README de Outcome.
3. Escribir `blocked_by` con paths relativos explícitos.
4. Validar cada archivo y postcheck completo.
5. Detectar drift entre dry-run y apply.

**Out**:
- Commit/push.
- Auto-fix si postcheck falla.
- Implementar código de tasks.

## Estado inicial esperado

- Dry-run genera cambios estables.

## Criterios de Aceptación

- Apply en temp fixture crea estructura canónica.
- Postcheck falla si se generaría estructura inválida.
- No se genera ningún `*-tasks.md`.

## Fuente de verdad

- `.claude/skills/roadmap/plan-subcommand.md`
- `.claude/skills/roadmap/task-guide.md`
- `.claude/skills/roadmap/outcome-guide.md`
