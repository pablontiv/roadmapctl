---
estado: Completed
tipo: task
---
# T011: Añadir CI inicial y outline de release cross-platform

**Outcome**: [O01 roadmapctl MVP obligatorio para roadmaps](README.md)
**Contribuye a**: CE1 y CE2

[[blocked_by:./T002-create-go-cli-skeleton.md]]
[[blocked_by:./T009-add-fixtures-and-golden-tests.md]]

## Preserva

- INV1: Rootline permanece como DBMS/constraint engine genérico.
  - Verificar: release de `roadmapctl` no modifica release de `rootline`.

## Contexto

`roadmapctl` debe correr en Linux, macOS y Windows. El MVP necesita CI que pruebe build/test cross-platform y un outline de release aunque la publicación completa pueda quedar para una fase posterior.

## Alcance

**In**:
1. Añadir workflow CI con matriz Linux, macOS y Windows.
2. Ejecutar `go test ./...` y `go build ./cmd/roadmapctl`.
3. Definir plan de release con GoReleaser para linux/darwin/windows amd64/arm64.
4. Documentar instalación inicial vía `go install`.
5. Documentar compatibilidad mínima esperada con Rootline.

**Out**:
- Publicar release real si no está aprobado.
- Homebrew/Scoop/Winget.
- Signing/SBOM obligatorio en MVP.

## Estado inicial esperado

- Tests y fixtures existen.
- Skeleton Go compila localmente.

## Criterios de Aceptación

- CI corre en al menos Linux, macOS y Windows.
- `go test ./...` se ejecuta en la matriz.
- `go build ./cmd/roadmapctl` se ejecuta en la matriz.
- Docs indican cómo instalar con `go install`.
- Release outline menciona checksums y matriz futura.

## Fuente de verdad

- `.github/workflows/ci.yml`
- `.goreleaser.yml` o `docs/release.md`
- `README.md`
