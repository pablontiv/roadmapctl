---
estado: Completed
tipo: task
---
# T009: Poblar style fields en .roadmapctl.toml de este repo

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: el propio repo sirve de referencia de TOML bien configurado

[[blocked_by:./T001-config-gitflow-style-fields.md]]

## Preserva

- INV1: todos los campos existentes del TOML siguen presentes y con sus valores actuales
  - Verificar: `roadmapctl bootstrap --repo /home/shared/roadmapctl --output json` sale exit 0 tras el cambio

## Contexto

El TOML de este repo (`docs/roadmap/.roadmapctl.toml`) no tiene los 3 style fields nuevos. Hay que:
2. Agregar `branch_style`, `pr_title_style`, `pr_body_style` bajo `[gitflow]` con los valores reales del repo (inferidos de git log + commits existentes)

El estilo real de este repo: conventional commits (`feat/fix/docs/chore/refactor(OXX): ...`), branches directos a master, sin PRs (pr_mode=false).

## Alcance

**In**:
2. Agregar sección `[gitflow]` con `base_branch = "master"`, `branch_style`, `pr_title_style`, `pr_body_style` describiendo el estilo real del repo

**Out**:
- No modificar ningún otro campo del TOML
- No modificar otros archivos de configuración

## Estado inicial esperado

- T001 completada: el binario parsea los nuevos campos
- `docs/roadmap/.roadmapctl.toml` existe

## Criterios de Aceptación

- `grep 'branch_style\|pr_title_style\|pr_body_style' docs/roadmap/.roadmapctl.toml` → 3 matches con valores no-vacíos
- `roadmapctl bootstrap --repo /home/shared/roadmapctl --output json` sale exit 0 y no incluye `RMC_GITFLOW_NOT_CONFIGURED`

## Fuente de verdad

- `docs/roadmap/.roadmapctl.toml`
