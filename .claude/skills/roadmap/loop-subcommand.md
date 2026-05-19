# /roadmap loop [--filter PATTERN] [--max N]

Ejecutar tasks pendientes usando la configuración efectiva devuelta por `roadmapctl bootstrap`. El loop acepta solo `--filter`, `--max` y el flag global `--repo` (workspace mode).

Ruta normal autosuficiente: este archivo contiene el procedimiento operativo completo. No leer `common-logic.md` ni documentación de integración para ejecutar el loop; esos documentos son referencia de mantenimiento/troubleshooting.

## Opciones CLI permitidas

- `--filter PATTERN`: filtrar por path (`O01`, `T003`, slug, etc.).
- `--max N`: límite de esta ejecución. Tiene precedencia sobre `loop_max_tasks`.

Los flags de comportamiento históricos `--parallel`, `--worktree`, `--self-pace`, `--skip-reviews`, `--checkpoint-interval` y `--pr` están obsoletos; no documentarlos ni aceptarlos como comportamiento activo. Usar los campos de configuración `parallel`, `autonomy`, `compact_after_task_commit`, `pr_mode`, `pr_merge_strategy`, `commit_style`, `auto_push` y `outcome_close_verify` expuestos por `roadmapctl bootstrap`.

## Config efectiva

Del JSON de bootstrap/context leer:

- `loop_max_tasks`: límite repo-local; `0` significa sin límite.
- `parallel`: permite waves oportunistas cuando sea seguro.
- `autonomy`: `manual`, `supervised` o `until_done`.
- `compact_after_task_commit`: compactar contexto tras una task durable.
- `pr_mode`: activar workflow de PR por scope.
- `pr_merge_strategy`, `commit_style`, `auto_push`, `outcome_close_verify`.

Calcular `effective_max` así:

1. Si `--max N` está presente, `effective_max = N`.
2. Si no, `effective_max = loop_max_tasks`.
3. Si `effective_max == 0`, no limitar la cola.

## Autonomy

- `manual`: ejecutar una task/wave y preguntar antes de continuar. Si se descubre dependencia faltante, sugerir el `blocked_by` requerido y detenerse.
- `supervised`: continuar entre tasks/waves sin preguntar; preguntar antes de ediciones estructurales del roadmap como agregar `blocked_by`.
- `until_done`: continuar hasta agotar ready queue o `effective_max`. Puede aplicar reparaciones estructurales seguras de `blocked_by`, pero cada mutación debe ir seguida de `roadmapctl check --strict` antes de continuar. Si no hay ruta determinística segura para editar, detenerse y reportar la dependencia requerida.

## Workspace mode

El loop opera en un repo a la vez.

- Con `--repo <name>`: usar ese repo.
- Sin `--repo`: contar pendientes por repo con `roadmapctl pending --repo <repo-path> --output json` y pedir selección.

## Fase 1: Discovery

### Preflight obligatorio roadmapctl

Antes de consultar o ejecutar tasks pendientes:

```bash
command -v roadmapctl
roadmapctl doctor --repo <repo-path> --output json --strict
roadmapctl check --repo <repo-path> --output json --strict
```

Si `roadmapctl` falta o cualquier comando sale non-zero, detenerse antes de seleccionar o ejecutar tasks. Reportar comando, exit code y diagnostic IDs si hubo JSON. No ejecutar tasks ni mutar estados.

**Staged async auto-update** — `internal/updater` implementa el patrón staged async: `FetchAndStage` descarga en background la versión más nueva a `~/.cache/roadmapctl/staged/<version>/`; `ApplyStagedIfAvailable` detecta el binario staged en la siguiente invocación, lo reemplaza atómicamente (os.Rename en Unix, copy+swap en Windows) y re-exec (syscall.Exec en Unix). El CLI wiring está en `internal/cli/cli.go`: `updater.CurrentVersion = version`, `ApplyStagedIfAvailable()`, y `go FetchAndStage(version)`. Errores de red y permisos son silenciosos. Para desactivar: `ROADMAPCTL_NO_UPDATE=1`. Cobertura del paquete: 85.4%. Gosec: `httpClient.Do(req)` requiere `//nolint:gosec` (G704 SSRF taint); `isNewer` no usa range loop sobre dos arrays para evitar falso positivo G602.

