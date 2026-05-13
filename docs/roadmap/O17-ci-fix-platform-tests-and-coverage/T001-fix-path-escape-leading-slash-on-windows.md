---
estado: Completed
tipo: task
---
# T001: Fix ErrPathEscape para paths con slash inicial en Windows

## Descripción

**RC1 — `TestResolveInsideAbsolutePathUnix` y `TestResolveInsideLeadingSlash` fallan en Windows**

En `internal/fsx/path.go`, `ResolveInside` maneja paths que empiezan con `/`:

```go
if strings.HasPrefix(normalized, "/") {
    target = filepath.Clean(filepath.FromSlash(normalized))
}
```

En Linux: `filepath.FromSlash("/absolute/path")` → `/absolute/path` (absoluto).
`filepath.Rel(root, "/absolute/path")` → `"../../absolute/path"` → `ErrPathEscape` ✓

En Windows: `filepath.FromSlash("/absolute/path")` → `\absolute\path` (relativo al drive root).
`filepath.Rel("C:\\Temp\\...", "\\absolute\\path")` → error `"Rel: can't make \absolute\path relative to C:\..."`.
Ese error no es `ErrPathEscape`, así que `errors.Is(err, ErrPathEscape)` falla. ✗

**Fix:** Retornar `ErrPathEscape` directamente para cualquier path con slash inicial,
sin delegar en `filepath.Rel`. Un path que empieza con `/` siempre escapa el repo root.

```go
if strings.HasPrefix(normalized, "/") {
    return "", "", fmt.Errorf("%w: %s", ErrPathEscape, candidate)
}
target := filepath.Join(absRoot, filepath.FromSlash(normalized))
```

## Criterios de Aceptación

- `TestResolveInsideAbsolutePathUnix` pasa en Linux, macOS y Windows
- `TestResolveInsideLeadingSlash` pasa en Linux, macOS y Windows
- `TestResolveInsideUNCPath` sigue pasando (usa `//`, detectado antes del slash check)
- `go test ./internal/fsx/...` pasa sin errores
