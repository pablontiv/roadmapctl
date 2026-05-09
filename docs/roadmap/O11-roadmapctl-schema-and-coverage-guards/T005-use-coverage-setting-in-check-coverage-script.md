---
estado: Specified
tipo: task
---
# T005: Usar required_code_coverage en check-coverage.sh

**Outcome**: [Roadmapctl schema and coverage guards](README.md)

[[blocked_by:./T004-add-required-code-coverage-config.md]]

## Preserva

- El script no depende de Rootline para correr.
- COVERAGE_THRESHOLD conserva precedencia.

## Contexto

El umbral actual vive como fallback hardcodeado 85.0 en scripts/check-coverage.sh; debe quedar sincronizado con la configuración repo-local.

## Alcance

**In**:
1. Leer required_code_coverage desde ${ROADMAP_ROOT:-docs/roadmap}/.roadmapctl.toml si existe.
2. Mantener precedencia: env COVERAGE_THRESHOLD, TOML, fallback 85.0.
3. Actualizar documentación de release/coverage.

**Out**:
1. Agregar dependencia obligatoria de roadmapctl context al script.
2. Cambiar el cálculo de cobertura Go.

## Estado inicial esperado

scripts/check-coverage.sh usa COVERAGE_THRESHOLD:-85.0 sin consultar .roadmapctl.toml.

## Criterios de Aceptación

- El script honra COVERAGE_THRESHOLD cuando está seteado.
- El script usa required_code_coverage de TOML cuando no hay env override.
- El fallback sigue siendo 85.0 si no hay TOML o key.

## Fuente de verdad

- scripts/check-coverage.sh
- docs/release.md
- docs/cli-contract.md
