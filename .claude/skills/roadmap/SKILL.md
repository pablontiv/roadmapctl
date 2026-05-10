---
source: pablontiv/praxis
name: roadmap
description: |
  Usar cuando el usuario sabe QUÉ construir y necesita planificar CÓMO —
  descomponiendo trabajo en Outcomes opcionales y Tasks ejecutables por agentes,
  con criterios de aceptación, dependencias y validación. También para ver
  progreso, trabajo pendiente o ejecutar tasks en secuencia. Usar si el usuario
  describe capacidades a construir, pregunta "cómo estructuro esto", lista
  requerimientos, quiere ver pendientes/progreso, o dice "next task",
  "planificar" o "descomponer".
argument-hint: "<texto libre> | [pending|loop|plan] [args]"
allowed-tools:
  - Write
  - Read
  - Grep
  - Glob
  - Bash
  - TaskCreate
  - TaskList
  - TaskUpdate
  - TaskGet
  - Skill
  - AskUserQuestion
  - ExitPlanMode
  - Agent
effort: xhigh
execution-model: sonnet
worktree-per-outcome: false
parallel-independent-tasks: false
hooks:
  Stop:
    - type: agent
      prompt: "Verify that the critic agent was invoked during this skill execution and its evaluation passed. If no critic evaluation occurred and the skill produced artifacts, return {ok: false, reason: 'Critic evaluation was skipped'}. If work is still in progress, return {ok: true}."
      timeout: 60
---

# /roadmap — Planificación AI-Native Simple

Modelo canónico:

```text
Outcome/Objetivo  (opcional)
└── Task          (unidad ejecutable)
```

Para trabajo chico, usar solo tasks. El skill produce únicamente Outcomes y Tasks.

## Invariante de materialización

Cuando el usuario pide crear, generar o materializar tareas, el resultado NO puede
ser un único archivo resumen tipo `*-tasks.md`.

Materializar tareas significa crear archivos canónicos:

```text
<roadmap-root>/OXX-slug/README.md
<roadmap-root>/OXX-slug/TXXX-task.md
```

o tasks directas:

```text
<roadmap-root>/TXXX-task.md
```

Si no se puede crear esa estructura, detenerse y explicar el bloqueo. No hacer
fallback a markdown libre.

## Invariante de escritura segura

Además del criterio canónico, el skill no puede crear/reescribir archivos roadmap manualmente ni hacer dumps multi-file fuera de `roadmapctl`.

Prohibido:

- `bash` con múltiples heredocs/cats dirigidos a archivos distintos.
- loops de shell que llamen `rootline new` o escritura en varios paths.
- `write`/`edit` directo para crear o reescribir archivos canónicos del roadmap.

Permitido: `roadmapctl materialize --plan <plan-json> --apply` o `roadmapctl materialize --changes <dry-run-json> --apply` puede aplicar múltiples archivos en una ejecución, porque roadmapctl owns canonical writes, per-file diagnostics, validation, ordering and postcheck.

## Bootstrap mínimo y autosuficiente

Ejecutar solo el bootstrap necesario para despachar el flujo. No leer documentación adicional para resolver configuración, schema, helpers, pending, next o decision; `roadmapctl` es la API autosuficiente para esos datos.

### Paso 0: Detectar modo

```bash
test -d .git
```

- Sí → single-repo mode.
- No → workspace mode.

### Fuente primaria de contexto

`roadmapctl context` resuelve configuración efectiva, schema, helpers y comportamiento operacional. `.roadmapctl.toml` dentro de `<roadmap-root>/` es la configuración canónica; cualquier config local legacy es solo input de migración gestionado opacamente por roadmapctl.

Gate inicial para flujos implementados:

```bash
command -v roadmapctl
```

Ejecutar para cada repo objetivo:

```bash
roadmapctl context --repo <repo-path> --roadmap-root <roadmap-root-si-se-conoce> --output json
```

Usar el JSON devuelto como fuente de verdad para:

