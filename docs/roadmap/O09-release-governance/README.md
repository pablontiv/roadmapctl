---
estado: Pending
tipo: outcome
---
# O09: Release, CI, distribución y compatibilidad

## Objetivo

`roadmapctl` expandido se distribuye y valida de forma confiable: contrato CLI post-MVP, compatibilidad Rootline, CI/goldens, releases, integración CI y checklist de gobernanza.

## Criterios de Éxito

- CE1: El contrato CLI documenta comandos post-MVP y JSON/diagnostics.
  - Verificar: docs actualizadas.
- CE2: La compatibilidad Rootline está definida y testeada.
  - Verificar: CI/matrix o diagnostics.
- CE3: Hay estrategia de release/install con checksums.
  - Verificar: GoReleaser/install docs.
- CE4: Cada release/cutover tiene checklist de evidencia.
  - Verificar: release governance checklist.

## Invariantes

- INV1: Rootline sigue siendo dependencia externa, no import interno.
  - Verificar: `go list`/grep imports.
- INV2: CI no oculta cambios de goldens.
  - Verificar: docs/golden-tests.

## Alcance

**In**:
- CLI contract post-MVP.
- Compatibilidad Rootline.
- CI/goldens/release/installers.
- Action/output opcional.
- Checklist de gobernanza.

**Out**:
- Cambiar Rootline release.
- Publicar package managers sin aprobación.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-update-cli-contract-for-post-mvp-commands.md) | Actualizar contrato CLI post-MVP |
| [T002](T002-define-rootline-compatibility-policy.md) | Definir política compatibilidad Rootline |
| [T003](T003-add-ci-matrix-and-golden-stability.md) | Fortalecer CI matrix/goldens |
| [T004](T004-add-goreleaser-and-installers.md) | Añadir GoReleaser/installers/checksums |
| [T005](T005-add-github-action-or-ci-output.md) | Agregar integración CI/GitHub Action opcional |
| [T006](T006-release-governance-checklist.md) | Checklist release/governance |
