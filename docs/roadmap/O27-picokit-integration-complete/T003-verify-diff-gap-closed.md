---
estado: Specified
tipo: task
---
# T003: Verify diff package gap is closed (no orphan callers)

**Outcome**: [O27 Complete picokit integration](README.md)
**Contribuye a**: Confirma que la extracción de `internal/diff` en O25-T003 no dejó callsites huérfanos y que roadmapctl efectivamente no necesita `picokit/diff`.

## Preserva

- INV1: roadmapctl compila sin necesidad de `picokit/diff`.
  - Verificar: `go build ./...` pasa sin agregar el import.

## Contexto

O25-T003 borró `internal/diff/` por dead code. La hipótesis era que no había callers; esta task lo verifica formalmente para descartar deuda oculta.

Si la verificación encuentra callsites huérfanos, esta task se promociona a fix (importar `picokit/diff` y restaurar la funcionalidad). Si no, queda documentado el cierre del gap.

## Alcance

**In**:
1. Grep exhaustivo: `grep -rn "internal/diff\|picokit/diff" /home/shared/roadmapctl --include="*.go"`.
2. Si vacío: documentar en `CLAUDE.md` o comentario en `go.mod` que diff fue extraído y no consumido.
3. Si no vacío: convertir esta task en migración a `picokit/diff` (scope expand documentado en este archivo).

**Out**:
- No reintroducir `internal/diff/`.

## Estado inicial esperado

- `internal/diff/` no existe.
- roadmapctl no importa `picokit/diff`.

## Criterios de Aceptación

- `grep -rn "internal/diff\|picokit/diff" /home/shared/roadmapctl --include="*.go"` retorna vacío.
- `go build ./...` pasa.
- Si hay deuda detectada: scope expand registrado y nueva task creada en el outcome.
