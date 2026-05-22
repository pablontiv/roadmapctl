---
estado: Completed
tipo: task
---
# T001: Eliminar `pr_merge_strategy` del código Go

**Outcome**: [O31 Fix `/integrate` defaults y limpieza de config muerta](README.md)
**Contribuye a**: superficie de config sin campos muertos; bootstrap JSON refleja sólo lo aplicado

## Preserva

- INV1: build verde sin errores ni warnings nuevos
  - Verificar: `go build ./...` exit 0
- INV2: tests existentes siguen pasando (sin contar los que se eliminan en esta task)
  - Verificar: `go test ./internal/... -count=1` exit 0

## Contexto

El campo `pr_merge_strategy` está en `Config`, `tomlConfig`, default `"squash"`, lectura TOML, comparación en `equal()`, constante `DiagnosticGitflowPRMergeStrategyDeprecated`, y emisión de warning de deprecación. Está también en el output JSON de bootstrap. Y en el template `DefaultRoadmapctlTOML` y en sus tests. Y en golden tests + fixtures.

Ningún consumidor lo aplica: la Fase 6 del skill `/integrate` ejecuta `gh pr merge --auto` sin pasarlo. Decisión: eliminación total sin fallback ni migración. Quien tenga el campo en su TOML local debe quitarlo manualmente; lo que hoy es warning, en el futuro será error de parsing TOML.

## Alcance

**In**:
1. `internal/config/config.go`: borrar campo `PRMergeStrategy` de `Config` y `tomlConfig`; borrar lectura TOML, default `"squash"`, comparación en `equal()`, constante `DiagnosticGitflowPRMergeStrategyDeprecated`, y todo el bloque que emite el warning deprecado.
2. `internal/cli/bootstrap.go`: borrar campo `PRMergeStrategy` de `bootstrapConfigReport` y su asignación en `newBootstrapConfigReport`.
3. `internal/templates/bootstrap.go`: borrar línea `pr_merge_strategy = "squash"` de `DefaultRoadmapctlTOML`.
4. `internal/templates/bootstrap_test.go`: quitar `"pr_merge_strategy"` de la lista de strings esperados.
5. `internal/config/config_test.go`: borrar `TestPRMergeStrategyDeprecatedWarning` y todos los TOMLs sintéticos que escriben `pr_merge_strategy = ...`.
6. `testdata/golden/bootstrap-valid-legacy-config-fallback.json`: regenerar (o editar) quitando el par `"pr_merge_strategy": "squash"`.
7. `testdata/golden/context-valid-workspace.json` y `testdata/golden/context-valid-legacy-config-fallback.json`: borrar ambos archivos (orphans — no hay test que los use; `context` no es subcomando).
8. `testdata/fixtures/invalid-root-escape/.claude/roadmap.local.md`, `invalid-missing-outcome-readme/...`, `invalid-duplicate-ids/...`, `invalid-extra-nesting/...`: quitar línea `pr-merge-strategy: 'squash'` del frontmatter YAML de los 4 fixtures.

**Out**:
- Documentación markdown (skills, docs, CHANGELOG, históricos): T002.
- Cambios al skill `/integrate`: T003.

## Estado inicial esperado

- `grep -rn 'pr_merge_strategy\|PRMergeStrategy' internal/ testdata/` retorna múltiples matches.
- `go build ./...` y `go test ./internal/...` están verdes en master.

## Criterios de Aceptación

- `go build ./...` exit 0 desde `/home/shared/roadmapctl`.
- `go test ./internal/config/... ./internal/cli/... ./internal/templates/... -count=1` exit 0.
- `grep -rn 'pr_merge_strategy\|PRMergeStrategy\|PR_MERGE_STRATEGY_DEPRECATED' internal/ testdata/` retorna vacío.
- `roadmapctl bootstrap --repo <tmp> --output json` en un directorio temporal recién bootstrapped no contiene la clave `pr_merge_strategy` en el JSON.
- `testdata/golden/context-valid-workspace.json` y `testdata/golden/context-valid-legacy-config-fallback.json` no existen.

## Fuente de verdad

- `internal/config/config.go`
- `internal/cli/bootstrap.go`
- `internal/templates/bootstrap.go`
- `internal/templates/bootstrap_test.go`
- `internal/config/config_test.go`
- `testdata/golden/bootstrap-valid-legacy-config-fallback.json`
- `testdata/golden/context-valid-workspace.json` (borrar)
- `testdata/golden/context-valid-legacy-config-fallback.json` (borrar)
- `testdata/fixtures/invalid-root-escape/.claude/roadmap.local.md`
- `testdata/fixtures/invalid-missing-outcome-readme/.claude/roadmap.local.md`
- `testdata/fixtures/invalid-duplicate-ids/.claude/roadmap.local.md`
- `testdata/fixtures/invalid-extra-nesting/.claude/roadmap.local.md`
