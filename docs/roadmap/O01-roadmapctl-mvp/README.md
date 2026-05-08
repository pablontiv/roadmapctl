---
estado: Pending
tipo: outcome
---
# O01: roadmapctl MVP obligatorio para roadmaps

## Objetivo

Existe un CLI cross-platform `roadmapctl` que valida roadmaps gobernados por Rootline y actúa como guard obligatorio para los comandos `/roadmap` que materializan, mutan o ejecutan trabajo.

## Criterios de Éxito

- CE1: `roadmapctl doctor` diagnostica entorno, repo, Rootline y configuración de roadmap.
  - Verificar: `go test ./...` y `go run ./cmd/roadmapctl doctor --repo testdata/fixtures/valid-outcome-with-tasks --output json`
- CE2: `roadmapctl check` detecta estructuras inválidas, incluyendo el fallback de un único `*-tasks.md`.
  - Verificar: `go run ./cmd/roadmapctl check --repo testdata/fixtures/invalid-single-summary-file --output json` sale con código `1`.
- CE3: `/roadmap` documenta que `roadmapctl` es obligatorio para comandos implementados que escriben, mutan o ejecutan.
  - Verificar: revisar docs/skill integration y fixtures de validación.
- CE4: El repo `roadmapctl` contiene la fuente canónica del skill `roadmap` y un git hook lo instala en el user scope.
  - Verificar: `.claude/skills/roadmap/` existe en el repo, `.githooks/post-merge` ejecuta `scripts/sync-roadmap-skill.sh --install`, y `scripts/sync-roadmap-skill.sh --check` confirma que la copia instalada coincide con la fuente.

## Invariantes

- INV1: Rootline permanece como DBMS/constraint engine genérico; no se agregan subcomandos roadmap a `rootline`.
  - Verificar: no hay cambios en `cmd/rootline` ni imports de `github.com/pablontiv/rootline/internal/*`.
- INV2: `roadmapctl` no invoca subprocess con shell strings; usa argumentos explícitos y timeouts.
  - Verificar: tests de `internal/rootlinecli` y revisión de `exec.CommandContext`.
- INV3: El MVP no materializa ni corrige automáticamente; solo diagnostica y valida.
  - Verificar: comandos disponibles limitados a `doctor` y `check`.

## Alcance

**In**:
- CLI Go separado con `doctor` y `check`.
- Salida human-readable y JSON estable.
- Fixtures de roadmaps válidos e inválidos.
- Integración documentada para que `/roadmap` bloquee si `roadmapctl` falla.
- Skill `roadmap` versionado en este repo e instalado al user scope por hook.

**Out**:
- Materialización automática de planes.
- Fix automático de roadmaps.
- Subcomandos roadmap dentro de `rootline`.
- Publicación a package managers en el MVP.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-define-cli-contract.md) | Definir contrato CLI, flags, JSON y exit codes |
| [T002](T002-create-go-cli-skeleton.md) | Crear esqueleto Go cross-platform para `roadmapctl` |
| [T003](T003-implement-diagnostics-model.md) | Implementar modelo de diagnostics y renderers |
| [T004](T004-load-roadmap-config.md) | Cargar `.claude/roadmap.local.md` y resolver roots |
| [T005](T005-wrap-rootline-cli.md) | Encapsular llamadas seguras al binario `rootline` |
| [T006](T006-implement-doctor-command.md) | Implementar `roadmapctl doctor` |
| [T007](T007-implement-structure-checks.md) | Implementar checks estructurales canónicos |
| [T008](T008-implement-rootline-backed-checks.md) | Integrar validaciones basadas en Rootline JSON |
| [T009](T009-add-fixtures-and-golden-tests.md) | Crear fixtures y golden tests cross-platform |
| [T010](T010-document-roadmap-skill-integration.md) | Documentar integración obligatoria con `/roadmap` |
| [T011](T011-add-ci-and-release-outline.md) | Añadir CI inicial y outline de release cross-platform |
| [T012](T012-version-roadmap-skill-and-install-hook.md) | Versionar skill `roadmap` en el repo y sincronizarlo al user scope con hook |
