---
estado: Completed
tipo: task
---
# T025: Purgar toda referencia a `roadmap.local.md` del código

**Contribuye a**: eliminar definitivamente el legacy `roadmap.local.md` del codebase; la única configuración soportada es `<roadmap-root>/.roadmapctl.toml`.

## Preserva

- INV1: `roadmapctl bootstrap`, `doctor` y `check` siguen funcionando sobre repos que ya migraron a TOML.
  - Verificar: `roadmapctl bootstrap --repo . --output json` + `roadmapctl doctor --repo . --output json --strict` + `roadmapctl check --repo . --output json --strict` retornan `status: ok`.
- INV2: `go build ./...`, `go test ./...` y `golangci-lint run ./...` siguen pasando.
  - Verificar: ejecutar los tres comandos desde la raíz del repo.

## Contexto

El usuario ha pedido reiteradamente eliminar toda referencia a `roadmap.local.md` del código. Aún quedan referencias en 5 archivos de producción y tests:

- `internal/config/config.go` — `Load()` con rama legacy + migración + parser YAML completo + `LegacyMigrationPlan`.
- `internal/cli/bootstrap.go:476` — `configSource()` ramifica por basename.
- `internal/roadmap/status.go:68` — fallback path en diagnóstico.
- `internal/config/config_test.go` — 4 tests legacy + helper `writeConfig` + 4 tests que usan ese helper.
- `internal/diagnostics/report_test.go:41,50` — path literal en fixtures.

El código de migración legacy ya no es necesario: cualquier repo activo ya migró hace tiempo. No se requiere backwards-compat ni shims.

Funciones que solo existen para soportar legacy y deben eliminarse junto al código que las usa:
- `MigrationPlan` (tipo), `LegacyMigrationPlan`, `loadLegacyFields`
- `applyFields`, `parseFrontmatter`, `parseYAMLLines`, `parseScalar`, `parseInlineStringList`, `unquote`
- `stringValue`, `stringValueOK`, `stringSliceValue`, `boolValue`, `intValue`, `floatValue`
- `renderTOMLConfig` (solo lo usaba la migración)
- `WarnConfigConflict` (dead code)

Después de eliminar legacy, `Load()` queda con tres ramas: TOML existe → cargar; roadmap-root no existe → `ErrConfigMissing`; default → usar defaults con `ConfigPath = tomlPath`.

## Alcance

**In**:
1. Eliminar de `internal/config/config.go` toda la lógica legacy: rama del switch en `Load()`, `LegacyMigrationPlan`, `loadLegacyFields`, `applyFields`, `parseFrontmatter` y helpers asociados, `renderTOMLConfig`, `MigrationPlan`, `WarnConfigConflict`.
2. En `Load()`, cambiar el `Path: legacyPath` del error `ErrConfigMissing` por `Path: tomlPath`.
3. Simplificar `configSource()` en `internal/cli/bootstrap.go` para que solo retorne `"toml"` o `"defaults"`.
4. En `internal/roadmap/status.go`, eliminar el fallback `path = ".claude/roadmap.local.md"`; si `status.Path` está vacío, dejarlo vacío en el `Diagnostic`.
5. En `internal/config/config_test.go`: eliminar `TestLoadLegacyOnlyMigratesToTOMLAndDeletesLegacy`, `TestLoadExistingTOMLDeletesLegacyWithoutConflictWarning`, `TestLoadInvalidTOMLDoesNotFallbackToLegacy`, `TestLegacyMigrationPlanGeneratesTOMLWithoutWriting`, y el helper `writeConfig`.
6. Migrar tests no-legacy que usan `writeConfig`:
   - `TestLoadResolvesValidRoadmapRootInsideRepo`: usar `writeRoadmapctlTOML` y borrar la aserción que comprueba ausencia de legacy file.
   - `TestLoadRejectsParentEscape`: como `roadmap-root` ya no se lee de YAML (es `roadmapRootDir = "docs/roadmap"` hardcoded), **eliminar el test**.
   - `TestLoadAcceptsWindowsStyleSeparatorsInRoadmapRoot`: por la misma razón, **eliminar**.
   - `TestLoadAppliesDocumentedDefaultsAndParsesOverrides`: reescribir usando `writeRoadmapctlTOML` con overrides en TOML syntax; eliminar si se vuelve redundante con tests TOML existentes.
7. `TestConfigErrorFormatsPathAndUnwrapsCause`: cambiar el path literal del fixture y del mensaje esperado de `.claude/roadmap.local.md` a `docs/roadmap/.roadmapctl.toml`.
8. En `internal/diagnostics/report_test.go`: reemplazar `.claude/roadmap.local.md` por `docs/roadmap/.roadmapctl.toml` en ambas ocurrencias (líneas 41 y 50).

**Out**:
- NO tocar `docs/roadmap/**` ni `docs/superpowers/**` ni `context-gather/**`: los docs históricos son registro de migraciones ya hechas.
- NO añadir flags o comandos CLI nuevos.
- NO crear ningún archivo `*.go` nuevo (solo eliminar/modificar existentes).

## Estado inicial esperado

- `roadmapctl bootstrap --repo /home/shared/roadmapctl --output json` retorna `config_source: "toml"`.
- `grep -rln "roadmap\.local" /home/shared/roadmapctl --include="*.go"` lista los 5 archivos del Contexto.

## Criterios de Aceptación

- AC1: `find /home/shared/roadmapctl -name "*.go" -not -path "*/dist/*" | xargs grep -l "roadmap\.local"` no imprime nada (cero hits en código Go).
- AC2: En `internal/config/config.go` no existen los identificadores `MigrationPlan`, `LegacyMigrationPlan`, `loadLegacyFields`, `applyFields`, `parseFrontmatter`, `parseYAMLLines`, `parseScalar`, `parseInlineStringList`, `unquote`, `stringValue`, `stringValueOK`, `stringSliceValue`, `boolValue`, `intValue`, `floatValue`, `renderTOMLConfig`, `WarnConfigConflict`.
- AC3: En `internal/config/config_test.go` no existen los tests `TestLoadLegacyOnlyMigratesToTOMLAndDeletesLegacy`, `TestLoadExistingTOMLDeletesLegacyWithoutConflictWarning`, `TestLoadInvalidTOMLDoesNotFallbackToLegacy`, `TestLegacyMigrationPlanGeneratesTOMLWithoutWriting`, ni la función helper `writeConfig`.
- AC4: `configSource()` en `internal/cli/bootstrap.go` solo retorna `"toml"` o `"defaults"`; no contiene la cadena `"roadmap.local.md"` ni el valor `"legacy"`.
- AC5: `internal/roadmap/status.go` no contiene la cadena `".claude/roadmap.local.md"`.
- AC6: `internal/diagnostics/report_test.go` no contiene la cadena `".claude/roadmap.local.md"`.
- AC7: `cd /home/shared/roadmapctl && go build ./...` exit code 0.
- AC8: `cd /home/shared/roadmapctl && go test ./...` exit code 0.
- AC9: `cd /home/shared/roadmapctl && golangci-lint run ./...` exit code 0.
- AC10: `roadmapctl doctor --repo /home/shared/roadmapctl --output json --strict` retorna `summary.status: ok` post-cambios.

## Fuente de verdad

- `/home/shared/roadmapctl/internal/config/config.go`
- `/home/shared/roadmapctl/internal/config/config_test.go`
- `/home/shared/roadmapctl/internal/cli/bootstrap.go`
- `/home/shared/roadmapctl/internal/roadmap/status.go`
- `/home/shared/roadmapctl/internal/diagnostics/report_test.go`
