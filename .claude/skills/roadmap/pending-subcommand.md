# /roadmap pending

Vista filtrada de trabajo pendiente. Muestra tasks pendientes usando `roadmapctl` como capa determinística de roadmap.

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
roadmapctl pending --repo <repo> --roadmap-root <roadmap-root> --output json
```

Reglas:

- No llamar `rootline tree` directamente para pending.
- No parsear tablas.
- No ejecutar `rootline stats`.
- No postprocesar JSON crudo de Rootline para producir la vista pending.
- No recalcular `done_statuses`, `leaf_filter` o agrupación en prompt; esa lógica pertenece a `roadmapctl pending`.