- `<repo-path>` = `root`
- `<abs-roadmap-root>` = `roadmap_root`
- `<roadmap-root>` = path relativo desde `root` a `roadmap_root`
- `<where-leaf>` = `helpers.where_leaf`
- `<where-not-done>` = `helpers.where_not_done`
- `<where-active>` = `helpers.where_active`
- schema/status/config operacional = campos `schema`, `status_values`, `done_statuses`, `active_statuses`, `outcome_close_verify`, `pr_merge_strategy`, `commit_style`, `auto_push`, `required_code_coverage`, `loop_max_tasks`, `parallel`, `autonomy`, `compact_after_task_commit`, `pr_mode` y cualquier campo adicional expuesto por `roadmapctl context`

`roadmapctl doctor` y `roadmapctl check` no forman parte del bootstrap read-only. Ejecutarlos solo antes de escribir, mutar, ejecutar tasks o declarar validez del roadmap, y como postcheck después de materializar o mutar.

Si `roadmapctl context` falla o `roadmapctl` no existe:

- Para flujos implementados read-only, writes, mutaciones, ejecución o declaraciones de validez: detenerse; no fallback.
- Para planificación conceptual sin writes/mutaciones/ejecución/validez: se permite usar defaults explícitos solo como ayuda conceptual, dejando claro que los guards faltan para materializar/ejecutar. El skill no migra ni parsea legacy para flujos implementados.

### Workspace mode

1. Escanear subdirectorios inmediatos con `.git` + config roadmap (`<roadmap-root>/.roadmapctl.toml`; señales legacy solo habilitan que roadmapctl migre, no que el skill las lea).
2. Para cada repo, ejecutar `roadmapctl context` si está disponible y calcular helpers desde su JSON.
3. Imprimir checkpoint con repos detectados.

### Single-repo mode

1. Resolver repo actual.
2. Ejecutar `roadmapctl context` si está disponible.
3. Si context no está disponible y el flujo es conceptual/no-write, preguntar dónde vive el roadmap o usar defaults explícitos marcados como no verificados.
4. Imprimir checkpoint desde JSON de `roadmapctl context` o desde defaults conceptuales explícitamente marcados.

Template mínimo:

```yaml
---
roadmap-root: # preguntar al usuario
done-statuses: ['Completed', 'Obsolete']
active-statuses: ['Pending', 'Specified', 'In Progress']
status-values:
  pending: 'Pending'
  specified: 'Specified'
  in-progress: 'In Progress'
  completed: 'Completed'
  blocked: 'Blocked'
  obsolete: 'Obsolete'
leaf-filter: 'isIndex == false'
outcome-close-verify: []
pr-merge-strategy: 'squash'
commit-style: 'conventional'
auto-push: true
---
```

## Configuración

Fuente de configuración:

1. `<roadmap-root>/.roadmapctl.toml` vía `roadmapctl context`.
2. La config local legacy es solo input de migración para roadmapctl, no fuente durable que el skill deba parsear en flujos implementados.
3. defaults solo para modo conceptual/no-write.

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

Helpers:

- `<where-not-done>`: `not (estado in <done-statuses>)`
- `<where-active>`: `estado in <active-statuses>`
- `<where-leaf>`: valor de `leaf-filter`

Checkpoint obligatorio:

```text
Bootstrap:
  roadmap-root: docs/roadmap
  <where-leaf>:     isIndex == false
  <where-not-done>: not (estado in ["Completed", "Obsolete"])
  <where-active>:   estado in ["Pending", "Specified", "In Progress"]
```

## Validación de configuración

`roadmapctl context`, `roadmapctl doctor` y `roadmapctl check` son los únicos validadores de configuración para flujos implementados. El skill no debe leer ni validar archivos de config legacy directamente.

No ejecutar `rootline describe` como paso normal de validación de configuración: el schema y los status efectivos vienen en el JSON de `roadmapctl context`, y `roadmapctl doctor/check` validan los flujos que escriben, mutan, ejecutan o declaran validez.

## Dependencias CLI

### rootline

Rootline es dependencia interna de `roadmapctl` para validar y consultar roadmaps. En flujos normales del skill, no hacer gate directo con `command -v rootline`; dejar que `roadmapctl context`, `roadmapctl doctor` o `roadmapctl check` reporten `RMC_ENV_ROOTLINE_MISSING` con diagnostics estables.

Usar comandos `rootline` directos solo para troubleshooting explícito después de que el flujo principal de `roadmapctl` lo requiera o reporte diagnostics.

