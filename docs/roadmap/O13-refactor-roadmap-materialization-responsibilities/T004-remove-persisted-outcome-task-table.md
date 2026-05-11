---
estado: Completed
tipo: task
---
# T004: Eliminar ## Tasks como dato persistido

**Outcome**: [Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill](README.md)

[[blocked_by:./T001-record-responsibility-separation-decision.md]]

## Preserva

- La estructura canónica OXX/README.md + TXXX-*.md sigue siendo la fuente de verdad.

## Contexto

La tabla ## Tasks duplica información derivable y obliga a updates del README que pueden quedar stale.

## Alcance

**In**:
1. Eliminar renderer/update de ## Tasks.
2. Actualizar lint/check para no exigir tabla.
3. Asegurar que read-model lista tasks hijas desde estructura canónica.
4. Actualizar fixtures/goldens afectados.

**Out**:
1. No eliminar vistas de tasks; deben calcularse en comandos read-only.

## Estado inicial esperado

Outcome guide, lint contract y materialize renderer esperan y mantienen ## Tasks.

## Criterios de Aceptación

- Outcome README materializado o escrito por skill no contiene ## Tasks por defecto.
- roadmapctl lint/check no reportan falta de tabla ## Tasks.
- pending/next/decision siguen encontrando tasks hijas.
- Existing outcome append ya no requiere update de README por tabla.

## Fuente de verdad

- internal/materialize/dryrun.go
- internal/lint
- internal/roadmap
- docs/cli-contract.md
- .claude/skills/roadmap/outcome-guide.md
