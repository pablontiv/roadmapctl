---
estado: Pending
tipo: task
---
# T005: Validar schema compatibility y nombres cross-platform

**Outcome**: [O05 Lint semántico determinístico](README.md)
**Contribuye a**: CE3

[[blocked_by:./T001-define-lint-taxonomy-severity-json.md]]
[[blocked_by:../O02-post-mvp-foundations/T002-make-stem-authoritative-for-document-schema.md]]

## Preserva

- INV1: `.stem` sigue siendo autoridad de schema; lint solo verifica compatibilidad requerida por roadmapctl.
  - Verificar: describe/validate Rootline.

## Contexto

Además de estructura básica, `roadmapctl` debe detectar problemas portables y schema incompatibles temprano: colisiones case-insensitive, nombres reservados y ausencia de fields/links necesarios para comandos.

## Alcance

**In**:
1. Detectar colisiones case-insensitive en filenames dentro de un scope.
2. Detectar nombres problemáticos para Windows si aplica.
3. Verificar que schema efectivo expone fields/links mínimos que comandos necesitan: `estado`, `tipo`, `blocked_by`.
4. Permitir extensiones de schema del proyecto.

**Out**:
- Rechazar schema extendido solo por tener campos extra.
- Reimplementar validaciones Rootline.

## Estado inicial esperado

- Structure checks validan regex de filenames pero no todas las portabilidad/case collisions.

## Criterios de Aceptación

- Fixture case collision emite diagnostic.
- Schema con extensiones pasa.
- Schema sin `blocked_by` o sin `estado` produce diagnostic útil para comandos dependientes.

## Fuente de verdad

- `internal/roadmap/structure.go`
- `docs/roadmap/.stem`
- `internal/roadmap/schema` nuevo
