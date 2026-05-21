---
estado: Specified
tipo: task
---
# T008: pending-subcommand.md — footer con gitflow resuelto

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: el usuario ve los settings de gitflow sin tener que abrir el TOML

[[blocked_by:./T002-bootstrap-serialize-and-diagnostic.md]]

## Preserva

- INV1: el output principal de pending (tabla de tasks) no cambia
  - Verificar: pending sigue mostrando la tabla de tasks antes del footer

## Contexto

Hoy `/roadmap pending` muestra el estado de las tasks pero no hay visibilidad del gitflow configurado. El usuario descubre cómo está configurado el branch naming / PR mode / commit style solo cuando falla algo.

El footer debe mostrar los valores resueltos del TOML (no propuestos, sino los que se van a usar) para que el usuario pueda validar antes de que el loop empiece a commitear.

## Alcance

**In**:
1. Agregar al final del output de `pending-subcommand.md` un bloque footer:
```
── Gitflow del proyecto ──────────────
Base branch:    <base_branch>
Modo:           <descripción de pr_mode + auto_push>
Autonomía:      <autonomy>
commit_style:   <commit_style>

branch_style:   "<valor del TOML>"
pr_title_style: "<valor o 'no aplica — pr_mode=false'>"
pr_body_style:  "<valor o 'no aplica — pr_mode=false'>"

Para refrescar:  /roadmap bootstrap
```
2. Si `RMC_GITFLOW_NOT_CONFIGURED` está en los diagnostics de bootstrap: mostrar aviso en el footer ("⚠ gitflow no configurado — ejecutar /roadmap bootstrap")

**Out**:
- No modificar la lógica de tabla de tasks

## Estado inicial esperado

- T002 completada: bootstrap JSON incluye los style fields

## Criterios de Aceptación

- `pending-subcommand.md` contiene sección de footer con los campos mencionados
- El footer incluye manejo del caso `RMC_GITFLOW_NOT_CONFIGURED`
- El footer aparece después de la tabla de tasks, no antes

## Fuente de verdad

- `.claude/skills/roadmap/pending-subcommand.md`
