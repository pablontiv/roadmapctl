---
estado: Completed
tipo: task
---
# T002: Activar fake rootline en TestMain cuando rootline no está en PATH

**Outcome**: [O14 CI verde](README.md)
**Contribuye a**: job `ci / Test` verde en CI

## Preserva

- INV1: Tests con rootline real siguen funcionando cuando rootline está en PATH.
  - Verificar: `go test ./...` (con rootline instalado localmente)

## Contexto

El job `ci / Test` de crossbeam corre `go test -race ./...` sin instalar rootline.
Unos 20 tests fallan con `RMC_ENV_ROOTLINE_MISSING` porque los comandos `check`,
`doctor`, `pending`, `next`, `decision` invocan rootline internamente.

Ya existe un mecanismo de fake rootline: si `ROADMAPCTL_FAKE_ROOTLINE=1` y
`ROOTLINE_BIN=<path>`, `TestMain` re-ejecuta el binario de tests como fake
rootline que responde con JSON stub a `validate`, `describe`, `query`, `graph`,
`tree`, `set`, `new`. Este mecanismo se activa manualmente en
`TestCheckUsesRootlineBinEnvironmentOverride`.

La solución es detectar en `TestMain` si rootline no está disponible y activar
el fallback automáticamente, sin afectar ejecuciones locales donde rootline sí
está instalado.

## Alcance

**In**:
1. Agregar `"os/exec"` a los imports de `internal/cli/golden_test.go`
2. En `TestMain`, después de crear los dirs `.git`, llamar
   `exec.LookPath("rootline")`; si falla, hacer
   `os.Setenv("ROADMAPCTL_FAKE_ROOTLINE", "1")` y
   `os.Setenv("ROOTLINE_BIN", os.Args[0])`

**Out**:
- No modificar `fakeRootline()` ni los stubs existentes
- No instalar rootline en el job de crossbeam

## Estado inicial esperado

- `TestMain` ya crea los dirs `.git` de fixtures (commit `8971a1b`)
- El mecanismo `ROADMAPCTL_FAKE_ROOTLINE` existe y funciona
- ci / Test falla con ~20 tests `RMC_ENV_ROOTLINE_MISSING`

## Criterios de Aceptación

- `go test ./...` pasa localmente (con rootline real en PATH)
- job `ci / Test` verde en CI (sin rootline instalado)
- `TestCheckUsesRootlineBinEnvironmentOverride` sigue pasando

## Fuente de verdad

- `internal/cli/golden_test.go` (TestMain + fakeRootline)
