---
estado: Specified
tipo: task
---
# T032: Loop — frontera de verificación local/global para agentes paralelos

**Contribuye a**: Agentes despachados en paralelo no reportan como fallas los rompimientos transitorios del árbol que producen otros agentes de la misma wave; la verificación global queda exclusivamente en manos del loop principal.

## Contexto

En la wave 4 de O25 el agente T003 (delete dead diff) reportó `go build ./... exit 1` como "pre-existing failure" mientras T002 (fsx → pathsec) aún estaba editando imports. Esa salida confundió al loop: técnicamente T003 había hecho su trabajo (eliminar el paquete), pero su AC global (`go build` exit 0) no podía evaluarse durante la wave.

El loop ya tiene la regla T026 ("re-verificación directa post-subagente paralelo") que cubre la verificación correcta en el loop principal. Falta la contraparte: instruir al agente para que NO ejecute comandos globales — solo locales a sus archivos.

## Alcance

**In**:
1. En `.claude/skills/roadmap/loop-subcommand.md`, Fase 3 sección "Parallel waves" (antes del párrafo "Pre-dispatch: serializar tasks con Fuente de verdad solapada"):
   - Definir explícitamente la **frontera de verificación** del agente paralelo:
     - **Local**: grep/find sobre el path tocado, lint sobre el paquete tocado, tests sobre el paquete tocado, `wc -l`/`ls` sobre archivos producidos.
     - **Global** (prohibido al agente): `go build ./...`, `go test ./...`, `golangci-lint run ./...`, cualquier verificación que cubra todo el repo.
   - Aclarar que durante la wave el árbol puede estar transitoriamente roto si los agentes tocan imports cruzados aunque sus fuentes de verdad sean disjuntas — eso no es falla del scope del agente.
   - Toda verificación global es responsabilidad del loop principal en el paso 6 de re-verificación directa (regla T026).
2. Agregar al template de prompt para agents paralelos (si existe canónico) una instrucción explícita: "limitá tus ACs verificables a comandos locales a tus archivos; no ejecutes builds/tests/lints globales".

**Out**:
- No cambiar el mecanismo de dispatch paralelo (waves, fuente de verdad overlap, etc.).
- No agregar worktrees ni isolation adicional — la regla es a nivel de instrucción al agente, no de aislamiento del filesystem.
- No relajar T026 ni la verificación post-agente del loop.

## Criterios de Aceptación

- La sección "Parallel waves" especifica la frontera local/global con ejemplos concretos.
- La regla menciona explícitamente que los ACs globales son responsabilidad del loop principal post-wave, no del agente.
- El skill instruye que el agente debe tolerar (no reportar como falla) los rompimientos transitorios del árbol durante la wave si están fuera de su scope.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
