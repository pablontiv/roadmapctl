# /roadmap pending

Vista filtrada de trabajo pendiente. Muestra tasks pendientes usando `roadmapctl` como capa determinística de roadmap.

Ruta normal autosuficiente: usar el `roadmapctl bootstrap` ya obtenido en bootstrap y luego `roadmapctl pending`. No leer `common-logic.md`, documentación de integración ni archivos Rootline para este flujo. No ejecutar `roadmapctl doctor`/`check`: pending es read-only y `roadmapctl pending` es la fuente canónica.

## Workspace mode

Si `<repos>` existe o bootstrap detectó workspace:

```bash
roadmapctl pending --workspace --repo <workspace-root> --output json
```

Si `--repo <name>` ya fue resuelto en bootstrap, ejecutar single-repo sobre ese repo.

Renderizar desde el JSON:

- `kind` debe ser `roadmapctl/pending`.
- `repos[]` agrupa por repo en workspace.
- `count` es el total pendiente.
- `tasks[]` contiene `path`, `outcome_path` y `status`.
- Si `summary.status != "ok"`, detenerse y reportar `diagnostics`.

## Single-repo

```bash
roadmapctl pending --repo <repo> --output json
```

Reglas:

- Si `summary.status != "ok"`, detenerse y reportar `diagnostics`.
- No llamar `roadmapctl doctor` ni `roadmapctl check` para pending.
- No llamar `rootline tree` directamente para pending.
- No parsear tablas.
- No ejecutar `rootline stats`.
- No postprocesar JSON crudo de Rootline para producir la vista pending.
- No recalcular `done_statuses`, `leaf_filter` o agrupación en prompt; esa lógica pertenece a `roadmapctl pending`.

## Rendering: Output with tasks table

La tabla de tasks se genera desde el JSON de `roadmapctl pending`:

```
Pending tasks (count: <count>)
┌─────┬───────────┬──────────┐
│ ID  │ Outcome   │ Status   │
├─────┼───────────┼──────────┤
│ ... │ ...       │ ...      │
└─────┴───────────┴──────────┘
```

## Footer gitflow

Después de la tabla de tasks, mostrar el footer de configuración gitflow obtenido desde `roadmapctl bootstrap`:

```
── Gitflow del proyecto ──────────────
Base branch:    <base_branch>
Modo:           <descripción de pr_mode + auto_push>
Autonomía:      <autonomy>
commit_style:   <commit_style>

branch_style:   "<valor del TOML o 'no configurado'>"
pr_title_style: "<valor o 'no aplica — pr_mode=false'>"
pr_body_style:  "<valor o 'no aplica — pr_mode=false'>"

Para refrescar:  /roadmap bootstrap
```

Si `RMC_GITFLOW_NOT_CONFIGURED` está presente en los diagnostics del bootstrap anterior, mostrar aviso al final del footer:

```
⚠ gitflow no configurado — ejecutar /roadmap bootstrap para configurar los style fields
```

El footer es informativo: no bloquea la ejecución de `pending` ni reemplaza la tabla de tasks. Leer los valores desde el JSON serializado de `roadmapctl bootstrap` en el workspace o repo actual.
