---
estado: Completed
tipo: task
---
# T004: Añadir GoReleaser, installers y checksums

**Outcome**: [O09 Release/governance](README.md)
**Contribuye a**: CE3

[[blocked_by:./T003-add-ci-matrix-and-golden-stability.md]]

## Preserva

- INV1: Release de roadmapctl es separado de Rootline.
  - Verificar: repo/module y artifacts.

## Contexto

`docs/release.md` hoy describe outline de release, pero no hay GoReleaser/package-manager flow completo.

## Alcance

**In**:
1. Agregar `.goreleaser.yml` para Linux/macOS/Windows amd64/arm64.
2. Generar checksums.
3. Documentar install script o `go install` recomendado.
4. Asegurar que artifacts no incluyen Rootline.
5. Actualizar release docs.

**Out**:
- Publicar package managers sin aprobación adicional.
- Cambiar Rootline installers.

## Estado inicial esperado

- CI build existe; release outline existe.

## Criterios de Aceptación

- Dry-run de GoReleaser pasa si se usa.
- Docs explican dependencia Rootline.
- Checksums se generan.

## Fuente de verdad

- `docs/release.md`
- `.github/workflows/ci.yml`
- `go.mod`
