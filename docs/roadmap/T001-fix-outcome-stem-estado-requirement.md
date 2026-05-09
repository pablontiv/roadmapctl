---
estado: Pending
tipo: task
---
# T001: Fix outcome .stem estado requirement

## Preserva

- roadmapctl check passes after fix

## Contexto

renderOutcome was fixed to not emit estado but the production .stem still required it for O*, causing roadmapctl check to fail with RMC_ROOTLINE_VALIDATE_FAILED.

## Alcance

**In**:
1. docs/roadmap/.stem required.match: T* only
2. remove validate non_empty for estado

**Out**:
1. no Go code changes

## Estado inicial esperado

.stem required estado on O* and T*

## Criterios de Aceptación

- docs/roadmap/.stem required.match only targets T*
- validate rule non_empty for estado removed
- rootline validate --all passes with no errors on outcome READMEs
- roadmapctl check passes

## Fuente de verdad

- docs/roadmap/.stem
