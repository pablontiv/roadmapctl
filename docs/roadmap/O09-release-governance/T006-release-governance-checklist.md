---
estado: Completed
tipo: task
---
# T006: Checklist de release y gobernanza

**Outcome**: [O09 Release/governance](README.md)
**Contribuye a**: CE4

[[blocked_by:./T001-update-cli-contract-for-post-mvp-commands.md]]
[[blocked_by:./T002-define-rootline-compatibility-policy.md]]
[[blocked_by:./T003-add-ci-matrix-and-golden-stability.md]]
[[blocked_by:../O08-skill-cutover/T005-automate-pi-headless-verification-evidence.md]]

## Preserva

- INV1: Cambios de skill/guard no se liberan sin evidencia.
  - Verificar: checklist exige Pi headless.

## Contexto

La expansión post-MVP toca CLI, skill, docs y compatibilidad Rootline. Necesita checklist explícito antes de releases/cutovers.

## Alcance

**In**:
1. Checklist de tests: `go test`, build, goldens, fixture smoke.
2. Checklist de Rootline compatibility.
3. Checklist de skill sync y Pi headless evidence.
4. Checklist de docs/changelog/release notes.
5. Definir qué bloquea release.

**Out**:
- Automatizar todo en CI si no es viable todavía.

## Estado inicial esperado

- CI/goldens y Pi evidence procedure existen.

## Criterios de Aceptación

- Checklist está documentado y usado en release docs.
- Cada cutover de skill tiene evidencia requerida.
- No quedan pasos críticos solo en memoria de la sesión.

## Fuente de verdad

- `docs/release.md`
- `docs/roadmap-skill-integration.md`
- `scripts/sync-roadmap-skill.sh`
