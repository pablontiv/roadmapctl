---
estado: Completed
tipo: task
---
# T006: Implementar workspace discovery

**Outcome**: [O03 Config/context/workspace](README.md)
**Contribuye a**: CE4

[[blocked_by:./T004-implement-context-command.md]]

## Preserva

- INV1: Si un workspace es ambiguo, roadmapctl falla con diagnostic en vez de adivinar.
  - Verificar: fixture de repo ambiguo.

## Contexto

El flag `--workspace` existe pero no está implementado. El skill describe workspace mode: detectar repos, leer roadmap roots y calcular helpers por repo.

## Alcance

**In**:
1. Implementar discovery de repos con `.git` y roadmap config/root.
2. Soportar selección por `--repo <name>` en workspace mode si se aprueba.
3. Agregar output agrupado por repo en `context`.
4. Definir diagnostics para config faltante, repo ambiguo, root escape y repo inválido.

**Out**:
- Cross-repo dependencies.
- Ejecución paralela de tasks.

## Estado inicial esperado

- `Options.Workspace` existe pero no se usa.

## Criterios de Aceptación

- Fixture workspace con dos repos produce contexto por repo.
- Repo ambiguo falla con diagnostic estable.
- Single-repo mode no cambia.

## Fuente de verdad

- `internal/cli/cli.go`
- `.claude/skills/roadmap/SKILL.md`
- `testdata/fixtures/*workspace*`
