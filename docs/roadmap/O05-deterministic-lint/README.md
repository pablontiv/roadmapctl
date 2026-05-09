---
estado: Pending
tipo: outcome
---
# O05: Lint semántico determinístico

## Objetivo

`roadmapctl lint` valida reglas semánticas determinísticas de roadmap (secciones, tablas, consistencia y portabilidad) sin introducir juicio AI ni auto-fix implícito.

## Criterios de Éxito

- CE1: `roadmapctl lint` detecta inconsistencias de tabla `## Tasks`.
  - Verificar: fixtures stale/missing rows.
- CE2: `lint` detecta secciones estructurales faltantes y ACs ausentes de forma determinística.
  - Verificar: fixtures de tasks incompletas.
- CE3: `lint` detecta problemas cross-platform de nombres.
  - Verificar: fixtures de colisiones case-insensitive/reserved names.

## Invariantes

- INV1: Lint no decide calidad AI subjetiva.
  - Verificar: diagnostics se basan en estructura observable.
- INV2: Lint no reescribe archivos salvo futuro fix explícito.
  - Verificar: comando read-only.

## Alcance

**In**:
- Taxonomía de lint.
- Parser Markdown con goldmark.
- Tablas, secciones, ACs, schema compatibility y nombres cross-platform.

**Out**:
- Auto-fix.
- Materialización.
- Evaluación AI de calidad.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-define-lint-taxonomy-severity-json.md) | Definir taxonomía de lint y contrato JSON |
| [T002](T002-implement-markdown-section-table-parser.md) | Implementar parser Markdown/secciones/tablas |
| [T003](T003-implement-outcome-task-table-consistency.md) | Validar consistencia de tabla `## Tasks` |
| [T004](T004-implement-task-section-and-ac-lint.md) | Validar secciones y ACs de tasks |
| [T005](T005-implement-schema-and-cross-platform-name-lints.md) | Validar compatibilidad schema y nombres cross-platform |
| [T006](T006-add-lint-fixtures-and-goldens.md) | Agregar fixtures/goldens de lint |
