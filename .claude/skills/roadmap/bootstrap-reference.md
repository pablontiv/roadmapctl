> Cargar cuando el subcomando va a escribir, mutar o ejecutar tasks.

# Bootstrap Reference

## Paso 0: Detectar modo

```bash
test -d .git
```

- Sí → single-repo mode.
- No → workspace mode.

## Fuente primaria de contexto

`roadmapctl bootstrap` resuelve configuración efectiva, helpers y comportamiento operacional. `.roadmapctl.toml` dentro de `<roadmap-root>/` es la configuración canónica.

Gate inicial:

```bash
command -v roadmapctl
```

Ejecutar para cada repo objetivo:

```bash
roadmapctl bootstrap --repo <repo-path> --output json
```

Usar el JSON devuelto como fuente de verdad para:

- `<repo-path>` = `root`
- `<abs-roadmap-root>` = `roadmap_root`
- `<roadmap-root>` = path relativo desde `root` a `roadmap_root`
- `<where-leaf>` = `helpers.where_leaf`
- `<where-not-done>` = `helpers.where_not_done`
- `<where-active>` = `helpers.where_active`
- status/config operacional = campos `status_values`, `done_statuses`, `active_statuses`, `outcome_close_verify`, `pr_merge_strategy`, `commit_style`, `auto_push`, `required_code_coverage`, `loop_max_tasks`, `parallel`, `autonomy`, `compact_after_task_commit`, `pr_mode` y cualquier campo adicional expuesto

Para el schema completo de `estado`, llamar `rootline describe` directamente si se necesita.

`roadmapctl doctor` y `roadmapctl check` no forman parte del bootstrap read-only. Ejecutarlos solo antes de escribir, mutar, ejecutar tasks o declarar validez del roadmap, y como postcheck después de materializar o mutar.

Si `roadmapctl bootstrap` falla o `roadmapctl` no existe:

- Para flujos implementados (read-only, writes, mutaciones, ejecución, validez): detenerse; no fallback.
- Para planificación conceptual sin writes: se permite usar defaults explícitos, marcando claramente que los guards faltan.

## Workspace mode

1. Escanear subdirectorios inmediatos con `.git` + config roadmap (`<roadmap-root>/.roadmapctl.toml`).
2. Para cada repo, ejecutar `roadmapctl bootstrap` y calcular helpers desde su JSON.
3. Imprimir checkpoint con repos detectados.

Cada repo mantiene su propio roadmap bajo `<repo>/docs/roadmap/`. El skill opera por repo:

- `/roadmap loop` se invoca sobre un repo a la vez. Usar `--repo <name>` o ejecutar dentro del repo.
- Cada repo gestiona su propio commit/push según `auto_push` en su `.roadmapctl.toml`.
- No existe routing de commits cross-repo.

## Single-repo mode

1. Resolver repo actual.
2. Ejecutar `roadmapctl bootstrap`.
3. Si bootstrap no está disponible y el flujo es conceptual, usar defaults explícitos marcados como no verificados.
4. Imprimir checkpoint desde JSON de `roadmapctl bootstrap`.

## Template mínimo `.roadmapctl.toml`

```yaml
---
roadmap-root: # preguntar al usuario
done-statuses: ['Completed', 'Obsolete']
active-statuses: ['Pending', 'Specified', 'In Progress']
status-values:
  pending: 'Pending'
  specified: 'Specified'
  in-progress: 'In Progress'
  completed: 'Completed'
  blocked: 'Blocked'
  obsolete: 'Obsolete'
leaf-filter: 'isIndex == false'
outcome-close-verify: []
pr-merge-strategy: 'squash'
commit-style: 'conventional'
auto-push: true
---
```
