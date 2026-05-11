---
estado: Completed
tipo: task
---
# T003: Diseñar path planning y guard en roadmapctl

**Outcome**: [Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill](README.md)

[[blocked_by:./T001-record-responsibility-separation-decision.md]]

## Preserva

- roadmapctl conserva diagnostics estables y no asume generación creativa de contenido.

## Contexto

roadmapctl debe ayudar al skill con determinismo: numbering, colisiones, destinos y validación, no con generación de Markdown semántico.

## Alcance

**In**:
1. Diseñar comando o modo para proponer paths canónicos para Outcomes y Tasks.
2. Detectar colisiones, paths no canónicos y destinos inesperados.
3. Definir JSON de salida compacto para el skill.
4. Definir cómo se valida que los paths aprobados coinciden con lo propuesto.

**Out**:
1. No renderizar README/Task body.
2. No mantener ## Tasks en README.

## Estado inicial esperado

internal/materialize renderiza Markdown completo desde un JSON semántico.

## Criterios de Aceptación

- La nueva superficie propone paths OXX-slug/README.md y TXXX-task.md sin escribir contenido.
- La salida permite detectar divergencia antes de Pi write.
- El diseño elimina la necesidad de que roadmapctl reciba estructuras semánticas largas.
- El diseño mantiene check/lint/read-model/transition como responsabilidades roadmapctl.

## Fuente de verdad

- internal/materialize/dryrun.go
- internal/roadmap/numbering.go
- internal/cli/materialize.go
- docs/cli-contract.md
