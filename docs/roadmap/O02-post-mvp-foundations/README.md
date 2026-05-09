---
estado: Pending
tipo: outcome
---
# O02: Fundaciones post-MVP y frontera de schema

## Objetivo

`roadmapctl` tiene una base técnica sólida para crecer más allá del MVP sin invadir responsabilidades de Rootline: `.stem` vía Rootline define el schema documental, la config local define roles operacionales y el CLI usa librerías comunitarias para infraestructura genérica.

## Criterios de Éxito

- CE1: `estado: On Hold` es válido cuando el `.stem` efectivo lo permite.
  - Verificar: `go test ./...` y fixture `valid-status-on-hold`.
- CE2: roles configurados inexistentes en schema producen diagnóstico de config, no rechazo de documentos válidos.
  - Verificar: fixture `invalid-config-role-not-in-schema`.
- CE3: los nuevos comandos pueden construirse sobre Cobra, TOML, Rootline JSON y modelos internos testeados.
  - Verificar: tests de `internal/cli`, `internal/config`, `internal/rootlinecli`, `internal/roadmap`.

## Invariantes

- INV1: Rootline + `.stem` es la autoridad de schema documental.
  - Verificar: no hay enums documentales hardcodeados como fuente primaria en `roadmapctl`.
- INV2: `roadmapctl` conserva subprocess seguro: argumentos explícitos, timeout, stdout/stderr separados.
  - Verificar: tests de `internal/rootlinecli`.
- INV3: JSON stdout sigue siendo un único objeto parseable en modo JSON.
  - Verificar: golden tests.

## Alcance

**In**:
- Migración de infraestructura CLI/config hacia paquetes comunitarios.
- Fix de autoridad schema/config y bug `On Hold`.
- Hardening del wrapper Rootline.
- Modelo read-domain y fixtures fundacionales.

**Out**:
- Implementar comandos read-only post-MVP.
- Materializar o mutar roadmaps.
- Cambiar Rootline.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-adopt-cobra-and-community-packages.md) | Adoptar Cobra y paquetes comunitarios para CLI/config/Markdown/diff |
| [T002](T002-make-stem-authoritative-for-document-schema.md) | Usar `.stem` vía Rootline como autoridad de schema documental |
| [T003](T003-validate-operational-status-roles-separately.md) | Validar roles config separados del enum documental |
| [T004](T004-harden-rootline-json-on-nonzero-exit.md) | Parsear JSON Rootline incluso con exit non-zero |
| [T005](T005-add-roadmap-domain-model-and-tree-wrapper.md) | Agregar modelo read-domain y wrapper `rootline tree` |
| [T006](T006-add-raw-blocked-by-validation.md) | Complementar validación raw de `blocked_by` |
| [T007](T007-expand-foundation-fixtures-and-goldens.md) | Agregar fixtures/goldens fundacionales |
