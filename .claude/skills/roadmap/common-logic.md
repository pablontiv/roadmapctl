# Lógica Común — Materialización y Ejecución

Referencia de mantenimiento/troubleshooting para el skill. Los subcomandos implementados deben ser autosuficientes en su ruta normal y no depender de leer este archivo para operaciones simples.

## Configuración de Campos

El vocabulario de campos del roadmap es configurable sin recompilar. Ver `.roadmapctl.toml` sección `[fields]` para `lifecycle`, `record_type`, `task_value`, `outcome_value`, `display_name`, `dependency_link`. Defaults retrocompatibles preservan comportamiento actual.

## Extracción de campos de bootstrap

`roadmapctl bootstrap --field <dot-path>` extrae un valor escalar del JSON de bootstrap sin necesitar jq o python3. Soporta paths simples (`roadmap_root`) y anidados (`helpers.where_leaf`). Exit 1 si el campo no existe o es un objeto/array. Compatible con y sin `--output json`.

## Título de tasks en next output

`roadmapctl next --output json` incluye el campo `title` en cada task (cuando disponible). El título se extrae del campo derivado configurado en `[fields].display_name` (default: `titulo`), que viene del H1 del documento vía `source: body.h1` en el `.stem`.

> En workspace mode, `<roadmap-root>` se reemplaza por `<abs-roadmap-root>` y los comandos git usan `git -C <repo-path>`.

## Guard obligatorio: roadmapctl

Para cualquier flujo que escriba, mute, ejecute tasks o declare validez del roadmap, `roadmapctl` es obligatorio además de Rootline.

Antes de escribir, mutar o ejecutar:

```bash
command -v roadmapctl
roadmapctl --version                                                                  # diagnóstico: version del binario instalado
roadmapctl doctor --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Después de cualquier materialización o mutación del roadmap:

```bash
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si `roadmapctl` no existe o cualquier comando sale non-zero, detenerse antes de continuar. Reportar comando, exit code y diagnostic IDs si hubo JSON. No auto-fix, no fallback a markdown libre, no ejecutar tasks y no commitear mutaciones del roadmap.

Si una materialización escribió algunos cambios y luego falló postcheck, reportar los paths aplicados, validar los paths afectados (`rootline validate <path>` o `roadmapctl check --strict`), y elegir explícitamente corregir o revertir antes de volver a ejecutar `roadmapctl check --strict`. No declarar éxito ni commitear mientras el postcheck falle.

Los diagnostics `RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED` y `RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY` indican `.stem` stale que impide materializar Outcomes sin `estado`; bloquear writes hasta corregir el schema mediante roadmapctl o intervención explícita. Los diagnostics de config deben tratar el path reportado por roadmapctl como fuente efectiva (`.roadmapctl.toml` en repos TOML-backed; legacy solo en migración).

La planificación conceptual que no escribe, no muta, no ejecuta y no declara validez puede continuar sin `roadmapctl`.

## Modelo

El roadmap usa máximo dos niveles:

```text
<roadmap-root>/
├── O01-outcome/
│   ├── README.md
│   └── T001-task.md
└── T001-task-directa.md
```

## Prohibición de fallback

Nunca representar múltiples tasks como una lista dentro de un único archivo:

```text
<roadmap-root>/algo-tasks.md
```

Eso no es una materialización válida del roadmap.

Si una operación pretende crear N tasks, debe crear N archivos `TXXX-*.md`.
Si las tasks pertenecen a un Outcome, debe existir también el README del Outcome.

Si falta schema, `.stem`, `rootline`, permisos o estructura para crear archivos
canónicos, detenerse. No usar `Write` directo para inventar una estructura
alternativa.

## Auto-numbering

El skill usa `rootline describe` para obtener numeración determinística:

```bash
# Retorna el siguiente O y T simultáneamente
rootline describe <roadmap-root> --field schema.id.next_by_pattern --output json
# → {"O*": "O14", "T*": "T014"}

# Retorna solo el siguiente T dentro de un outcome existente
rootline describe <roadmap-root>/OXX-slug/ --field schema.id.next_by_pattern --output json
# → {"T*": "T009"}
```

Preferir `next_by_pattern` (mapa) sobre `next` (string) en schemas con múltiples patrones de secuencia. `next` es determinístico post-fix pero retorna solo el primer patrón alfabético que coincide con entries existentes.

## Cascading links

El skill no edita manualmente la tabla `## Tasks`. `roadmapctl materialize` actualiza el README del Outcome y mantiene la tabla sin columna Estado; el estado se lee desde frontmatter.

## Dependencias duras

`blocked_by` es el nombre de link de dependencia dura por defecto. Se puede cambiar via `[fields].dependency_link` en `.roadmapctl.toml`. Las instrucciones a continuación usan el nombre por defecto.

Una dependencia dura no es una sugerencia de orden. Declarar `blocked_by` en la task bloqueada solo cuando la task actual no pueda ejecutarse o validarse si el target no está completado.

Guardrail obligatorio antes de serializar `blocked_by`:

```text
¿Qué fallaría objetivamente si ejecuto esta task antes?
```

- Si hay una falla concreta: usar `blocked_by` con path relativo explícito.
- Si solo hay orden sugerido, secuencia narrativa, relación temática, provenance, “conviene después de” o “usar su output si existe”: no usar `blocked_by`; ponerlo en `Contexto`, `Fuente de verdad` o prose.
- No usar `blocked_by` para orden sugerido ni para forzar serialización artificial dentro o entre Outcomes.
- Misma carpeta/Outcome: `[[blocked_by:./T001-prerequisite.md]]`
- Otro Outcome: `[[blocked_by:../O01-setup/T001-prerequisite.md]]`
- No usar targets bare como `[[blocked_by:T001-prerequisite]]`; rootline solo los resuelve por basename único y pueden romperse con duplicados.

## Comandos Rootline de Referencia (troubleshooting/legacy)

Estos comandos son solo para inspección de bajo nivel, troubleshooting o reparación puntual después de que los flujos obligatorios de `roadmapctl` hayan pasado. No son la ruta primaria para descubrir pendientes, elegir siguiente task, explicar blockers ni construir el árbol de decisión; para eso usar `roadmapctl pending`, `roadmapctl next` y `roadmapctl decision`.

| Comando | Cuándo usarlo |
|---------|---------------|
| `rootline validate <path>` | Troubleshooting después de crear/editar `.md` |
| `rootline fix <path>` | Reparación puntual si validate falla y la propuesta es segura |
| `rootline query <path> --where "expr"` | Inspección legacy de metadata, no selección pending/next/decision |
| `rootline tree <path> --where "expr" --output json` | Inspección legacy de jerarquía, no conteos principales del skill |
| `rootline graph <path> --where "expr" --check` | Inspección legacy de dependencias, no readiness/blockers del skill |

No usar `rootline stats`; `tree` ya incluye conteos. No postprocesar JSON crudo de Rootline para reconstruir lógica que pertenece a `roadmapctl`.

## Workspace context — fixture .git dirs

Los tests de workspace context requieren directorios `.git` en los fixtures. Git no rastrea directorios vacíos ni archivos dentro de un path component llamado `.git`. `TestMain` en `internal/cli/golden_test.go` los crea con `os.MkdirAll` al arrancar los tests. No agregar `.gitkeep` dentro de esos directorios.
