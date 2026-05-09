---
estado: Pending
tipo: task
---
# T007: Cortar plan-subcommand hacia `roadmapctl materialize`

**Outcome**: [O07 Controlador de materialización](README.md)
**Contribuye a**: CE2, CE3

[[blocked_by:./T006-add-materialization-fixtures-and-goldens.md]]

## Preserva

- INV1: El skill sigue haciendo descomposición conceptual y aprobación humana.
  - Verificar: `plan-subcommand.md` conserva fase de aprobación.

## Contexto

Cuando `materialize` esté estable, el skill debe dejar de hacer numbering, rootline new, writes y table updates manuales. Debe producir/usar input estructurado y delegar al CLI.

## Alcance

**In**:
1. Actualizar `plan-subcommand.md` para generar/pasar plan estructurado.
2. Reemplazar writes manuales con `roadmapctl materialize --dry-run` y luego `--apply` cuando aprobado.
3. Mantener preflight/postcheck y no fallback.
4. Ejecutar Pi headless verification.

**Out**:
- Cambiar pending/transition si ya fue cortado.
- Implementar code tasks.

## Estado inicial esperado

- Materialize dry-run/apply testeado.

## Criterios de Aceptación

- Skill no contiene instrucciones shell de numbering/materialization duplicadas salvo referencia histórica mínima.
- Headless materialization preflight pasa.
- `scripts/sync-roadmap-skill.sh --check` pasa.

## Fuente de verdad

- `.claude/skills/roadmap/plan-subcommand.md`
- `docs/roadmap-skill-integration.md`
