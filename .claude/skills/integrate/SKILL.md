---
name: integrate
description: |
  Ejecuta el gitflow per-task: commit, push, branch por scope, PR create y merge.
  Invocar cuando el caller (típicamente /roadmap loop) necesita integrar una task
  completada. También invocable ad-hoc cuando el usuario pide "integrate", "commit
  y push", "crear PR", "mergear PR del scope", o "gitflow".
argument-hint: "task_path=<path> scope=<scope> previous_scope=<scope> repo_path=<path> pr_mode=<bool> commit_style=<style> auto_push=<bool> pr_merge_strategy=<strategy> autonomy=<mode> [commit_files=<files>] [commit_message=<msg>] [is_last_in_scope=<bool>]"
allowed-tools:
  - Bash
  - Read
  - AskUserQuestion
---

# /integrate — Gitflow Per-Task

Encapsula commit, push, branch por scope, PR creation/merge y cleanup.
Invocable por `/roadmap loop` (paso 9) y ad-hoc por el usuario.

## Inputs

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `task_path` | string | Path relativo de la task recién completada (e.g. `docs/roadmap/O24/T001.md`) |
| `scope` | string | Outcome activo o `direct-tasks` (e.g. `O24-slug`, `direct-tasks`) |
| `previous_scope` | string | Scope anterior; vacío si es la primera task del loop |
| `repo_path` | string | Path absoluto al repo (e.g. `/home/shared/myrepo`) |
| `config` | JSON | Objeto con seis campos de `roadmapctl bootstrap`: `commit_style`, `auto_push`, `pr_mode`, `pr_merge_strategy`, `autonomy`, `base_branch` |
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
  "branch": "<nombre de branch o null>",
  "pr": <número entero o null>,
  "scope_changed": <true|false>,
  "diagnostics": ["<RMC_INTEGRATE_* si hubo error>"]
}
```

El caller consume este bloque para actualizar `prs_created`, `current_scope_branch` y decidir continuación.

## Gate previo

El caller **debe** haber ejecutado `roadmapctl transition complete --apply` antes de invocar este skill. El skill no revalida el estado de la task, pero verifica que haya algo para integrar:

```bash
git -C <repo_path> status --porcelain
```

Si el output está vacío → emitir `RMC_INTEGRATE_NOOP` en diagnostics y retornar con `commit_hash: null`. Si hay cambios staged o unstaged, continuar.

## Fase 1: Scope change

```
scope_changed = (scope != previous_scope && previous_scope != "")
```

Si `scope_changed == true` y `pr_mode == true`:
- El scope anterior puede tener un PR abierto pendiente de cierre. Detectar:
  ```bash
  gh pr list --head feat/<previous_scope> --state open --json number,url
  ```
- Si existe PR previo abierto, registrar en diagnostics como informativo (`PR anterior abierto para <previous_scope>: #N`). No cerrarlo automáticamente aquí; el caller decide según `autonomy`.

## Fase 2: Branch setup (si `pr_mode == true`)

1. Detectar branch actual:
   ```bash
   git -C <repo_path> rev-parse --abbrev-ref HEAD
   ```

2. Derivar branch target:
   - Scope = Outcome: `feat/<scope>` (e.g. `feat/O24-slug`)
   - Scope = `direct-tasks`: `feat/direct-roadmap-tasks`

3. Si el branch actual difiere del target:
   ```bash
   git -C <repo_path> fetch origin <base_branch>
   git -C <repo_path> checkout <base_branch>
   git -C <repo_path> pull --ff-only
   git -C <repo_path> checkout -B feat/<scope>
   ```

Si `pr_mode == false`, omitir esta fase; commitear en el branch actual.

## Fase 3: Commit

```bash
git -C <repo_path> add <commit_files>
# si commit_files omitido:
git -C <repo_path> add -A   # warning: staging todo
```

Derivar mensaje de commit según `commit_style`:

