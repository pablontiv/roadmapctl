---
estado: Pending
tipo: task
---
# T003: Implementar fallback y migración legacy

**Outcome**: [O03 Config/context/workspace](README.md)
**Contribuye a**: CE2, INV2

[[blocked_by:./T002-implement-config-discovery-and-toml-loader.md]]

## Preserva

- INV1: Repos existentes con `.claude/roadmap.local.md` siguen funcionando.
  - Verificar: fixture legacy.

## Contexto

La config legacy vive en `.claude/roadmap.local.md`. La nueva config preferida vive en `<roadmap-root>/.roadmapctl.toml`. La migración debe ser explícita, reversible y no destructiva.

## Alcance

**In**:
1. Implementar fallback legacy cuando no existe TOML.
2. Si TOML y legacy coexisten, usar TOML y emitir warning/info si difieren.
3. Diseñar e implementar flujo de migración dry-run/apply si se aprueba el nombre del comando.
4. Generar TOML equivalente desde config legacy.

**Out**:
- Borrar `.claude/roadmap.local.md` automáticamente.
- Cambiar `.claude/.stem` sin aprobación.

## Estado inicial esperado

- Loader TOML existe.
- Parser legacy aún disponible.

## Criterios de Aceptación

- `valid-legacy-config-fallback` pasa.
- `warning-config-conflict` muestra config source y warning determinístico.
- Migración dry-run muestra path destino y contenido propuesto sin escribir.

## Fuente de verdad

- `.claude/roadmap.local.md`
- `.claude/.stem`
- `internal/config/config.go`
- `docs/roadmap-skill-integration.md`
