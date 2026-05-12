---
tipo: outcome
---
# Refactorizar materialización y responsabilidades roadmapctl/Rootline/skill

Separar generación semántica, escritura, path planning, validación y vistas calculadas para evitar que roadmapctl, Rootline y el skill dupliquen responsabilidades.

## Criterios de Aceptación

- El flujo /roadmap plan conserva una aprobación humana antes de escribir archivos.
- Outcome README no persiste listas de Tasks ni criterios de aceptación obligatorios.
- roadmapctl no renderiza contenido semántico de Tasks y Rootline conserva las responsabilidades genéricas de Markdown/.stem.

## Tasks

| Task | Descripción |
|------|-------------|
| [T001](T001-record-responsibility-separation-decision.md) | Reemplazar la decisión previa de materialize como writer semántico por la separación acordada entre skill, Pi write, roadmapctl y Rootline. |
| [T002](T002-redefine-roadmap-plan-skill-contract.md) | Actualizar el skill para proponer semántica, pedir aprobación y materializar con Pi write usando templates después de validar paths. |
| [T003](T003-design-roadmapctl-path-planning-guard.md) | Reemplazar el rol writer de materialize por una superficie determinística que calcule y valide paths canónicos sin renderizar Markdown semántico. |
| [T004](T004-remove-persisted-outcome-task-table.md) | Hacer que la lista de Tasks de un Outcome sea una vista calculada desde filesystem + Rootline en lugar de una tabla en README. |
| [T005](T005-move-acceptance-criteria-to-tasks-only.md) | Quitar AC obligatorios de Outcomes y asegurar que cada Task sea la unidad verificable con sus propios criterios. |
| [T006](T006-verify-rootline-generic-record-operations.md) | Validar o ajustar Rootline para soportar el flujo genérico de .stem, new, set, secciones, validate, query, tree y graph sin lógica roadmap-specific. |
| [T007](T007-migrate-docs-fixtures-and-goldens-to-new-boundary.md) | Actualizar contrato, integración, fixtures y golden tests para reflejar Pi write, path planning, views calculadas y Rootline como DB genérica. |
| [T008](T008-add-headless-approval-and-parallel-write-verification.md) | Probar el flujo del skill desde propuesta hasta materialización aprobada, incluyendo aprobación obligatoria, escritura paralela y validación final. |
