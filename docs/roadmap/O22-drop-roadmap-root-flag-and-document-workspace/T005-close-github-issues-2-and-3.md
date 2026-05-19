---
estado: Completed
tipo: task
---
# T005: Close GitHub issues #2 and #3 with explanatory comments

**Outcome**: [O22 Drop --roadmap-root flag and document workspace](README.md)
**Contribuye a**: dos issues abiertos en GitHub cerrados con trazabilidad a la decisión de diseño documentada

[[blocked_by:./T001-drop-roadmap-root-flag-from-go.md]]
[[blocked_by:./T002-sweep-skill-markdown-for-roadmap-root-flag.md]]
[[blocked_by:./T003-document-workspace-multi-repo-in-readme-and-skill.md]]
[[blocked_by:./T004-update-cli-contract-doc-after-flag-drop.md]]

## Preserva

- INV1: los comentarios de cierre referencian la sección Workspace mode del README (debe existir tras T003)
- INV2: no se cierran issues sin que los cambios documentados estén mergeados/push a la rama default

## Contexto

Los GitHub issues fueron analizados en una investigación documentada en `/home/pones/.claude/plans/tenmos-2-issues-en-iridescent-naur.md`. La conclusión:

- **#2 (multi-repo workspace commit routing)**: WONTFIX como feature; el modelo correcto es que cada repo tenga su propio roadmap. Workspace mode actual (vía `workspaceRepoRoots()`) ya itera repos. No se añade `code_repos = [...]` ni routing cross-repo.
- **#3 (`.roadmapctl.toml` en repo root)**: WONTFIX; `docs/roadmap/` permanece como convención fija. Backscroll confirmó cero casos reales de layouts alternativos.

Cerrar los issues vía `gh issue close` con comentarios explicativos que referencien:

- La nueva sección "Workspace mode" del `README.md` (creada en T003)
- La actualización de `docs/cli-contract.md` (T004) que registra la convención fija
- El cambio de código (drop del flag en T001)

Comentario sugerido para issue #2:

```
Closed as WONTFIX. After investigation, the correct model is that each repo
participating in a workspace maintains its own complete roadmap under
`<repo>/docs/roadmap/`. There is no "code repo without roadmap" scenario;
each repo is autonomous. Cross-repo commit routing is intentionally not
supported — each repo's loop only touches files in its own repo.

The existing workspace mode (`workspaceRepoRoots()` in `internal/cli/pending.go`)
already discovers sibling `.git` repos and loads each repo's config
independently. The skill `/roadmap` is invoked per repo.

This convention is now documented in the README's "Workspace mode" section
and clarified in the `/roadmap` skill's Paso 0.

If you have a scenario where a single roadmap genuinely needs to drive
commits across multiple repos, please open a new issue describing the
specific workflow — the current closure is based on no real-world case
appearing in our investigation.
```

Comentario sugerido para issue #3:

```
Closed as WONTFIX. Investigation with prior sessions confirmed zero
real-world cases where `--roadmap-root` was used with a non-default value;
every observed invocation passed `docs/roadmap` redundantly. The flag has
been dropped (see #<PR-number>) and `docs/roadmap/` is now a fixed
convention, not configurable.

This means `.roadmapctl.toml` lives only at `<repo>/docs/roadmap/.roadmapctl.toml`.
Each repo in a workspace has its own such layout — see the README's new
"Workspace mode" section.

If you have a layout that genuinely requires a different roadmap-root
location, please open a new issue describing the case.
```

Adaptar `#<PR-number>` al PR real que mergee este outcome.

## Alcance

**In**:
1. Verificar que los cambios de T001-T004 están en la rama default (`master`) y push (`git log origin/master` debe contener los commits)
2. `gh issue view 2` y `gh issue view 3` para verificar que siguen abiertos
3. `gh issue close 2 --comment "..."` con el comentario adaptado al issue #2
4. `gh issue close 3 --comment "..."` con el comentario adaptado al issue #3
5. Verificar con `gh issue view <n>` que ambos quedaron cerrados

**Out**:
- Cambios al código o docs — fuera de scope, son tareas anteriores
- Crear PRs adicionales o features nuevas — fuera de scope

## Estado inicial esperado

- T001-T004 Completed
- Cambios de T001-T004 mergeados a `master` y push a `origin`
- `gh issue list --state open` muestra #2 y #3 como abiertos

## Criterios de Aceptación

- `gh issue view 2 --json state -q .state` retorna `CLOSED`
- `gh issue view 3 --json state -q .state` retorna `CLOSED`
- El comentario de cierre de cada issue referencia la sección "Workspace mode" del README y/o el PR que dropeó el flag
- `gh issue list --state open` no incluye #2 ni #3

## Fuente de verdad

- `gh issue` (no es un archivo del repo, es estado en GitHub)
- Comentarios de cierre adaptados a partir del template arriba
