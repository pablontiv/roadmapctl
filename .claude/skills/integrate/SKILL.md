---
name: integrate
description: |
  Ejecuta el gitflow per-task: commit, push, branch opcional por scope, PR create y merge.
  Invocar cuando el caller (típicamente /roadmap loop) necesita integrar una task
  completada. También invocable ad-hoc cuando el usuario pide "integrate", "commit
  y push", "crear PR", "mergear PR del scope", o "gitflow".
argument-hint: "task_path=<path> scope=<scope> previous_scope=<scope> repo_path=<path> branch_mode=<direct_push|scope_branch> pr_create=<never|manual|auto> commit_style=<style> auto_push=<bool> branch_style=<template> pr_title_style=<style> pr_body_style=<style> base_branch=<branch> autonomy=<mode> [commit_files=<files>] [commit_message=<msg>] [is_last_in_scope=<bool>]"
allowed-tools:
  - Bash
  - Read
  - AskUserQuestion
---

# /integrate — Gitflow Per-Task

Encapsula commit, push, branch opcional por scope, PR creation/merge y cleanup.
Invocable por `/roadmap loop` (paso 9) y ad-hoc por el usuario.

## Inputs

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `task_path` | string | Path relativo de la task recién completada |
| `scope` | string | Outcome activo o `direct-tasks` |
| `previous_scope` | string | Scope anterior; vacío si es la primera task del loop |
| `repo_path` | string | Path absoluto al repo |
| `branch_mode` | enum | `direct_push` \| `scope_branch` — leído de `gitflow.branch_mode` |
| `pr_create` | enum | `never` \| `manual` \| `auto` — leído de `gitflow.pr_create` |
| `commit_style` | string | Leído de `gitflow.commit_style` |
| `auto_push` | bool | Leído de `gitflow.auto_push` |
| `base_branch` | string | Leído de `gitflow.base_branch` — no se auto-detecta |
| `branch_style` | string | Template para nombre de branch (e.g. `feat/{scope}`). Solo relevante cuando `branch_mode=scope_branch` |
| `pr_title_style` | string | Template para título de PR. Requerido cuando `pr_create=auto` |
| `pr_body_style` | string | Template para body de PR. Requerido cuando `pr_create=auto` |
| `commit_files[]` | string[] | (opcional) Lista de archivos a `git add`; si omitido, usar `-A` con warning |
| `commit_message` | string | (opcional) Override del mensaje de commit; si omitido, derivar desde `commit_style` y `task_path` |
| `is_last_in_scope` | bool | (opcional) `true` si es la última task del scope |

## Salida

Al terminar (éxito o fallo parcial documentado), imprimir el bloque:

```
INTEGRATE_RESULT: {
  "commit_hash": "<hash o null>",
  "branch": "<branch efectivo o null>",
  "pr": <número entero o null>,
  "scope_changed": <true|false>,
  "diagnostics": ["<RMC_INTEGRATE_* si hubo error>"]
}
```

El caller consume este bloque para actualizar `prs_created`, `current_scope_branch`
y decidir continuación. En direct-push, no hay `current_scope_branch` por scope;
`branch` solo reporta la rama base efectiva.

## Gate previo

El caller **debe** haber ejecutado `roadmapctl transition complete --apply` antes de invocar este skill. El skill no revalida el estado de la task, pero verifica que haya algo para integrar:

```bash
git -C <repo_path> status --porcelain
```

Si el output está vacío → emitir `RMC_INTEGRATE_NOOP` en diagnostics y retornar con `commit_hash: null`. Si hay cambios staged o unstaged, continuar.

## Fase 1: Mode determination y scope change

El modo se lee directamente de los campos enum — no se clasifica texto libre.

1. Derivar `effective_branch_mode` y `effective_pr_create` directamente desde los inputs:
   - `effective_branch_mode = branch_mode` (siempre explícito)
   - `effective_pr_create = pr_create` (siempre explícito)

