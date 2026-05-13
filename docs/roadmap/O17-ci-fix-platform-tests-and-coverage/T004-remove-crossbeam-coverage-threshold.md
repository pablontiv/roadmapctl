---
estado: Completed
tipo: task
---
# T004: Eliminar coverage threshold del job crossbeam

## Descripción

**RC4 — `ci/Test` (crossbeam) falla con coverage ~78% usando fake rootline**

El workflow `pablontiv/crossbeam/.github/workflows/go-ci.yml@v1` no instala
rootline. Cuando `go test ./... -race` corre sin rootline disponible, `TestMain`
en `internal/cli` detecta la ausencia y activa fake rootline automáticamente:

```go
if _, err := exec.LookPath("rootline"); err != nil {
    _ = os.Setenv("ROADMAPCTL_FAKE_ROOTLINE", "1")
    _ = os.Setenv("ROOTLINE_BIN", os.Args[0])
}
```

Con fake rootline activo, los tests marcados `requiresRealRootline` se saltean.
Esto baja `internal/cli` de 81.4% a 56.9% y el total cae a ~78% < 85%.

**Solución:** Poner `coverage-threshold: 0` en el job `ci` de `ci.yml`.
La cobertura ya es validada por el job `smoke` en 3 plataformas (Ubuntu, macOS,
Windows) usando `check-coverage.sh`. El job `ci/Test` de crossbeam sigue corriendo
`go test ./... -race` como gate de race conditions y corrección general, solo sin
el gate de cobertura duplicado.

```yaml
jobs:
  ci:
    uses: pablontiv/crossbeam/.github/workflows/go-ci.yml@v1
    with:
      go-version-file: go.mod
      coverage-threshold: 0   # cobertura validada por smoke job en 3 plataformas
```

## Criterios de Aceptación

- `ci/Test` (crossbeam) pasa aunque la cobertura sea < 85% con fake rootline
- `smoke/ubuntu-latest`, `smoke/macos-latest`, `smoke/windows-latest` validan cobertura con `check-coverage.sh`
- El job `release` llega a ejecutarse cuando los demás jobs pasan
