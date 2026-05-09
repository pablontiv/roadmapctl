---
estado: Pending
tipo: task
---
# T006: Complementar validación raw de `blocked_by`

**Outcome**: [O02 Fundaciones post-MVP](README.md)
**Contribuye a**: INV1

[[blocked_by:./T002-make-stem-authoritative-for-document-schema.md]]

## Preserva

- INV1: La definición de links permitidos vive en `.stem`; roadmapctl solo agrega governance donde Rootline no expone suficiente detalle.
  - Verificar: tests con `.stem` y graph Rootline.

## Contexto

El contrato exige que `blocked_by` use path relativo explícito y apunte a task files. Rootline validate/graph cubre gran parte, pero graph resuelto puede ocultar detalles raw si hay fallback por basename o shapes cambian.

## Alcance

**In**:
1. Investigar qué información expone Rootline sobre links raw vs resueltos.
2. Agregar validación complementaria solo si Rootline no basta para diagnosticar targets bare/no explícitos.
3. Mantener el `.stem` como autoridad de patrón de link.
4. Agregar fixtures para bare target, broken target, target no-task y links válidos.

**Out**:
- Reimplementar todo el parser de links de Rootline si el JSON ya provee suficiente información.
- Cambiar semántica de `blocked_by` en `.stem`.

## Estado inicial esperado

- `graphDiagnostics` usa `broken_links` de Rootline.
- Fixture `invalid-bare-blocked-by` existe.

## Criterios de Aceptación

- `RMC_GRAPH_INVALID_BLOCKED_BY` se emite para links inválidos o no explícitos.
- Links válidos same-outcome y cross-outcome pasan.
- La implementación documenta si depende de Rootline validate, graph o scan raw.

## Fuente de verdad

- `internal/roadmap/dependencies.go`
- `internal/roadmap/status.go`
- `docs/roadmap/.stem`
- `testdata/fixtures/invalid-bare-blocked-by`
