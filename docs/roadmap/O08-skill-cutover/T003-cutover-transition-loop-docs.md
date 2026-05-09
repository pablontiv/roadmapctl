---
estado: Pending
tipo: task
---
# T003: Cortar loop/status docs hacia transition

**Outcome**: [O08 Cutover de skill](README.md)
**Contribuye a**: CE2

[[blocked_by:../O06-transition-controller/T007-cutover-loop-status-transitions-in-skill.md]]

## Preserva

- INV1: El agente sigue implementando código y ejecutando ACs; roadmapctl solo gobierna estado.
  - Verificar: loop docs.

## Contexto

Después del cutover de O06, este task asegura que la documentación del skill no siga enseñando `rootline set` directo como ruta primaria.

## Alcance

**In**:
1. Revisar `loop-subcommand.md`.
2. Mantener Rootline como referencia baja solo si hace falta troubleshooting.
3. Documentar `roadmapctl transition` como ruta primaria.
4. Ejecutar Pi headless loop scenario.

**Out**:
- Materialization docs.
- Cambiar implementación de transition.

## Estado inicial esperado

- O06/T007 completado.

## Criterios de Aceptación

- Loop docs no duplican reglas de can-start/can-complete.
- Headless loop evidencia preflight y transition route.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
- `docs/roadmap-skill-integration.md`
