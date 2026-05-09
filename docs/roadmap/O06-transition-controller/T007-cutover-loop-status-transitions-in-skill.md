---
estado: Pending
tipo: task
---
# T007: Cortar loop/status del skill hacia transition

**Outcome**: [O06 Controlador de transiciones](README.md)
**Contribuye a**: CE3

[[blocked_by:./T006-add-transition-fixtures-and-goldens.md]]

## Preserva

- INV1: Cambios de skill requieren Pi headless verification.
  - Verificar: escenarios documentados.

## Contexto

`loop-subcommand.md` actualmente llama `rootline set` para In Progress/Completed. Una vez estable `roadmapctl transition`, el skill debe delegar esas mutaciones.

## Alcance

**In**:
1. Actualizar `loop-subcommand.md` para usar `transition can-start/start/complete`.
2. Remover o degradar instrucciones directas `rootline set` donde roadmapctl ya gobierna.
3. Mantener ejecución de código/ACs en el agente.
4. Ejecutar sync del skill y Pi headless verification.

**Out**:
- Cutover de materialización.
- Cambiar pending/decision si ya fue hecho.

## Estado inicial esperado

- Transition commands testeados.

## Criterios de Aceptación

- Skill no duplica reglas de transición cubiertas por roadmapctl.
- Headless loop preflight muestra uso correcto de roadmapctl.
- `scripts/sync-roadmap-skill.sh --check` pasa.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
- `docs/roadmap-skill-integration.md`
- `scripts/sync-roadmap-skill.sh`
