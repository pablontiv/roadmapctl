---
estado: Completed
tipo: task
---
# T002: Skip test de chmod en Windows

## Descripción

**RC2 — `TestBootstrapInitApplyReportsDiagnosticsOnFileError` falla en Windows**

El test crea un directorio, aplica `os.Chmod(path, 0o555)` y espera que
`bootstrap init --apply` falle al intentar crear subdirectorios.

En Linux/macOS, `chmod 0o555` impide escritura → el proceso sale con código ≠ 0 ✓

En Windows, `os.Chmod` no tiene el mismo efecto sobre la creación de subdirectorios.
El directorio sigue siendo escribible → el proceso sale con código 0 → test falla:
```
bootstrap_test.go:212: bootstrap init with permission error should fail, got exit = 0
```

**Fix:** Skipear el test en Windows con `runtime.GOOS == "windows"`.

```go
func TestBootstrapInitApplyReportsDiagnosticsOnFileError(t *testing.T) {
    if runtime.GOOS == "windows" {
        t.Skip("chmod 0o555 does not prevent directory creation on Windows")
    }
    ...
}
```

## Criterios de Aceptación

- `TestBootstrapInitApplyReportsDiagnosticsOnFileError` pasa en Linux y macOS
- En Windows el test se reporta como SKIP (no FAIL)
- `go test ./internal/cli/...` pasa sin errores
