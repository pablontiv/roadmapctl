---
estado: Pending
tipo: task
---
# T003: Validar consistencia de tabla `## Tasks`

**Outcome**: [O05 Lint semántico determinístico](README.md)
**Contribuye a**: CE1

[[blocked_by:./T002-implement-markdown-section-table-parser.md]]

## Preserva

- INV1: La tabla no incluye estado; el estado vive en frontmatter.
  - Verificar: lints no exigen columna Estado.

## Contexto

Los Outcomes deben tener README con tabla `## Tasks` que enlace a tasks hijas. La tabla puede quedar stale cuando se crean/renombran tasks.

## Alcance

**In**:
1. Comparar archivos `TXXX-*.md` dentro del Outcome contra links de tabla.
2. Detectar filas faltantes y stale.
3. Validar links relativos de filas.
4. Emitir diagnostics con README path y task/link afectado.

**Out**:
- Auto-fix de tabla.
- Validar contenido semántico de descripciones.

## Estado inicial esperado

- Parser Markdown de tablas existe.

## Criterios de Aceptación

- Fixture con task sin fila falla/warn según política.
- Fixture con fila stale detecta link inexistente.
- Fixture válido pasa.

## Fuente de verdad

- `.claude/skills/roadmap/outcome-guide.md`
- `.claude/skills/roadmap/common-logic.md`
- `docs/roadmap/*/README.md`
