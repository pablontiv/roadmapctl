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

```text
Outcome/Objetivo  (opcional)
└── Task          (unidad ejecutable)
```
Para trabajo chico, usar solo tasks.

## Invariante de materialización

Materializar tareas significa crear archivos canónicos, no un único archivo resumen:

```text
<roadmap-root>/OXX-slug/README.md
<roadmap-root>/OXX-slug/TXXX-task.md
<roadmap-root>/TXXX-task.md  (Task directa)
```

Si no se puede crear esa estructura, detenerse. No fallback a markdown libre.

## Invariante de escritura segura

El skill es el único writer vía Write tool, después de aprobación explícita y preflight pasado.

Prohibido: heredocs/cat> en shell para múltiples archivos; loops con `rootline new`; escribir sin aprobación; escribir si `roadmapctl doctor` o `check --strict` retornan non-zero.

Permitido: Write tool por archivo canónico tras aprobación y preflight exitoso.

## Bootstrap

Ejecutar `roadmapctl bootstrap --repo <repo> --output json`. Su JSON es fuente de verdad para config, helpers y comportamiento. Gate inicial: `command -v roadmapctl`. Detalle: [bootstrap-reference.md](bootstrap-reference.md) | [config-reference.md](config-reference.md).

## Gates CLI

Antes de escribir/mutar/ejecutar/declarar validez:

```bash
roadmapctl doctor --repo <repo> --output json --strict
roadmapctl check --repo <repo> --output json --strict
```

Si sale non-zero, detenerse. No auto-fix, no fallback. Verificación del skill: [verification-reference.md](verification-reference.md).

## Routing por subcomando

| `$ARGUMENTS` | Archivo |
|---|---|
| vacío / `pending` | [pending-subcommand.md](pending-subcommand.md) |
| `decision` / `next` | [decision-tree-subcommand.md](decision-tree-subcommand.md) |
| `plan` | [plan-subcommand.md](plan-subcommand.md) |
| `loop [--filter] [--max]` | [loop-subcommand.md](loop-subcommand.md) |
| texto libre | [autonomous-mode.md](autonomous-mode.md) |

## Flag global `--repo`

Solo workspace mode: `--repo <name>` resuelve un repo. Single-repo: se ignora.

## Regla de dispatch

1. Vacío → `pending`.
2. `pending`, `loop`, `plan` → subcomando directo.
3. `decision`/`next`/priorización → `decision-tree`.
4. Estado/progreso/pendientes → `pending`.
5. "crea las tareas"/"materializa"/"genera archivos" → `plan`.
6. Texto libre sin pedir archivos → modo autónomo.

Ambigüedad: "descompón/planifica" = propuesta (no escribe); "crea/materializa" = `plan`. Si no hay plan suficiente, preguntar antes de escribir.

## Referencia

- Modelo completo: [framework-reference.md](framework-reference.md)
- Outcomes: [outcome-guide.md](outcome-guide.md)
- Tasks: [task-guide.md](task-guide.md)
- `.stem` base: [base.stem](base.stem)
- Bootstrap detallado: [bootstrap-reference.md](bootstrap-reference.md)
- Config completa: [config-reference.md](config-reference.md)
- Verificación skill: [verification-reference.md](verification-reference.md)
