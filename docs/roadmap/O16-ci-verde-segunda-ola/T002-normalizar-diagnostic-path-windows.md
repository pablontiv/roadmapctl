---
estado: Completed
tipo: task
---
# T002: Normalizar Diagnostic.Path a forward slashes en Windows

## Descripción

**RC4 — backslash en `Diagnostic.Path` (`smoke/windows-latest`)**

En `internal/diagnostics/report.go`, `NewReport()` normaliza `Root` y `RoadmapRoot`
con `filepath.ToSlash()` (líneas 108–109), pero NO normaliza el campo `Path` de cada
`Diagnostic`. En Windows, rutas relativas como `O01-work\T001-task.md` retienen el
backslash nativo.

`AssertNoBackslashes` en `golden_test.go:99` falla con:
```
value contains backslash path separator: ... "path":"O01-work\\T001-task.md" ...
```

Fix: en `NewReport()`, después del `copy(copied, diagnostics)`, iterar sobre `copied`
y aplicar `filepath.ToSlash()` a cada `copied[i].Path`.

```go
for i := range copied {
    copied[i].Path = filepath.ToSlash(copied[i].Path)
}
```

## Criterios de Aceptación

- `TestCheckGoldenJSONFixtures/invalid_status_bogus` pasa en Windows (no `AssertNoBackslashes` failure)
- Otros tests con `Diagnostic.Path` no regresionan en Linux/macOS
- Ningún golden file nuevo tiene backslashes en el campo `path`
