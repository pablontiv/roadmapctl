---
estado: Specified
tipo: task
---
# T007: Consolidar detección RMC_GITFLOW_NOT_CONFIGURED en todos los subcomandos

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: cualquier subcomando de /roadmap ejecuta el wizard si faltan los style fields

[[blocked_by:./T005-bootstrap-skill-open-scan-wizard.md]]

## Preserva

- INV1: cada subcomando sigue ejecutando bootstrap en su Fase 1
  - Verificar: `grep -l 'bootstrap-reference' .claude/skills/roadmap/*.md` incluye loop, plan, pending, decision-tree

## Contexto

El wizard de adopción (escaneo + pre-llenado + escritura TOML) vive en `bootstrap-reference.md` (T005). Pero los subcomandos existentes (`plan`, `pending`, `decision-tree`) necesitan referenciar ese archivo y reaccionar al diagnostic `RMC_GITFLOW_NOT_CONFIGURED` en su Fase 1.

`loop-subcommand.md` ya referencia `bootstrap-reference.md`; verificar que el detector está ahí también.

## Alcance

**In**:
1. Verificar y actualizar `plan-subcommand.md`: Fase 1 referencia `bootstrap-reference.md` + detector `RMC_GITFLOW_NOT_CONFIGURED` → wizard → re-bootstrap → continúa
2. Verificar y actualizar `pending-subcommand.md`: ídem
3. Verificar y actualizar `decision-tree-subcommand.md`: ídem
4. Verificar `loop-subcommand.md`: asegurar que el detector también está en Fase 1 (no solo en paso 9)

**Out**:
- No modificar el contenido específico de cada subcomando más allá de la Fase 1

## Estado inicial esperado

- T005 completada: `bootstrap-reference.md` tiene la sección del wizard

## Criterios de Aceptación

- `grep -l 'RMC_GITFLOW_NOT_CONFIGURED\|bootstrap-reference' .claude/skills/roadmap/plan-subcommand.md .claude/skills/roadmap/pending-subcommand.md .claude/skills/roadmap/decision-tree-subcommand.md` retorna los 3 archivos
- Cada uno de los 3 describe la acción ante `RMC_GITFLOW_NOT_CONFIGURED`: wizard → re-bootstrap → continuar

## Fuente de verdad

- `.claude/skills/roadmap/plan-subcommand.md`
- `.claude/skills/roadmap/pending-subcommand.md`
- `.claude/skills/roadmap/decision-tree-subcommand.md`
- `.claude/skills/roadmap/loop-subcommand.md`
