---
estado: Completed
tipo: task
---
# T001: Actualizar contrato CLI post-MVP

**Outcome**: [O09 Release/governance](README.md)
**Contribuye a**: CE1

[[blocked_by:../O03-config-context-workspace/T004-implement-context-command.md]]

## Preserva

- INV1: Contrato JSON/exit codes existente no se rompe sin versionado.
  - Verificar: golden tests.

## Contexto

`docs/cli-contract.md` hoy describe el MVP `doctor`/`check`. Debe evolucionar para documentar comandos post-MVP cuando existan.

## Alcance

**In**:
1. Documentar comandos directos: context, pending, next, decision, lint, transition, materialize.
2. Documentar JSON `kind`, fields y diagnostics por comando.
3. Separar contrato MVP histórico de contrato vigente.
4. Documentar `.stem` authority y `.roadmapctl.toml` config.

**Out**:
- Documentar comandos no implementados como disponibles antes de tiempo.

## Estado inicial esperado

- Al menos `context` implementado; otros pueden documentarse como roadmap/preview si aplica.

## Criterios de Aceptación

- Docs distinguen implemented vs planned.
- No hay mención de `roadmapctl roadmap ...`.
- `go test ./...` y docs checks pasan.

## Fuente de verdad

- `docs/cli-contract.md`
- `README.md`
- `docs/release.md`
