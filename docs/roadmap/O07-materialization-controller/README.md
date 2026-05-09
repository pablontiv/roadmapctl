---
tipo: outcome
---
# O07: Controlador de materialización

## Objetivo

`roadmapctl materialize` convierte un plan estructurado aprobado en archivos canónicos de roadmap con dry-run/apply, numbering cross-platform, path containment, updates de README y postcheck obligatorio.

## Criterios de Éxito

- CE1: `materialize --dry-run` muestra archivos/cambios sin escribir.
  - Verificar: golden dry-run.
- CE2: `materialize --apply` crea Outcomes/Tasks canónicos y actualiza tablas.
  - Verificar: temp fixture + postcheck.
- CE3: nunca se genera un único `*-tasks.md` como fallback.
  - Verificar: tests anti-regresión.

## Invariantes

- INV1: Materializar no implementa código.
  - Verificar: solo archivos roadmap `.md`, `.stem`, `.roadmapctl.toml` cuando corresponda.
- INV2: Writes requieren apply explícito, preflight y postcheck.
  - Verificar: tests.

## Alcance

**In**:
- Plan input estructurado.
- Numbering/path planning.
- Dry-run/diff.
- Apply/postcheck.
- Bootstrap materialization aprobado.
- Cutover plan-subcommand.

**Out**:
- Auto-fix general.
- Commit/push automático.
- Decomposición AI dentro de roadmapctl.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-define-structured-materialize-plan-schema.md) | Definir schema de plan materializable |
| [T002](T002-implement-cross-platform-numbering-path-planning.md) | Implementar numbering/path planning cross-platform |
| [T003](T003-implement-materialize-dry-run-and-diff.md) | Implementar dry-run y diff |
| [T004](T004-implement-materialize-apply-with-postcheck.md) | Implementar apply con postcheck |
| [T005](T005-support-bootstrap-materialization.md) | Soportar bootstrap materialization aprobado |
| [T006](T006-add-materialization-fixtures-and-goldens.md) | Agregar fixtures/goldens materialization |
| [T007](T007-cutover-plan-materialization-in-skill.md) | Cortar plan-subcommand hacia materialize |
