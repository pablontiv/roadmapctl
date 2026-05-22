---
name: integrate
description: |
  Ejecuta el gitflow per-task: commit, push, branch opcional por scope, PR create y merge.
  Invocar cuando el caller (típicamente /roadmap loop) necesita integrar una task
  completada. También invocable ad-hoc cuando el usuario pide "integrate", "commit
  y push", "crear PR", "mergear PR del scope", o "gitflow".
argument-hint: "task_path=<path> scope=<scope> previous_scope=<scope> repo_path=<path> pr_mode=<bool> commit_style=<style> auto_push=<bool> branch_style=<style> pr_title_style=<style> pr_body_style=<style> autonomy=<mode> [commit_files=<files>] [commit_message=<msg>] [is_last_in_scope=<bool>]"
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
| `task_path` | string | Path relativo de la task recién completada (e.g. `docs/roadmap/O24/T001.md`) |
| `scope` | string | Outcome activo o `direct-tasks` (e.g. `O24-slug`, `direct-tasks`) |
| `previous_scope` | string | Scope anterior; vacío si es la primera task del loop |
| `repo_path` | string | Path absoluto al repo (e.g. `/home/shared/myrepo`) |
| `config` | JSON | Objeto con campos de `roadmapctl bootstrap`: `commit_style`, `auto_push`, `pr_mode`, `branch_style`, `pr_title_style`, `pr_body_style`, `autonomy`, `base_branch` |
| `branch_style` | string | Estilo gitflow del repo. Puede pedir branch por scope (`feat/<scope>`, `<scope>`, `branch per scope`) o direct-push/no feature branch. Si es ambiguo, solo crear branch cuando `pr_mode=true` |
| `pr_title_style` | string | Plantilla para título de PR (generado por LLM desde `commit_style` y scope) |
| `pr_body_style` | string | Plantilla para body de PR (generado por LLM) |
| `base_branch` | string | Rama base para PR (ej: `main`, `master`) |
| `commit_files[]` | string[] | (opcional) Lista de archivos a `git add`; si omitido, usar `-A` con warning |
| `commit_message` | string | (opcional) Override del mensaje de commit; si omitido, derivar desde `commit_style` y `task_path` |
| `is_last_in_scope` | bool | (opcional) `true` si esta es la última task del scope actual; activa merge de PR |

`base_branch` dentro de `config` se detecta así si no viene explícito:

```bash
git -C <repo_path> symbolic-ref refs/remotes/origin/HEAD 2>/dev/null \
  | sed 's@^refs/remotes/origin/@@'
# fallback: main, luego master
```

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

## Fase 1: Gitflow mode y scope change

Antes de ejecutar comandos de branch, clasificar `branch_style` y derivar el modo
efectivo. Si un texto matchea señales de `scope-branch` y `direct-push`, gana
`scope-branch` solo cuando hay patrón explícito de scope.

1. Clasificar `branch_style`:
   - `scope-branch`: contiene un patrón explícito (`<scope>`, `{scope}`) o frases
     como `branch per scope`, `feature branch`, `topic branch`.
   - `direct-push`: contiene señales como `direct`, `push directly`,
     `no feature branch`, `no feature branches`, `no branch`, `no branches` o
     `trunk-based`, y no contiene un patrón explícito de scope.
   - `ambiguous`: vacío o no clasificable.

2. Derivar `branch_mode`, `branch_target` y `effective_pr_mode`:
   - Si `branch_style` es `direct-push`: `branch_mode=direct-push`,
     `branch_target=<base_branch>`, `effective_pr_mode=false`. Emitir warning
     informativo si `pr_mode=true`: `RMC_INTEGRATE_PR_DISABLED_BY_DIRECT_PUSH`.
   - Si `branch_style` es `scope-branch`: `branch_mode=scope-branch` y generar
     `branch_target` desde el patrón indicado.
   - Si `branch_style` es `ambiguous` y `pr_mode=true`: `branch_mode=scope-branch`
     con fallback `feat/<scope>` (e.g. `feat/O24-slug` para Outcome,
     `feat/direct-roadmap-tasks` para `direct-tasks`).
   - Si `branch_style` es `ambiguous` y `pr_mode=false`: `branch_mode=direct-push`,
     `branch_target=<base_branch>`, `effective_pr_mode=false`.

   Para `previous_scope`, derivar `previous_branch_target` con las mismas reglas.
   No hardcodear `feat/<scope>` salvo en el fallback ambiguo con PR.

