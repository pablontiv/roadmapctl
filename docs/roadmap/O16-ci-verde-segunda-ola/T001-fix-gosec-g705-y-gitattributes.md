---
estado: Specified
tipo: task
---
# T001: Fix gosec G705 y añadir .gitattributes

## Descripción

Dos fixes de infraestructura que desbloquean `ci/Lint` y el job `smoke/windows-latest`.

**RC1 — gosec G705 (`ci/Lint`)**
`golangci-lint v2.10.1` con gosec activo reporta G705 (XSS taint analysis) en
`internal/cli/golden_test.go:230`:
```go
fmt.Fprintf(stderr, "unknown fake rootline command %q\n", args[0])
```
Es un falso positivo en código de test. Fix: añadir `//nolint:gosec` en esa línea.

**RC5 — CRLF en goldens Windows (`smoke/windows-latest`)**
Sin `.gitattributes`, git con `core.autocrlf=true` checkoutea los archivos `*.json`
con CRLF. `AssertGoldenJSON` hace `bytes.Equal(want, normalized)` — `want` tiene
`\r\n`, `normalized` tiene `\n`. Visualmente idénticos, bytes distintos → golden
mismatch en `invalid_single_summary_file`, `valid_outcome_with_tasks`, etc.

Fix: crear `.gitattributes` en la raíz del repo con normalización LF para todos
los archivos de texto.

## Criterios de Aceptación

- `golangci-lint run ./...` no reporta G705 ni errores en `golden_test.go`
- Existe `.gitattributes` con `* text=auto eol=lf` y reglas para `.go`, `.json`, `.md`, `.toml`, `.yml`, `.sh`
- `git diff --check` no reporta CRLF en los golden files existentes tras re-checkout
