---
estado: Pending
tipo: task
---
# T005: Agregar modelo read-domain y wrapper `rootline tree`

**Outcome**: [O02 Fundaciones post-MVP](README.md)
**Contribuye a**: CE3

[[blocked_by:./T002-make-stem-authoritative-for-document-schema.md]]
[[blocked_by:./T004-harden-rootline-json-on-nonzero-exit.md]]

## Preserva

- INV1: Rootline sigue siendo fuente de datos genérica; roadmapctl normaliza para su dominio.
  - Verificar: wrappers Rootline siguen usando CLI JSON.

## Contexto

Los comandos `pending`, `next` y `decision` necesitan un modelo común de tasks, outcomes, status roles y dependencias. Rootline ya provee `tree`, `query` y `graph`; falta wrapper y normalización interna.

## Alcance

**In**:
1. Agregar wrapper `Tree(ctx, root, wheres...)` en `internal/rootlinecli`.
2. Definir tipos internos para Task, Outcome, Dependency, StatusRole y RoadmapContext.
3. Normalizar paths y campos desde Rootline JSON.
4. Agregar tests de parsing para shapes conocidos.

**Out**:
- Implementar comandos `pending`/`next`/`decision`.
- Mutar estados o archivos.

## Estado inicial esperado

- `rootlinecli.Client` tiene wrappers para validate/describe/query/graph, no tree.

## Criterios de Aceptación

- Tests demuestran que `rootline tree ... --output json` se invoca con args explícitos.
- Modelo interno soporta tasks directas y tasks dentro de outcomes.
- No se introducen hardcodes de valores documentales fuera de fallback explícito.

## Fuente de verdad

- `internal/rootlinecli/client.go`
- `internal/rootlinecli/client_test.go`
- `internal/roadmap/*`
