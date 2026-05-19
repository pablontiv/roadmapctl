---
estado: Specified
tipo: task
---
# T001: Crear skill `/integrate` con prosa de gitflow per-task

**Outcome**: [O24 Extraer skill /integrate y achicar el skill /roadmap](README.md)
**Contribuye a**: existencia de un skill aparte que encapsula commit, push, branch por scope, PR creation/merge y cleanup, invocable por el loop y ad-hoc.

## Preserva

- INV1: el skill no reimplementa `git`/`gh` en Go ni en otro binario; solo
  prescribe los comandos en orden.
  - Verificar: `grep -E "(git |gh )" .claude/skills/integrate/SKILL.md | wc -l` retorna > 0; no se crea ningún paquete `internal/gitops` ni similar.
- INV2: el skill convive con el skill `/roadmap` en este repo siguiendo el mismo
  patrón que `retrospective`.
  - Verificar: `ls .claude/skills/integrate/SKILL.md` existe; `scripts/sync-roadmap-skill.sh --install --skill integrate` instala sin error.
- INV3: el skill declara un contrato de salida parseable que el caller
  consume para bookkeeping (`prs_created`, `current_scope`).
  - Verificar: `grep -c "INTEGRATE_RESULT" .claude/skills/integrate/SKILL.md` ≥ 1; el bloque documenta `commit_hash`, `branch`, `pr`, `scope_changed`, `diagnostics`.

## Contexto

Hoy `loop-subcommand.md` paso 9 (~30 líneas) y `pr-workflow.md` (97 líneas)
prescriben gitflow per-task: commit + push + opcional creación/merge de PR. La
prescripción depende de seis campos de `roadmapctl bootstrap`: `commit_style`,
`auto_push`, `pr_mode`, `pr_merge_strategy`, `autonomy`, y `base_branch` (este
último derivado de `git symbolic-ref refs/remotes/origin/HEAD`).

Variables canónicas que el caller mantiene hoy en el loop: `prs_created`,
`current_scope_branch`, `base_branch`. El nuevo skill las recibe como input y
devuelve actualizaciones via su contrato de salida.

Scope = Outcome activo o `direct-tasks`. Branch = `feat/<scope>`. Hay un PR
por scope, no por task; el skill debe detectar si el scope ya tiene PR abierto
para no crear duplicados.

Comando de creación de PR usa heredoc estilo `gh pr create --title ...
--body "$(cat <<'EOF' ... EOF)"` y `gh pr merge <n> --auto --<strategy> --delete-branch`
para auto-merge.

Frontmatter del skill: `name: integrate`, `description` con triggers (qué
condiciones lo invocan), `argument-hint` describiendo inputs canónicos,
`allowed-tools: [Bash, Read, AskUserQuestion]`.

## Alcance

**In**:
1. Crear `.claude/skills/integrate/` con `SKILL.md` único conteniendo:
   - Frontmatter (`name`, `description`, `argument-hint`, `allowed-tools`).
   - Sección "Inputs" tabulando `task_path`, `scope`, `previous_scope`, `config` (JSON con los seis campos), `repo_path`, `commit_files[]?`, `commit_message?`, `is_last_in_scope?`.
   - Sección "Salida" documentando el bloque `INTEGRATE_RESULT: { ... }` que el skill debe imprimir al terminar.
   - Sección "Gate previo" requiriendo que el caller ya haya ejecutado `roadmapctl transition complete --apply` y verificando con `git status --porcelain` que hay algo para integrar.
   - Fase 1 "Scope change": comparar `scope` con `previous_scope` y setear `scope_changed`.
   - Fase 2 "Branch setup (si `pr_mode==true`)": `git rev-parse --abbrev-ref HEAD`; si difiere de `feat/<scope>`, `git fetch origin <base_branch>`, `git checkout <base_branch>`, `git pull --ff-only`, `git checkout -B feat/<scope>`.
   - Fase 3 "Commit": `git add <commit_files>` (o `-A` fallback con warning), `git commit -m "<mensaje>"` con mensaje derivado de `commit_style` (`conventional` → `<type>(<scope-corto>): <título>`) o override de `commit_message`. Capturar `commit_hash = git rev-parse HEAD`.
   - Fase 4 "Push": si `auto_push==true`, `git push -u origin <branch>`. Manejo de push rejected según `autonomy` (manual/supervised parar; until_done rebase y reintentar 1 vez).
   - Fase 5 "PR (si `pr_mode==true && auto_push==true`)": `gh pr list --head feat/<scope> --state open` para detectar; `gh pr create --base <base_branch> --head feat/<scope>` con body heredoc si no existe.
   - Fase 6 "Merge (si `pr_mode==true && is_last_in_scope==true`)": `gh pr merge <n> --auto --<pr_merge_strategy> --delete-branch` según autonomy (manual/supervised preguntan; until_done auto-merge). Post-merge: `git checkout <base_branch> && git pull --ff-only`.
   - Sección "Errores comunes" tabulando diagnostics `RMC_INTEGRATE_NOOP`, `RMC_INTEGRATE_PUSH_REJECTED`, `RMC_INTEGRATE_GH_AUTH`, `RMC_INTEGRATE_NO_GIT`, `RMC_INTEGRATE_NO_GH` con causa y acción.
   - Sección "Verificación al modificar este skill" con los dos comandos `pi --skill` de los AC4-AC5.
