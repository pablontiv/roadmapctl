---
estado: Specified
tipo: task
---
# T034: Loop — fallback de re-verificación textual cuando `grep` está interceptado

**Contribuye a**: La re-verificación directa post-agente nunca aborta porque un hook del proyecto intercepta `grep`/`git grep`; siempre hay un camino verificable o se detiene con reporte explícito.

## Contexto

En la sesión O25, un hook del repo (cartyx) interceptó `grep` y `git grep` redirigiéndolos a `cartyx query` (que indexa estructura de código, no texto libre). Esto chocó con la regla T026 de re-verificación directa, que asume que `grep` funciona para validar ACs como "el archivo contiene la cadena X".

Los fallbacks usados in-flight: `find -name '*.go' -exec grep -l <patrón> {} +`, leer el archivo con la herramienta Read y matchear en memoria, abrir el archivo con cat/awk. Funcionaron pero fueron improvisados, no parte del skill.

El comando alternativo del hook puede ser útil para algunas verificaciones (búsqueda estructural), pero NO sustituye búsqueda textual — si la AC pide "esta cadena exacta aparece en este archivo", indexación estructural no aplica.

## Alcance

**In**:
1. En `.claude/skills/roadmap/loop-subcommand.md`, Fase 3 paso 6, dentro de la regla "Re-verificación directa post-subagente paralelo":
   - Agregar sub-regla "Fallback de re-verificación textual": si el shell del loop tiene `grep`/`git grep` interceptado por un hook local del proyecto, no abortar la re-verificación.
   - Listar al menos tres fallbacks portables:
     1. `find <path> -name '*.<ext>' -exec grep -l <patrón> {} +` (o `+` final para batching).
     2. Leer el archivo con la herramienta Read y matchear textualmente en la conversación.
     3. Usar el comando alternativo que el hook sugiera, SOLO si indexa contenido textual; si solo indexa estructura/símbolos, NO cuenta como AC verificado.
   - La re-verificación nunca aborta por hook intercept; o encuentra un camino verificable o reporta y detiene el loop.

**Out**:
- No agregar lógica de detección automática del hook (basta con que el modelo reconozca el error stderr típico: "use cartyx instead", "repo is indexed", etc.).
- No tocar la regla T026 base — solo agregar la sub-regla de fallback dentro de ella.
- No documentar el hook específico del repo en el skill (eso va al CLAUDE.md del proyecto).

## Criterios de Aceptación

- Fase 3 paso 6 incluye la sub-regla "Fallback de re-verificación textual" dentro de la regla de re-verificación post-agente.
- La sub-regla lista ≥3 fallbacks portables (no específicos a un hook del proyecto).
- La sub-regla establece que indexación estructural (no textual) NO sustituye búsqueda textual y no cuenta como AC verificado.
- La re-verificación se documenta como no-abortable por hook intercept.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
