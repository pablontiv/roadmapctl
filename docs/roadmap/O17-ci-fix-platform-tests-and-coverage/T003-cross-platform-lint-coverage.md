---
estado: Completed
tipo: task
---
# T003: Cobertura cross-platform para funciones de lint de portabilidad

## Descripción

**RC3 — `internal/lint` cae a 75.7% en macOS y Windows → total 83.5% < 85%**

El task O15/T002 skipea `TestCheckFilenamePortabilityReportsCaseCollisionAndReservedName`
en FS insensible a mayúsculas. Esto fue necesario para evitar fallos, pero dejó
`CheckFilenamePortability`, `checkCaseCollisionsInDir`, `reservedWindowsName` y
`lintNameDiagnostic` sin cobertura en macOS y Windows.

**Gap de cobertura en macOS/Windows:**
- `CheckFilenamePortability` → 0%
- `checkCaseCollisionsInDir` → 0%
- `reservedWindowsName` → 0%
- `lintNameDiagnostic` → 0%

**Fix:** Añadir 3 tests en `internal/lint/schema_portability_test.go`:

### Test 1: directorio limpio (todas las plataformas)

```go
func TestCheckFilenamePortabilityNoIssuesOnCleanDir(t *testing.T) {
    root := t.TempDir()
    for _, name := range []string{"T001-task.md", "T002-feature.md"} {
        if err := os.WriteFile(filepath.Join(root, name), nil, 0o644); err != nil {
            t.Fatal(err)
        }
    }
    found, err := CheckFilenamePortability(root)
    if err != nil {
        t.Fatal(err)
    }
    if len(found) != 0 {
        t.Fatalf("expected no diagnostics, got %v", found)
    }
}
```

Cubre: `CheckFilenamePortability` (walk completo, ruta sin diagnósticos),
`checkCaseCollisionsInDir` (ruta sin colisión), `sortDiagnostics`.

### Test 2: unidad pura de reservedWindowsName (todas las plataformas)

```go
func TestReservedWindowsNameDetectsAndIgnores(t *testing.T) {
    if got := reservedWindowsName("CON.md"); got != "CON" {
        t.Fatalf("reservedWindowsName(CON.md) = %q, want CON", got)
    }
    if got := reservedWindowsName("T001-task.md"); got != "" {
        t.Fatalf("reservedWindowsName(T001-task.md) = %q, want empty", got)
    }
    if got := reservedWindowsName("NUL"); got != "NUL" {
        t.Fatalf("reservedWindowsName(NUL) = %q, want NUL", got)
    }
}
```

Cubre: `reservedWindowsName` — ambas ramas (reservado y no reservado).
Sin syscalls de FS, funciona en todas las plataformas.

### Test 3: nombre reservado real (Linux y macOS, skip en Windows)

```go
func TestCheckFilenamePortabilityDetectsReservedName(t *testing.T) {
    if runtime.GOOS == "windows" {
        t.Skip("cannot create files named CON on Windows")
    }
    root := t.TempDir()
    for _, name := range []string{"CON.md", "T001-task.md"} {
        if err := os.WriteFile(filepath.Join(root, name), nil, 0o644); err != nil {
            t.Fatal(err)
        }
    }
    found, err := CheckFilenamePortability(root)
    if err != nil {
        t.Fatal(err)
    }
    assertLintDiagnostic(t, found, diagnostics.DiagnosticLintFilenameReserved, "CON.md", "CON")
}
```

Cubre: `lintNameDiagnostic`, rama "reserved name found" de `CheckFilenamePortability`.

## Criterios de Aceptación

- `go test ./internal/lint/...` pasa en Linux, macOS y Windows
- Cobertura de `internal/lint` ≥ 85% en macOS y Windows (compensando el skip del collision test)
- `./scripts/check-coverage.sh` reporta ≥ 85.0% en Linux
- CI smoke/macos-latest pasa el coverage gate
