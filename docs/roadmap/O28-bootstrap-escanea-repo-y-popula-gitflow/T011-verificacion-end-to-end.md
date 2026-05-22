---
estado: Completed
tipo: task
---
# T011: Verificación end-to-end de O28

**Outcome**: [O28 Bootstrap escanea repo y popula gitflow](README.md)
**Contribuye a**: confirma que todos los cambios de O28 funcionan integrados

[[blocked_by:./T010-docs-cli-contract-changelog.md]]

## Preserva

- INV1: el repo sigue compilando y los tests pasan
  - Verificar: `go build ./...` y `go test ./... -count=1` salen exit 0

## Contexto

Task de verificación final. Confirma que los cambios de T001-T010 funcionan como un sistema integrado, no solo individualmente.

## Alcance

**In**:
Ejecutar y confirmar cada verificación:

1. `roadmapctl bootstrap --output json` incluye `branch_style`, `pr_title_style`, `pr_body_style`
2. En repo sin TOML: `roadmapctl bootstrap` devuelve `RMC_GITFLOW_NOT_CONFIGURED`
4. **Fix Causa B**: `grep -F 'omitir esta fase; commitear en el branch actual' .claude/skills/integrate/SKILL.md` → vacío
5. **Fix Causa A**: `grep -i 'única puerta\|único gate' .claude/skills/roadmap/loop-subcommand.md` → match en paso 9
6. `integrate/SKILL.md` no contiene `Co-authored-by`, `Signed-off-by`, `--squash`, `--delete-branch`
7. Config snapshot en loop paso 9 tiene `branch_style`
8. `bootstrap-reference.md` no contiene lista rígida de comandos fijos para el escaneo
9. `bootstrap-subcommand.md` existe
10. SKILL.md routing tiene fila para `bootstrap`
11. plan-subcommand.md, pending-subcommand.md, decision-tree-subcommand.md referencian `RMC_GITFLOW_NOT_CONFIGURED`
12. `docs/roadmap/.roadmapctl.toml` tiene los 3 style fields

**Out**:
- No modificar ningún archivo de implementación

## Estado inicial esperado

- T010 completada (todas las tasks anteriores completadas)

## Criterios de Aceptación

- Las 12 verificaciones del Alcance pasan sin errores
- `go test ./... -count=1` sale exit 0
- `golangci-lint run ./...` sale exit 0
- `roadmapctl check --repo /home/shared/roadmapctl --strict` sale exit 0

## Fuente de verdad

- Todos los archivos modificados en O28 (T001-T010)
