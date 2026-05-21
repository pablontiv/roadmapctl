> Cargar cuando el subcomando va a escribir, mutar o ejecutar tasks.

# Bootstrap Reference

## Fuente primaria de contexto

`roadmapctl bootstrap` detecta automáticamente si `<repo>` es un workspace root (contiene siblings con `<sibling>/docs/roadmap/.roadmapctl.toml`) o un single-repo, y resuelve configuración efectiva, helpers y comportamiento operacional. No hace falta `test -d .git` ni escaneo manual de subdirectorios desde el skill.

`.roadmapctl.toml` dentro de `<roadmap-root>/` es la configuración canónica por repo.

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

Cuando `roadmapctl bootstrap --repo <workspace-root>` detecta workspace, el JSON devuelve:

- `config_source: "workspace"` (marcador inequívoco).
- `repos[]` con cada miembro: `name`, `root`, `roadmap_root`, `config_path`, `config_source`.
- `roadmap_root` y `config_path` a nivel top vacíos.

El flag explícito `--workspace` fuerza este modo aunque la auto-detección no lo dispare.

Diagnostics relevantes:

- `RMC_WORKSPACE_EMPTY` (error) — el path es workspace root pero ningún miembro tiene config válido.
- `RMC_WORKSPACE_MEMBER_SKIPPED` (info) — miembro inválido, ignorado.

## Diagnósticos gitflow

- `RMC_GITFLOW_NOT_CONFIGURED` (info) — los campos `[gitflow]` (`branch_style`, `pr_title_style`, `pr_body_style`) están vacíos. El skill puede detectar esto para ofrecer el wizard de adopción de gitflow.
- `RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED` (warning) — `pr_merge_strategy` está seteado en el TOML; migrar a `[gitflow]` fields.

## Wizard de adopción gitflow (RMC_GITFLOW_NOT_CONFIGURED)

Cuando bootstrap devuelve `RMC_GITFLOW_NOT_CONFIGURED`, el skill debe ejecutar el wizard de adopción:

1. **Escaneo abierto del repo**: El LLM decide qué leer según lo que encuentra — cualquier combinación de `.md`, `CLAUDE.md`, `AGENTS.md`, `git log`, `gh pr list --state merged`, `gh pr list --state closed`, branch protection, CI workflows (`.github/workflows/`), código fuente, o cualquier otro artefacto relevante. No hay lista fija de comandos; el LLM usa su criterio para entender las convenciones del proyecto.

2. **Redactar propuesta de TOML**: Con la información recopilada, el LLM propone valores para los 3 style fields:
   - `branch_style`: descripción de la convención de branches
   - `pr_title_style`: descripción del formato de PR titles
   - `pr_body_style`: descripción del template de PR body (o "no aplica" si pr_mode=false)
   
   Los campos determinísticos (`done_statuses`, `active_statuses`, etc.) se deducen del TOML existente o se usan defaults.

3. **Wizard pre-llenado**: Presentar cada campo al usuario para confirmar, editar o saltar. No escribir hasta que el usuario confirme el bloque completo.

4. **Escribir el TOML**: Solo tras confirmación final, actualizar `<roadmap-root>/.roadmapctl.toml` con los 3 style fields bajo `[gitflow]`.

5. **Re-ejecutar bootstrap**: Invocar `roadmapctl bootstrap --repo <repo-path> --output json` y continuar desde el JSON actualizado.

El subcomando `/roadmap bootstrap` fuerza re-escaneo on-demand aunque el TOML ya tenga valores, sobrescribiendo los style fields si el usuario lo confirma.

Cada repo mantiene su propio roadmap bajo `<repo>/docs/roadmap/`. El skill opera por repo:

- `/roadmap loop` se invoca sobre un repo a la vez. Usar `--repo <name>` (resuelto contra `repos[].name`) o ejecutar dentro del repo.
- Cada repo gestiona su propio commit/push según `auto_push` en su `.roadmapctl.toml`.
- No existe routing de commits cross-repo.

## Single-repo mode

`roadmapctl bootstrap` lo detecta cuando `<repo>/docs/roadmap/.roadmapctl.toml` existe. El JSON devuelve `roadmap_root`, `config_path`, helpers y todos los campos operacionales directamente (sin `repos[]`).

Si bootstrap no está disponible y el flujo es conceptual, usar defaults explícitos marcados como no verificados.

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
