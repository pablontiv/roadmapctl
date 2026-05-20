---
estado: Specified
tipo: task
---
# T031: Loop — cross-check de coherencia output↔filesystem en discovery

**Contribuye a**: El loop detecta cuando la salida de `roadmapctl next`/`pending` no coincide con el filesystem (binario stale, índice desactualizado) y se detiene antes de ejecutar tasks con datos incorrectos.

## Contexto

En la sesión O25 el binario `roadmapctl` estaba construido desde un commit anterior al HEAD local (`24b7f66` vs HEAD `1192e78`). Eso ocurre porque el pre-push hook del repo solo rebuildea el binario en `git push` — commits locales sin push lo dejan stale.

El síntoma observado fue que `roadmapctl pending` reportó slugs (`T028-loop-check-ci-green-before-ci-invariant-tasks.md`) que no existían en el filesystem; los archivos reales tenían slugs distintos (`T028-loop-ci-green-prereq-check.md`). El loop investigó cada incoherencia in-flight y perdió tiempo.

Rebuildear no le toca al skill (esa es responsabilidad del CLI orquestador o del entorno). Pero detectar la incoherencia y detenerse sí está dentro del scope del skill.

## Alcance

**In**:
1. En `.claude/skills/roadmap/loop-subcommand.md`, Fase 1 paso 2 (después de `roadmapctl pending`):
   - Para cada `path` en `next.ready[]` y `pending.tasks[]`, verificar que el archivo existe como tracked en el filesystem: `git ls-files <roadmap_root>/<path>` retorna ese path exactamente.
   - Si algún path reportado no aparece tracked (o aparece un slug que no matchea ningún archivo en `git ls-files`), emitir mensaje `STALE TOOLCHAIN OUTPUT` listando los slugs sospechosos vs los archivos reales más cercanos.
   - Detener el loop antes de pasar a Fase 2 (TodoList).
2. Documentar explícitamente que el skill NO intenta rebuildear ningún binario — la reparación queda al operador (recargar shell, `make install`, rebuildear con `go build`, etc.).

**Out**:
- No rebuildear binarios ni ejecutar `install-user.sh` ni equivalentes.
- No relajar el check (es bloqueante por diseño — datos incoherentes producen ejecuciones incorrectas).
- No agregar cross-check sobre `blocked[]` (esos paths pueden referenciar tasks que también están stale; basta con bloquear desde `ready[]`/`pending`).

## Criterios de Aceptación

- `.claude/skills/roadmap/loop-subcommand.md` Fase 1 paso 2 incluye la regla de cross-check.
- El criterio de "existe tracked" usa `git ls-files`, no `test -f` (un archivo untracked no es canónico).
- El mensaje `STALE TOOLCHAIN OUTPUT` enumera explícitamente paths reportados vs archivos reales y propone como acción "refrescar herramientas y reintentar"; no propone bypass.
- El skill establece explícitamente que rebuildear no es responsabilidad del skill.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
