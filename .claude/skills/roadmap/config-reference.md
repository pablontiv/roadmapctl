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
| `pr-merge-strategy` | `'squash'` | `<pr-merge-strategy>` |
| `commit-style` | `'conventional'` | `<commit-style>` |
| `auto-push` | `true` | `<auto-push>` |
| `required-code-coverage` | `85.0` | `<required-code-coverage>` |
| `loop-max-tasks` | `0` | `<loop-max-tasks>` |
| `parallel` | `true` | `<parallel>` |
| `autonomy` | `'until_done'` | `<autonomy>` |
| `compact-after-task-commit` | `true` | `<compact-after-task-commit>` |
| `pr-mode` | `false` | `<pr-mode>` |

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
