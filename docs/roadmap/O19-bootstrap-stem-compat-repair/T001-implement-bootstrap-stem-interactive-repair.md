---
estado: Specified
tipo: task
---
# T001: Implement bootstrap stem interactive repair

**Outcome**: [O19 Bootstrap stem compat repair](README.md)
**Contribuye a**: bootstrap detecta .stem incompatible y ofrece repararlo con confirmación del usuario

## Preserva

- INV1: `roadmapctl doctor` y `roadmapctl check` siguen siendo read-only; la reparación es un paso explícito adicional, no una auto-corrección silenciosa.
  - Verificar: `roadmapctl doctor --strict` sin flags de repair no modifica ningún archivo.
- INV2: La reparación solo toca `<roadmap-root>/.stem`; no modifica outcomes, tasks ni ningún otro archivo del roadmap.
  - Verificar: git diff después de repair muestra solo `.stem`.

## Contexto

Cuando un repo tiene un `.stem` con `estado.required.match: ["O*", "T*"]` o una regla
`validate estado non_empty` global, `roadmapctl check --strict` falla con:
- `RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED`
- `RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY`

Esto bloquea todas las escrituras del `/roadmap` skill. Actualmente la única solución es
copiar manualmente el `.stem` canónico. Este task agrega a `roadmapctl bootstrap` la
capacidad de detectar esos diagnostics específicos y ofrecer la reparación de forma segura
e interactiva.

El `.stem` canónico tiene `required.match: ["T*"]` solo y no tiene `validate estado non_empty`.

## Alcance

**In**:
1. En `internal/cli/bootstrap.go` (o el handler de bootstrap): después de correr doctor/check internamente, si los diagnostics incluyen `RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED` o `RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY`, mostrar el problema y el diff propuesto.
2. Preguntar al usuario: `Update .stem to canonical schema? [y/N]`.
3. Aceptar `--yes` como flag para confirmar sin prompt interactivo (útil en CI/headless).
4. Aplicar solo si el usuario confirma: escribir el `.stem` canónico sobre el existente.
5. Si el `.stem` tiene campos custom no reconocidos (no es el template legacy conocido), emitir `RMC_BOOTSTRAP_REPAIR_UNSUPPORTED_STEM` y no modificar nada.
6. Después de aplicar, re-correr `roadmapctl check --strict` internamente y reportar si pasó.

**Out**:
- No reparar ningún otro diagnostic fuera de los dos mencionados.
- No tocar outcomes, tasks ni otros archivos del roadmap.
- No cambiar el comportamiento de `roadmapctl doctor` ni `roadmapctl check` standalone.

## Estado inicial esperado

- `roadmapctl bootstrap` existe y corre doctor/check internamente.
- Los diagnostics `RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED` y `RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY` están implementados (O11, Completed).

## Criterios de Aceptación

- `roadmapctl bootstrap --repo <repo-con-stem-legacy>` reporta los diagnostics bloqueantes, muestra el diff del `.stem` y pregunta `[y/N]`.
- Respondiendo `y` (o usando `--yes`): el `.stem` se actualiza al canónico.
- `roadmapctl check --strict` pasa en el repo afectado después de aplicar la reparación.
- Respondiendo `N`: no se modifica ningún archivo; bootstrap reporta el bloqueo como antes.
- Si el `.stem` tiene campos no reconocibles: diagnostic `RMC_BOOTSTRAP_REPAIR_UNSUPPORTED_STEM`, sin modificaciones.
- `golangci-lint run ./...` reporta 0 issues.

## Fuente de verdad

- `internal/cli/bootstrap.go`
- `internal/cli/schema_compatibility.go` (lógica de detección existente de O11)
- Canonical `.stem` template: `docs/roadmap/.stem` de este repo
