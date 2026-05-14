---
estado: Pending
tipo: task
---
# T006: Add codeql.yml workflow

**Contribuye a**: add static security analysis (CodeQL) to roadmapctl CI — the only repo in the ecosystem currently missing it.

## Alcance

**In**:
- Create `/.github/workflows/codeql.yml` using `pablontiv/crossbeam/.github/workflows/codeql.yml@v1` with `language: go`
- Schedule: nightly cron (e.g. `0 3 * * *`) as done in backscroll and rootline

**Out**:
- No changes to existing ci.yml

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/.github/workflows/codeql.yml` passes
- File calls crossbeam codeql.yml@v1 with `language: go`
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/.github/workflows/codeql.yml (new)
- /home/shared/rootline/.github/workflows/codeql.yml (reference)
