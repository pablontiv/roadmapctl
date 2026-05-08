---
estado: Completed
tipo: task
---
# T001: Definir contrato CLI de roadmapctl

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE1, CE2 y CE3

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: no se proponen subcomandos roadmap dentro de `rootline`.
- INV3: El MVP no materializa ni corrige automáticamente.
  - Verificar: contrato limita comandos a `doctor` y `check`.

## Contexto

`roadmapctl` será el guard obligatorio para comandos `/roadmap` que escriben, mutan o ejecutan. Esta task define la interfaz pública antes de implementar código, para evitar acoplar decisiones de producto al diseño interno.

## Alcance

**In**:
1. Definir comandos MVP: `roadmapctl doctor` y `roadmapctl check`.
2. Definir flags: `--repo`, `--roadmap-root`, `--workspace`, `--output json|text`, `--strict`, `--rootline`, `--timeout`.
3. Definir exit codes: `0`, `1`, `2`, `3`, `4`.
4. Definir JSON report schema con `version`, `kind`, `summary`, `root`, `roadmap_root`, `diagnostics`.
5. Definir diagnostic ID convention con prefijo `RMC_`.

**Out**:
- Implementar comandos.
- Diseñar materialización o fix automático.
- Cambiar Rootline.

## Estado inicial esperado

- Repo Go nuevo o skeleton mínimo disponible.
- No existe contrato CLI formal versionado.

## Criterios de Aceptación

- Existe documentación del contrato CLI en `docs/cli-contract.md` o equivalente.
- El contrato especifica stdout/stderr para modo JSON y text.
- El contrato declara que `roadmapctl` es obligatorio para comandos `/roadmap` implementados que escriben, mutan o ejecutan.
- El contrato incluye al menos tres diagnostics: fallback `*-tasks.md`, missing `rootline`, y invalid `blocked_by`.

## Fuente de verdad

- `README.md`
- `docs/cli-contract.md`
- `docs/roadmap/O01-roadmapctl-mvp/README.md`