- `conventional`: `<type>(<scope-corto>): <título-tarea>` — el `type` se infiere del prefijo de la task (`feat`, `fix`, `docs`, `chore`, etc.); el scope-corto es el código del Outcome (e.g. `O24`) o `direct`.
- Cualquier otro valor o override explícito de `commit_message`: usar el texto directo.

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
git -C <repo_path> push -u origin <branch>
```

Si el push es rechazado (exit ≠ 0):

- `manual` / `supervised`: emitir `RMC_INTEGRATE_PUSH_REJECTED`, reportar al usuario y detenerse. No reintentar.
- `until_done`: intentar rebase y reintentar una vez:
  ```bash
  git -C <repo_path> pull --rebase origin <branch>
  git -C <repo_path> push -u origin <branch>
  ```
  Si aún falla: emitir `RMC_INTEGRATE_PUSH_REJECTED` y detenerse.

Si `auto_push == false`, omitir push. `branch` en `INTEGRATE_RESULT` refleja el branch local.

## Fase 5: PR (si `pr_mode == true && auto_push == true`)

Detectar si ya existe PR abierto para el scope:

```bash
gh pr list --head feat/<scope> --state open --json number,url
```

Si no existe, crear:

```bash
gh pr create \
  --base <base_branch> \
  --head feat/<scope> \
  --title "<commit_style-title para el scope>" \
  --body "$(cat <<'EOF'
## Scope
<scope>

## Cambios
- Lista de tasks completadas con sus commits

## Verificación
- ACs: passed
- Invariantes preservadas
EOF
)"
```

Registrar número de PR en `INTEGRATE_RESULT.pr`.

Si `gh` no está disponible o `gh auth status` falla:
- `manual`: emitir `RMC_INTEGRATE_GH_AUTH` o `RMC_INTEGRATE_NO_GH`, preguntar si continuar sin PR.
- `supervised` / `until_done`: degradar a modo sin PR; advertir; continuar.

## Fase 6: Merge (si `pr_mode == true && is_last_in_scope == true`)

Por `autonomy`:

- `manual`: preguntar al usuario si mergear ahora o dejar abierto.
- `supervised`: preguntar antes de mergear.
- `until_done`: ejecutar auto-merge si branch protection lo permite:
  ```bash
  gh pr merge <pr_number> --auto --<pr_merge_strategy> --delete-branch
  ```

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
| `RMC_INTEGRATE_PUSH_REJECTED` | Push rechazado (remote tiene commits adelante) | Sincronizar con `git pull --rebase origin <branch>` manualmente y reinvocar |
| `RMC_INTEGRATE_GH_AUTH` | `gh auth status` falla | Ejecutar `gh auth login` y reinvocar |
| `RMC_INTEGRATE_NO_GIT` | `git` no encontrado en PATH | Instalar git o verificar entorno |
| `RMC_INTEGRATE_NO_GH` | `gh` no encontrado en PATH | Instalar GitHub CLI (`gh`) o degradar a `pr_mode=false` |

## Verificación al modificar este skill

Ejecutar desde el repo canónico después de cualquier cambio:

```bash
./scripts/sync-roadmap-skill.sh --install --skill integrate
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con pr_mode=false, autonomy=until_done, task=docs/roadmap/T020-x.md, scope=direct-tasks. Listar los comandos que correrías, SIN ejecutar git/gh ni modificar archivos.'
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/integrate/SKILL.md --tools read,bash \
  -p 'HEADLESS: invocar integrate con pr_mode=true, scope=O22-slug, previous_scope=O21-slug, is_last_in_scope=false. Listar comandos SIN ejecutar ni modificar archivos.'
```

Escenario A debe listar `git add`, `git commit`, `git push`; NO debe listar `gh pr create` ni `gh pr merge`; debe imprimir `INTEGRATE_RESULT` con `pr: null`.

Escenario B debe mencionar detección de scope change, branch `feat/O22-slug`, y plan de detectar PR previo de `O21-slug` (sin ejecutarlo).