**Rootline binary staleness** — Cuando el loop ejecuta tasks en el repo `rootline` (o en repos que modifican `cmd/rootline/` o `internal/`), verificar que el binario instalado refleja los cambios recientes. Si el binario es stale, `roadmapctl next` puede devolver JSON formato v1 (sin `frontmatter` map) produciendo títulos vacíos y otros fallos silenciosos.

```bash
rootline --version                              # versión del binario instalado
git -C /home/shared/rootline log --oneline -1  # último commit de fuente
```

Si la fuente es más nueva que el binario, reconstruir antes de continuar:

```bash
go build -o $(which rootline) /home/shared/rootline/cmd/rootline
```

Nota CI: `go test ./...` funciona sin rootline instalado — `TestMain` activa el fake rootline automáticamente cuando `exec.LookPath("rootline")` falla (`ROADMAPCTL_FAKE_ROOTLINE=1`). El fake `describe` retorna el envelope completo `rootline/describe` (versión 1, schema, links, validate[]). Tests que requieren rootline real deben llamar `requiresRealRootline(t)` para saltearse automáticamente (ciclos, broken blocked_by, query/graph/tree, can-start/can-complete, decision scoring). La cobertura se verifica con `./scripts/check-coverage.sh` (umbral: 85.0%) en el job `smoke` (Ubuntu, macOS, Windows); el job `ci/Test` de crossbeam corre `go test ./... -race` sin gate de cobertura (instala fake rootline, no el real). Áreas de cobertura reciente: `bootstrap.go` (bootstrapApplyDiagnostic, renderBootstrap), `fsx/path.go` (symlink containment, prefix eval, ErrPathEscape para paths `/`-prefixed en Windows), `lint/schema_portability.go` (CheckFilenamePortability, reservedWindowsName, lintNameDiagnostic, arrayValue — cobertura cross-platform con tests que no dependen de filesystem case-sensitive ni chmod). Regla invariante: todo test que skippea con `runtime.GOOS` o filesystem case-insensitivo debe tener un gemelo cross-platform que cubra el mismo código desde otro ángulo.

1. Obtener estado determinístico de ejecución:
   ```bash
   roadmapctl next --repo <repo-path> --output json
   ```
   - Si `summary.status != "ok"` o el comando sale non-zero: reportar diagnostics y parar.
   - Usar `ready[]` como cola ejecutable; usar `blocked[]` solo para explicar skips/bloqueos.
   - `roadmapctl next`/`blocked_by` es la única fuente de dependencias para readiness y parallel waves.
   - No ejecutar `rootline graph`, `rootline query` o `rootline tree` ni postprocesar JSON crudo de Rootline para reconstruir la cola.
2. Obtener listado activo para tabla y conteos:
   ```bash
   roadmapctl pending --repo <repo-path> --output json
   ```
   - Si `summary.status != "ok"` o el comando sale non-zero: reportar diagnostics y parar.
3. Aplicar `--filter` por path sobre `ready[]` si existe.
4. Mantener el orden determinístico devuelto por `roadmapctl next`; no hacer topological sort manual.
5. Aplicar `effective_max` si es mayor que cero.
6. Renderizar tabla desde JSON (`ready[]`, `blocked[]` y `pending.tasks[]`).
7. Si no hay tasks en `ready[]` después del filtro: informar pendientes bloqueadas y parar.

## Fase 2: TodoList

Para cada task:

- subject: `TXXX: título`
- description: `Path: <filepath>`
- activeForm: `Implementando TXXX`

Mostrar `TaskList`.

## Fase 2.5: PR mode

Si `pr_mode == true`, leer [pr-workflow.md](pr-workflow.md) y ejecutar Branch & PR Detection. Si `pr_mode == false`, omitir workflow de PR.

## Observabilidad de procesos largos

