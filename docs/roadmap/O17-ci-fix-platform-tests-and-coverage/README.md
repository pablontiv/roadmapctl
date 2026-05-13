---
tipo: outcome
---
# O17: CI verde — fixes de plataforma y cobertura cruzada

O14/O15/O16 cierran la mayoría de los bugs de CI, pero la pipeline sigue roja.
Investigación del run `25831261949` revela 4 bugs restantes en 3 plataformas.

**Jobs que deben quedar verdes:** `ci/Test`, `smoke/macos-latest`, `smoke/windows-latest`

## Bugs activos

| RC | Job | Síntoma |
|----|-----|---------|
| RC1 | smoke/windows-latest | `TestResolveInsideAbsolutePathUnix` y `TestResolveInsideLeadingSlash` fallan — `filepath.Rel` no puede relacionar paths `/absolute` con rutas `C:\...` |
| RC2 | smoke/windows-latest | `TestBootstrapInitApplyReportsDiagnosticsOnFileError` falla — `chmod 0o555` no impide crear directorios en Windows |
| RC3 | smoke/macos-latest + smoke/windows-latest | Cobertura total 83.5% < 85% — `internal/lint` cae a 75.7% porque el test de case-collision se saltea en FS insensible y ningún otro test cubre `CheckFilenamePortability` / `checkCaseCollisionsInDir` / `reservedWindowsName` |
| RC4 | ci/Test (crossbeam) | Cobertura ~78% con fake rootline — crossbeam no instala rootline, `requiresRealRootline` tests se saltean, `internal/cli` cae a 56.9% |
