---
estado: Specified
tipo: task
---
# T002: Reemplazar gitflow inline del loop por invocación a `/integrate`

**Outcome**: [O24 Extraer skill /integrate y achicar el skill /roadmap](README.md)
**Contribuye a**: el loop deja de prescribir `git`/`gh` directamente; toda la prosa de gitflow vive en el skill `/integrate`, no duplicada.

[[blocked_by:./T001-create-integrate-skill.md]]

## Preserva

- INV1: invariante de escritura segura del skill `/roadmap` se mantiene (skill es único writer de archivos roadmap; no se introducen heredocs ni loops shell para gitflow).
  - Verificar: `grep -E "(cat\s*>|<<EOF)" .claude/skills/roadmap/loop-subcommand.md | wc -l` retorna 0 fuera de bloques de documentación.
- INV2: el loop sigue respetando los gates `roadmapctl transition can-start`, `transition start --apply`, `transition complete --apply` antes y después de cada task.
  - Verificar: `grep -c "roadmapctl transition" .claude/skills/roadmap/loop-subcommand.md` retorna ≥ 3.
- INV3: ningún comando `git ` ni `gh ` ejecutivo queda en `loop-subcommand.md` fuera de la sección "Observabilidad de procesos largos" (la cual sí menciona `gh run view` como ejemplo).
  - Verificar: contar líneas con `git ` o `gh ` (espacio después) fuera de "Observabilidad de procesos largos" y "Outcome close check"; debe ser 0.

## Contexto

El paso 9 actual de `loop-subcommand.md` mezcla `roadmapctl transition complete --apply`
(gate de modelo) con `git add` + `git commit` + `git push` + bookkeeping de PR.
La fase 2.5 carga condicionalmente `pr-workflow.md` cuando `pr_mode==true`.

T001 creó el skill `/integrate` que recibe `task_path`, `scope`, `previous_scope`,
`config`, `repo_path`, `commit_files[]?`, `is_last_in_scope?` y devuelve un
bloque `INTEGRATE_RESULT: { commit_hash, branch, pr?, scope_changed, diagnostics[] }`.

El loop debe invocar `/integrate` vía la tool `Skill` después de que el gate
`transition complete --apply` haya pasado, y antes de actualizar UI / compactar.
El skill `/integrate` reemplaza la fase 2.5 entera (branch setup ahora vive en
fase 2 del skill, no en una fase aparte del loop).

`current_scope`, `prs_created` y `base_branch` siguen siendo variables que el
loop mantiene, pero ahora se actualizan con datos que devuelve `/integrate`.

## Alcance

**In**:
1. En `.claude/skills/roadmap/loop-subcommand.md`:
   - Reemplazar el paso 9 actual ("Complete + commit") por un nuevo paso 9 "Integrate" que:
     - Mantenga la pre-condición de ACs/invariantes pasados.
     - Mantenga `roadmapctl transition complete <task> --apply` como gate del modelo.
     - Después del gate, invoque `Skill("integrate", task_path=<task>, scope=<current_scope>, previous_scope=<last_scope>, config=<bootstrap snapshot relevante>, repo_path=<repo>, commit_files=<archivos modificados por la task>, is_last_in_scope=<true|false según ready[] post-completion>)`.
     - Capture `INTEGRATE_RESULT` y actualice `current_scope` si `scope_changed`, `prs_created` si hubo `pr`.
   - Eliminar completamente la sección "Fase 2.5: PR mode" (líneas 104-106).
   - Eliminar la referencia a `pr-workflow.md` de cualquier otra parte del archivo.
2. En `.claude/skills/roadmap/SKILL.md`:
   - Eliminar cualquier mención a `pr-workflow.md` (revisar la sección "Referencia" y "Routing").
3. Borrar `.claude/skills/roadmap/pr-workflow.md` (97 líneas).
4. Re-sincronizar el skill `/roadmap` a user-scope: `scripts/sync-roadmap-skill.sh --install --skill roadmap` debe ser idempotente y reflejar los cambios.

**Out**:
- No modificar contenido del skill `/integrate`. T001 lo creó; este task no lo toca.
- No refactorizar SKILL.md por progressive disclosure. Eso es T003.
- No actualizar docs/README. Eso es T004.

## Estado inicial esperado

- `.claude/skills/integrate/SKILL.md` existe y pasó los AC de T001.
- `.claude/skills/roadmap/loop-subcommand.md` contiene la prosa actual del paso 9 (líneas ~194-198 con `transition complete --apply` + `git add` + `git commit` + `git push` + bookkeeping).
- `.claude/skills/roadmap/pr-workflow.md` existe con 97 líneas.

## Criterios de Aceptación

- AC1: `loop-subcommand.md` contiene un paso 9 reescrito que invoca explícitamente al skill `/integrate`. Verificar: `grep -E "Skill.*integrate" .claude/skills/roadmap/loop-subcommand.md` retorna ≥ 1 línea.
- AC2: ninguna línea ejecutiva con `git ` ni `gh ` permanece en el paso 9 del loop. Verificar: extraer el rango del paso 9 (entre `9. ` y `10. `) y `grep -cE "(^|\s)(git|gh) " <rango>` retorna 0.
- AC3: la sección "Fase 2.5: PR mode" no existe en `loop-subcommand.md`. Verificar: `grep -c "Fase 2.5" .claude/skills/roadmap/loop-subcommand.md` retorna 0.
- AC4: `.claude/skills/roadmap/pr-workflow.md` no existe. Verificar: `test -f .claude/skills/roadmap/pr-workflow.md` retorna 1 (no existe).
- AC5: ninguna referencia a `pr-workflow.md` queda en ningún archivo del skill `/roadmap`. Verificar: `grep -rl "pr-workflow" .claude/skills/roadmap/` retorna vacío.
- AC6: pi headless del loop autónomo confirma que el agente carga `loop-subcommand.md` y planea invocar `/integrate` sin emitir comandos `git`/`gh` directamente. Comando: `./scripts/sync-roadmap-skill.sh --install --skill roadmap && PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/roadmap/SKILL.md --tools read,bash -p 'HEADLESS: el usuario dice "loop autonomo" en este repo. Hacer bootstrap + preflight. Cuando llegues al paso 9 de una task ficticia, listar EXACTAMENTE qué tool/skill invocarías y con qué argumentos, sin ejecutar git/gh.'`. La salida debe mencionar `Skill` con `integrate`, no `Bash` con `git`/`gh`.
- AC7: `roadmapctl check --repo /home/shared/roadmapctl --output json --strict` retorna exit 0 después de los cambios (sanity: no se rompieron archivos del roadmap).

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md` (paso 9, fase 2.5).
- `.claude/skills/roadmap/SKILL.md` (referencias a pr-workflow.md).
- `.claude/skills/roadmap/pr-workflow.md` (a eliminar).
- `.claude/skills/integrate/SKILL.md` (creado en T001, no se modifica aquí).
