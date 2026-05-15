---
estado: Specified
tipo: task
---
# T020: Add rootline binary staleness check to loop-subcommand.md preflight

**Contribuye a**: el loop detecta proactivamente cuando el binario de rootline está desactualizado respecto a su fuente, antes de ejecutar tasks que dependan del formato de su output

## Preserva

- INV1: `roadmapctl check --strict` verde
  - Verificar: `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict`

## Contexto

Cuando rootline source cambia (e.g. refactor de tree.go en O14/T004), el binario instalado en `/home/pones/.local/bin/rootline` puede quedar stale. `roadmapctl next/pending/decision` invocan ese binario en runtime — si el binario es pre-T004, devuelve JSON v1 sin `frontmatter` map. El resultado son títulos vacíos y otros fallos silenciosos.

Durante el loop de O20, este problema causó que T003 fallara la verificación de AC (`ready titles: ['<MISSING>']`) aunque el código era correcto. La detección fue manual y tardía.

El loop skill es el lugar correcto para documentar esta verificación — es quien ejecuta tasks de rootline y quien sufre el síntoma.

## Alcance

**In**:
1. Agregar un bloque de verificación en `## Fase 1: Discovery → Preflight obligatorio` de `/home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md`
2. El bloque debe:
   - Indicar cuándo aplica (loop con tasks que cambian `cmd/rootline/` o `internal/` de rootline)
   - Mostrar cómo detectar discrepancia (comparar `rootline --version` / fecha de `which rootline` con último commit de fuente)
   - Dar el comando de reconstrucción: `go build -o $(which rootline) ./cmd/rootline` desde el repo de rootline
   - Mencionar el síntoma concreto: `roadmapctl next` retorna títulos vacíos o devuelve JSON formato v1

**Out**:
- No cambiar código
- No agregar lógica automática al loop — solo documentar el check manual

## Estado inicial esperado

`loop-subcommand.md § Fase 1` no tiene ninguna mención de binary staleness.

## Criterios de Aceptación

- `grep -n "staleness\|stale\|binary\|binario\|go build" /home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md` retorna al menos una línea relevante
- El bloque menciona el síntoma (títulos vacíos / JSON v1) y el fix (`go build`)
- `roadmapctl check --repo /home/shared/roadmapctl --roadmap-root docs/roadmap --output json --strict` verde

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md` — archivo a modificar