### roadmapctl

Requerido para comandos implementados de `/roadmap` que escriben, mutan,
ejecutan tasks o declaran validez del roadmap.

Gate antes de escribir/mutar/ejecutar/declarar validez:

```bash
command -v roadmapctl
roadmapctl doctor --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Postcheck obligatorio después de materializar o mutar:

```bash
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si `roadmapctl` falta o sale non-zero, detenerse y reportar diagnostics. No
hacer auto-fix, no fallback a markdown libre, no commitear y no ejecutar tasks.

La generación conceptual de planes puede continuar sin rootline/roadmapctl solo
si no escribe, no muta, no ejecuta y no declara validez.

### Verificación obligatoria al modificar este skill

Todo cambio al skill `/roadmap` o a sus guards `roadmapctl` debe probarse con Pi headless antes de commit/release. No alcanza con grep.

Ejecutar desde el repo canónico:

```bash
./scripts/sync-roadmap-skill.sh --install
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/roadmap/SKILL.md --tools read,bash -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: the user asks "loop autonomo" in this repository. Do not modify files and do not run git commit/push. Perform only the bootstrap and the required preflight checks from the skill, then stop. In your final answer, list the exact commands you ran and whether roadmapctl doctor/check were required and passed.'
PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/roadmap/SKILL.md --tools read,bash -p 'HEADLESS VERIFICATION TEST. Use the roadmap skill. Scenario: there is an already approved plan to materialize one direct task, and the user says "crea las tareas". Do not create or modify files and do not run git commit/push. Perform only bootstrap and the required preflight checks that must happen before any roadmap write, then stop. In your final answer, list exact commands run and whether roadmapctl doctor/check were required and passed.'
```

La evidencia debe mostrar que `roadmapctl doctor` y `roadmapctl check` fueron requeridos y pasaron antes de loop/materialización, sin modificar archivos.

## Routing por subcomando

Después del bootstrap:

| `$ARGUMENTS` | Archivo | Descripción |
|--------------|---------|-------------|
| *(sin argumentos)* | [pending-subcommand.md](pending-subcommand.md) | Trabajo pendiente por defecto |
| `pending` | [pending-subcommand.md](pending-subcommand.md) | Trabajo pendiente |
| `decision`, `next`, `prioriza`, `qué sigue` | [decision-tree-subcommand.md](decision-tree-subcommand.md) | Priorizar qué ejecutar |
| `plan` | [plan-subcommand.md](plan-subcommand.md) | Materializar plan como `.md` |
| `loop [--filter] [--max]` | [loop-subcommand.md](loop-subcommand.md) | Ejecutar tasks pendientes |
| *(texto libre)* | [autonomous-mode.md](autonomous-mode.md) | Descomponer en Outcome/Tasks |

## Flag global `--repo`

Solo workspace mode:

- `--repo <name>` resuelve un único repo y remueve el flag antes del dispatch.
- En single-repo mode se ignora.

## Regla de dispatch

1. Si vacío → `pending`.
2. Si empieza con `pending`, `loop`, `plan` → subcomando directo.
3. Si empieza con `decision`/`next` o pide priorización/“qué sigue” → `decision-tree`.
4. Si pide estado/progreso/pendientes → `pending`.
5. Si dice "crea las tareas", "materializa", "genera los archivos",
   "pasalo al roadmap", "crea el roadmap" o equivalente → `plan`.
6. Si describe algo a construir o descomponer sin pedir archivos → modo autónomo.

Ambigüedad crítica:

- "descompón/planifica" = proponer estructura, no escribir archivos.
- "crea/materializa/genera tareas" = `plan-subcommand.md`.
- Si no hay plan previo suficiente para materializar, preguntar antes de escribir.

## Lógica común

[common-logic.md](common-logic.md) es referencia de mantenimiento. Los subcomandos implementados deben ser autosuficientes en su ruta normal: no leer lógica común para `pending`/`decision`/`next`; leerla solo para troubleshooting o cambios al skill.

## Referencia

- Modelo completo: [framework-reference.md](framework-reference.md)
- Outcomes: [outcome-guide.md](outcome-guide.md)
- Tasks: [task-guide.md](task-guide.md)
- `.stem` base: [base.stem](base.stem)
