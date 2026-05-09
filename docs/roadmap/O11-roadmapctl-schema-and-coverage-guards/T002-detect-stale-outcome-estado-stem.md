---
estado: Specified
tipo: task
---
# T002: Detectar .stem legacy que exige estado en outcomes

**Outcome**: [Roadmapctl schema and coverage guards](README.md)

[[blocked_by:./T001-centralize-roadmapctl-bootstrap-templates.md]]

## Preserva

- Rootline sigue siendo la fuente del schema efectivo.
- roadmapctl no parsea YAML .stem directamente cuando rootline describe alcanza.

## Contexto

Pinata demostró que un .stem stale puede pasar check mientras los Outcome README tienen estado manual, pero falla al materializar outcomes nuevos sin estado.

## Alcance

**In**:
1. Crear helper compartido de schema compatibility basado en rootline describe --output json.
2. Detectar schema.estado.required_match que incluye O* o required global sin scope.
3. Detectar validate global/unscoped field=estado rule=non_empty.
4. Definir diagnostics accionables para ambas condiciones.

**Out**:
1. Rechazar Outcome README existentes solo por tener estado manual.
2. Modificar Rootline o su semántica de schema.

## Estado inicial esperado

lint solo valida presencia de campos/links y check no inspecciona required_match/validate para compatibilidad de outcomes.

## Criterios de Aceptación

- Unit tests cubren schema correcto sin diagnostics.
- Unit tests cubren required_match [O*, T*] con diagnostic de error.
- Unit tests cubren validate estado non_empty global con diagnostic de error.

## Fuente de verdad

- internal/lint/schema_portability.go
- internal/roadmap/dependencies.go
- internal/diagnostics/report.go
