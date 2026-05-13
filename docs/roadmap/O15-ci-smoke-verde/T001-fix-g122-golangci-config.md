---
estado: Completed
tipo: task
---
# T001: Remover G122 de gosec excludes y añadir nolint en sitios

**Outcome**: [O15 CI smoke verde](README.md)
**Contribuye a**: job `ci / Lint` verde en CI

## Contexto

`G122` no existe en el enum de gosec aceptado por golangci-lint v2.10.1 (el que
usa crossbeam). Se agregó en v2.11+. La validación de config falla con:

```
jsonschema: "linters.settings.gosec.excludes.0" does not validate [...] value must be one of 'G101', ...
```

Las dos líneas afectadas son patrones legítimos que merecen supresión local:
- `internal/roadmap/structure.go`: `os.ReadFile(path)` en callback de `filepath.WalkDir`
- `internal/cli/transition_test.go`: `os.ReadFile`/`os.WriteFile` en helper de copia de fixtures

## Alcance

**In**:
1. Remover `G122` de `linters.settings.gosec.excludes` en `.golangci.yml`
2. Añadir `//nolint:gosec` en `internal/roadmap/structure.go:159`
3. Añadir `//nolint:gosec` en `internal/cli/transition_test.go:293` y `:297`

**Out**:
- No modificar lógica de producción
- No cambiar otras reglas de golangci-lint

## Criterios de Aceptación

- `golangci-lint config verify` sale 0 (válido para v2.10.1)
- `golangci-lint run ./...` reporta 0 issues
- `go test ./...` pasa sin errores