2. Verificar que `scripts/sync-roadmap-skill.sh --install --skill integrate` sincroniza el nuevo skill a `~/.claude/skills/integrate/` sin tocar otros skills.

**Out**:
- No modificar `loop-subcommand.md` ni `pr-workflow.md` ni `SKILL.md` del skill `/roadmap`. Esos cambios son de T002.
- No modificar `docs/roadmap-skill-integration.md` ni `README.md`. Eso es de T004.
- No escribir código Go.

## Estado inicial esperado

- `.claude/skills/integrate/` no existe en `git ls-files`.
- `scripts/sync-roadmap-skill.sh` ya soporta `--skill NAME` (verificado).
- Bootstrap del repo retorna `commit_style=conventional`, `auto_push=true`, `pr_mode=false`, `autonomy=until_done`, `pr_merge_strategy=squash`.

## Criterios de Aceptación

- AC1: `test -f .claude/skills/integrate/SKILL.md` retorna 0.
- AC2: frontmatter del skill contiene exactamente `name: integrate`, una `description` con la palabra "integrate" y al menos un trigger explícito, un `argument-hint` no vacío y `allowed-tools` con `Bash` y `Read`. Verificar parseando frontmatter con `head -30 .claude/skills/integrate/SKILL.md`.
- AC3: `grep -c "INTEGRATE_RESULT" .claude/skills/integrate/SKILL.md` ≥ 1; el bloque referencia los cinco campos `commit_hash`, `branch`, `pr`, `scope_changed`, `diagnostics`.
- AC4: el skill describe explícitamente las 6 fases (Scope change, Branch setup, Commit, Push, PR create, Merge) y al menos los 5 diagnostics `RMC_INTEGRATE_*` listados arriba. Verificar con `grep -c "## Fase" .claude/skills/integrate/SKILL.md` ≥ 6 y `grep -c "RMC_INTEGRATE_" .claude/skills/integrate/SKILL.md` ≥ 5.
- AC5: `./scripts/sync-roadmap-skill.sh --install --skill integrate` retorna exit 0; `test -f ~/.claude/skills/integrate/SKILL.md` retorna 0; `diff -r .claude/skills/integrate ~/.claude/skills/integrate` no muestra diferencias.
- AC6: pi headless scenario A pasa — `PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash -p 'HEADLESS: invocar integrate con pr_mode=false, autonomy=until_done, task=docs/roadmap/T020-x.md, scope=direct-tasks. Listar los comandos que correrías, SIN ejecutar git/gh ni modificar archivos.'` — la salida del agente debe listar `git add`, `git commit -m`, y `git push`; NO debe listar `gh pr create` ni `gh pr merge`; el skill debe imprimir un bloque `INTEGRATE_RESULT` con `pr: null`.
- AC7: pi headless scenario B pasa — invocar el skill con `pr_mode=true, scope=O22-x, previous_scope=O21-x` y verificar que la salida menciona detección de scope change, branch `feat/O22-x`, y plan de cerrar PR previo (sin ejecutarlo). No debe modificar archivos.

## Fuente de verdad

- `.claude/skills/roadmap/pr-workflow.md` (97 líneas) — fuente original de la lógica de branch/PR que se traslada.
- `.claude/skills/roadmap/loop-subcommand.md` líneas 194-198 (paso 9 actual) — fuente de la lógica de commit/push.
- `.claude/skills/retrospective/` — patrón estructural de skill aparte que convive con `/roadmap` en este repo.
- `/home/pones/.claude/plans/tiene-sentido-que-roadmap-modular-stroustrup.md` — plan aprobado.
- `scripts/sync-roadmap-skill.sh` — ya soporta `--skill NAME`.
