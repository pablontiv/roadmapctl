# /roadmap loop [--filter PATTERN] [--max N] [--pr]

> Pre-requisito: leer [common-logic.md](common-logic.md).

Ejecutar tasks pendientes en loop con confirmación entre cada una.

## Opciones

- `--filter PATTERN`: filtrar por path (`O01`, `T003`, slug, etc.).
- `--max N`: limitar a N tasks.
- `--checkpoint-interval N`: checkpoint de calidad cada N tasks (default 5).
- `--skip-reviews`: desactivar quality gates.
- `--pr`: crear branch/PR por Outcome o grupo de tasks directas.
- `--worktree`: crear git worktree aislado por Outcome via `EnterWorktree`/`ExitWorktree` si esas tools están disponibles. Al entrar a un Outcome se crea un worktree; al cerrar se limpia con `ExitWorktree`. Requiere que el repo soporte worktrees.
- `--self-pace`: usar `ScheduleWakeup` si está disponible entre tasks para loops de larga duración. Mantiene el cache caliente (delay ~270s) sin bloquear la sesión.
- `--parallel`: ejecutar tasks independientes dentro de un Outcome en paralelo via `Agent` tool usando el `execution-model` declarado en el frontmatter del skill.

## Workspace mode

El loop opera en un repo a la vez.

- Con `--repo <name>`: usar ese repo.
- Sin `--repo`: contar pendientes por repo con `roadmapctl pending --repo <repo-path> --roadmap-root <roadmap-root> --output json` y pedir selección.

## Fase 1: Discovery

### Preflight obligatorio roadmapctl

Antes de consultar o ejecutar tasks pendientes:

```bash
command -v roadmapctl
roadmapctl doctor --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
roadmapctl check --repo <repo-path> --roadmap-root <roadmap-root> --output json --strict
```

Si `roadmapctl` falta o cualquier comando sale non-zero, detenerse antes de seleccionar o ejecutar tasks. Reportar comando, exit code y diagnostic IDs si hubo JSON. No ejecutar tasks ni mutar estados.

1. Obtener estado determinístico de ejecución:
   ```bash
   roadmapctl next --repo <repo-path> --roadmap-root <roadmap-root> --output json
   ```
   - Si `summary.status != "ok"` o el comando sale non-zero: reportar diagnostics y parar.
   - Usar `ready[]` como cola ejecutable; usar `blocked[]` solo para explicar skips/bloqueos.
   - No recalcular dependencias en prompt: `roadmapctl next` es la fuente canónica de `blocking_dependencies`/blockers.
2. Obtener listado activo para tabla y conteos:
   ```bash
   roadmapctl pending --repo <repo-path> --roadmap-root <roadmap-root> --output json
   ```
   - Si `summary.status != "ok"` o el comando sale non-zero: reportar diagnostics y parar.
3. Aplicar `--filter` por path sobre `ready[]` si existe.
4. Mantener el orden determinístico devuelto por `roadmapctl next`; no hacer topological sort manual.
5. Aplicar `--max`.
6. Renderizar tabla desde JSON (`ready[]`, `blocked[]` y `pending.tasks[]`).
7. Si no hay tasks en `ready[]` después del filtro: informar pendientes bloqueadas y parar.

## Fase 2: TodoList

Para cada task:

- subject: `TXXX: título`
- description: `Path: <filepath>`
- activeForm: `Implementando TXXX`

Mostrar `TaskList`.

## Fase 2.5: PR mode

Si `--pr`, leer [pr-workflow.md](pr-workflow.md) y ejecutar Branch & PR Detection.

## Fase 2.6: Worktree setup

Solo si `--worktree` (o `worktree-per-outcome: true` en el frontmatter del skill). Sin este flag → skip.

Al entrar a un Outcome en el loop:
- `EnterWorktree` con nombre derivado del Outcome ID (ej: `outcome-O01`).
- Todos los commits de ese Outcome ocurren dentro del worktree.
- Al cerrar el Outcome (última task completada, o loop interrumpido): `ExitWorktree`.

## Fase 3: Loop

Variables:

- `checkpoint_commit`: HEAD inicial.
- `checkpoint_task_count`: 0.
- `current_scope`: Outcome actual o `direct-tasks`.
- `checkpoint_interval`: default 5.

Para cada task ordenada:

