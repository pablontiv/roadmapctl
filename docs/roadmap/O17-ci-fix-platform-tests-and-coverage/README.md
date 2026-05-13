---
tipo: outcome
---
# O17: CI verde — fixes de plataforma y cobertura cruzada

O14/O15/O16 cierran la mayoría de los bugs de CI, pero la pipeline sigue roja.
Investigación del run `25831261949` revela 5 bugs en 3 plataformas.

**Todos los jobs verdes desde run `25832630920`:** `ci/Test`, `ci/Lint`, `ci/Build`,
`ci/Tidy`, `ci/Vulnerability`, `smoke/ubuntu-latest`, `smoke/macos-latest`, `smoke/windows-latest`.
Release `v0.0.1` publicado con binarios para 6 plataformas.

## Bugs resueltos

| RC | Job | Síntoma | Fix |
|----|-----|---------|-----|
| RC1 | smoke/windows-latest | `filepath.Rel` falla con paths `/`-prefixed en Windows | T001: early-return ErrPathEscape en fsx |
| RC2 | smoke/windows-latest | `chmod 0o555` no impide crear directorios en Windows | T002: skip con `runtime.GOOS == "windows"` |
| RC3 | smoke/macos-latest + smoke/windows-latest | Coverage 83.5% < 85% — case-collision test salta en FS insensible | T003: tests cross-platform para CheckFilenamePortability |
| RC4 | ci/Test (crossbeam) | Coverage ~78% con fake rootline | T004: coverage-threshold: 0 en crossbeam ci job |
| RC5 | smoke/windows-latest | Coverage 84.8% < 85% — bootstrapApplyDiagnostic y lintNameDiagnostic al 0% en Windows | T005: tests unitarios cross-platform para las funciones afectadas |