2. Generar `branch_target`:
   - Si `effective_branch_mode = "direct_push"`: `branch_target = base_branch`
   - Si `effective_branch_mode = "scope_branch"`: LLM genera `branch_target` sustituyendo
     `{scope}` o `<scope>` en `branch_style` con el valor de `scope` (slugificado).
     Si `branch_style` está vacío o no contiene patrón de scope: emitir
     `RMC_INTEGRATE_BRANCH_STYLE_MISSING` y detener.

3. Detectar cambio de scope:
   ```
   scope_changed = (scope != previous_scope && previous_scope != "")
   ```
   Si `scope_changed == true` y `effective_pr_create != "never"`:
   ```bash
   gh pr list --head <previous_branch_target> --state open --json number,url
   ```
   Si existe PR previo abierto: registrar en diagnostics como informativo.

## Fase 2: Branch setup

Fase 2 siempre sincroniza `base_branch` antes de commitear, garantizando que HEAD queda
en la posición correcta independientemente del estado previo.

1. Detectar branch actual:
   ```bash
   git -C <repo_path> rev-parse --abbrev-ref HEAD
   ```

2. Sincronizar `base_branch` siempre:
   ```bash
   git -C <repo_path> fetch origin <base_branch>
   git -C <repo_path> checkout <base_branch>
   git -C <repo_path> pull --ff-only
   ```
   Si el `checkout` falla porque hay cambios locales en conflicto: emitir
   `RMC_INTEGRATE_CHECKOUT_BLOCKED` y detener. El operador debe resolver el estado
   antes de reinvocar.

3. Crear branch de scope (solo si `effective_branch_mode == "scope_branch"`):
   - Si `branch_style` está vacío o no contiene un patrón de scope explícito
     (`<scope>`, `{scope}`): emitir `RMC_INTEGRATE_BRANCH_STYLE_MISSING` y detener.
     Recovery: correr `/roadmap bootstrap` para popular `[gitflow].branch_style`.
   - Si `branch_style` provee un patrón de scope válido: LLM genera `branch_target`
     desde el patrón y ejecuta:
     ```bash
     git -C <repo_path> checkout -B <branch_target>
     ```

## Fase 3: Commit

```bash
git -C <repo_path> add <commit_files>
# si commit_files omitido:
git -C <repo_path> add -A   # warning: staging todo
```

Derivar mensaje de commit según `commit_style`:

- LLM genera commit message desde `commit_style` usando conocimiento de training, sin tabla local.
- Si `commit_style` es `conventional`: LLM genera `<type>(<scope-corto>): <título-tarea>` analizando el contenido de la task (docs, chore, refactor, feat, fix) según convención.
- El `<scope-corto>` es el código del Outcome (`O24`, etc.) o `direct` para `direct-tasks`.
- Override explícito vía `commit_message`: usar el texto directo. El override por `commit_message` siempre tiene precedencia.

```bash
git -C <repo_path> commit -m "$(cat <<'EOF'
<mensaje derivado>
EOF
)"
```

Capturar hash:

```bash
commit_hash=$(git -C <repo_path> rev-parse HEAD)
```

## Fase 4: Push (si `auto_push == true`)

```bash
git -C <repo_path> push -u origin <branch_target>
```

Si el push es rechazado (exit ≠ 0), diferenciar la causa según stderr:

### Rechazo por pre-push hook

Si stderr contiene las palabras `hook` o `pre-push`:

- **Prohibido** usar `--no-verify` o `--force-with-lease` automáticamente. El hook es del proyecto; el skill respeta su decisión.
- Intentar parsear stderr buscando paths o nombres de archivo mencionados como requeridos. Heurística: líneas que contienen un path-like (segmento con `/` o extensión de archivo) junto a un verbo de coordinación (`update`, `change`, `required`, `not updated`, `missing`).
- Si se detectan paths candidatos: emitir diagnostic informativo listando esos paths y proponer al caller un commit complementario que los cubra dentro del mismo push range. Permitir un reintento de push después de ese commit.
- Si no se puede parsear nada estructurado: emitir `INTEGRATE_HOOK_REJECTED` con el mensaje literal del hook y detenerse. El operador decide cómo proceder.

