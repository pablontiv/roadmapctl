---
estado: Completed
tipo: task
---
# T008: Integrar checks respaldados por Rootline JSON

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE2

[[blocked_by:./T005-wrap-rootline-cli.md]]
[[blocked_by:./T007-implement-structure-checks.md]]

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: se consumen comandos genéricos JSON.
- INV2: `roadmapctl` no invoca subprocess con shell strings.
  - Verificar: todos los comandos pasan por `internal/rootlinecli`.

## Contexto

Los checks estructurales cubren filesystem, pero Rootline ya sabe validar `.stem`, frontmatter, query y graph. `roadmapctl` debe usar esos contratos genéricos para agregar garantías de negocio.

## Alcance

**In**:
1. Ejecutar y parsear `rootline validate --all <roadmap-root> -o json`.
2. Ejecutar y parsear `rootline describe <roadmap-root>/ -o json`.
3. Ejecutar y parsear `rootline query <roadmap-root> --where <leaf-filter> --where 'tipo == "task"' -o json`.
4. Ejecutar y parsear `rootline graph <roadmap-root> --where <leaf-filter> -o json`.
5. Convertir errores Rootline a diagnostics `RMC_ROOTLINE_*`.
6. Detectar ciclos y links bloqueantes rotos.
7. Validar `estado` y `tipo` contra schema/config.

**Out**:
- Usar `rootline graph --check` como parser machine-readable.
- Intentar reparar datos.
- Modificar schema Rootline.

## Estado inicial esperado

- `internal/rootlinecli` existe.
- Checks estructurales existen.

## Criterios de Aceptación

- Fixture con ciclo falla con diagnostic de ciclo y exit `1`.
- Fixture con link `blocked_by` roto falla con diagnostic bloqueante.
- Fixture con status fuera de enum falla.
- Si Rootline falta, `check` falla con exit `3` salvo modo explícito futuro.

## Fuente de verdad

- `internal/rootlinecli/`
- `internal/roadmap/dependencies.go`
- `internal/roadmap/status.go`
- `testdata/fixtures/invalid-cycle/`
