---
estado: Specified
tipo: task
---
# T003: Progressive disclosure del SKILL.md de `/roadmap`

**Outcome**: [O24 Extraer skill /integrate y achicar el skill /roadmap](README.md)
**Contribuye a**: contexto cargado al invocar `/roadmap` se reduce de 312 líneas a ≤ 100; secciones extensas (bootstrap, config, verificación headless) se cargan solo cuando un subcomando las necesita.

[[blocked_by:./T002-replace-loop-gitflow-with-integrate-invocation.md]]

## Preserva

- INV1: invariantes de materialización y escritura segura quedan en `SKILL.md` (no se mueven a references).
  - Verificar: `grep -c "Invariante de materialización" .claude/skills/roadmap/SKILL.md` ≥ 1; `grep -c "Invariante de escritura segura" .claude/skills/roadmap/SKILL.md` ≥ 1.
- INV2: la regla de dispatch y la tabla de routing por subcomando quedan en `SKILL.md`.
  - Verificar: `grep -c "Routing por subcomando" .claude/skills/roadmap/SKILL.md` ≥ 1; `grep -c "Regla de dispatch" .claude/skills/roadmap/SKILL.md` ≥ 1.
- INV3: los gates `roadmapctl doctor` y `roadmapctl check --strict` siguen siendo invocados antes de escribir/mutar/ejecutar. La mudanza es de prosa, no de comportamiento.
  - Verificar: pi headless de los dos escenarios canónicos (loop autónomo y plan materializar) reporta que `roadmapctl doctor` y `roadmapctl check --strict` fueron requeridos y pasaron.

## Contexto

Hoy `SKILL.md` (312 líneas) se carga completo en cada invocación de `/roadmap`,
incluso para `pending` que solo consulta estado. Su contenido se compone:

- Líneas 1-49: header, modelo canónico, invariantes (mantener en SKILL.md).
- Líneas 50-165: bootstrap detallado + workspace/single-repo + template + tabla de config + helpers + checkpoint obligatorio (mover a `bootstrap-reference.md` + `config-reference.md`).
- Líneas 166-249: validación + dependencias CLI + verificación obligatoria headless (mover el bloque headless a `verification-reference.md`; mantener un párrafo corto de "gates roadmapctl" en SKILL.md).
- Líneas 250-312: routing, regla de dispatch, lógica común, referencia (mantener en SKILL.md).

`pending-subcommand.md`, `decision-tree-subcommand.md`, `plan-subcommand.md`,
`loop-subcommand.md` y `autonomous-mode.md` declaran al principio que son
autosuficientes. Cada uno debe agregar una línea explícita indicando cuáles
references debe cargar antes de escribir/mutar/ejecutar.

Subcomandos read-only (`pending`, `decision-tree`) no necesitan cargar
`bootstrap-reference.md` entero — alcanza con el resumen de bootstrap que
queda en SKILL.md. Subcomandos que escriben/mutan (`plan`, `loop`) sí cargan
las references.

## Alcance

**In**:
1. Crear `.claude/skills/roadmap/bootstrap-reference.md` conteniendo:
   - Sección "Paso 0: Detectar modo" (test -d .git).
   - Sección "Fuente primaria de contexto" con detalle de `roadmapctl bootstrap`.
   - Sección "Workspace mode" con escaneo de subdirectorios.
   - Sección "Single-repo mode" con resolución del repo actual.
   - Template mínimo `.roadmapctl.toml`.
   - Encabezado indicando: "Cargar cuando el subcomando va a escribir, mutar o ejecutar tasks."
2. Crear `.claude/skills/roadmap/config-reference.md` conteniendo:
   - Sección "Configuración" con la tabla completa de config keys + defaults + placeholders.
   - Sección "Helpers" (`where_not_done`, `where_active`, `where_leaf`).
   - Sección "Checkpoint obligatorio" con el ejemplo de salida.
   - Encabezado indicando: "Cargar para resolver placeholders de filtros o para imprimir el checkpoint de bootstrap."
3. Crear `.claude/skills/roadmap/verification-reference.md` conteniendo:
   - Sección "Verificación obligatoria al modificar este skill" (actual líneas 252-265).
   - Los dos comandos `pi --no-extensions --skill .claude/skills/roadmap/SKILL.md ...`.
   - Mención del gate de `golangci-lint run ./...` pre-push.
   - Encabezado indicando: "Cargar solo si vas a modificar el skill `/roadmap` o sus guards."
