# /roadmap plan

Materializa el plan de la conversación como archivos `.md` del roadmap. No implementa código.

Ruta normal autosuficiente: este archivo contiene el procedimiento operativo completo. No leer `common-logic.md` ni documentación de integración para ejecutar el flujo; esos documentos son referencia de mantenimiento/troubleshooting.

Materializar es una operación estructural. Está prohibido crear un único archivo
con una lista de tareas. Cada task debe tener su propio archivo `TXXX-*.md`.

## Fuente del plan

1. Contexto actual de conversación.

Si no hay plan, informar: "No hay plan en esta conversación. Primero planificar, luego ejecutar `/roadmap plan`." y parar.

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

### Paso 1: Path planning y propuesta visual

Determinar rutas canónicas para Outcomes y Tasks:

1. Usar `roadmapctl plan-paths` para proponer rutas canónicas:
   ```bash
   # input.json: {"version":1,"kind":"roadmapctl/path-plan","items":[{"type":"outcome","slug":"slug"},{"type":"task","slug":"slug","outcome_slug":"slug"}]}
   roadmapctl plan-paths --input input.json --repo <repo-path> --roadmap-root <roadmap-root> --output json
   ```
   Revisar la propuesta: `summary.status == "ok"` y `paths[]` contiene rutas canónicas `OXX-slug/README.md` y `OXX-slug/TXXX-*.md` o `TXXX-*.md` directas.

Presentar la propuesta visual al usuario en formato legible:

```
Propuesta aprobada:

O01-nombre-del-objetivo/
├── README.md
└── T001-task-corta.md
    - AC1: ...
    - AC2: ...
└── T002-task-larga.md
    - AC1: ...
```

**STOP obligatorio hasta aprobación explícita del usuario.** Solo si el usuario aprueba el árbol exacto propuesto, pasar a Paso 2.

### Paso 2: Re-pregunta y divergencia

Re-preguntar **solo si**:

- Hay divergencia entre las rutas propuestas en Paso 1 y lo que el usuario aprobó (ej: cambió cantidad de tasks, destino, o slugs).
- El usuario pide cambios estructura respecto a la propuesta visual presentada.

Si el usuario aprueba explícitamente el árbol visual propuesto (rutas, slugs, tasks, ACs), no re-preguntar; pasar directo a Paso 3.

### Paso 3: Escritura directa de archivos aprobados

Después del `preflight` y **aprobación explícita** del árbol visual propuesto en Paso 1:

1. **Crear directorios padre**: si una task pertenece a un Outcome (ej: `OXX-slug/TXXX-task.md`), crear el directorio `OXX-slug/` si no existe.

2. **Escribir archivos en paralelo** usando Write tool:
   - Outcome `README.md` con template `outcome-guide.md`: frontmatter `tipo: outcome`, título, descripción/contexto (SIN `## Criterios de Aceptación` persistidos, SIN `## Tasks` persistida — son vistas calculadas).
   - Task `TXXX-*.md` con template `task-guide.md`: frontmatter `estado: Specified`, título, descripción, ACs en `## Criterios de Aceptación`, contexto, scope, hard blockers si aplican.

3. **Validación local post-write**:
   ```bash
   rootline validate <path-creado>
   ```
   Ejecutar después de escribir cada archivo crítico; si `rootline validate` falla, reportar error antes de continuar.

4. **Postcheck obligatorio**: tras escribir todos los archivos,
   ```bash
   roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
   ```
   Si falla, detener y reportar diagnostics. No declarar éxito ni commitear hasta resolverlo.

### Paso 4: Commit + push

- `git add` solo archivos `.md` y `.stem` creados/modificados del roadmap.
- `git commit -m "chore(roadmap): create planning docs"`
- `git push` si `<auto-push>` es true.

STOP obligatorio. Informar: "Archivos de planificación creados. Ejecutar `/roadmap loop` cuando esté listo para implementar."