### Rechazo por divergencia con remote

Si stderr NO contiene `hook` ni `pre-push` (remote tiene commits adelante):

- `manual` / `supervised`: emitir `RMC_INTEGRATE_PUSH_REJECTED`, reportar al usuario y detenerse. No reintentar.
- `until_done`: intentar rebase y reintentar una vez:
  ```bash
  git -C <repo_path> pull --rebase origin <branch_target>
  git -C <repo_path> push -u origin <branch_target>
  ```
  Si aún falla: emitir `RMC_INTEGRATE_PUSH_REJECTED` y detenerse.

Si `auto_push == false`, omitir push. `branch` en `INTEGRATE_RESULT` refleja
`branch_target`; en direct-push es `<base_branch>`, no un branch por scope.

## Fase 5: PR (si `effective_branch_mode = "scope_branch"`)

### pr_create = "never"

No ejecutar ningún comando PR. Continuar directamente a Fase 6 (que no hace nada en este caso).

### pr_create = "manual"

Detectar si ya existe PR abierto para el scope:

```bash
gh pr list --head <branch_target> --state open --json number,url
```

Si no existe, imprimir el comando sugerido (NO ejecutarlo):

```bash
echo "PR sugerido (ejecutar manualmente):"
echo "gh pr create \\"
echo "  --base <base_branch> \\"
echo "  --head <branch_target> \\"
echo "  --title \"<LLM-generated desde pr_title_style>\" \\"
echo "  --body \"<LLM-generated desde pr_body_style>\""
```

Registrar `pr: null` en `INTEGRATE_RESULT` (no se creó PR automáticamente).

### pr_create = "auto"

Detectar si ya existe PR abierto:

```bash
gh pr list --head <branch_target> --state open --json number,url
```

Si no existe, crear:

```bash
gh pr create \
  --base <base_branch> \
  --head <branch_target> \
  --title "<LLM-generated desde pr_title_style>" \
  --body "<LLM-generated desde pr_body_style>"
```

LLM genera PR title desde `pr_title_style` y body desde `pr_body_style`, incluyendo contexto de scope y lista de tasks completadas.

Registrar número de PR en `INTEGRATE_RESULT.pr`.

Si `gh` no está disponible o `gh auth status` falla:
- `manual`: emitir `RMC_INTEGRATE_GH_AUTH` o `RMC_INTEGRATE_NO_GH`, preguntar si continuar sin PR.
- `supervised` / `until_done`: degradar a modo sin PR; advertir; continuar.

## Fase 6: Merge (si `effective_pr_create == "auto" && is_last_in_scope == true`)

Por `autonomy`:

- `manual`: preguntar al usuario si mergear ahora o dejar abierto.
- `supervised`: preguntar antes de mergear.
- `until_done`: ejecutar auto-merge si branch protection lo permite:
  ```bash
  gh pr merge <pr_number> --auto
  ```
  GitHub decide la merge strategy desde sus branch protection rules.

Post-merge cleanup:

```bash
git -C <repo_path> checkout <base_branch>
git -C <repo_path> pull --ff-only
```

Registrar `{number, url, scope, status: "merged"}` para el caller.

## Errores comunes

