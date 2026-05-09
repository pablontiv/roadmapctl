---
estado: Pending
tipo: task
---
# T001: Definir schema de plan materializable

**Outcome**: [O07 Controlador de materialización](README.md)
**Contribuye a**: CE1, CE2

[[blocked_by:../O03-config-context-workspace/T004-implement-context-command.md]]
[[blocked_by:../O05-deterministic-lint/T001-define-lint-taxonomy-severity-json.md]]

## Preserva

- INV1: `roadmapctl` no hace descomposición AI; recibe input estructurado aprobado.
  - Verificar: command contract.

## Contexto

El skill puede razonar y proponer un plan, pero la materialización determinística debe recibir datos estructurados: outcomes/tasks/dependencies/ACs.

## Alcance

**In**:
1. Definir formato JSON para outcomes, direct tasks, tasks dentro de outcomes, dependencies y ACs.
2. Definir validación de input y diagnostics.
3. Documentar cómo el skill producirá o pasará ese input.
4. Decidir versionado del input.

**Out**:
- Parsear chat libre.
- Escribir archivos.

## Estado inicial esperado

- `plan-subcommand.md` materializa desde contexto conversacional.

## Criterios de Aceptación

- Existe spec de input con ejemplos.
- Input inválido produce diagnostics claros.
- No hay dependencia de LLM dentro de roadmapctl.

## Fuente de verdad

- `.claude/skills/roadmap/plan-subcommand.md`
- `.claude/skills/roadmap/task-guide.md`
- `.claude/skills/roadmap/outcome-guide.md`
