# Lógica Común — Materialización y Ejecución

Referencia de mantenimiento/troubleshooting para el skill. Los subcomandos implementados deben ser autosuficientes en su ruta normal y no depender de leer este archivo para operaciones simples.

> En workspace mode, `<roadmap-root>` se reemplaza por `<abs-roadmap-root>` y los comandos git usan `git -C <repo-path>`.

## Guard obligatorio: roadmapctl

Para cualquier flujo que escriba, mute, ejecute tasks o declare validez del roadmap, `roadmapctl` es obligatorio además de Rootline.

Antes de escribir, mutar o ejecutar:

```bash
command -v roadmapctl
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

## Invariante de escritura segura

Materializar implica delegar las escrituras canónicas a `roadmapctl`; el skill no escribe markdown del roadmap directamente ni arma dumps multi-file manuales.

**Prohibido en una misma tool call del skill:**

- `bash`/`sh` con múltiples `cat >`, heredocs multi-file o loops de escritura.
- una sola llamada que genere `rootline new` para múltiples rutas.
- `write`/`edit` directo para crear o reescribir archivos canónicos del roadmap.
- cualquier escritura manual que cree/reescriba más de un archivo canónico del roadmap.

Permitido: un único comando `roadmapctl materialize --plan <plan-json> --apply` o `roadmapctl materialize --changes <dry-run-json> --apply` puede aplicar múltiples archivos, porque `roadmapctl` valida el plan/change-set, ordena padres antes de hijos, reporta diagnostics por path, ejecuta validaciones y corre postcheck antes de éxito.

En `/roadmap plan`:

1. Guardar el plan aprobado en un archivo temporal (`plan_json`); no pegar JSON grande en el prompt o respuesta.
2. Revisar siempre un dry-run determinístico guardado en archivo temporal (`dry_run_json`).
3. Para la revisión normal, mostrar solo `summary`, `diagnostics`, y por cambio `path`, `operation`, `applied`, `preconditions`; no volcar `changes[].content` ni diffs completos salvo pedido explícito o troubleshooting puntual.
4. Guardar el dry-run JSON como change-set congelado cuando se use `--changes`.
5. Preferir batch apply owned by roadmapctl cuando `parallel = true`; usar `--target` granular solo para recuperación puntual, selección humana de un único archivo, o troubleshooting.
6. Bootstrap explícito (`.` / `.stem` / `.roadmapctl.toml`) conserva su flujo canónico propio.

## Materialización determinística

La ruta primaria para crear archivos del roadmap es:

```bash
roadmapctl materialize --plan <plan-json> --dry-run --repo <repo-path> --roadmap-root <roadmap-root> --output json
roadmapctl materialize --plan <plan-json> --apply --repo <repo-path> --roadmap-root <roadmap-root> --output json
```

El skill no debe duplicar numbering, `rootline new`, writes, actualización de tablas ni escritura final de `blocked_by`; debe producir plan estructurado, revisar dry-run y delegar en `roadmapctl materialize`. El skill sí debe decidir si una relación es un hard blocker antes de serializarla como `blocked_by`.

## Auto-numbering

El skill no calcula números `OXX`/`TXXX`. `roadmapctl materialize` asigna numbering determinístico y reporta las rutas propuestas en `changes[]` durante dry-run. Si el dry-run no produce rutas canónicas, detenerse y reportar diagnostics.

## Verificación de padre

La verificación normal de padres, rutas, allowlist y orden de creación pertenece a `roadmapctl materialize --dry-run` y `--apply`. No ejecutar `rootline describe` como paso primario antes de crear archivos.

Si `roadmapctl materialize` reporta que falta un padre, schema, `.stem`, permisos o estructura, informar al usuario y no crear archivos fuera del roadmap.

Excepción permitida: `plan-subcommand.md` puede crear `<roadmap-root>/` y `<roadmap-root>/.stem` solo mediante el bootstrap explícito gobernado por `roadmapctl bootstrap init`. Fuera de ese flujo, no crear directorios ad-hoc.

## Cascading links

El skill no edita manualmente la tabla `## Tasks`. `roadmapctl materialize` actualiza el README del Outcome y mantiene la tabla sin columna Estado; el estado se lee desde frontmatter.

## Dependencias duras

`blocked_by` es una dependencia dura, no una sugerencia de orden. Declararlo en la task bloqueada solo cuando la task actual no pueda ejecutarse o validarse si el target no está completado.

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
