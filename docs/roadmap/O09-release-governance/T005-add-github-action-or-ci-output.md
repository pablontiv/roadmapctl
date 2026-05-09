---
estado: Pending
tipo: task
---
# T005: Agregar GitHub Action o output CI opcional

**Outcome**: [O09 Release/governance](README.md)
**Contribuye a**: CE3

[[blocked_by:../O05-deterministic-lint/T006-add-lint-fixtures-and-goldens.md]]
[[blocked_by:./T001-update-cli-contract-for-post-mvp-commands.md]]

## Preserva

- INV1: JSON/text siguen siendo outputs principales hasta aprobar formatos extra.
  - Verificar: docs.

## Contexto

Proyectos podrían querer ejecutar `roadmapctl check/lint` en CI y ver anotaciones. Esta integración debe basarse en JSON estable.

## Alcance

**In**:
1. Evaluar GitHub Action composite o docs CI.
2. Evaluar formatos adicionales como SARIF/JUnit solo si aporta valor.
3. Documentar uso recomendado en repos con Rootline instalado.
4. Mantener JSON como fuente de verdad.

**Out**:
- Cambiar comandos core para depender de GitHub.
- Requerir CI para uso local.

## Estado inicial esperado

- `check` y `lint` tienen JSON estable.

## Criterios de Aceptación

- Hay ejemplo CI funcional.
- Action/output no rompe uso local.
- Tests/docs explican exit codes.

## Fuente de verdad

- `docs/cli-contract.md`
- `docs/release.md`
- `.github/workflows/ci.yml`
