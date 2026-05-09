---
estado: Pending
tipo: outcome
---
# O06: Controlador de transiciones de estado

## Objetivo

`roadmapctl` gobierna preflight, dry-run y apply de transiciones de estado del roadmap usando roles operacionales y Rootline `set`, con postcheck obligatorio y sin ejecutar implementación de código.

## Criterios de Éxito

- CE1: `transition can-start` explica si una task puede iniciar según dependencias y status roles.
  - Verificar: fixtures ready/blocked.
- CE2: `transition start/complete --dry-run` produce cambios planeados sin escribir.
  - Verificar: golden JSON.
- CE3: `transition ... --apply` muta vía Rootline y corre postcheck.
  - Verificar: tests en temp fixtures.

## Invariantes

- INV1: No ejecutar código, ACs, commits ni PRs desde roadmapctl.
  - Verificar: command scope.
- INV2: Apply requiere intención explícita y postcheck.
  - Verificar: tests de dry-run/apply.

## Alcance

**In**:
- `can-start`, `can-complete`, `start`, `complete`, `set-status` gobernado.
- Wrappers Rootline `set`/validate-one.
- Fixtures/goldens y cutover loop/status.

**Out**:
- Implementar tasks de código.
- Commit/push automático.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-design-transition-status-role-model.md) | Diseñar modelo de roles/transiciones |
| [T002](T002-add-rootline-set-and-validateone-wrappers.md) | Agregar wrappers Rootline `set` y validate-one |
| [T003](T003-implement-can-start-can-complete.md) | Implementar can-start/can-complete |
| [T004](T004-implement-transition-dry-run.md) | Implementar dry-run de transición |
| [T005](T005-implement-transition-apply-with-postcheck.md) | Implementar apply con postcheck |
| [T006](T006-add-transition-fixtures-and-goldens.md) | Agregar fixtures/goldens transición |
| [T007](T007-cutover-loop-status-transitions-in-skill.md) | Cortar loop/status del skill hacia transition |
