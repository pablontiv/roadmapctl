---
estado: Specified
tipo: task
---
# T003: integrate/SKILL.md — eliminar bug Causa B y migrar a style fields

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: `/integrate` nunca commitea en branch actual; genera artefactos desde style fields

[[blocked_by:./T002-bootstrap-serialize-and-diagnostic.md]]

## Preserva

- INV1: las 6 fases de integrate siguen existiendo en el mismo orden
  - Verificar: el archivo tiene secciones Fase 1 a Fase 6

## Contexto

Bug documentado en `integrate/SKILL.md` línea 99:
> "Si `pr_mode == false`, omitir esta fase; commitear en el branch actual."

Esta instrucción literal causa que con `pr_mode=false` el LLM commitee en el branch actual, que puede ser `main`. Es la **Causa B** del incidente wcc-deviceas.

Además, el skill tiene hardcodeado:
- Branch naming: `feat/<scope>`, `feat/direct-roadmap-tasks`
- Tabla determinística de tipos de commit (docs/chore/refactor/feat/fix)
- Template fijo de PR body
- Trailers `Co-authored-by`, `Signed-off-by`
- `gh pr merge --auto --squash --delete-branch` (flags que decide GitHub, no el CLI)

Todo eso debe reemplazarse por generación LLM desde los style fields del TOML.

## Alcance

**In**:
1. Eliminar literalmente la frase *"Si `pr_mode == false`, omitir esta fase; commitear en el branch actual"* de Fase 2
2. Reemplazar por: *"Fase 2 corre siempre. Branch local por scope es invariante. `pr_mode` solo controla si Fase 5 y Fase 6 se ejecutan."*
3. Fase 2: LLM lee `branch_style` del config snapshot y genera nombre de branch
4. Fase 3: LLM genera commit message desde `commit_style` (conocimiento de training, sin tabla local)
5. Fase 5: LLM genera PR title/body desde `pr_title_style` / `pr_body_style`
6. Fase 6: `gh pr merge --auto` sin `--squash`, `--rebase`, `--merge`, `--delete-branch`
7. Eliminar tabla determinística de tipos de commit
8. Eliminar toda mención de trailers (`Co-authored-by`, `Signed-off-by`)
9. Config snapshot en el frontmatter/inputs: agregar `branch_style`, `pr_title_style`, `pr_body_style`, `base_branch`; remover `pr_merge_strategy`

**Out**:
- No modificar loop-subcommand.md (T004)
- No modificar bootstrap-reference.md (T005)

## Estado inicial esperado

- T002 completada: bootstrap JSON incluye los 3 style fields
- `.claude/skills/integrate/SKILL.md` existe

## Criterios de Aceptación

- `grep -F 'omitir esta fase; commitear en el branch actual' .claude/skills/integrate/SKILL.md` → exit 1 (sin matches)
- El archivo no contiene `Co-authored-by` ni `Signed-off-by`
- El archivo no contiene `--squash` ni `--delete-branch` en el comando `gh pr merge`
- Fase 2 no contiene condición `if pr_mode`; el texto dice que corre siempre
- Inputs del skill incluyen `branch_style`, `pr_title_style`, `pr_body_style`; no incluyen `pr_merge_strategy`

## Fuente de verdad

- `.claude/skills/integrate/SKILL.md`
