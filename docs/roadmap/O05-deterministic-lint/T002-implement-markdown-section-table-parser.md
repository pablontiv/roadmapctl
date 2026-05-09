---
estado: Pending
tipo: task
---
# T002: Implementar parser Markdown de secciones y tablas

**Outcome**: [O05 Lint semántico determinístico](README.md)
**Contribuye a**: CE1, CE2

[[blocked_by:./T001-define-lint-taxonomy-severity-json.md]]
[[blocked_by:../O02-post-mvp-foundations/T001-adopt-cobra-and-community-packages.md]]

## Preserva

- INV1: Parser no reescribe documentos.
  - Verificar: tests read-only.

## Contexto

Se necesita inspeccionar headings, secciones y tablas de Markdown. Usar `github.com/yuin/goldmark` con extensiones GFM/table evita mantener parsers propios.

## Alcance

**In**:
1. Agregar goldmark y extensión de tablas.
2. Extraer headings, rangos y tablas relevantes.
3. Normalizar links de tabla sin modificar Markdown.
4. Tests con README/task samples.

**Out**:
- Round-trip completo de Markdown.
- Reformatting.

## Estado inicial esperado

- No existe parser Markdown en roadmapctl.

## Criterios de Aceptación

- Tests detectan `## Tasks`, filas y headings esperados.
- Parser preserva posiciones suficientes para diagnostics.
- No se hace string-grep frágil para tablas complejas si goldmark lo cubre.

## Fuente de verdad

- `go.mod`
- `internal/lint` nuevo
- `.claude/skills/roadmap/outcome-guide.md`
- `.claude/skills/roadmap/task-guide.md`
