---
estado: Completed
tipo: task
---
# T008: Align outcome guide with renderer

## Preserva

- Outcome README remains an index with derived status.
- The Tasks table remains maintained consistently and without manual Estado column.

## Contexto

The guide template includes sections not emitted by the current materialize renderer.

## Alcance

**In**:
1. Compare outcome-guide.md with renderOutcome output and goldens.
2. Either simplify the guide or extend schema/renderer intentionally.
3. Update tests/goldens if renderer output changes.

**Out**:
1. Do not add manual estado to outcome frontmatter.
2. Do not introduce unsupported prose-only roadmap formats.

## Estado inicial esperado

Outcome guide and generated README structure differ.

## Criterios de Aceptación

- outcome-guide.md and materialized outcome README output are consistent.
- Goldens/tests are updated if output changes.
- roadmapctl check remains strict-clean.

## Fuente de verdad

- .claude/skills/roadmap/outcome-guide.md
- internal/materialize/dryrun.go
- internal/materialize/dryrun_test.go
