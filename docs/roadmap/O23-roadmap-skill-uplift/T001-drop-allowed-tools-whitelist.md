---
estado: Completed
tipo: task
---
# T001: Eliminar bloque allowed-tools del frontmatter de SKILL.md

**Outcome**: [O23 Roadmap skill: tool access, observability docs, README outputs](README.md)
**Contribuye a**: dar al skill `/roadmap` acceso a la tool roster completa de la sesión (incluyendo `Monitor`, `ScheduleWakeup`, y futuras herramientas) sin necesidad de mantener un whitelist explícito.

## Preserva

- INV1: La sección "Invariante de escritura segura" de `SKILL.md` (líneas 70-81 en el estado pre-cambio) sigue siendo binding por narrativa.
  - Verificar: leer `SKILL.md` después del cambio y confirmar que la sección "Invariante de escritura segura" no fue tocada; las prohibiciones de heredocs múltiples / loops shell siguen presentes textualmente.
- INV2: La verificación headless con Pi del skill (`SKILL.md` sección "Verificación obligatoria al modificar este skill") sigue pasando.
  - Verificar: `./scripts/sync-roadmap-skill.sh --install` y luego correr los dos comandos `pi --no-extensions --skill ... -p '...'` documentados en `SKILL.md`; deben mostrar bootstrap y preflight ejecutados sin errores.

## Contexto

El frontmatter de `/home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` actualmente enumera tools permitidas en un bloque `allowed-tools:` (líneas 13-26). La lista omite `Monitor` y `ScheduleWakeup`, lo que bloquea invocaciones a esas tools incluso cuando el usuario las pide explícitamente — incidente confirmado en Backscroll session `9d8a66cd-01eb-4415-8e33-04f3b9cec020.jsonl` (el usuario pidió "usa el tool monitor" y el modelo nunca pudo invocarlo).

La convención del ecosistema de skills es omitir el campo: confirmado en `playground`, `skill-creator`, `claude-automation-recommender`, `verification-before-completion`, `using-superpowers`, y la mayoría de skills bajo `superpowers/`. Cuando se omite, el skill hereda la tool roster completa de la sesión. Mantener un whitelist explícito significa que cada herramienta nueva del harness requiere edición del skill — costo de mantenimiento sin beneficio de seguridad real, porque la "Invariante de escritura segura" del skill ya restringe el uso de Bash a operaciones puntuales vía narrativa.

Después del cambio, el skill se sincroniza al user scope con `scripts/sync-roadmap-skill.sh --install`. La parity post-sync se verifica con `scripts/sync-roadmap-skill.sh --check`.

## Alcance

**In**:
1. Eliminar el bloque `allowed-tools:` y todas sus entradas (Write, Read, Grep, Glob, Bash, TaskCreate, TaskList, TaskUpdate, TaskGet, Skill, AskUserQuestion, ExitPlanMode, Agent) del frontmatter YAML de `SKILL.md`.
2. Correr `./scripts/sync-roadmap-skill.sh --install` para espejar el cambio a `~/.claude/skills/roadmap/`.
3. Validar parity con `./scripts/sync-roadmap-skill.sh --check`.

**Out**:
- Editar narrativa o cualquier otra sección de `SKILL.md` (otras tasks lo hacen si aplica; aquí sólo frontmatter).
- Tocar otros skills, sus configuraciones, o el script `sync-roadmap-skill.sh`.
- Agregar wildcard (`*`) literal — Claude Code no documenta esa sintaxis; el patrón es la omisión.

## Estado inicial esperado

- `SKILL.md` tiene el bloque `allowed-tools:` con 13 entradas explícitas (Write, Read, Grep, Glob, Bash, TaskCreate, TaskList, TaskUpdate, TaskGet, Skill, AskUserQuestion, ExitPlanMode, Agent).
- `diff /home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md ~/.claude/skills/roadmap/SKILL.md` exit 0 (parity actual).

## Criterios de Aceptación

- `grep -c '^allowed-tools:' /home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` retorna `0`.
- `./scripts/sync-roadmap-skill.sh --check` exit 0.
- `diff /home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md ~/.claude/skills/roadmap/SKILL.md` no produce salida.
- Verificación headless documentada en `SKILL.md` (sección "Verificación obligatoria al modificar este skill") pasa: ambos prompts de `pi --no-extensions ...` muestran bootstrap + preflight ejecutados sin errores.

## Fuente de verdad

- `/home/shared/roadmapctl/.claude/skills/roadmap/SKILL.md` (frontmatter, líneas 13-26 pre-cambio)
- `/home/shared/roadmapctl/scripts/sync-roadmap-skill.sh`
