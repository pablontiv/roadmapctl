# /roadmap plan

> Pre-requisito: leer [common-logic.md](common-logic.md).

Materializa el plan de la conversación como archivos `.md` del roadmap. No implementa código.

Materializar es una operación estructural. Está prohibido crear un único archivo
con una lista de tareas. Cada task debe tener su propio archivo `TXXX-*.md`.

## Fuente del plan

1. Contexto actual de conversación.
2. Fallback: `~/.claude/plans/${CLAUDE_SESSION_ID}.md`.

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
  - `.` / `.stem` / `.roadmapctl.toml` solo en bootstrap explícito,
  - `OXX-slug/README.md`,
  - `OXX-slug/TXXX-task.md`,
  - `TXXX-task.md`.
- `applied == false` para todo dry-run.
- No aparece ningún `*-tasks.md`.

Si el dry-run falla o propone rutas fuera del allowlist, detenerse y reportar diagnostics. No escribir archivos manualmente.

### Paso 3: Apply aprobado

Solo después del dry-run válido y aprobación humana explícita, ejecutar:

```bash
roadmapctl materialize --plan <plan-json> --apply --repo <repo-path> --roadmap-root <roadmap-root> --output json
```

`roadmapctl materialize --apply` es responsable de:

- crear `OXX/README.md` y `TXXX-*.md`,
- actualizar la tabla `## Tasks` del README de Outcome,
- escribir `blocked_by` con paths relativos explícitos,
- materializar bootstrap aprobado sin sobrescribir `.stem` existente,
- validar archivos creados y ejecutar postcheck.

Si `summary.status != "ok"`, detenerse antes de declarar éxito o commitear. Reportar diagnostics; no auto-fix y no fallback a markdown libre.

### Paso 4: Postcheck explícito

Aunque `materialize --apply` ejecuta postcheck, correr el guard común antes de declarar éxito:

```bash
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si falla, detenerse y reportar diagnostics.

### Paso 5: Commit + push

- `git add` solo archivos `.md` y `.stem` creados/modificados del roadmap.
- `git commit -m "chore(roadmap): create planning docs"`
- `git push` si `<auto-push>` es true.

STOP obligatorio. Informar: “Archivos de planificación creados. Ejecutar `/roadmap loop` cuando esté listo para implementar.”
