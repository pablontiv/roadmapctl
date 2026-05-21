---
estado: Completed
tipo: task
---
# T004: loop-subcommand.md — declarar gate único y ampliar config snapshot

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: el loop nunca ejecuta git/gh directamente; pasa style fields a /integrate

[[blocked_by:./T002-bootstrap-serialize-and-diagnostic.md]]

## Preserva

- INV1: el paso 9 del loop sigue invocando Skill("integrate", ...) como su acción principal
  - Verificar: el archivo contiene `Skill("integrate"` en el paso 9

## Contexto

**Causa A** del incidente: el paso 9 del loop está enterrado entre cómputos burocráticos (`transition complete`, recálculo de `is_last_in_scope`, recálculo de `next`) sin una declaración prominente de que `/integrate` es el único gate a `git`/`gh`. Un LLM que optimiza puede saltarse el skill y hacer git por su cuenta.

Además, el config snapshot pasado a `/integrate` en el paso 9 incluye `pr_merge_strategy` (deprecated) y no incluye los 3 nuevos style fields.

## Alcance

**In**:
1. En el paso 9 de `loop-subcommand.md`: agregar declaración explícita al inicio del bloque: *"`/integrate` es la única puerta a `git` y `gh` desde el loop. El loop NO ejecuta `git commit/push` ni `gh pr create/merge` directamente bajo ninguna condición, ni siquiera para `pr_mode=false`. Saltarlo es un bug del flujo."*
2. Actualizar el config snapshot en la invocación `Skill("integrate", ..., config=<snapshot>)`: agregar `branch_style`, `pr_title_style`, `pr_body_style`, `base_branch`; remover `pr_merge_strategy`

**Out**:
- No modificar bootstrap-reference.md (T005)
- No modificar integrate/SKILL.md (T003)

## Estado inicial esperado

- T002 completada
- `.claude/skills/roadmap/loop-subcommand.md` existe

## Criterios de Aceptación

- `grep -i 'única puerta\|único gate\|only gate' .claude/skills/roadmap/loop-subcommand.md` retorna al menos 1 match en el paso 9
- El config snapshot del `Skill("integrate",...)` en paso 9 contiene `branch_style`, `pr_title_style`, `pr_body_style`
- El config snapshot no contiene `pr_merge_strategy`

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md` — paso 9 (Integrate)
