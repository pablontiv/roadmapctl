---
estado: Pending
tipo: task
---
# T004: Validar secciones y ACs de tasks

**Outcome**: [O05 Lint semántico determinístico](README.md)
**Contribuye a**: CE2

[[blocked_by:./T002-implement-markdown-section-table-parser.md]]

## Preserva

- INV1: No evaluar calidad subjetiva; solo presencia/estructura observable.
  - Verificar: diagnostics no dicen “mal especificado” sin regla objetiva.

## Contexto

`task-guide.md` define secciones esperadas para tasks AI-ready. `roadmapctl lint` debe validar estructura mínima sin intentar juzgar el contenido como lo haría un LLM.

## Alcance

**In**:
1. Detectar secciones esperadas: Preserva, Contexto, Alcance, Estado inicial esperado, Criterios de Aceptación, Fuente de verdad.
2. Validar presencia de al menos un AC observable.
3. Validar Fuente de verdad no vacía.
4. Configurar severities iniciales como warnings o errores según taxonomía aprobada.

**Out**:
- Evaluar si una task cabe en una sesión.
- Interpretar lenguaje natural de ACs.

## Estado inicial esperado

- Parser de headings existe.

## Criterios de Aceptación

- Fixture con sección faltante emite diagnostic.
- Fixture con ACs ausentes emite diagnostic.
- Task con todas las secciones mínimas pasa.

## Fuente de verdad

- `.claude/skills/roadmap/task-guide.md`
- `.claude/skills/roadmap/framework-reference.md`
