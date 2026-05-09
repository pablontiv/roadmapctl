---
estado: Completed
tipo: task
---
# T005: Soportar bootstrap materialization aprobado

**Outcome**: [O07 Controlador de materialización](README.md)
**Contribuye a**: CE2

[[blocked_by:../O03-config-context-workspace/T005-implement-bootstrap-inspect-init.md]]
[[blocked_by:./T003-implement-materialize-dry-run-and-diff.md]]

## Preserva

- INV1: Bootstrap es explícito y acotado; no habilita writes arbitrarios.
  - Verificar: path allowlist.

## Contexto

Un roadmap nuevo puede no tener root, `.stem` o `.roadmapctl.toml`. La excepción bootstrap debe ser explícita, dry-run/apply y segura.

## Alcance

**In**:
1. Permitir que materialize integre bootstrap aprobado cuando falta root/schema/config.
2. Escribir solo archivos permitidos: `.roadmapctl.toml`, `.stem`, README/tasks canónicos.
3. Reusar base `.stem` y config defaults aprobados.
4. Postcheck completo.

**Out**:
- Crear estructuras fuera del roadmap root.
- Inventar schema documental en roadmapctl; usar template base versionado.

## Estado inicial esperado

- Bootstrap inspect/init existe.
- Dry-run materialize existe.

## Criterios de Aceptación

- Fixture missing root produce dry-run de bootstrap claro.
- Apply bootstrap no escribe fuera de paths permitidos.
- Existing `.stem` no se sobrescribe sin aprobación.

## Fuente de verdad

- `.claude/skills/roadmap/base.stem`
- `.claude/skills/roadmap/plan-subcommand.md`
- `internal/fsx/path.go`
