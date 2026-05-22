---
estado: Specified
tipo: task
---
# T002: Eliminar `pr_merge_strategy` de skills, docs e históricos

**Outcome**: [O31 Fix `/integrate` defaults y limpieza de config muerta](README.md)
**Contribuye a**: documentación coherente con la API real; sin instrucciones que apuntan a config muerta

## Preserva

- INV1: bootstrap sigue limpio (sin diagnostics nuevos)
  - Verificar: `roadmapctl bootstrap --repo /home/shared/roadmapctl --output json` summary.status="ok", warnings=0, errors=0
- INV2: la documentación editada sigue válida (sin enlaces internos rotos ni renderizado roto)
  - Verificar: revisión manual de los archivos editados

## Contexto

El usuario pidió "eliminar todo, sin excepciones": se tocan skills, docs de contrato, CHANGELOG y task records históricos del roadmap (`docs/roadmap/O*/T*.md`), y specs (`docs/superpowers/specs/`).

Quedan intactos los archivos auto-generables: `context-gather/`, `graphify-out/`, `cartyx-out/`.

## Alcance

**In**:

Skills:
1. `.claude/skills/roadmap/loop-subcommand.md`: quitar `pr_merge_strategy` de las dos listas de campos config.
2. `.claude/skills/roadmap/bootstrap-reference.md`: quitar de la lista de campos operacionales; borrar fila del diagnostic `RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED`; borrar `pr-merge-strategy: 'squash'` del ejemplo YAML.
3. `.claude/skills/roadmap/config-reference.md`: borrar fila `pr-merge-strategy` de la tabla "Config key".

Docs de contrato:
4. `docs/cli-contract.md`: quitar el campo del ejemplo TOML, la fila completa de la tabla TOML schema, `pr_merge_strategy` de la columna Outputs del `context` capability, y la fila del diagnostic deprecado.
5. `docs/roadmap-skill-integration.md`: quitar de la línea que describe el JSON config snapshot.

CHANGELOG e históricos:
6. `CHANGELOG.md`: borrar la entry "pr_merge_strategy TOML field: now emits..." y agregar nueva entry: `Removed pr_merge_strategy TOML field and RMC_GITFLOW_PR_MERGE_STRATEGY_DEPRECATED diagnostic. Merge strategy is now controlled exclusively by GitHub branch protection rules. No fallback or migration path — consumers must delete the field from their .roadmapctl.toml if present.`
7. `docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md`: quitar las 2 menciones.
8. `docs/roadmap/O03-config-context-workspace/T001-design-roadmapctl-toml-config-contract.md`: quitar mención.
9. `docs/roadmap/O10-repo-local-execution-settings/T002-expose-execution-settings-in-context.md`: quitar 2 menciones.
10. `docs/roadmap/O10-repo-local-execution-settings/T005-cutover-roadmap-skill-loop-to-config.md`: quitar mención.
11. `docs/roadmap/O24-extract-integrate-skill-and-shrink-roadmap-skill/README.md`: quitar mención.
12. `docs/roadmap/O24-extract-integrate-skill-and-shrink-roadmap-skill/T001-create-integrate-skill.md`: quitar 3 menciones.
13. `docs/roadmap/O24-extract-integrate-skill-and-shrink-roadmap-skill/T003-progressive-disclosure-of-roadmap-skill.md`: quitar AC del grep `pr-merge-strategy`.
14. `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T001-config-gitflow-style-fields.md`: quitar la lógica de deprecation warning.
15. `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T003-integrate-skill-fix-causa-b.md`: quitar mención.
16. `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T004-loop-integrate-gate-unico.md`: quitar 3 menciones.
17. `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T009-poblar-toml-este-repo.md`: quitar 4 menciones.
18. `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T010-docs-cli-contract-changelog.md`: quitar mención.
19. `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T011-verificacion-end-to-end.md`: quitar 3 menciones.

**Out**:
- Código Go y tests: T001.
- Skill `/integrate`: T003.
- Artefactos auto-generables (`context-gather/`, `graphify-out/`, `cartyx-out/`).

## Estado inicial esperado

- T001 puede estar pendiente o completada (independiente; orden no importa).
- `grep -rn 'pr_merge_strategy\|pr-merge-strategy' .claude/ docs/ CHANGELOG.md` retorna múltiples matches.

## Criterios de Aceptación

- `grep -rn 'pr_merge_strategy\|pr-merge-strategy\|PR_MERGE_STRATEGY_DEPRECATED' .claude/ docs/ CHANGELOG.md` retorna vacío.
- `CHANGELOG.md` contiene la nueva entry literal: `Removed pr_merge_strategy TOML field`.
- `roadmapctl bootstrap --repo /home/shared/roadmapctl --output json` summary.status="ok", warnings=0, errors=0, infos=0.

## Fuente de verdad

- `.claude/skills/roadmap/loop-subcommand.md`
- `.claude/skills/roadmap/bootstrap-reference.md`
- `.claude/skills/roadmap/config-reference.md`
- `docs/cli-contract.md`
- `docs/roadmap-skill-integration.md`
- `CHANGELOG.md`
- `docs/superpowers/specs/2026-05-09-roadmap-repo-settings-design.md`
- `docs/roadmap/O03-config-context-workspace/T001-design-roadmapctl-toml-config-contract.md`
- `docs/roadmap/O10-repo-local-execution-settings/T002-expose-execution-settings-in-context.md`
- `docs/roadmap/O10-repo-local-execution-settings/T005-cutover-roadmap-skill-loop-to-config.md`
- `docs/roadmap/O24-extract-integrate-skill-and-shrink-roadmap-skill/README.md`
- `docs/roadmap/O24-extract-integrate-skill-and-shrink-roadmap-skill/T001-create-integrate-skill.md`
- `docs/roadmap/O24-extract-integrate-skill-and-shrink-roadmap-skill/T003-progressive-disclosure-of-roadmap-skill.md`
- `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T001-config-gitflow-style-fields.md`
- `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T003-integrate-skill-fix-causa-b.md`
- `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T004-loop-integrate-gate-unico.md`
- `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T009-poblar-toml-este-repo.md`
- `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T010-docs-cli-contract-changelog.md`
- `docs/roadmap/O28-bootstrap-escanea-repo-y-popula-gitflow/T011-verificacion-end-to-end.md`
