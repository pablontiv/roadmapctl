---
estado: Completed
tipo: task
---
# T008: Añadir verificación headless del flujo completo

**Outcome**: [Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill](README.md)

[[blocked_by:./T002-redefine-roadmap-plan-skill-contract.md]]
[[blocked_by:./T003-design-roadmapctl-path-planning-guard.md]]
[[blocked_by:./T007-migrate-docs-fixtures-and-goldens-to-new-boundary.md]]
[[blocked_by:./T006-verify-rootline-generic-record-operations.md]]

## Preserva

- Los guards bloquean writes si roadmapctl/Rootline reportan errores.

## Contexto

El cambio afecta el comportamiento crítico del skill; necesita evidencia headless, no solo grep o tests unitarios.

## Alcance

**In**:
1. Añadir escenario que prueba no-write antes de aprobación.
2. Añadir escenario que prueba Pi write después de aprobación.
3. Añadir escenario que prueba escritura paralela de archivos independientes.
4. Añadir escenario de divergencia que re-pregunta antes de escribir.
5. Validar que no se crea *-tasks.md.

**Out**:
1. No ejecutar implementación de tasks de código; solo materialización de roadmap.

## Estado inicial esperado

Los headless actuales verifican preflight/materialize viejo, no el flujo Pi write aprobado.

## Criterios de Aceptación

- Headless evidencia pregunta obligatoria antes de materializar.
- Headless evidencia archivos canónicos sin ## Tasks en Outcome.
- Headless evidencia AC por Task.
- Headless evidencia writes paralelos cuando es seguro.
- Headless evidencia postcheck ok y divergencia bloqueante cuando corresponde.

## Fuente de verdad

- scripts/verify-roadmap-skill-headless.sh
- .claude/skills/roadmap/plan-subcommand.md
- docs/roadmap-skill-integration.md
