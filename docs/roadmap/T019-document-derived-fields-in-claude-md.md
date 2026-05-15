---
estado: Completed
tipo: task
---
# T019: Document `derived` vs `frontmatter` in roadmapctl CLAUDE.md

**Contribuye a**: agentes que implementen features sobre next/pending/decision no asumen que todos los campos de una task están en `frontmatter`

## Preserva

- INV1: `go test ./...` verde
  - Verificar: `cd /home/shared/roadmapctl && go test ./...`
- INV2: `roadmapctl check --strict` verde
  - Verificar: `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict`

## Contexto

`roadmapctl next/pending` llaman internamente a `rootline query --output json`. Cada row devuelta tiene dos mapas separados:

- `frontmatter`: campos del YAML del archivo (`estado`, `tipo`, etc.)
- `derived`: campos calculados por el stem (`titulo` via `source: body.h1`, `is_done`, `isIndex`)

Durante el loop de O20, agent-t003 implementó la extracción de `titulo` usando `stringField(frontmatter, ...)`. Como `titulo` está en `derived` y no en `frontmatter`, los títulos eran siempre vacíos. El fix fue agregar `effectiveFields(row)` que hace merge de ambos mapas con `derived` ganando.

Este invariante no está documentado en ningún artefacto duradero.

## Alcance

**In**:
1. Agregar una sección en `/home/shared/roadmapctl/CLAUDE.md` que documente:
   - La separación `frontmatter` / `derived` en rows de `rootline query`
   - Qué tipo de campos va a cada mapa (frontmatter = YAML raw; derived = `source:` rules, expr derivations)
   - El patrón `effectiveFields(row)` como forma correcta de acceder a cualquier campo
   - Ejemplos concretos: `titulo` en derived, `estado` en frontmatter

**Out**:
- No cambiar código
- No tocar tasks completadas

## Estado inicial esperado

`/home/shared/roadmapctl/CLAUDE.md` no tiene ninguna mención de `derived` ni de la separación de campos.

## Criterios de Aceptación

- `grep -n "derived" /home/shared/roadmapctl/CLAUDE.md` retorna al menos una línea relevante
- La sección menciona explícitamente que `titulo` (y otros campos con `source:`) aparecen en `derived`, no en `frontmatter`
- La sección menciona `effectiveFields(row)` como patrón de merge
- `go test ./...` verde (no hay código que cambiar)
- `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict` verde

## Fuente de verdad

- `/home/shared/roadmapctl/CLAUDE.md` — archivo a modificar
- `/home/shared/roadmapctl/internal/roadmap/model.go` — contiene `effectiveFields()` como referencia
