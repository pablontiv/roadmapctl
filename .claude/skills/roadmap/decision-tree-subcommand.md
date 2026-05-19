# /roadmap — Árbol de Decisión

Generar un árbol de decisión que muestre recomendaciones ejecutables, quick wins, dependencias y bloqueos usando `roadmapctl` como fuente determinística.

Ruta normal autosuficiente: usar el `roadmapctl bootstrap` ya obtenido en bootstrap y luego `roadmapctl decision` o `roadmapctl next`. No leer `common-logic.md`, documentación de integración ni archivos Rootline para este flujo. No ejecutar `roadmapctl doctor`/`check`: decision/next son read-only y sus comandos son la fuente canónica.

## Workspace mode

Si `<repos>` existe, ejecutar por repo y agrupar salida con prefijo `[repo]`. Si el usuario quiere focalizar, sugerir `/roadmap --repo <name>`.

Para cada repo resuelto:

```bash
roadmapctl decision --repo <repo> --output json
```

Opcionalmente, para mostrar solo readiness/blockers sin scoring:

```bash
roadmapctl next --repo <repo> --output json
```

## Single-repo

```bash
roadmapctl decision --repo <repo> --output json
```

Si el usuario pide específicamente "next task", "listas" o "bloqueadas":

```bash
roadmapctl next --repo <repo> --output json
```

## Interpretar `roadmapctl decision`

Usar el JSON sin recalcular lógica en prompt:

- `kind` debe ser `roadmapctl/decision`.
- `recommendations[]` ya viene ordenado de forma determinística.
- Cada recomendación incluye `score`, `reasons`, `path`, `status`, `outcome_path` y `unblocks` cuando aplica.
- `quick_wins[]` identifica progreso rápido.
- `critical_blockers[]` identifica tareas listas que desbloquean otras tareas.
- `blocked[]` explica blockers por path resuelto.
- Si `summary.status != "ok"`, detenerse y reportar `diagnostics`.

## Interpretar `roadmapctl next`

Usar `roadmapctl next` cuando se necesita separar explícitamente:

- `ready[]`: tasks activas cuyas dependencias están completadas según `done_statuses`.
- `blocked[]`: tasks activas con `blockers[]` insatisfechos.

## Renderizar

```text
ROADMAP DECISION TREE

Qué objetivo priorizar?
│
├─► RECOMENDACIONES
│   <path> [<status>] score=<score>
│   razones: <reasons>
│   desbloquea: <unblocks>
│
├─► QUICK WINS
│   <path> [<status>]
│
└─► BLOQUEADAS
    <path> bloqueada por: <blockers>
```

## Criterios de decisión

```text
CRITERIOS
├─ Hay trabajo In Progress? → roadmapctl decision lo prioriza en score/reasons
├─ Hay task que desbloquea muchas otras? → mirar critical_blockers[]
├─ Quiero progreso rápido? → mirar quick_wins[]
└─ Si no → primera entrada de recommendations[]
```

Reglas:

- Si `summary.status != "ok"`, detenerse y reportar `diagnostics`.
- No llamar `roadmapctl doctor` ni `roadmapctl check` para decision/next.
- No llamar `rootline tree`, `rootline graph` ni `rootline query` directamente para pending/next/decision.
- No postprocesar JSON crudo de Rootline para construir blockers, reverse dependencies, quick wins, scoring o recomendaciones.
- No buscar links con grep.
- No usar `git log` ni Backscroll como input de scoring read-only salvo que exista una futura opción explícita en `roadmapctl`.
- No recalcular dependencias, reverse dependencies, quick wins o scoring en prompt; esa lógica pertenece a `roadmapctl decision`/`roadmapctl next`.
