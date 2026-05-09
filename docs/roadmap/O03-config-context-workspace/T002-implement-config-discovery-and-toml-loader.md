---
estado: Completed
tipo: task
---
# T002: Implementar discovery y loader TOML

**Outcome**: [O03 Config/context/workspace](README.md)
**Contribuye a**: CE1, INV1

[[blocked_by:./T001-design-roadmapctl-toml-config-contract.md]]
[[blocked_by:../O02-post-mvp-foundations/T001-adopt-cobra-and-community-packages.md]]

## Preserva

- INV1: `--roadmap-root` y path containment siguen protegiendo escapes.
  - Verificar: root escape tests.

## Contexto

El loader actual hardcodea `.claude/roadmap.local.md`. Debe preferir `<roadmap-root>/.roadmapctl.toml` y usar una librería TOML mantenida.

## Alcance

**In**:
1. Agregar `github.com/pelletier/go-toml/v2`.
2. Descubrir config en `docs/roadmap/.roadmapctl.toml` por defecto.
3. Si se pasa `--roadmap-root`, buscar `<roadmap-root>/.roadmapctl.toml`.
4. Inferir roadmap root desde dirname del TOML preferido.
5. Usar defaults cuando el TOML no exista pero root/schema sí existan.
6. Mantener path containment con `fsx.ResolveInside`.

**Out**:
- Migración automática desde legacy.
- Workspace multi-repo.

## Estado inicial esperado

- `internal/config/config.go` parsea frontmatter custom.

## Criterios de Aceptación

- Fixture `valid-roadmapctl-toml-default` carga correctamente.
- Legacy behavior no se rompe cuando no hay TOML.
- Errores TOML tienen diagnostics claros y línea/campo cuando sea posible.

## Fuente de verdad

- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/fsx/path.go`
