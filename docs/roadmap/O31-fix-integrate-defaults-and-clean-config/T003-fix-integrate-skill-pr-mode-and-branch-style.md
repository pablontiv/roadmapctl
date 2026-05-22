---
estado: Completed
tipo: task
---
# T003: Fix `/integrate` Fase 2 (pr_mode condicional) y remover fallback `branch_style`

**Outcome**: [O31 Fix `/integrate` defaults y limpieza de config muerta](README.md)
**Contribuye a**: `pr_mode=false` default funciona correctamente; falta de `branch_style` falla temprano

## Preserva

- INV1: las 6 fases de `integrate/SKILL.md` siguen existiendo en el mismo orden
  - Verificar: el archivo tiene secciones Fase 1 a Fase 6
- INV2: ningún cambio toca código Go ni otros skills
  - Verificar: `git diff --name-only` después de la task lista solo `.claude/skills/integrate/SKILL.md` y (si aplica) la copia sincronizada en cache

## Contexto

Bug original en `integrate/SKILL.md` línea 86:
> "Fase 2 corre siempre. Branch local por scope es invariante. `pr_mode` solo controla si Fase 5 y Fase 6 se ejecutan."

Con `pr_mode=false` (default del TOML) y `auto_push=true`, el skill crea `feat/<scope>` igual; sin Fase 5/6, los commits quedan en el feature branch huérfanos. master nunca los recibe. Ese fue el origen del incidente que dejó T005/T006/T007 (O03) fuera de master.

Fix: el skill es el dueño único de la posición de HEAD. Antes de commitear, fuerza el checkout del branch target:
- `pr_mode == false` → target = `base_branch`. Skill hace `fetch + checkout base_branch + pull --ff-only`, después commitea ahí.
- `pr_mode == true` → target = derivado de `branch_style`. Skill hace `fetch + checkout base_branch + pull --ff-only + checkout -B <branch_target>`, después commitea ahí.

Adicionalmente, el skill actualmente tiene un fallback silencioso: "Si `branch_style` está vacío, usar fallback: `feat/<scope>`". Decisión del usuario: sin fallback. Si está vacío, emitir `RMC_INTEGRATE_BRANCH_STYLE_MISSING` y detenerse — el operador debe correr `/roadmap bootstrap` para popular el style field.

## Alcance

**In**:

1. Reescribir la sección "## Fase 2: Branch setup" de `.claude/skills/integrate/SKILL.md` siguiendo el contrato Opción B descrito arriba:
   - Paso 1: detectar branch actual (`git rev-parse --abbrev-ref HEAD`).
   - Paso 2: sincronizar base_branch siempre (`fetch + checkout + pull --ff-only`). Si el checkout falla por conflicto local, emitir `RMC_INTEGRATE_CHECKOUT_BLOCKED` y detener.
   - Paso 3 (sólo si `pr_mode == true`): si `branch_style` vacío → `RMC_INTEGRATE_BRANCH_STYLE_MISSING` y detener; si poblado, LLM genera nombre de branch y `git checkout -B <branch_target>`.
2. Eliminar la frase "Si `branch_style` está vacío, usar fallback: `feat/<scope>`" (y sus ejemplos como `feat/O24-slug`, `feat/direct-roadmap-tasks`).
3. Ajustar Fase 4 (Push): el target del push es el branch actual (HEAD) — `base_branch` con `pr_mode=false`, `<branch_target>` con `pr_mode=true`.
4. Agregar a la tabla "Errores comunes":
   - `RMC_INTEGRATE_BRANCH_STYLE_MISSING` — `branch_style` no configurado en `[gitflow]`; `pr_mode=true` lo requiere. Recovery: correr `/roadmap bootstrap` y popular el style field.
   - `RMC_INTEGRATE_CHECKOUT_BLOCKED` — conflicto local impide checkout de `base_branch`; operador debe resolver state antes de reinvocar.
5. Agregar Escenario C en la sección "Verificación al modificar este skill":
   - Comando headless pi con `pr_mode=false`, `autonomy=until_done`, scope=direct-tasks.
   - Debe listar: `git fetch`, `git checkout <base_branch>`, `git pull --ff-only`, `git add`, `git commit`, `git push` a `<base_branch>`.
   - NO debe listar: `git checkout -B`, `gh pr create`, `gh pr merge`.
   - `INTEGRATE_RESULT.branch` = `<base_branch>` (e.g. `"master"`).
6. Ejecutar `./scripts/sync-roadmap-skill.sh --install --skill integrate` para sincronizar el cache de skills.

**Out**:
- Código Go (T001).
- Otros skills (T002 cubre `.claude/skills/roadmap/*`).
- Cualquier mención de `pr_merge_strategy` (parte de T001/T002).

## Estado inicial esperado

- `.claude/skills/integrate/SKILL.md` contiene literalmente "Fase 2 corre siempre. Branch local por scope es invariante."
- `.claude/skills/integrate/SKILL.md` contiene literalmente "usar fallback".
- T001/T002 pueden estar pendientes o completadas (independencia).

## Criterios de Aceptación

- `grep -F 'Fase 2 corre siempre. Branch local por scope es invariante' .claude/skills/integrate/SKILL.md` retorna vacío.
- `grep -F 'usar fallback' .claude/skills/integrate/SKILL.md` retorna vacío.
- `grep -F 'RMC_INTEGRATE_BRANCH_STYLE_MISSING' .claude/skills/integrate/SKILL.md` retorna al menos 2 matches (Fase 2 + tabla errores).
- `grep -F 'RMC_INTEGRATE_CHECKOUT_BLOCKED' .claude/skills/integrate/SKILL.md` retorna al menos 2 matches.
- `grep -F 'Escenario C' .claude/skills/integrate/SKILL.md` retorna al menos 1 match.
- Verificación headless pi del Escenario A (pr_mode=false): el output NO contiene `checkout -B`, `gh pr create`, ni `gh pr merge`; SÍ contiene `git fetch`, `git checkout`, `git pull`, `git add`, `git commit`, `git push`.
- Verificación headless pi del Escenario B (pr_mode=true): el output SÍ contiene la cadena completa con `checkout -B <branch_target>` y `gh pr create`.
- Verificación headless pi del Escenario C (pr_mode=true + branch_style vacío): el output contiene `RMC_INTEGRATE_BRANCH_STYLE_MISSING` y NO contiene `checkout -B`.
- `./scripts/sync-roadmap-skill.sh --install --skill integrate` exit 0.

## Fuente de verdad

- `.claude/skills/integrate/SKILL.md`
