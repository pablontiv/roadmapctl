---
estado: Completed
tipo: task
---
# T005: Implementar bootstrap inspect/init

**Outcome**: [O03 Config/context/workspace](README.md)
**Contribuye a**: CE3

[[blocked_by:./T004-implement-context-command.md]]

## Preserva

- INV1: Bootstrap es la única excepción controlada para crear config/root/schema inicial.
  - Verificar: docs y tests de dry-run/apply.

## Contexto

Hoy el skill puede crear `<roadmap-root>/` y `.stem` durante `plan`. Esto debe convertirse en una operación explícita y diagnosticable de `roadmapctl`, no en writes implícitos del prompt.

## Alcance

**In**:
1. Diseñar comandos `bootstrap inspect` y `bootstrap init` o alternativa aprobada.
2. `inspect` es read-only y muestra faltantes.
3. `init --dry-run` muestra archivos propuestos.
4. `init --apply` escribe solo `.roadmapctl.toml`, `<roadmap-root>/` y `<roadmap-root>/.stem` si se aprueba.
5. Validar path containment y postcheck.

**Out**:
- Crear Outcomes/Tasks.
- Auto-fix de roadmaps existentes.

## Estado inicial esperado

- El skill tiene bootstrap prose en `plan-subcommand.md`.

## Criterios de Aceptación

- Bootstrap inspect funciona sin modificar archivos.
- Apply requiere flag explícito y no escribe fuera del roadmap root/config permitido.
- Tests cubren missing config, missing root, missing `.stem` y existing files.

## Fuente de verdad

- `.claude/skills/roadmap/plan-subcommand.md`
- `.claude/skills/roadmap/base.stem`
- `internal/fsx/path.go`
