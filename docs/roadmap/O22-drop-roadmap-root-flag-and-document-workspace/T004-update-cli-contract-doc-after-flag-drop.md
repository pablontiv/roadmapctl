---
estado: Completed
tipo: task
---
# T004: Update `docs/cli-contract.md` after flag drop

**Outcome**: [O22 Drop --roadmap-root flag and document workspace](README.md)
**Contribuye a**: contrato CLI público refleja la realidad post-T001 (sin flag, convención fija)

[[blocked_by:./T001-drop-roadmap-root-flag-from-go.md]]

## Preserva

- INV1: `docs/cli-contract.md` sigue siendo la fuente de verdad estable del contrato CLI; no se elimina ni se reorganiza más allá de lo necesario
- INV2: la integridad de tablas (flags globales, exit codes) se mantiene con la fila/sección eliminada

## Contexto

`docs/cli-contract.md` documenta el contrato público del CLI. Hoy menciona el flag `--roadmap-root` en al menos dos lugares identificados durante la investigación:

- **Línea ~55**: tabla de flags globales — fila `| --roadmap-root | path | inferred from <roadmap-root>/.roadmapctl.toml or one-time legacy migration input | Override configured roadmap root. The resolved path must stay inside the repo. |`
- **Línea ~66**: párrafo "Preferred post-MVP config lives at `<roadmap-root>/.roadmapctl.toml`... `--roadmap-root` remains a command-line override for **explicit inspection and migration workflows**."

Adicionalmente, bloques de código bash en el documento contienen ejemplos como `roadmapctl check --repo . --roadmap-root docs/roadmap --output json` que deben eliminarse o ajustarse para no pasar el flag.

Tras T001 el flag no existe. Hay que:

1. Eliminar la fila de la tabla de flags
2. Eliminar/reformular el párrafo de la línea 66
3. Reemplazar todos los ejemplos de comandos en bash code blocks que pasen `--roadmap-root`
4. Añadir una nota breve documentando que el roadmap root es convención fija `docs/roadmap/`, no configurable

## Alcance

**In**:
1. Buscar `grep -n "roadmap-root" docs/cli-contract.md` para enumerar matches
2. Editar la tabla de flags globales eliminando la fila de `--roadmap-root`
3. Editar/eliminar el párrafo que menciona "migration workflows" en relación al flag
4. Editar los bloques de código bash para eliminar el flag de invocaciones de ejemplo
5. Añadir una sección breve (1-2 líneas) o nota documentando: el roadmap root es siempre `docs/roadmap/`, fijo, no configurable vía flag; resuelto por `roadmapctl bootstrap` en el campo `roadmap_root` del JSON

**Out**:
- Cambios a otros docs (`README.md`, `SKILL.md`) — T003
- Cambios al binario o sus comandos — T001

## Estado inicial esperado

- T001 Completed
- `grep -n "roadmap-root" docs/cli-contract.md` retorna varias líneas (tabla + párrafo + ejemplos)

## Criterios de Aceptación

- `grep -cn "\-\-roadmap-root" docs/cli-contract.md` retorna `0` (cero menciones del flag CLI)
- `grep -n "roadmap-root" docs/cli-contract.md` solo retorna referencias al **concepto** (e.g. "el roadmap root es `docs/roadmap/`") o al campo JSON `roadmap_root` — no al flag CLI
- El documento contiene una afirmación clara de que `docs/roadmap/` es la convención fija
- Tablas (flags globales, exit codes) tienen sintaxis válida tras la eliminación
- Documento sigue siendo coherente (no quedan referencias colgantes a "Override configured roadmap root" sin contexto)

## Fuente de verdad

- `docs/cli-contract.md`
