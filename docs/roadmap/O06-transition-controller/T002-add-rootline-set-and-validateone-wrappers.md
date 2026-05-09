---
estado: Pending
tipo: task
---
# T002: Agregar wrappers Rootline `set` y validate-one

**Outcome**: [O06 Controlador de transiciones](README.md)
**Contribuye a**: CE3

[[blocked_by:../O02-post-mvp-foundations/T004-harden-rootline-json-on-nonzero-exit.md]]

## Preserva

- INV1: Subprocess seguro: args explícitos, timeout y sin shell strings.
  - Verificar: tests de rootlinecli.

## Contexto

Para aplicar transiciones, `roadmapctl` debe mutar frontmatter vía Rootline genérico `set` y validar el archivo. Rootline `set` puede no tener JSON estable, así que el wrapper debe tratar stdout/stderr cuidadosamente.

## Alcance

**In**:
1. Agregar wrapper `Set(ctx, file, assignments...)`.
2. Agregar helper validate-one si hace falta.
3. Tests de args, cwd, timeout, stderr y exit code.
4. Manejar comandos no JSON sin romper JSON stdout de roadmapctl.

**Out**:
- Implementar reglas de transición.
- Modificar archivos reales fuera de tests temp.

## Estado inicial esperado

- rootlinecli tiene validate/describe/query/graph.

## Criterios de Aceptación

- Tests demuestran que `rootline set <file> field=value` usa args explícitos.
- Errores de `set` producen diagnostics útiles.
- No se imprime stdout Rootline crudo en JSON mode de roadmapctl.

## Fuente de verdad

- `internal/rootlinecli/client.go`
- `internal/rootlinecli/client_test.go`
- `docs/release.md`
