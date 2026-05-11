# /roadmap plan

Materializa el plan de la conversación como archivos `.md` del roadmap. No implementa código.

Ruta normal autosuficiente: este archivo contiene el procedimiento operativo completo. No leer `common-logic.md` ni documentación de integración para ejecutar el flujo; esos documentos son referencia de mantenimiento/troubleshooting.

Materializar es una operación estructural. Está prohibido crear un único archivo
con una lista de tareas. Cada task debe tener su propio archivo `TXXX-*.md`.

## Fuente del plan

1. Contexto actual de conversación.

Si no hay plan, informar: “No hay plan en esta conversación. Primero planificar, luego ejecutar `/roadmap plan`.” y parar.

## Workspace mode

Resolver repo target:

1. `--repo <name>` si fue dado.
2. Repo mencionado en el plan.
3. Si ambiguo, preguntar.

Usar `<abs-roadmap-root>` y `git -C <repo-path>`.

## Fase 1: Descomposición

1. Identificar el plan más reciente.
2. Leer contexto existente relacionado bajo `<roadmap-root>/`.
3. Aplicar [framework-reference.md](framework-reference.md): máximo Outcome + Tasks.
4. Producir:
   - tasks directas, o
   - Outcome(s) + tasks.
5. Cada task debe tener nombre, descripción, ACs principales y, solo si aplica, `hard_blockers` explícitos. Un hard blocker es una dependencia objetiva: si no está completada, la task actual no debe ejecutarse.

## Fase 2: Aprobación

Presentar árbol completo y pedir aprobación con `AskUserQuestion`.

STOP hasta aprobación. No crear archivos antes.

## Fase 3: Materialización

**MATERIALIZAR ≠ IMPLEMENTAR.** Crear solo archivos `.md` y `.stem` dentro de `<roadmap-root>/`.

### Preflight obligatorio roadmapctl

Antes de crear o modificar cualquier archivo del roadmap:

```bash
command -v roadmapctl
roadmapctl doctor --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si `roadmapctl` falta o cualquier comando sale non-zero, detenerse antes de escribir. Reportar comando, exit code y diagnostic IDs si hubo JSON. No crear archivos, no usar fallback `*-tasks.md`, no auto-fix.

#### Excepción explícita de bootstrap

La única excepción al `doctor`/`check` previo es un bootstrap pedido explícitamente para un `<roadmap-root>` inexistente. Flujo permitido:

```bash
roadmapctl bootstrap inspect --repo <repo-path> --roadmap-root <roadmap-root> --output json
roadmapctl bootstrap init --repo <repo-path> --roadmap-root <roadmap-root> --dry-run --output json
# tras aprobación explícita del dry-run:
roadmapctl bootstrap init --repo <repo-path> --roadmap-root <roadmap-root> --apply --output json
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

`bootstrap init --apply` también ejecuta postcheck internamente; el `check --strict` posterior sigue siendo obligatorio como evidencia externa antes de declarar éxito. Si el root ya existe y `doctor`/`check` falla, no usar bootstrap como reparación: detenerse y reportar diagnostics.

Guardrail obligatorio antes de escribir:

1. Confirmar que se va a crear una de estas formas:
   - Outcome + tasks: `OXX-slug/README.md` + `OXX-slug/TXXX-*.md`
   - Tasks directas: `TXXX-*.md` en la raíz del roadmap.
2. Si el plan contiene varias tasks, no escribirlas en un archivo único.
3. Si no hay información suficiente para nombrar/separar tasks, preguntar.
4. Si falta `rootline` y no se puede crear estructura canónica, detenerse.

### Paso 1: Serializar plan estructurado

Convertir el árbol aprobado a JSON `roadmapctl/materialize-plan` versión 1 y guardarlo en un archivo temporal, por ejemplo:

```bash
plan_json="$(mktemp)"
```

Schema operativo mínimo autosuficiente:

```json
{
  "version": 1,
  "kind": "roadmapctl/materialize-plan",
  "items": []
}
```

`items[]` contiene Outcomes y/o Tasks aprobadas con `slug`, `title`, `description`, `preserves`, `context`, `scope_in`, `scope_out`, `acceptance_criteria`, `source_of_truth`, `initial_status` y `hard_blockers` cuando correspondan.

Fuente del schema: el contrato canónico versionado vive en el repositorio `roadmapctl` como `docs/materialize-plan-schema.md`; esa ruta es repo-relativa cuando se está trabajando dentro de `roadmapctl`, no relativa al directorio del skill instalado. En repos consumidores normales, usar el schema mínimo autosuficiente de este archivo y los diagnostics de `roadmapctl materialize`; no buscar un `docs/materialize-plan-schema.md` local salvo que el repo sea `roadmapctl`. Exposición CLI tipo `roadmapctl materialize schema --output json` queda diferida hasta que se apruebe una versión de schema exportable.

Reglas:

