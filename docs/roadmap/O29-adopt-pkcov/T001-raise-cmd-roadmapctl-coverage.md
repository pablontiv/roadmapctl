---
estado: Completed
tipo: task
---
# T001: Subir cmd/roadmapctl a ≥85% coverage

**Outcome**: [O29 Adoptar pkcov + reactivar gate](README.md)
**Contribuye a**: precondición para reactivar el gate sin que bloquee push (INV2)

## Preserva

- INV2 del outcome: el pre-work no relaja el contrato — sube el paquete al estándar
  - Verificar: `go test ./cmd/roadmapctl/ -cover` ≥85

## Contexto

`cmd/roadmapctl` está hoy en 50% (medido `2026-05-21`). Es el paquete más débil del repo y arrastra el total. El paquete contiene los runE de cobra para los subcomandos del CLI; los tests existentes prueban helpers internos por debajo del runE pero no los entry points completos.

Patrón de test cobra (mismo molde que rootline usa para sus tests de comandos):
```go
buf := new(bytes.Buffer)
rootCmd.SetOut(buf)
rootCmd.SetArgs([]string{"subcommand", "...args"})
err := rootCmd.Execute()
```

Subcomandos a cubrir: identificar con `go tool cover -func=cov.out | awk '$3<"85.0%"' | grep cmd/roadmapctl/`.

## Alcance

**In**:

1. Identificar runE/funciones de `cmd/roadmapctl` con cobertura <85.
2. Agregar tests cobra que invoquen los subcomandos con args válidos e inválidos.
3. Cubrir error paths (archivos inexistentes, JSON inválido, flags requeridos faltantes).
4. Validar que el paquete cruza 85.

**Out**:
- No refactor de la lógica de comandos.
- No cambiar la firma o output de los comandos.
- No tocar otros paquetes.

## Estado inicial esperado

- `go test ./cmd/roadmapctl/ -cover` reporta ~50%.

## Criterios de Aceptación

- `go test ./cmd/roadmapctl/ -cover` ≥85.0%.
- `go test ./... -race` verde.
- Diff incluye sólo archivos `_test.go` en `cmd/roadmapctl/`.
- `golangci-lint run` sin issues nuevos.

## Fuente de verdad

- `/home/shared/roadmapctl/cmd/roadmapctl/` — runE handlers
- `/home/shared/rootline/cmd/rootline/describe_test.go` — patrón cobra a imitar
