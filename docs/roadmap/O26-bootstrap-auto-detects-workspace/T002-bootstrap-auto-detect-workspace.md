---
estado: Completed
tipo: task
---
# T002: Auto-detect workspace mode in `roadmapctl bootstrap`

**Outcome**: [O26 Bootstrap auto-detects workspace mode](README.md)
**Contribuye a**: Que `roadmapctl bootstrap` invocado desde la raíz de un workspace devuelva la estructura agregada sin requerir `--workspace` explícito.

[[blocked_by:./T001-extract-workspace-detection.md]]

## Preserva

- INV1: Salida JSON para repos single (con `.git/` propio + config) idéntica a la actual.
  - Verificar: tests existentes en `bootstrap_test.go` pasan sin modificarse.
- INV2: Pasar `--workspace` explícito en un repo single sigue devolviendo modo workspace (no se rompe la rama actual `bootstrap inspect --workspace`).
- INV3: `go build ./... && go vet ./...` pasa desde `/home/shared/roadmapctl`.

## Contexto

`buildBootstrapConfig` (`internal/cli/bootstrap.go:244-281`) llama incondicionalmente a `config.Load(options.Repo)` y, ante error, emite `RMC_CONFIG_MISSING`. Nunca consulta `options.Workspace` ni inspecciona el filesystem.

La estructura `bootstrapConfigReport` actualmente no expone un slot `repos[]` para workspace (la rama de inspect lo construye en otro tipo). Hay dos opciones:
- (a) Agregar campo `Repos []RepoEntry` con `json:"repos,omitempty"` a `bootstrapConfigReport`.
- (b) Devolver un tipo nuevo `bootstrapWorkspaceReport` que comparta `kind`, `version`, `summary`, `diagnostics` con el report single.

Preferida: (a), porque mantiene un solo `kind` (`roadmapctl/bootstrap`) consistente con lo que ya consume el skill. La forma del JSON queda alineada con el golden de `bootstrap inspect --workspace` (`testdata/fixtures/valid-workspace/context-valid-workspace.json`).

## Alcance

**In**:
1. Importar `roadmapctl/internal/workspace` en `internal/cli/bootstrap.go`.
2. En `buildBootstrapConfig`, antes de `config.Load`:
   ```go
   if options.Workspace || workspace.IsWorkspaceRoot(root) {
       return buildBootstrapWorkspaceConfig(ctx, root, options)
   }
   ```
3. Implementar `buildBootstrapWorkspaceConfig(ctx, root, options) bootstrapConfigReport`:
   - Itera `workspace.MemberRoots(root)`.
   - Por cada member, llama `config.Load(member)`. Si carga OK, agrega una entrada al slice `Repos` con `name` (basename del path), `root` absoluto, `roadmap_root`, `config_path` relativo, `config_source="toml"`.
   - Si `config.Load` falla para un member, lo omite del slice `Repos` y emite un diagnóstico `info` con el path del member y el error subyacente (no aborta).
   - Si `len(Repos) == 0`, devuelve un report con un diagnóstico `error` con id `RMC_WORKSPACE_EMPTY`, severity `error`, mensaje `"workspace root detected but no member has a valid roadmap config"`, path `root`.
   - Si `len(Repos) >= 1`, `summary.status="ok"`, `config_source="workspace"`, `roadmap_root=""` (no aplica a nivel workspace).
4. Agregar campo `Repos []bootstrapRepoEntry \`json:"repos,omitempty"\`` (o nombre análogo siguiendo la convención del struct existente) a `bootstrapConfigReport`. Definir `bootstrapRepoEntry` con los campos descriptos en (3).
5. Registrar el nuevo error id `RMC_WORKSPACE_EMPTY` donde se registran los otros (`internal/config/config.go` o `internal/diagnostics/`, según convención del repo).

**Out**:
- No tocar `bootstrap inspect` ni `bootstrap init`. Comparten el código de detección a través del paquete `workspace` pero su salida y semántica no cambian en esta task.
- No agregar lógica de carga recursiva de defaults workspace-level. Cada member sigue siendo autosuficiente.
- No modificar `pending.go` (ya migrado en T001).

## Estado inicial esperado

- T001 mergeada: `internal/workspace/{detect.go,detect_test.go}` existen.
- `grep -n "options.Workspace" /home/shared/roadmapctl/internal/cli/bootstrap.go` retorna solo declaración/flag, sin uso real en `buildBootstrapConfig`.
- `roadmapctl bootstrap --repo /home/shared --output json` retorna `RMC_CONFIG_MISSING` (estado actual del bug).

## Criterios de Aceptación

- `roadmapctl bootstrap --repo /home/shared --output json | jq -r .summary.status` imprime `"ok"`.
- `roadmapctl bootstrap --repo /home/shared --output json | jq -r .config_source` imprime `"workspace"`.
- `roadmapctl bootstrap --repo /home/shared --output json | jq '.repos | length'` imprime un entero ≥7.
- `roadmapctl bootstrap --repo /home/shared/roadmapctl --output json | jq -r .config_source` sigue imprimiendo `"toml"` (no se rompe el single-repo).
- En un directorio temporal sin `.git/` y sin members con config: `roadmapctl bootstrap --repo <tmp> --output json | jq -r '.diagnostics[0].id'` imprime `"RMC_WORKSPACE_EMPTY"`.
- `go build ./... && go vet ./...` pasa.
- `go test ./internal/cli/...` pasa con los tests pre-existentes intactos.

## Fuente de verdad

- `/home/shared/roadmapctl/internal/cli/bootstrap.go` (modificado)
- `/home/shared/roadmapctl/internal/workspace/detect.go` (consumido, no modificado)
- `/home/shared/roadmapctl/internal/config/config.go` (posible: nuevo error id)
- `/home/shared/roadmapctl/testdata/fixtures/valid-workspace/` (referencia para shape del JSON workspace)
