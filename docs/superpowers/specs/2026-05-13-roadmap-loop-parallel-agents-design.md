# Design: `/roadmap loop` parallel Agent dispatch

**Date:** 2026-05-13
**Status:** Approved

## Problem

`/roadmap loop` ejecuta tasks una a la vez en serie aunque `parallel = true` esté configurado en `.roadmapctl.toml`. La config es correcta; el skill es ambiguo.

**Root cause:** La sección "Parallel waves" de `loop-subcommand.md` dice:

> "La ejecución paralela debe usar worktrees aislados o una ruta de integración equivalente con control de conflictos. Si no hay aislamiento/control seguro, ejecutar secuencialmente aunque `parallel == true`."

Sin worktrees configurados, el AI interpreta "no hay aislamiento" → cae a secuencial. El mecanismo real de ejecución paralela (llamadas paralelas al tool `Agent`) nunca se menciona explícitamente.

**Evidencia histórica:** En pinata (2026-05-12) el loop SÍ ejecutó con múltiples agentes en paralelo ("dispatching all four agents in parallel", waves 1-8). El AI interpretó Agent calls sobre archivos distintos como "ruta de integración equivalente". Ese comportamiento es el correcto y debe ser el default explícito.

## Solución

Reescritura quirúrgica de la sección "Parallel waves" en `loop-subcommand.md`.

### Cambio

**Sección actual:**
```
Si `parallel == true`, formar waves oportunistas desde `ready[]` usando solo la información canónica de `roadmapctl next` y `blocked_by`:

- Tasks en una misma wave no tienen dependencia explícita entre sí según `roadmapctl next`.
- No inferir dependencias por heurísticas de paths, nombres o secciones; si aparece un conflicto real durante integración, tratarlo como dependencia faltante.
- La ejecución paralela debe usar worktrees aislados o una ruta de integración equivalente con control de conflictos. Si no hay aislamiento/control seguro, ejecutar secuencialmente aunque `parallel == true`.
```

**Sección reemplazada:**
```
Si `parallel == true`, formar waves oportunistas desde `ready[]` y ejecutar
cada wave con llamadas paralelas al tool `Agent` — una por task de la wave.

Las tasks de una wave son independientes por definición (`roadmapctl next`
garantiza que no tienen `blocked_by` entre sí), por lo que Agent calls
paralelas son la ruta de integración correcta sin necesidad de worktrees.

- No inferir dependencias por heurísticas de paths, nombres o secciones.
- Si dos tasks de la misma wave modifican el mismo archivo y producen
  conflicto al integrar, tratar como dependencia faltante según el modo
  de autonomía, no intentar merge de worktrees.
```

### Qué no cambia

- El resto del loop (transitions, ACs, quality gates, commits, compaction) queda idéntico.
- `parallel = true` en `.roadmapctl.toml` ya está configurado en todos los repos.
- No requiere cambios en `roadmapctl` CLI ni en el toolset del skill.
- `parallel-independent-tasks: false` en `SKILL.md` no es relevante: ese campo controla paralelismo automático del harness, no el dispatch manual con `Agent` tool.

## Scope

Un único archivo: `.claude/skills/roadmap/loop-subcommand.md`

Cambio acotado a la sección "Parallel waves" (~5 líneas reemplazadas).

## Archivo a editar

Fuente canónica: `/home/shared/roadmapctl/.claude/skills/roadmap/loop-subcommand.md`

Después del cambio, sincronizar a user-scope:
```bash
./scripts/sync-roadmap-skill.sh --install
```

## Verificación

Después del cambio y sync:
1. Ejecutar `./scripts/sync-roadmap-skill.sh --check` para confirmar que source e installed coinciden.
2. Ejecutar headless verification definida en `SKILL.md` para confirmar comportamiento de bootstrap:
   ```bash
   PI_SKIP_VERSION_CHECK=1 pi --no-extensions --skill .claude/skills/roadmap/SKILL.md --tools read,bash -p 'HEADLESS VERIFICATION TEST...'
   ```
