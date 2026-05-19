---
estado: Completed
tipo: task
---
# T001: Drop `--roadmap-root` flag from Go code

**Outcome**: [O22 Drop --roadmap-root flag and document workspace](README.md)
**Contribuye a**: superficie CLI sin flag noise; `config.Load` simplificado a constante hardcodeada `docs/roadmap`

## Preserva

- INV1: `go test ./...` verde
  - Verificar: `cd /home/shared/roadmapctl && go test ./...`
- INV2: `go build ./...` limpio
  - Verificar: `cd /home/shared/roadmapctl && go build ./...`
- INV3: `golangci-lint run ./...` con 0 issues
  - Verificar: `cd /home/shared/roadmapctl && golangci-lint run ./...`
- INV4: `roadmapctl bootstrap --repo . --roadmap-root docs/roadmap` sigue funcionando para roadmap actual del propio repo (último uso documentado del default — debe seguir aceptándolo como argumento posicional? NO: con el drop, el flag no existe; el skill deberá invocar sin el flag — ver T002)
  - Verificar: tras T001+T002, `roadmapctl bootstrap --repo . --output json` (sin `--roadmap-root`) carga config desde `docs/roadmap/.roadmapctl.toml`

## Contexto

Investigación con backscroll confirmó cero usos reales del flag con valor distinto al default `docs/roadmap`. El flag es plumbing inútil: definido en `internal/cli/cli.go:69`, propagado a `Options.RoadmapRoot` (cli.go:25), reenviado vía `config.Options.RoadmapRoot` a `config.Load()` (`internal/config/config.go:103`) en 10+ call sites de CLI, y consumido en línea 112 de `config.Load` solo para construir `tomlPath`.

`config.Load` hoy:

```go
roadmapRoot := opts.RoadmapRoot
if strings.TrimSpace(roadmapRoot) == "" {
    roadmapRoot = filepath.ToSlash(filepath.Join("docs", "roadmap"))
}
```

Tras T001, hardcodear directamente:

```go
const roadmapRoot = "docs/roadmap"
```

La rama de migración legacy (`config.go:140-141, 169`) también usa `opts.RoadmapRoot`. Simplificar:

- Línea 140-142: eliminar el override desde `opts.RoadmapRoot` para `stringValue(fields["roadmap-root"])`
- Línea 169: eliminar la rama que distingue por presencia de `opts.RoadmapRoot`
- `LegacyMigrationPlan` (config.go:227-253) también consume `opts.RoadmapRoot` en líneas 241-243; simplificar análogamente

Después del drop, la firma `Load(repo string, opts Options)` queda con `Options` vacío. Decisión: eliminar el parámetro `Options` y la struct entera (`internal/config/config.go:41-43`), simplificando a `Load(repo string) (*Config, error)` y `LegacyMigrationPlan(repo string) (MigrationPlan, error)`.

Call sites a actualizar (10 sitios identificados con `awk` sobre fuentes Go):

- `internal/cli/doctor.go:28`
- `internal/cli/check.go:14`
- `internal/cli/lint.go:14`
- `internal/cli/pending.go:4763`, `pending.go:4785` (último ya pasa Options vacío)
- `internal/cli/next.go:2454`
- `internal/cli/decision.go:5839`
- `internal/cli/transition.go:4313`
- `internal/cli/bootstrap.go:2747`, `bootstrap.go:2771`, `bootstrap.go:3010`
- `internal/config/config_test.go:9217` (test de `LegacyMigrationPlan`)

Tests obsoletos a eliminar de `internal/config/config_test.go`:

- `TestLoadUsesRoadmapRootOverrideForTOMLDiscovery` (líneas ~2350)
- `TestLoadRoadmapRootOverride` (líneas ~2674)

Referencia residual en `internal/cli/cli.go:139`: `diagnostics.NewReport("roadmapctl/"+name, options.Repo, options.RoadmapRoot, nil)`. Reemplazar `options.RoadmapRoot` por `""` (el report es defensivo, sobreescrito por doctor/check antes de renderizar).

## Alcance

**In**:
1. Eliminar registro del flag en `internal/cli/cli.go:69`
2. Eliminar campo `RoadmapRoot string` de `cli.Options` (cli.go:25)
3. Eliminar referencia residual a `options.RoadmapRoot` en cli.go:139 (reemplazar por `""`)
4. Eliminar struct `config.Options` entera (config.go:41-43)
5. Simplificar `config.Load` para usar `const roadmapRoot = "docs/roadmap"` hardcodeado
6. Simplificar `LegacyMigrationPlan` quitando dependencia de `opts.RoadmapRoot`
7. Actualizar las 10 call sites de `config.Load` y `LegacyMigrationPlan` para quitar el segundo argumento
8. Eliminar tests obsoletos (`TestLoadUsesRoadmapRootOverrideForTOMLDiscovery`, `TestLoadRoadmapRootOverride`)
9. Actualizar `TestLegacyMigrationPlanGeneratesTOMLWithoutWriting` para quitar `Options{}` argument

**Out**:
- Cambios al skill markdown (`.claude/skills/roadmap/*.md`) — son responsabilidad de T002
- Cambios a `docs/cli-contract.md` — responsabilidad de T004
- Cambios a `README.md` y `SKILL.md` Paso 0 — responsabilidad de T003
- Cualquier feature nueva relacionada a workspace, code_repos, modo minimal: explícitamente fuera de scope

## Estado inicial esperado

- `roadmapctl --help` muestra `--roadmap-root` como persistent flag
- `internal/cli/cli.go:69` registra el flag
- `internal/config/config.go:42` declara `Options.RoadmapRoot string`
- 10 call sites pasan `config.Options{RoadmapRoot: options.RoadmapRoot}`
- `go test ./...` verde

## Criterios de Aceptación

- `roadmapctl --help 2>&1 | grep -c "roadmap-root"` retorna `0`
- `roadmapctl doctor --roadmap-root foo --output json 2>&1 | grep -c "unknown flag"` retorna `≥ 1`
- `roadmapctl doctor --repo . --output json` (sin `--roadmap-root`) carga config correctamente y retorna `summary.status == "ok"` contra el roadmap actual del repo
- `go test ./...` verde
- `go build ./...` limpio
- `golangci-lint run ./...` reporta 0 issues
- Buscar `RoadmapRoot:` en `internal/cli/` y `internal/config/` (excluyendo `RoadmapRoot:` como field de structs de output JSON/diagnostics y `RoadmapRoot string` en `Config`/`Report`) no retorna ninguna referencia a la struct `config.Options.RoadmapRoot` eliminada
- Buscar `config.Options{` en `internal/cli/` retorna 0 matches

## Fuente de verdad

- `internal/cli/cli.go` — registro de flag, struct `Options`
- `internal/config/config.go` — struct `Options`, función `Load`, función `LegacyMigrationPlan`
- `internal/config/config_test.go` — tests a eliminar/ajustar
- `internal/cli/{doctor,check,lint,pending,next,decision,transition,bootstrap}.go` — call sites