- No pasar prose libre a `roadmapctl materialize`.
- No pegar el JSON completo en la respuesta si es grande; guardar el plan en `plan_json` y reportar solo un resumen.
- Cada Outcome/task aprobado debe tener `slug`, `title`, `description`, ACs, `source_of_truth` y límites suficientes.
- Los Outcome v1 crean un Outcome nuevo cuando el `slug` no existe. Si el `slug` ya existe como `OXX-<slug>`, `roadmapctl materialize` debe appendear las tasks dentro de ese Outcome con un `update` gobernado del `README.md` y `create` de los `TXXX-*.md`; el skill no debe intentar append/update manual ni proponer un Outcome numerado duplicado.
- Serializar `blocked_by` **solo** desde hard blockers aprobados, con `ref` plan-local o `path` explícito; nunca targets bare.
- Antes de incluir cualquier `blocked_by`, responder: “¿Qué fallaría objetivamente si ejecuto esta task antes?”. Si la respuesta es “nada; solo es mejor orden/contexto”, no es hard blocker.
- No usar `blocked_by` para orden sugerido, secuencia narrativa, agrupación por Outcome, relación temática, provenance, “conviene después de”, ni “usar su output si existe”. Poner ese contexto en `context`, `source_of_truth` o prose de la task.
- Si falta información para poblar campos requeridos o justificar un hard blocker, preguntar antes de materializar.

### Paso 2: Dry-run determinístico

Ejecutar:

```bash
dry_run_json="$(mktemp)"
roadmapctl materialize --plan "$plan_json" --dry-run --repo <repo-path> --roadmap-root <roadmap-root> --output json >"$dry_run_json"
```

Revisar el JSON desde `dry_run_json` sin volcar contenido/diffs completos al contexto:

- `summary.status == "ok"`.
- `changes[]` contiene únicamente operaciones canónicas permitidas:
  - `.` / `.stem` / `.roadmapctl.toml` solo en bootstrap explícito aprobado (preferir `roadmapctl bootstrap init` para crear solo bootstrap),
  - `OXX-slug/README.md`,
  - `OXX-slug/TXXX-task.md`,
  - `TXXX-task.md`.
- `applied == false` para todo dry-run.
- No aparece ningún `*-tasks.md`.

Si el dry-run falla o propone rutas fuera del allowlist, detenerse y reportar diagnostics. No escribir archivos manualmente.

Reporte normal del dry-run: `summary`, `diagnostics`, `path`, `operation`, `applied` y preconditions relevantes. No usar `cat "$dry_run_json"` ni pegar el JSON completo en la respuesta; extraer esos campos de forma selectiva. Leer `changes[].content` o diffs completos solo si el usuario lo pide explícitamente o para troubleshooting dirigido.

### Paso 3: Aplicación batch gobernada por roadmapctl

Solo después del dry-run válido y aprobación humana explícita, aplicar con una única operación owned by roadmapctl:

```bash
roadmapctl materialize --plan "$plan_json" --apply --repo <repo-path> --roadmap-root <roadmap-root> --output json
```

Alternativa con change-set congelado revisado:

```bash
roadmapctl materialize --changes "$dry_run_json" --apply --repo <repo-path> --roadmap-root <roadmap-root> --output json
```

Reglas:

1. El skill no escribe `content` del dry-run manualmente y no usa shell heredocs/loops para crear archivos.
2. `roadmapctl` debe reportar `summary.status == "ok"`, `applied == true`, y `changes[]` con todos los paths aplicados.
3. Si `summary.status != "ok"`, detener la materialización y reportar diagnostics; no continuar ni hacer fallback.
4. Usar `--changes <dry-run-json> --target <target.path> --apply` solo para recuperación puntual o cuando se aprobó explícitamente aplicar un único archivo.

### Paso 4: Postcheck explícito

Tras completar todos los targets, ejecutar:

```bash
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si falla, detenerse y reportar diagnostics. No declarar éxito, no commitear y no continuar con más materialización hasta resolverlo.

#### Recuperación ante postcheck fallido o materialización parcial

1. Reportar `summary`, diagnostic IDs y `changes[]` con `applied=true` para dejar claro qué quedó escrito.
2. Inspeccionar solo los paths aplicados relevantes; no reescribirlos manualmente en batch.
3. Ejecutar `rootline validate <path>` sobre los `.md` aplicados o `roadmapctl check --strict` para confirmar el fallo actual.
4. Elegir explícitamente una de dos rutas: corregir con una edición puntual aprobada o revertir los paths aplicados.
5. Volver a ejecutar `roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict` y solo entonces permitir commit o declaración de éxito.

### Paso 5: Commit + push

- `git add` solo archivos `.md` y `.stem` creados/modificados del roadmap.
- `git commit -m "chore(roadmap): create planning docs"`
- `git push` si `<auto-push>` es true.

STOP obligatorio. Informar: “Archivos de planificación creados. Ejecutar `/roadmap loop` cuando esté listo para implementar.”
