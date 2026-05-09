---
estado: Completed
tipo: task
---
# T001: Centralizar templates bootstrap de roadmapctl

**Outcome**: [Roadmapctl schema and coverage guards](README.md)

## Preserva

- El .stem canónico mantiene estado requerido solo para T*.
- Outcome README materializados no incluyen estado manual.

## Contexto

Los templates base están duplicados en internal/cli/bootstrap.go e internal/materialize/dryrun.go, y agregar required_code_coverage duplicaría el riesgo de drift.

## Alcance

**In**:
1. Crear o mover los templates canónicos a una ubicación compartida interna.
2. Actualizar bootstrap y materialize para consumir el mismo template.
3. Agregar required_code_coverage = 85.0 al TOML default compartido.

**Out**:
1. Auto-reparar .stem existentes.
2. Cambiar el formato de tasks materializadas.

## Estado inicial esperado

baseStemContent y defaultRoadmapctlTOML existen duplicados en bootstrap.go y dryrun.go.

## Criterios de Aceptación

- Bootstrap y materialize usan un único origen para base .stem y default .roadmapctl.toml.
- Tests confirman que Outcome README materializado no contiene estado.
- Tests confirman que el TOML default contiene required_code_coverage = 85.0.

## Fuente de verdad

- internal/cli/bootstrap.go
- internal/materialize/dryrun.go
- internal/materialize/dryrun_test.go
- internal/cli/bootstrap_test.go
