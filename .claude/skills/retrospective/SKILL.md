---
source: pablontiv/roadmapctl
name: retrospective
description: |
  Análisis de errores de ejecución y comprensión al finalizar un ciclo de trabajo.
  Se invoca automáticamente al terminar /roadmap loop. También usar cuando el usuario
  pida "retrospectiva", "qué salió mal", "qué mejorar", "analizar errores",
  "revisar la sesión", o quiera identificar correcciones a skills/flujos/tasks.

  Produce una tabla de propuestas accionable como input para /roadmap plan + /roadmap loop.
argument-hint: "[opcional: área o skill específico a analizar]"
allowed-tools:
  - Read
  - Grep
  - Glob
  - Bash
---

# /retrospective — Análisis de Errores y Comprensión

## Filosofía

> "Los errores son información. Los malentendidos son señales de que el skill/flujo
> no fue lo suficientemente claro. Documentarlos donde se usan previene recurrencia."

Dos tipos de problema a identificar:
- **Ejecución**: algo falló técnicamente (comando, transición, AC, entorno)
- **Comprensión**: el agente malinterpretó scope, instrucciones o criterios — fue corregido
  por el usuario o produjo work que tuvo que rehacerse

## Detección de modo

Al inicio, determinar contexto disponible:

- **Post-loop**: hay `checkpoint_commit` en la conversación (pasado por el loop)
  → usar directamente
- **Manual**: sin `checkpoint_commit`
  → inferir desde `git log` y conversación activa

## Fase 0 — Recuperación de contexto

### Post-loop
```bash
# Commits producidos durante el loop
git log <checkpoint_commit>..HEAD --oneline

# Estado actual del roadmap
roadmapctl pending --repo <repo-path> --output json 2>/dev/null
```

### Manual
```bash
git log --oneline -20
```

### Recuperación de errores pre-compact (si sesión larga)

**Primario** — usar `backscroll search` para recuperar errores de sesiones anteriores:


```bash
backscroll search "error fail wrong issue incorrecto" --session last
```

**Fallback** — si `backscroll` no está disponible, leer el archivo de sesión directamente:

```bash
PROJECT_SLUG=$(pwd | tr '/' '-' | sed 's/^-//')
SESSION_FILE=$(ls -t ~/.claude/projects/${PROJECT_SLUG}/*.jsonl 2>/dev/null | grep -v agent | head -1)

if [ -n "$SESSION_FILE" ]; then
  python3 -c "
import json, re, sys
with open('$SESSION_FILE') as f:
    for line in f:
        try:
            text = json.dumps(json.loads(line))
            if re.search(r'error|fail|wrong|issue|no[,\s]|incorrecto|repetir', text, re.I):
                m = re.search(r'.{0,60}(error|fail|wrong|issue|incorrecto).{0,60}', text, re.I)
                if m: print(m.group(0)[:140])
        except: pass
  " 2>/dev/null | sort -u | head -30
fi
```

Combinar con la conversación visible.

## Fase 1 — Errores de ejecución

Revisar la conversación e identificar fallos técnicos:

| # | Error | Skill / Comando | Comportamiento esperado | Comportamiento real |
|---|-------|-----------------|------------------------|---------------------|

Clasificar cada error:
- **Transición**: `roadmapctl transition` falló o produjo estado inesperado
- **AC**: Criterio de aceptación no pasó, era ambiguo o era falso positivo
- **Entorno**: Path, permisos, variable de entorno, dependencia faltante
- **Flujo**: Orden de operaciones incorrecto, pre-check omitido
- **Externo**: Red, API, servicio externo

## Fase 2 — Errores de comprensión

Identificar momentos donde el agente se desvió del objetivo real o fue corregido.

Señales a buscar:
- Mensajes del usuario con "no", "eso no", "en realidad", "no era eso"
- Tasks que tuvieron que rehacerse
- Scope que se expandió o contrajo respecto a la task
- Instrucciones del skill que resultaron ambiguas o incorrectas en este contexto
- Patrones aplicados que no correspondían

| # | Momento | Skill / Regla involucrada | Qué entendió el agente | Qué era correcto |
|---|---------|--------------------------|------------------------|------------------|

## Fase 3 — Verificación pre-propuesta

Antes de proponer, buscar en los artefactos candidatos para confirmar que la corrección
no existe ya:

```bash
rg -n "<keyword>" .claude/skills/ .claude/rules/ 2>/dev/null | head -10
```

Si `rg` no está disponible, usar `grep -r` como fallback. Esto previene proponer reglas
que ya existen o duplicar retrospectivas anteriores.

### Clasificación de scope

Antes de poner una propuesta en la tabla, clasificar `Scope` y `Repo destino`:

- **local**: el artefacto pertenece al repo actual. Incluye `CLAUDE.md`, `AGENTS.md`,
  `.claude/rules/`, docs locales y roadmap local. También aplica a skills cuyo
  `SKILL.md` no tiene `source:` o cuyo `source:` coincide con el repo actual.
- **upstream-skill**: el artefacto está bajo `.claude/skills/<skill>/...` y el
  frontmatter de `.claude/skills/<skill>/SKILL.md` tiene `source:` apuntando a otro
  repo. `Repo destino` debe ser ese `source:` (por ejemplo `pablontiv/roadmapctl`).

Para candidatos bajo `.claude/skills/<skill>/...`, leer siempre el frontmatter de
`.claude/skills/<skill>/SKILL.md` antes de proponer. Si el scope es `upstream-skill`
y existe un clone local en `/home/shared/<repo-name>`, materializar la corrección en
ese repo fuente, no en la copia instalada del usuario ni en el repo del proyecto actual.

## Fase 4 — Tabla de propuestas (output final)

Presentar tabla de correcciones propuestas, estructurada como spec para `/roadmap plan`:

| # | Tipo | Scope | Repo destino | Artefacto | Sección | Cambio propuesto | Previene error # |
|---|------|-------|--------------|-----------|---------|-----------------|-----------------|
| 1 | Comprensión | upstream-skill | `pablontiv/roadmapctl` | `.claude/skills/roadmap/loop-subcommand.md` | Fase 3.X | Aclarar que Y debe hacerse antes de Z | E2 |
| 2 | AC | local | repo actual | `docs/roadmap/O01/TXXX-task.md` | Criterios de Aceptación | Especificar condición exacta del AC | E1 |

**El skill termina aquí.** Para materializar correcciones:
```
/roadmap plan  → convierte propuestas en tasks de corrección
/roadmap loop  → ejecuta las correcciones
```

Las propuestas `upstream-skill` deben planearse y ejecutarse en `Repo destino`.
No corregirlas editando copias instaladas en user-scope ni mezclándolas en el repo actual.

## Mapeo de artefactos por tipo de error

| Tipo de Error | Artefacto candidato |
|---------------|---------------------|
| Instrucción de skill incorrecta / ambigua | `.claude/skills/<skill>/SKILL.md` o subcomando; clasificar `local` vs `upstream-skill` desde frontmatter `source:` |
| Regla operativa mal aplicada | `.claude/rules/<rule>.md` |
| AC de task incompleto o ambiguo | `<roadmap-root>/<outcome>/<task>.md` sección ACs |
| Criterio de completitud del outcome | `<roadmap-root>/<outcome>/README.md` |
| Pre-check de entorno faltante | Skill o rule donde se ejecuta el comando |
| Pattern de proyecto no documentado | `CLAUDE.md` del proyecto |