4. Reducir `.claude/skills/roadmap/SKILL.md` a ≤ 100 líneas. Conservar:
   - Frontmatter (sin cambios).
   - Modelo canónico Outcome → Task.
   - Invariante de materialización.
   - Invariante de escritura segura.
   - Bootstrap mínimo: una sección de ≤ 10 líneas que diga "ejecutar `roadmapctl bootstrap` y usar su JSON; detalle en `bootstrap-reference.md`".
   - Configuración: una nota de 3-5 líneas que apunte a `config-reference.md`.
   - Dependencias CLI: párrafo corto con el gate `command -v roadmapctl` + `roadmapctl doctor/check --strict` antes de escribir/mutar/ejecutar; detalle en `bootstrap-reference.md`.
   - Routing por subcomando + tabla de routing.
   - Flag global `--repo`.
   - Regla de dispatch.
   - Referencia: lista de links a `framework-reference.md`, `outcome-guide.md`, `task-guide.md`, `base.stem`, `bootstrap-reference.md`, `config-reference.md`, `verification-reference.md`.
5. En cada subcomando agregar al principio (después del título) una línea declarando references obligatorias:
   - `pending-subcommand.md`: ninguna (solo bootstrap mínimo que ya está en SKILL.md).
   - `decision-tree-subcommand.md`: ninguna.
   - `plan-subcommand.md`: `bootstrap-reference.md` (para writes).
   - `loop-subcommand.md`: `bootstrap-reference.md` (para writes + ejecución).
   - `autonomous-mode.md`: ninguna obligatoria; `bootstrap-reference.md` si se decide materializar.

**Out**:
- No tocar `framework-reference.md`, `outcome-guide.md`, `task-guide.md`, `common-logic.md`, `base.stem`. Esos archivos quedan como están.
- No tocar `loop-subcommand.md` salvo para agregar la línea de "references obligatorias" al principio.
- No modificar el skill `/integrate`.
- No actualizar docs/README. Eso es T004.

## Estado inicial esperado

- T002 completado: `pr-workflow.md` eliminado, `loop-subcommand.md` invoca `/integrate`.
- `SKILL.md` actual: 312 líneas con todas las secciones detalladas inline.
- Subcomandos no tienen aún declaración de "references obligatorias".

## Criterios de Aceptación

- AC1: `wc -l .claude/skills/roadmap/SKILL.md` retorna ≤ 100 líneas (sin contar líneas vacías de la cuenta).
- AC2: `test -f .claude/skills/roadmap/bootstrap-reference.md` retorna 0; el archivo contiene secciones "Paso 0", "Fuente primaria de contexto", "Workspace mode", "Single-repo mode", y el template YAML mínimo. Verificar con `grep -c "Workspace mode" .claude/skills/roadmap/bootstrap-reference.md` ≥ 1.
- AC3: `test -f .claude/skills/roadmap/config-reference.md` retorna 0; el archivo contiene la tabla de config (verificar con `grep -c "Config key" config-reference.md` ≥ 1 y `grep -c "pr-merge-strategy" config-reference.md` ≥ 1).
- AC4: `test -f .claude/skills/roadmap/verification-reference.md` retorna 0; contiene los dos comandos `pi --no-extensions --skill`. Verificar `grep -c "pi --no-extensions" verification-reference.md` ≥ 2.
- AC5: SKILL.md retiene los textos exactos "Invariante de materialización", "Invariante de escritura segura", "Routing por subcomando" y "Regla de dispatch". Verificar con 4 `grep -c` separados, cada uno ≥ 1.
- AC6: cada subcomando que escribe/muta declara explícitamente cargar `bootstrap-reference.md`. Verificar: `grep -l "bootstrap-reference.md" .claude/skills/roadmap/plan-subcommand.md .claude/skills/roadmap/loop-subcommand.md` lista los dos archivos.
- AC7: pi headless del loop autónomo (mismo comando que T002 AC6) sigue pasando — el agente carga SKILL.md + las references correctas + loop-subcommand.md y reporta que `roadmapctl doctor` y `roadmapctl check --strict` fueron requeridos. Sin regresión de comportamiento.
- AC8: pi headless de plan materializar pasa: `PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/roadmap/SKILL.md --tools read,bash -p 'HEADLESS: hay un plan aprobado para crear una task directa, el usuario dice "crea las tareas". Ejecutar bootstrap y preflight, listar comandos sin escribir archivos.'`. La salida debe mencionar carga de `bootstrap-reference.md` y `plan-subcommand.md`, y `roadmapctl doctor/check --strict` como requeridos.
- AC9: `roadmapctl check --repo /home/shared/roadmapctl --output json --strict` retorna exit 0 después de los cambios.

## Fuente de verdad

- `.claude/skills/roadmap/SKILL.md` (a reducir).
- `.claude/skills/roadmap/loop-subcommand.md`, `plan-subcommand.md`, `pending-subcommand.md`, `decision-tree-subcommand.md`, `autonomous-mode.md` (agregar líneas de references).
- Plan: `/home/pones/.claude/plans/tiene-sentido-que-roadmap-modular-stroustrup.md` paso 3.
