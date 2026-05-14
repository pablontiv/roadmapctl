---
estado: Completed
tipo: task
---
# T003: Update skill and docs for bootstrap stem repair

**Outcome**: [O19 Bootstrap stem compat repair](README.md)
**Contribuye a**: agentes y docs reflejan el nuevo repair path de bootstrap

[[blocked_by:./T002-add-tests-for-bootstrap-stem-repair.md]]

## Preserva

- INV1: El flujo normal de `/roadmap plan` (sin stem legacy) no cambia.
  - Verificar: headless Pi test del flujo normal sigue pasando.
- INV2: El guard de preflight (`doctor/check --strict` antes de escribir) sigue siendo obligatorio.
  - Verificar: `plan-subcommand.md` mantiene el gate.

## Contexto

Con T001 y T002 implementados, el skill y la documentación deben reflejar el nuevo comportamiento:
cuando el preflight falla por stem legacy, el agente puede invocar `roadmapctl bootstrap` (con
`--yes` en modo autónomo) para resolver el bloqueo antes de reintentar el preflight.

## Alcance

**In**:
1. `plan-subcommand.md` §3.2 (Preflight obligatorio): agregar nota — si doctor/check falla con `RMC_LINT_SCHEMA_OUTCOME_ESTADO_*`, ejecutar `roadmapctl bootstrap --repo <repo> --roadmap-root <root> --yes` y reintentar el preflight.
2. `docs/cli-contract.md`: documentar el nuevo comportamiento de bootstrap (repair path, diagnostics, `--yes` flag).
3. Headless Pi test: verificar que bootstrap detecta y reporta los diagnostics en un fixture con stem legacy (no requiere confirm real, solo detección).

**Out**:
- No cambiar el contrato de doctor/check.
- No agregar auto-repair silencioso al skill; el repair sigue siendo explícito (--yes o prompt).

## Estado inicial esperado

- T001 y T002 completados: bootstrap repair implementado y con tests.
- `plan-subcommand.md` existe en `.claude/skills/roadmap/`.
- `docs/cli-contract.md` existe.

## Criterios de Aceptación

- `plan-subcommand.md` §3.2 describe el repair path con el comando exacto.
- `docs/cli-contract.md` tiene sección de bootstrap repair con diagnostics y flag `--yes`.
- El headless Pi test pasa, confirmando detección del stem legacy.
- `./scripts/sync-roadmap-skill.sh --install` (si existe) actualiza la skill distribuida.

## Fuente de verdad

- `.claude/skills/roadmap/plan-subcommand.md`
- `docs/cli-contract.md`
- `internal/cli/bootstrap.go`
