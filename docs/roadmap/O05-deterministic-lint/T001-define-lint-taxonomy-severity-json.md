---
estado: Completed
tipo: task
---
# T001: Definir taxonomía de lint y contrato JSON

**Outcome**: [O05 Lint semántico determinístico](README.md)
**Contribuye a**: CE1, CE2, CE3

[[blocked_by:../O02-post-mvp-foundations/T002-make-stem-authoritative-for-document-schema.md]]

## Preserva

- INV1: JSON diagnostics siguen formato `RMC_*` estable.
  - Verificar: golden tests.

## Contexto

`lint` agrega checks más semánticos que `check`, pero deben ser determinísticos y observables. Hay que definir qué es error vs warning y cómo interactúa con `--strict`.

## Alcance

**In**:
1. Definir IDs `RMC_LINT_*` y severities.
2. Definir `kind: roadmapctl/lint` y summary.
3. Documentar diferencia entre `check` y `lint`.
4. Definir política de `--strict`.

**Out**:
- Implementar parser Markdown.
- Auto-fix.

## Estado inicial esperado

- `check` cubre estructura/Rootline, no secciones ni tablas.

## Criterios de Aceptación

- Contrato lint queda documentado.
- Hay tests de exit code warning vs strict.
- No se reclasifican checks MVP sin compatibilidad.

## Fuente de verdad

- `docs/cli-contract.md`
- `internal/diagnostics/report.go`
- `docs/golden-tests.md`