1. **Verificar transición de inicio**
   ```bash
   roadmapctl transition can-start <task.md> --repo <repo-path> --roadmap-root <roadmap-root> --output json
   ```
   - Usar el JSON de `roadmapctl`; no recalcular reglas de dependencias en prompt.
   - Si `allowed=false`, skip con `blocking_dependencies[]`/`diagnostics[]`.
   - No llamar `rootline set` directamente para iniciar tasks.

2. **Scope change**
   - Si cambia Outcome/direct scope y `--pr`, cerrar PR anterior si corresponde y ejecutar Outcome Setup.
   - Sin `--pr`, solo actualizar `current_scope`.

3. **Marcar inicio**
   ```bash
   roadmapctl transition start <task.md> --apply --repo <repo-path> --roadmap-root <roadmap-root> --output json
   ```
   Si `allowed=false`, `summary.status="error"`, o el comando sale non-zero, detenerse antes de ejecutar la task o commitear. `roadmapctl transition start --apply` es responsable de `rootline set`, `rootline validate` y postcheck; no duplicar esas reglas en prompt. Actualizar UI con `TaskUpdate` solo después de pasar.

4. **Leer task**
   Leer el archivo completo. La task debe ser suficiente para implementar.

5. **Implementar**
   Ejecutar exactamente el alcance de la task. Si hay una sección `## Especificación Técnica`, seguirla.

5.5. **Paralelismo** (solo si `--parallel`):
   Solo paralelizar si se puede probar que las tasks son independientes:
   - no hay dependencia `blocked_by` entre ellas,
   - sus secciones `## Fuente de verdad` no solapan paths,
   - ninguna toca archivos globales/sensibles,
   - cada subagente trabaja en worktree aislado o en archivos disjuntos.

   Si cualquiera de esas condiciones no puede probarse, ejecutar secuencialmente. Si procede, invocarlas como subagentes en paralelo via `Agent` tool con `model` igual al `execution-model` del frontmatter del skill. Consolidar resultados antes de continuar a verificación de ACs.

6. **Verificar ACs e invariantes**
   - Ejecutar cada AC.
   - Ejecutar cada verificación en `## Preserva` si existe.
   - Si falla algo: parar y reportar.

7. **Outcome close check**
   Si es la última task pendiente del Outcome, ejecutar comandos de `<outcome-close-cmds>` si existen. Warning informativo, no bloqueo automático.

8. **Security review selectivo**
   Si se tocaron archivos sensibles (`secret`, `credentials`, `.env`, `auth`, `crypto`) o la task lo pide, ejecutar review de seguridad. Findings HIGH bloquean.

9. **Commit**
   ```bash
   roadmapctl transition complete <task.md> --apply --repo <repo-path> --roadmap-root <roadmap-root> --output json
   ```
   Ejecutar este comando solo después de que ACs e invariantes pasaron. Si `allowed=false`, `summary.status="error"`, o el comando sale non-zero, reportar diagnostics y detenerse antes de declarar completada la iteración o commitear. Si pasa: `git add` específico, commit según `<commit-style>`, push según `<auto-push>` y `--pr`.

10. **Actualizar UI y resumen**
   ```bash
   TaskUpdate <id> status: completed
   TaskOutput <id> "ACs: N/M passed | Commit: <hash>"
   ```
   Mostrar resultado de iteración.

10.5. **Self-pace** (solo si `--self-pace`):
   Si quedan más de 3 tasks en la cola:
   ```
   ScheduleWakeup(delaySeconds: 270, reason: "loop roadmap: <N> tasks restantes — <siguiente task>")
   ```
   Mantiene el cache caliente (< 5 min TTL) entre iteraciones largas.

11. **Checkpoint**
   Activar si:
   - `checkpoint_task_count >= checkpoint_interval`,
   - cambia scope,
   - usuario decide parar.

   Revisar diff acumulado, reportar findings informativos y resetear checkpoint.

12. **Confirmar continuación**
   Preguntar: continuar, saltar siguiente, o parar.

13. **Reintentar bloqueadas**
   Al final, reintentar tasks cuyas dependencias pasaron a done. Si no progresa ninguna, parar por deadlock.

## Fase 4: Resumen final

```text
RESUMEN LOOP
├─ Tasks completadas: N/TOTAL
├─ Tasks saltadas: M
├─ ACs: passed/total
├─ Security reviews: N
├─ Quality checkpoints: N
├─ PRs: ... (si --pr)
├─ Commits: ...
└─ Tasks restantes: ...
```
