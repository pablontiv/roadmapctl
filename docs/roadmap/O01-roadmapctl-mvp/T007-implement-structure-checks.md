---
estado: Completed
tipo: task
---
# T007: Implementar checks estructurales canónicos

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE2

[[blocked_by:./T003-implement-diagnostics-model.md]]
[[blocked_by:./T004-load-roadmap-config.md]]

## Preserva

- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: checks solo producen diagnostics.

## Contexto

El bug original fue que un agente materializó múltiples tasks como un único archivo markdown. Esta task implementa la protección mínima que debe bloquear ese resultado en cualquier plataforma.

## Alcance

**In**:
1. Validar que direct tasks en raíz sean `TXXX-*.md`.
2. Validar que outcomes sean directorios `OXX-*` con `README.md`.
3. Validar que tasks dentro de outcome sean `TXXX-*.md`.
4. Rechazar `*-tasks.md` como fallback de múltiples tasks.
5. Rechazar nesting extra bajo outcomes.
6. Detectar duplicados `OXX` en raíz y `TXXX` por scope.
7. Emitir diagnostics con path relativo.

**Out**:
- Validación de frontmatter.
- Validación de dependencias con Rootline.
- Fix automático.

## Estado inicial esperado

- Config loader resuelve roadmap-root.
- Diagnostics model existe.

## Criterios de Aceptación

- Fixture `invalid-single-summary-file` falla con diagnostic `RMC_STRUCTURE_SINGLE_FILE_FALLBACK`.
- Fixture `invalid-missing-outcome-readme` falla con diagnostic específico.
- Fixture `valid-direct-tasks` pasa checks estructurales.
- Fixture `valid-outcome-with-tasks` pasa checks estructurales.
- Tests normalizan path separators para Windows.

## Fuente de verdad

- `internal/roadmap/structure.go`
- `internal/roadmap/numbering.go`
- `testdata/fixtures/invalid-single-summary-file/`
- `docs/roadmap/O01-roadmapctl-mvp/README.md`
