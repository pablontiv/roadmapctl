---
estado: Specified
tipo: task
---
# T004: Documentar la nueva arquitectura `/roadmap` ↔ `/integrate`

**Outcome**: [O24 Extraer skill /integrate y achicar el skill /roadmap](README.md)
**Contribuye a**: la separación de responsabilidades queda registrada en docs canónicos del repo y en el README, con métricas verificables del refactor.

[[blocked_by:./T003-progressive-disclosure-of-roadmap-skill.md]]

## Preserva

- INV1: la documentación refleja exactamente el estado final post-T003 (no menciona `pr-workflow.md`, no describe el paso 9 antiguo, no asume `SKILL.md` de 312 líneas).
  - Verificar: `grep -r "pr-workflow" docs/ README.md` retorna vacío.
- INV2: las métricas de tamaño documentadas son verificables y se corresponden con `wc -l` real al momento del commit.
  - Verificar: ejecutar `wc -l .claude/skills/roadmap/SKILL.md` y `wc -l .claude/skills/integrate/SKILL.md`; los números mencionados en docs deben coincidir (± 5 líneas para tolerar formateo).

## Contexto

`docs/roadmap-skill-integration.md` documenta hoy la separación
skill/`roadmapctl`/Rootline y contiene una tabla "Thin adapter audit" (líneas
~278-289) listando responsabilidades por área (Bootstrap, Pending/decision,
Loop, Plan). Esa tabla y la sección de Loop integration deben actualizarse
para reflejar la presencia del skill `/integrate` como callee del paso 9.

`README.md` lista skills disponibles del repo (`/roadmap`, `/retrospective`).
Debe agregar `/integrate` con una descripción corta de su rol.

El skill `/integrate` no aparece en `cli-contract.md` porque no introduce
nuevos comandos CLI; tampoco en `release.md` porque no es un cambio de
binario. Solo se documenta en `roadmap-skill-integration.md` y `README.md`.

## Alcance

**In**:
1. En `docs/roadmap-skill-integration.md`:
   - Agregar nueva sección "## `/roadmap loop` ↔ `/integrate` integration" después de "## `/roadmap loop` integration snippet". Describir el contrato de entrada/salida del skill `/integrate` (los inputs canónicos y el bloque `INTEGRATE_RESULT`) y explicar que el loop no contiene prosa de `git`/`gh`.
   - Actualizar la tabla "Thin adapter audit" agregando una fila "Loop integration (git/PR)" con la columna "Remains in skill" = "Skill `/integrate` (callee del loop)" y "Owned by roadmapctl" = "—".
   - Modificar la fila "Loop" existente: cambiar "commit/push according to config" por "delegate commit/push/PR to `/integrate`".
2. En `README.md`:
   - En la sección que lista skills del repo, agregar entrada `/integrate` con descripción de una línea: "Encapsula commit/push/branch/PR per-task; invocado por `/roadmap loop` y disponible ad-hoc."
   - Si existe sección de "Installing skills", mencionar que `scripts/sync-roadmap-skill.sh --install --skill integrate` sincroniza el nuevo skill.
3. Capturar y commitear como parte del mismo commit las métricas finales del refactor en una sección breve "Refactor metrics" dentro de `docs/roadmap-skill-integration.md` (al final del documento, debajo de "Relationship to Rootline"):
   - `SKILL.md` /roadmap: 312 → N líneas (donde N es el conteo real post-T003).
   - Total prosa `/roadmap`: 1458 → M líneas.
   - `/integrate` SKILL.md: K líneas.
   - Fecha del refactor (YYYY-MM-DD).

**Out**:
- No modificar el skill `/integrate` ni los archivos del skill `/roadmap`. Esa parte ya está congelada al entrar a T004.
- No tocar `docs/cli-contract.md`, `docs/release.md`, `docs/auto-update.md`. Sin cambios.
- No agregar al `CHANGELOG.md` (es responsabilidad del flujo de release, no de este refactor de skills).

## Estado inicial esperado

- T003 completado: SKILL.md /roadmap ≤ 100 líneas; bootstrap/config/verification references existen; loop-subcommand.md invoca /integrate.
- `pr-workflow.md` no existe en el repo.
- `docs/roadmap-skill-integration.md` describe `/roadmap` end-to-end pero sin mencionar `/integrate`.
- `README.md` lista `/roadmap` y `/retrospective` pero no `/integrate`.

## Criterios de Aceptación

- AC1: `docs/roadmap-skill-integration.md` contiene una sección con título "/roadmap loop ↔ /integrate integration" (o variante con el carácter ↔ o "->"). Verificar: `grep -c "/integrate" docs/roadmap-skill-integration.md` ≥ 3.
- AC2: la tabla "Thin adapter audit" en `docs/roadmap-skill-integration.md` menciona explícitamente `/integrate` o `Skill /integrate`. Verificar: extraer el bloque de la tabla y `grep -c integrate` ≥ 1.
- AC3: `README.md` contiene una línea mencionando `/integrate` con descripción no vacía. Verificar: `grep -c "/integrate" README.md` ≥ 1.
- AC4: `docs/roadmap-skill-integration.md` contiene una sección "Refactor metrics" (o equivalente) con tres números verificables: SKILL.md /roadmap N líneas, total /roadmap M líneas, /integrate K líneas. Verificar: `grep -c "Refactor metrics" docs/roadmap-skill-integration.md` ≥ 1.
- AC5: los números de líneas mencionados en "Refactor metrics" coinciden con `wc -l` real al momento del commit (tolerar ± 5 por formateo). Verificar manualmente comparando `wc -l .claude/skills/roadmap/SKILL.md`, `find .claude/skills/roadmap -name "*.md" | xargs wc -l | tail -1`, y `wc -l .claude/skills/integrate/SKILL.md` contra los números documentados.
- AC6: `grep -r "pr-workflow" docs/ README.md` retorna vacío (no quedaron referencias colgadas).
- AC7: `roadmapctl check --repo /home/shared/roadmapctl --output json --strict` retorna exit 0 (sanity).

## Fuente de verdad

- `docs/roadmap-skill-integration.md` (sección a agregar + tabla a actualizar + métricas).
- `README.md` (mención del nuevo skill).
- `.claude/skills/integrate/SKILL.md` y `.claude/skills/roadmap/SKILL.md` (para extraer métricas finales con `wc -l`).
