---
estado: Completed
tipo: task
---
# T001: Registrar la nueva decisión arquitectónica

**Outcome**: [Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill](README.md)

## Preserva

- Rootline sigue siendo genérico y no aprende semántica roadmap-specific.

## Contexto

La sesión acordó que no hay compat obligatoria, que Outcome no persiste ## Tasks, que los AC viven en Tasks y que Pi write puede escribir Markdown aprobado.

## Alcance

**In**:
1. Actualizar o reemplazar la decisión sobre materialize writer.
2. Documentar el boundary final entre skill, Pi write, roadmapctl y Rootline.
3. Declarar explícitamente que no se mantiene compat obligatoria con el contrato materialize actual.

**Out**:
1. No implementar todavía cambios de CLI o skill.

## Estado inicial esperado

docs/decisions/materialize-writer-vs-guard-flow.md todavía acepta roadmapctl materialize como writer determinístico.

## Criterios de Aceptación

- La decisión declara que Outcome README no persiste ## Tasks.
- La decisión declara que Outcome no requiere AC y que los AC viven en Tasks.
- La decisión declara que Pi write escribe Markdown aprobado tras aprobación humana.
- La decisión declara que roadmapctl queda como guard/path planner/validator/view/policy layer.

## Fuente de verdad

- docs/decisions/materialize-writer-vs-guard-flow.md
- README.md
- docs/cli-contract.md
- docs/roadmap-skill-integration.md
