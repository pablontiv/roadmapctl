---
estado: Pending
tipo: task
---
# T001: Adoptar Cobra y paquetes comunitarios

**Outcome**: [O02 Fundaciones post-MVP](README.md)
**Contribuye a**: CE3, INV2, INV3

## Preserva

- INV1: JSON stdout sigue siendo un único objeto parseable.
  - Verificar: golden tests existentes.
- INV2: `Execute(args, stdout, stderr) int` sigue testeable.
  - Verificar: tests de `internal/cli`.

## Contexto

La CLI actual usa `flag` estándar y solo despacha `doctor`/`check`. La superficie aprobada crecerá a `context`, `pending`, `next`, `decision`, `lint`, `transition ...`, `materialize`, por lo que conviene adoptar paquetes mantenidos por la comunidad antes de multiplicar routing propio.

## Alcance

**In**:
1. Evaluar y adoptar `github.com/spf13/cobra` para estructura de comandos.
2. Preparar dependencias aprobadas: `github.com/pelletier/go-toml/v2`, `github.com/yuin/goldmark`, `github.com/aymanbagabas/go-udiff` donde corresponda.
3. Mantener compatibilidad de `doctor`/`check`, flags globales, exit codes y stdout/stderr.
4. Agregar tests/goldens de help, errores de parsing y JSON mode.

**Out**:
- Implementar comandos post-MVP.
- Cambiar semántica de diagnostics.

## Estado inicial esperado

- `internal/cli/cli.go` usa `flag` estándar.
- `go.mod` no contiene dependencias externas.

## Criterios de Aceptación

- `go test ./...` pasa.
- `go build ./cmd/roadmapctl` pasa.
- `roadmapctl doctor/check` mantienen flags y JSON contract.
- Help lista comandos existentes sin introducir `roadmapctl roadmap ...`.

## Fuente de verdad

- `internal/cli/cli.go`
- `internal/cli/*_test.go`
- `go.mod`
- `docs/cli-contract.md`