3. Detectar cambio de scope:
   ```text
   scope_changed = (scope != previous_scope && previous_scope != "")
   ```

Si `scope_changed == true` y `effective_pr_mode == true`:
- El scope anterior puede tener un PR abierto pendiente de cierre. Detectar:
  ```bash
  gh pr list --head <previous_branch_target> --state open --json number,url
  ```
- Si existe PR previo abierto, registrar en diagnostics como informativo (`PR anterior abierto para <previous_scope>: #N`). No cerrarlo automáticamente aquí; el caller decide según `autonomy`.

## Fase 2: Branch setup

Fase 2 corre siempre, pero branch por scope **no** es invariante. `branch_style`
define si se crea branch de scope o si se integra directo en la rama base.
`pr_mode` solo fuerza branch cuando `branch_style` es ambiguo; si `branch_style`
declara direct-push/no feature branch, no crear branch por scope aunque
`pr_mode=true`.

1. Detectar branch actual:
   ```bash
   git -C <repo_path> rev-parse --abbrev-ref HEAD
   ```

2. Si `branch_mode=scope-branch` y el branch actual difiere del target:
   ```bash
   git -C <repo_path> fetch origin <base_branch>
   git -C <repo_path> checkout <base_branch>
   git -C <repo_path> pull --ff-only
   git -C <repo_path> checkout -B <branch_target>
   ```

3. Si `branch_mode=direct-push`:
   - No ejecutar `checkout -B` ni crear branch por scope.
   - Si el branch actual ya es `<base_branch>`, continuar.
   - Si el branch actual difiere y `git status --porcelain` está vacío:
     ```bash
     git -C <repo_path> fetch origin <base_branch>
     git -C <repo_path> checkout <base_branch>
     git -C <repo_path> pull --ff-only
     ```
   - Si el branch actual difiere y hay cambios sin commitear, detenerse con
     `RMC_INTEGRATE_DIRECT_BRANCH_MISMATCH`. No commitear cambios direct-push en
     una rama distinta a `<base_branch>`.

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

## Fase 5: PR (si `effective_pr_mode == true && auto_push == true`)

Detectar si ya existe PR abierto para el scope:

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

## Fase 6: Merge (si `effective_pr_mode == true && is_last_in_scope == true`)

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
| `RMC_INTEGRATE_PR_DISABLED_BY_DIRECT_PUSH` | `pr_mode=true` pero `branch_style` declara direct-push/no feature branch | Continuar sin PR o cambiar `branch_style` a un patrón de branch por scope |
| `RMC_INTEGRATE_PUSH_REJECTED` | Push rechazado por divergencia con remote (remote tiene commits adelante) | Sincronizar con `git pull --rebase origin <branch>` manualmente y reinvocar |
| `RMC_INTEGRATE_GH_AUTH` | `gh auth status` falla | Ejecutar `gh auth login` y reinvocar |
| `RMC_INTEGRATE_NO_GIT` | `git` no encontrado en PATH | Instalar git o verificar entorno |
| `RMC_INTEGRATE_NO_GH` | `gh` no encontrado en PATH | Instalar GitHub CLI (`gh`) o degradar a `pr_mode=false` |

## Verificación al modificar este skill

Ejecutar desde el repo canónico después de cualquier cambio:

```bash
./scripts/sync-roadmap-skill.sh --install --skill integrate
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con pr_mode=false, branch_style="push directly to master; no feature branches", autonomy=until_done, task=docs/roadmap/T020-x.md, scope=direct-tasks. Listar los comandos que correrías, SIN ejecutar git/gh ni modificar archivos.'
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con pr_mode=true, branch_style="", scope=O22-slug, previous_scope=O21-slug, is_last_in_scope=false. Listar comandos SIN ejecutar ni modificar archivos.'
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con pr_mode=true, branch_style="push directly to master; no feature branches", scope=O22-slug, previous_scope=O21-slug, is_last_in_scope=false. Listar comandos SIN ejecutar ni modificar archivos.'
```

Escenario A debe listar `git add`, `git commit`, `git push -u origin master`; NO debe listar `checkout -B`, `gh pr create` ni `gh pr merge`; debe imprimir `INTEGRATE_RESULT` con `pr: null`.

Escenario B debe mencionar detección de scope change, branch `feat/O22-slug`, y plan de detectar PR previo de `O21-slug` (sin ejecutarlo).

Escenario C debe degradar a direct-push, emitir warning `RMC_INTEGRATE_PR_DISABLED_BY_DIRECT_PUSH`, no crear branch por scope y no crear PR.
