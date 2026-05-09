# /roadmap plan

> Pre-requisito: leer [common-logic.md](common-logic.md).

Materializa el plan de la conversación como archivos `.md` del roadmap. No implementa código.

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
5. Cada task debe tener nombre, descripción, dependencias `blocked_by` con paths relativos explícitos y ACs principales.

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

Convertir el árbol aprobado a JSON `roadmapctl/materialize-plan` versión 1 según `docs/materialize-plan-schema.md`:

```json
{
  "version": 1,
  "kind": "roadmapctl/materialize-plan",
  "items": []
}
```

Reglas:

- No pasar prose libre a `roadmapctl materialize`.
- Cada Outcome/task aprobado debe tener `slug`, `title`, `description`, ACs, `source_of_truth` y límites suficientes.
- Las dependencias deben representarse como `blocked_by` con `ref` plan-local o `path` explícito; nunca targets bare.
- Si falta información para poblar campos requeridos, preguntar antes de materializar.

### Paso 2: Dry-run determinístico

Ejecutar:

```bash
roadmapctl materialize --plan <plan-json> --dry-run --repo <repo-path> --roadmap-root <roadmap-root> --output json
```

Revisar el JSON:

- `summary.status == "ok"`.
- `changes[]` contiene únicamente operaciones canónicas permitidas:
  - `.` / `.stem` / `.roadmapctl.toml` solo en bootstrap explícito aprobado (preferir `roadmapctl bootstrap init` para crear solo bootstrap),
  - `OXX-slug/README.md`,
  - `OXX-slug/TXXX-task.md`,
  - `TXXX-task.md`.
- `applied == false` para todo dry-run.
- No aparece ningún `*-tasks.md`.

Si el dry-run falla o propone rutas fuera del allowlist, detenerse y reportar diagnostics. No escribir archivos manualmente.

### Paso 3: Aplicación granular por archivo

El comportamiento estricto es materializar **por archivo** (no batch): ninguna operación puede escribir más de un archivo roadmap canónico.

Solo después del dry-run válido y aprobación humana explícita:

1. Guardar el JSON completo del dry-run como change-set congelado (`<dry-run-json>`).
2. Derivar `targets[]` desde `changes[]` del dry-run (solo archivos permitidos: `OXX-slug/README.md`, `OXX-slug/TXXX-*.md`, `TXXX-*.md`).
3. Si `targets` está vacío, detenerse con error.
4. Para cada target aprobado, ejecutar una operación independiente:
   ```bash
   roadmapctl materialize --changes <dry-run-json> --target <target.path> --apply --repo <repo-path> --roadmap-root <roadmap-root> --output json
   ```
5. Cada invocación debe reportar exactamente un `changes[]` aplicado. Si `summary.status != "ok"`, detener la materialización y reportar diagnostics.

Este paso no permite ejecutar un único `roadmapctl materialize --plan <plan-json> --apply` sobre el plan completo dentro del skill. Tampoco permite escribir el `content` del dry-run manualmente desde el prompt o un subagente.

Si algún target falla, detener la materialización y reportar su diagnóstico; no continuar ni hacer fallback.

### Paso 4: Postcheck explícito

Tras completar todos los targets, ejecutar:

```bash
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si falla, detenerse y reportar diagnostics.

### Paso 5: Commit + push

- `git add` solo archivos `.md` y `.stem` creados/modificados del roadmap.
- `git commit -m "chore(roadmap): create planning docs"`
- `git push` si `<auto-push>` es true.

STOP obligatorio. Informar: “Archivos de planificación creados. Ejecutar `/roadmap loop` cuando esté listo para implementar.”
