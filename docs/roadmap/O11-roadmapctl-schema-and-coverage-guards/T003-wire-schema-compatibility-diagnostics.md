---
estado: Completed
tipo: task
---
# T003: Conectar diagnostics de schema compatibility en comandos

**Outcome**: [Roadmapctl schema and coverage guards](README.md)

[[blocked_by:./T002-detect-stale-outcome-estado-stem.md]]

## Preserva

- Los comandos siguen siendo determinísticos y reportan JSON estable.
- Bootstrap/materialize no sobrescriben .stem existentes sin aprobación explícita.

## Contexto

La validación debe ocurrir en las rutas que hoy declaran validez o escriben roadmap, especialmente antes de materialize apply.

## Alcance

**In**:
1. Invocar schema compatibility desde check después de describe.
2. Extender lint para incluir las nuevas reglas.
3. Extender doctor para reportar schema incompatible cuando .stem existe.
4. Extender bootstrap inspect/init para reportar .stem stale y bloquear apply si corresponde.
5. Agregar preflight en materialize apply para .stem existente incompatible.

**Out**:
1. Migración automática de schemas legacy.
2. Cambios grandes en contrato CLI fuera de diagnostics.

## Estado inicial esperado

check/lint/doctor/bootstrap/materialize no fallan temprano por .stem legacy que fuerza estado en outcomes.

## Criterios de Aceptación

- Fixture stale falla check y lint con diagnostics accionables.
- doctor reporta el schema incompatible.
- materialize apply contra .stem stale bloquea antes de escribir planned files.
- bootstrap init --apply con .stem stale no escribe archivos adyacentes y reporta diagnostics.

## Fuente de verdad

- internal/cli/check.go
- internal/cli/lint.go
- internal/cli/doctor.go
- internal/cli/bootstrap.go
- internal/cli/materialize.go
- internal/roadmap/dependencies.go
