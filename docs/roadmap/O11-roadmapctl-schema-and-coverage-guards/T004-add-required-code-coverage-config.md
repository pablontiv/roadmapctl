---
estado: Completed
tipo: task
---
# T004: Agregar setting required_code_coverage

**Outcome**: [Roadmapctl schema and coverage guards](README.md)

[[blocked_by:./T001-centralize-roadmapctl-bootstrap-templates.md]]

## Preserva

- COVERAGE_THRESHOLD de entorno sigue pudiendo sobreescribir el script de coverage.
- La configuración operacional sigue viviendo en .roadmapctl.toml.

## Contexto

El workflow de release ya usa un umbral de coverage, pero está hardcodeado en scripts/check-coverage.sh y no está declarado en roadmapctl context/config.

## Alcance

**In**:
1. Agregar Config.RequiredCodeCoverage float64 con default 85.0.
2. Parsear/renderizar required_code_coverage en .roadmapctl.toml.
3. Validar rango 0..100.
4. Exponer required_code_coverage en roadmapctl context single-repo y workspace.
5. Actualizar docs de contrato/config.

**Out**:
1. Cambiar comandos Go de coverage o motor de tests.
2. Remover soporte de override por variable de entorno.

## Estado inicial esperado

Config no tiene campo required_code_coverage y context no lo expone.

## Criterios de Aceptación

- TOML custom carga required_code_coverage.
- Default config usa 85.0.
- Valores menores a 0 o mayores a 100 fallan validación.
- roadmapctl context JSON incluye required_code_coverage.

## Fuente de verdad

- internal/config/config.go
- internal/config/config_test.go
- internal/cli/context.go
- internal/cli/context_test.go
- docs/cli-contract.md