- **Usar `Monitor`** (no `Bash` foreground bloqueante) cuando un proceso corre en background y queremos surfacear stdout línea-por-línea (tests, builds, agent dispatches durante una task). Patrón canónico: lanzar con `Bash` + `run_in_background: true` teando a `/tmp/roadmap-<task-id>.log`, luego `Monitor` con `grep -E --line-buffered` filtrando hitos (`PASS|FAIL|ERROR|heartbeat`).
- **Usar `ScheduleWakeup`** (no `bash sleep` loops) cuando hay que esperar estado externo que el harness no notifica: GitHub Actions runs (`gh run watch` bloquea; preferir wakeup + `gh run view --json status`), deploys, queues remotas.
- **Prohibido** encadenar `Bash sleep` para polling: elegir `Monitor` (stdout streamable) o `ScheduleWakeup` (poll interval externo) según el caso.
- **Instrucción directa del usuario**: si dice "monitorea" / "use monitor" / "watch this", invocar `Monitor` inmediatamente en el siguiente paso de proceso largo — no sustituir silenciosamente por `Bash background + sleep`.

Ejemplo del patrón completo:

```bash
# 1. Lanzar tests en background, tee a log
Bash(command="go test ./... -v 2>&1 | tee /tmp/roadmap-T001.log",
     run_in_background=true)

# 2. Streamear hitos del log
Monitor(command="grep -E --line-buffered '(PASS|FAIL|ok|---)' /tmp/roadmap-T001.log",
        description="T001 test run")
```

## Fase 3: Loop

Variables:

- `checkpoint_commit`: HEAD inicial.
- `checkpoint_task_count`: 0.
- `current_scope`: Outcome actual o `direct-tasks`.
- `checkpoint_interval`: 5 (quality gates siempre activos).

### Parallel waves

Si `parallel == true`, formar waves oportunistas desde `ready[]` usando solo la información canónica de `roadmapctl next` y `blocked_by`:

- Tasks en una misma wave no tienen dependencia explícita entre sí según `roadmapctl next`.
- No inferir dependencias por heurísticas de paths, nombres o secciones; si aparece un conflicto real durante integración, tratarlo como dependencia faltante.

Si `parallel == true`, ejecutar cada wave despachando llamadas paralelas al tool `Agent` — una por task de la wave. Las tasks de una wave son independientes por definición (`roadmapctl next` garantiza ausencia de `blocked_by` entre ellas), por lo que Agent calls paralelas sobre archivos distintos son la ruta correcta sin necesidad de worktrees.

Si dos tasks de la misma wave producen conflicto al integrar, tratar como dependencia faltante según el modo de autonomía — no usar worktrees para forzar el merge.

### Pre-dispatch: serializar tasks con `## Fuente de verdad` solapada

Antes de despachar la wave en paralelo, leer la sección `## Fuente de verdad` de cada task de `ready[]`. Si dos o más tasks declaran el mismo path como fuente de verdad, **serializarlas dentro de la wave** — ejecutar una, integrar, recalcular `roadmapctl next`, ejecutar la siguiente. Esta es ordenación de ejecución, no una dependencia estructural del roadmap: no agregar `blocked_by`, no mutar el grafo.

Sin este pre-check, dos agentes en paralelo editan el mismo archivo y producen conflicto al integrar; el handler "Conflictos por dependencia faltante" descrito abajo es reactivo y ya desperdició cómputo del agente. La detección por `Fuente de verdad` evita el desperdicio sin contradecir la regla de no inferir `blocked_by` por heurísticas — solo reordena la cola intra-wave.

Conflictos por dependencia faltante:

- `manual`: reportar el `blocked_by` recomendado y detenerse.
- `supervised`: pedir aprobación antes de aplicar `blocked_by`; luego `roadmapctl check --strict`.
- `until_done`: aplicar solo si la edición es determinística y segura; ejecutar `roadmapctl check --strict`; recalcular con `roadmapctl next`. Si no es seguro, detenerse y reportar.

Si `parallel == false`, ejecutar tasks en orden secuencial de `ready[]`.

Para cada task o wave ordenada:

1. **Verificar transición de inicio**
   ```bash
   roadmapctl transition can-start <task.md> --repo <repo-path> --output json
   ```
   - Usar el JSON de `roadmapctl`; no recalcular reglas de dependencias en prompt.
   - Si `allowed=false`, skip con `blocking_dependencies[]`/`diagnostics[]`.
   - No llamar `rootline set` directamente para iniciar tasks.

2. **Scope change**
   - Si cambia Outcome/direct scope y `pr_mode == true`, cerrar PR anterior si corresponde y ejecutar Outcome Setup.
   - Sin PR mode, solo actualizar `current_scope`.

