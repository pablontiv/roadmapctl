---
tipo: outcome
---
# Roadmapctl schema and coverage guards

Roadmapctl debe validar la compatibilidad del .stem que gobierna el roadmap y exponer un setting repo-local de cobertura requerida. Los Outcome README no deben requerir estado manual; el estado de outcomes se deriva de sus tasks.

## Criterios de Aceptación

- roadmapctl detecta schemas legacy que requieren estado en outcomes antes de materializar o declarar check/lint/doctor ok.
- Los templates/bootstrap/materialize siguen generando Outcome README sin estado manual.
- required_code_coverage está en configuración efectiva, contexto JSON, templates, documentación y script de coverage.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-centralize-roadmapctl-bootstrap-templates.md) | Eliminar drift entre los templates base de .stem y .roadmapctl.toml usados por bootstrap y materialize. |
| [T002](T002-detect-stale-outcome-estado-stem.md) | Agregar validación de compatibilidad del schema efectivo para detectar cuando estado es requerido en Outcome README o existe validate estado non_empty global. |
| [T003](T003-wire-schema-compatibility-diagnostics.md) | Hacer que check, lint, doctor, bootstrap y materialize reporten schemas incompatibles antes de declarar éxito o escribir archivos. |
| [T004](T004-add-required-code-coverage-config.md) | Agregar required_code_coverage a la configuración efectiva de roadmapctl, con default 85.0 y validación de rango. |
| [T005](T005-use-coverage-setting-in-check-coverage-script.md) | Actualizar el script de coverage para leer el umbral desde .roadmapctl.toml cuando COVERAGE_THRESHOLD no está seteado. |
| [T006](T006-update-fixtures-docs-and-goldens.md) | Alinear fixtures normales con el .stem canónico, mantener fixtures stale explícitos y actualizar documentación/goldens para los nuevos diagnostics y setting de coverage. |
