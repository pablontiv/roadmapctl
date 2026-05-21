---
estado: Specified
tipo: task
---
# T006: SKILL.md — agregar routing para subcomando bootstrap

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: `/roadmap bootstrap` se despacha correctamente

[[blocked_by:./T005-bootstrap-skill-open-scan-wizard.md]]

## Preserva

- INV1: el routing table de SKILL.md sigue siendo la única fuente de verdad de dispatch
  - Verificar: SKILL.md contiene la tabla de routing con los subcomandos existentes intactos

## Contexto

El subcomando `bootstrap` creado en T005 no tiene entrada en el routing table de `SKILL.md`. Sin esa entrada, `/roadmap bootstrap` caería en "texto libre → autonomous-mode.md" en lugar de despachar al archivo correcto.

## Alcance

**In**:
1. Agregar fila `bootstrap` → `bootstrap-subcommand.md` en la tabla de routing de `.claude/skills/roadmap/SKILL.md`
2. Actualizar el argument-hint si aplica

**Out**:
- No modificar el contenido de los subcomandos existentes

## Estado inicial esperado

- T005 completada: `bootstrap-subcommand.md` existe
- `.claude/skills/roadmap/SKILL.md` existe

## Criterios de Aceptación

- La tabla de routing en SKILL.md contiene una fila para `bootstrap` apuntando a `bootstrap-subcommand.md`
- Los demás subcomandos (plan, loop, pending, decision) siguen presentes e intactos

## Fuente de verdad

- `.claude/skills/roadmap/SKILL.md` — sección "Routing por subcomando"
