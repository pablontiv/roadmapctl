---
estado: Completed
tipo: task
---
# T018: Merge del bump setup-go 6.3→6.4 (PR #1)

**Contribuye a**: Mantener actions de CI actualizadas

## Contexto

PR #1: `actions/setup-go from 6.3.0 to 6.4.0`. CI verde (+1/-1 línea en workflow).

## Alcance

**In**:
1. `gh pr merge 1 --repo pablontiv/roadmapctl --merge`
2. `git -C /home/shared/roadmapctl pull --rebase origin master`

**Out**:
- No otros cambios

## Estado inicial esperado

- PR #1 abierto con CI verde

## Criterios de Aceptación

- `gh pr list --repo pablontiv/roadmapctl --state open` retorna 0 PRs

## Fuente de verdad

- `gh pr list --repo pablontiv/roadmapctl --state open`
