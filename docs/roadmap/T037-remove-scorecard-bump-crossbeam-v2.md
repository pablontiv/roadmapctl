---
estado: Specified
tipo: task
---
# T037: Remover workflow Scorecard local y bumpear a `crossbeam@v2`

**Contribuye a**: eliminar 100% de los startup_failure de Scorecard en roadmapctl (8/8 recientes) y heredar los cambios saneados de crossbeam (coverage default=0).

## Preserva

- INV1: CodeQL, gitleaks y multi-platform smoke tests siguen funcionando.
  - Verificar: `gh run list --repo pablontiv/roadmapctl --workflow codeql.yml --limit 3` retorna success post-merge.
- INV2: el smoke job multi-plataforma sigue validando coverage (patrón O17 T004 ya validado en este repo).

## Contexto

roadmapctl usa `pablontiv/crossbeam@v1` y mantiene un workflow local de Scorecard que ha fallado 8/8 veces en startup. Una vez `crossbeam@v2` publicado (sin scorecard en el set reusable), roadmapctl debe:
1. Eliminar el workflow scorecard local (o la referencia al reusable dentro de `ci.yml`).
2. Bumpear `@v1` → `@v2`.

Dependencia cross-repo: requiere `crossbeam@v2` publicado. Se valida en AC.

## Alcance

**In**:
1. Eliminar `.github/workflows/scorecard.yml` si existe (o la referencia al reusable scorecard dentro de `ci.yml`).
2. Bumpear todas las referencias `pablontiv/crossbeam/...@v1` a `@v2`.
3. Actualizar documentación si lista scorecard como workflow activo.

**Out**:
- No tocar `internal/`, `cmd/`, ni los smoke tests.
- No bajar el coverage gate del smoke job (sigue siendo la fuente de verdad de cobertura).

## Estado inicial esperado

- `crossbeam@v2` publicado.
- `.github/workflows/` referencia `@v1` y/o llama scorecard.

## Criterios de Aceptación

- `grep -rE 'pablontiv/crossbeam/.*@v1' .github/workflows/` retorna 0 matches.
- `grep -rE '(^| )scorecard' .github/workflows/` retorna 0 matches en workflows activos (sí permitido en comentarios histórico).
- Próximo push a master termina sin startup_failure: `gh run list --repo pablontiv/roadmapctl --branch master --limit 5` no muestra `startup_failure`.

## Fuente de verdad

- `/home/shared/roadmapctl/.github/workflows/`
- `/home/shared/roadmapctl/CLAUDE.md` (si lista CI workflows)
