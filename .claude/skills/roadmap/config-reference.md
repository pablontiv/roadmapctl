> Cargar para resolver placeholders de filtros o para imprimir el checkpoint de bootstrap.

# Config Reference

## Configuración

Fuente de configuración:

1. `<roadmap-root>/.roadmapctl.toml` vía `roadmapctl bootstrap`.
2. Defaults solo para modo conceptual/no-write.

| Config key | Default | Placeholder |
|------------|---------|-------------|
| `done-statuses` | `['Completed', 'Obsolete']` | `<done-statuses>` |
| `active-statuses` | `['Pending', 'Specified', 'In Progress']` | `<active-statuses>` |
| `status-values.pending` | `'Pending'` | `<status-pending>` |
| `status-values.specified` | `'Specified'` | `<status-specified>` |
| `status-values.in-progress` | `'In Progress'` | `<status-in-progress>` |
| `status-values.completed` | `'Completed'` | `<status-completed>` |
| `status-values.blocked` | `'Blocked'` | `<status-blocked>` |
| `status-values.obsolete` | `'Obsolete'` | `<status-obsolete>` |
| `leaf-filter` | `'isIndex == false'` | `<where-leaf>` |
| `outcome-close-verify` | `[]` | `<outcome-close-cmds>` |
| `required-code-coverage` | `85.0` | `<required-code-coverage>` |
| `loop-max-tasks` | `0` | `<loop-max-tasks>` |
| `parallel` | `true` | `<parallel>` |
| `autonomy` | `'until_done'` | `<autonomy>` |
| `compact-after-task-commit` | `true` | `<compact-after-task-commit>` |
| `gitflow.branch_mode` | `'direct_push'` | `<branch-mode>` |
| `gitflow.pr_create` | `'never'` | `<pr-create>` |
| `gitflow.commit_style` | `'conventional'` | `<commit-style>` |
| `gitflow.auto_push` | `true` | `<auto-push>` |
| `gitflow.base_branch` | `''` | `<base-branch>` |
| `gitflow.branch_style` | `''` | `<branch-style>` |
| `gitflow.pr_title_style` | `''` | `<pr-title-style>` |
| `gitflow.pr_body_style` | `''` | `<pr-body-style>` |

TOML schema canónico:

```toml
# DEPRECATED (still parsed for migration):
# commit_style = "conventional"
# auto_push = true
# pr_mode = false

[gitflow]
base_branch    = "master"
branch_mode    = "direct_push"    # direct_push | scope_branch
branch_style   = ""               # required when branch_mode=scope_branch
pr_create      = "never"          # never | manual | auto
pr_title_style = ""               # required when pr_create=auto
pr_body_style  = ""               # required when pr_create=auto
commit_style   = "conventional"
auto_push      = true
```

## Helpers

- `<where-not-done>`: `not (estado in <done-statuses>)`
- `<where-active>`: `estado in <active-statuses>`
- `<where-leaf>`: valor de `leaf-filter`

## Checkpoint obligatorio

```text
Bootstrap:
  roadmap-root: docs/roadmap
  <where-leaf>:     isIndex == false
  <where-not-done>: not (estado in ["Completed", "Obsolete"])
  <where-active>:   estado in ["Pending", "Specified", "In Progress"]
```