3. **Marcar inicio**
   ```bash
   roadmapctl transition start <task.md> --apply --repo <repo-path> --output json
   ```
   Si `allowed=false`, `summary.status="error"`, o el comando sale non-zero, detenerse antes de ejecutar la task o commitear. `roadmapctl transition start --apply` es responsable de `rootline set`, `rootline validate` y postcheck; no duplicar esas reglas en prompt. Actualizar UI con `TaskUpdate` solo después de pasar.

4. **Leer task**
   Leer el archivo completo. La task debe ser suficiente para implementar.

5. **Implementar**
   Ejecutar exactamente el alcance de la task. Si hay una sección `## Especificación Técnica`, seguirla.

   Prohibido añadir trabajo fuera del spec de la task activa, aunque sea relacionado, conveniente o evidente. Si durante la implementación se detecta trabajo útil fuera del spec: anotarlo en el contexto de la sesión como candidata a nueva task, no implementarlo.

6. **Verificar ACs e invariantes**
   - Ejecutar cada AC.
   - Ejecutar cada verificación en `## Preserva` si existe.
   - Si falla algo: parar y reportar.

7. **Outcome close check**
   Si es la última task pendiente del Outcome, ejecutar comandos de `outcome_close_verify` si existen. Warning informativo, no bloqueo automático.

8. **Security review selectivo**
   Si se tocaron archivos sensibles (`secret`, `credentials`, `.env`, `auth`, `crypto`) o la task lo pide, ejecutar review de seguridad. Findings HIGH bloquean.

9. **Complete + commit**
   ```bash
   roadmapctl transition complete <task.md> --apply --repo <repo-path> --output json
   ```
   Ejecutar este comando solo después de que ACs e invariantes pasaron. Si `allowed=false`, `summary.status="error"`, o el comando sale non-zero, reportar diagnostics y detenerse antes de declarar completada la iteración o commitear. Si pasa: `git add` específico, commit según `commit_style`, push según `auto_push`, y PR bookkeeping según `pr_mode`.

10. **Actualizar UI y resumen**
   ```bash
   TaskUpdate <id> status: completed
   TaskOutput <id> "ACs: N/M passed | Commit: <hash>"
   ```
   Mostrar resultado de iteración.

11. **Compaction opcional**
   Si `compact_after_task_commit == true`, compactar solo después de que la task sea durable:
   1. ACs e invariantes pasaron.
   2. `roadmapctl transition complete --apply` pasó.
   3. Commit creado.
   4. Push/PR bookkeeping terminado o bloqueo reportado.

   Preferir la herramienta `compact_roadmap_context` con `task_path`, `commit_hash`, `validation_summary`, `next_work` y `config_summary`. Si no está disponible, usar `/compact <instrucciones roadmap>` como fallback. Fallar al compactar debe advertir claramente, pero no invalida una task ya completada y commiteada.

12. **Checkpoint**
   Activar si:
   - `checkpoint_task_count >= checkpoint_interval`,
   - cambia scope,
   - autonomía `manual` solicita pausa,
   - usuario decide parar.

   Revisar diff acumulado, reportar findings informativos y resetear checkpoint.

13. **Continuación**
   - `manual`: preguntar continuar, saltar siguiente o parar después de cada task/wave.
   - `supervised` y `until_done`: no preguntar entre tasks/waves; recalcular cola con `roadmapctl next` y continuar hasta agotar ready queue o `effective_max`.

14. **Reintentar bloqueadas**
   Al final, reintentar tasks cuyas dependencias pasaron a done. Si no progresa ninguna, parar por deadlock.

## Fase 4: Resumen final

```text
RESUMEN LOOP
├─ Tasks completadas: N/TOTAL
├─ Tasks saltadas: M
├─ ACs: passed/total
├─ Security reviews: N
├─ Quality checkpoints: N
├─ PRs: ... (si pr_mode)
├─ Commits: ...
└─ Tasks restantes: ...
```

→ Invocar skill `/retrospective` pasando en contexto:
  - `checkpoint_commit` (capturado al inicio de Fase 3)
  - Tasks completadas, saltadas y ACs del resumen anterior