| ID | Causa | Acción recomendada |
|----|-------|-------------------|
| `RMC_INTEGRATE_NOOP` | `git status --porcelain` vacío; nada que commitear | Verificar que `roadmapctl transition complete --apply` fue ejecutado y los cambios fueron staged antes de invocar integrate |
| `INTEGRATE_HOOK_REJECTED` | Pre-push hook del proyecto exige cambios coordinados que no fueron incluidos en el commit | Leer el mensaje literal del hook, identificar los paths requeridos, crear un commit complementario sobre esos paths y reinvocar integrate |
| `RMC_INTEGRATE_DIRECT_BRANCH_MISMATCH` | Direct-push configurado pero el worktree con cambios está en una rama distinta a la base | Cambiar a la rama base antes de completar la task o mover los cambios explícitamente |
| `RMC_INTEGRATE_PUSH_REJECTED` | Push rechazado por divergencia con remote (remote tiene commits adelante) | Sincronizar con `git pull --rebase origin <branch>` manualmente y reinvocar |
| `RMC_INTEGRATE_GH_AUTH` | `gh auth status` falla | Ejecutar `gh auth login` y reinvocar |
| `RMC_INTEGRATE_NO_GIT` | `git` no encontrado en PATH | Instalar git o verificar entorno |
| `RMC_INTEGRATE_NO_GH` | `gh` no encontrado en PATH | Instalar GitHub CLI (`gh`) o degradar a `pr_create=never` |
| `RMC_INTEGRATE_BRANCH_STYLE_MISSING` | `branch_mode=scope_branch` pero `branch_style` vacío o no contiene patrón de scope | Correr `/roadmap bootstrap` y popular `[gitflow].branch_style` en `.roadmapctl.toml` |
| `RMC_INTEGRATE_CHECKOUT_BLOCKED` | Checkout de `base_branch` fallido por cambios locales en conflicto | Resolver state del worktree (commit, stash o discard) y reinvocar |
| `RMC_INTEGRATE_BRANCH_MODE_MISSING` | `branch_mode` no fue provisto | Verificar que bootstrap JSON fue leído y `gitflow.branch_mode` fue pasado al skill |
| `RMC_INTEGRATE_BASE_BRANCH_MISSING` | `base_branch` vacío | Configurar `[gitflow].base_branch` en `.roadmapctl.toml` y re-ejecutar bootstrap |

## Verificación al modificar este skill

Ejecutar desde el repo canónico después de cualquier cambio:

```bash
./scripts/sync-roadmap-skill.sh --install --skill integrate

# Scenario A: direct_push
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con branch_mode=direct_push, pr_create=never, base_branch=master, commit_style=conventional, auto_push=true, autonomy=until_done, task=docs/roadmap/T020-x.md, scope=direct-tasks. Listar los comandos que correrías, SIN ejecutar git/gh ni modificar archivos.'

# Scenario B: scope_branch + pr_create=auto
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con branch_mode=scope_branch, pr_create=auto, base_branch=main, branch_style="feat/{scope}", scope=O22-slug, previous_scope=O21-slug, is_last_in_scope=false. Listar comandos SIN ejecutar ni modificar archivos.'

# Scenario C: scope_branch + pr_create=manual
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con branch_mode=scope_branch, pr_create=manual, base_branch=master, branch_style="feat/{scope}", scope=O22-slug, previous_scope="", is_last_in_scope=false. Listar comandos SIN ejecutar ni modificar archivos. Verificar que NO ejecuta gh pr create sino que imprime el comando sugerido.'

# Scenario D: scope_branch + branch_style vacío
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con branch_mode=scope_branch, pr_create=auto, base_branch=master, branch_style="", scope=O22-slug. Listar comandos SIN ejecutar. Verificar que emite RMC_INTEGRATE_BRANCH_STYLE_MISSING y detiene.'
```

Escenario A debe listar `git fetch`, `git checkout master`, `git pull --ff-only`, `git add`, `git commit`, `git push -u origin master`; NO debe listar `checkout -B`, `gh pr create` ni `gh pr merge`; debe imprimir `INTEGRATE_RESULT` con `branch: "master"` y `pr: null`.

Escenario B debe listar `git checkout -B feat/O22-slug`, `gh pr create`; PR registrado.

Escenario C debe listar branch creado, push, imprime `gh pr create ...` sin ejecutar; `INTEGRATE_RESULT.pr = null`.

Escenario D debe contener `RMC_INTEGRATE_BRANCH_STYLE_MISSING`, no listar `checkout -B` ni `gh pr create`.
