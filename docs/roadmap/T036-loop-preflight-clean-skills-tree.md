---
estado: Specified
tipo: task
---
# T036: Loop — preflight bloquea si `.claude/skills/` tiene cambios uncommitted

**Contribuye a**: El loop no arranca con fuentes del skill modificadas sin integrar; obliga al operador a decidir cómo cerrar el patch antes de ejecutar.

## Contexto

En la sesión O25 el working tree tenía un patch uncommitted en `.claude/skills/roadmap/loop-subcommand.md` (la formalización de `Skill("retrospective", ...)`) que provenía de una sesión previa de debugging (`33b6199f-...`). Esa sesión editó el skill pero no commiteó ni pusheó — porque no era una sesión de loop, era debug y el modelo se fue sin integrar.

Cuando MI sesión arrancó el loop, ese patch coincidía en el mismo archivo (`loop-subcommand.md`) que iban a editar tres tasks pendientes (T026/T027/T028). Riesgo: mezclar scopes en commits. Lo resolví in-flight preguntando al usuario y commiteando standalone, pero esa disciplina debería ser obligatoria del skill.

## Alcance

**In**:
1. En `.claude/skills/roadmap/loop-subcommand.md`, Fase 1 sección "Preflight obligatorio roadmapctl", antes de `doctor`/`check`:
   - Ejecutar `git status --porcelain .claude/skills/` (limitado a esa ruta).
   - Si la salida NO está vacía (cualquier `M`, `A`, `D`, `??`), detener el loop antes de continuar.
   - Mensaje al operador: listar los archivos dirty + ofrecer tres acciones explícitas: (a) commit standalone (sugerencia de mensaje), (b) `git stash` con instrucción para reaplicar después, (c) descartar con `git restore .claude/skills/<archivo>`.
   - El loop no continúa hasta que el operador resuelva (el chequeo se repite en el siguiente intento).
2. El check aplica solo a `.claude/skills/` — no bloquea por otras rutas (e.g., `dist/`, `docs/roadmap/` con tasks in-progress, código fuente del repo).

**Out**:
- No bloquear por cambios en otras rutas (eso es del operador, no del skill).
- No agregar reparación automática (commit/stash/discard).
- No bloquear por archivos ignorados (`??` que matchea `.gitignore`).

## Criterios de Aceptación

- Fase 1 "Preflight obligatorio" incluye el check `git status --porcelain .claude/skills/` antes de `doctor`/`check`.
- Si hay cambios uncommitted/untracked en esa ruta, el loop se detiene con mensaje listando archivos.
- El mensaje propone explícitamente las tres acciones (a/b/c) al operador.
- El check NO aplica fuera de `.claude/skills/`.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
