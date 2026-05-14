---
estado: Pending
tipo: task
---
# T007: Add scorecard.yml workflow

**Contribuye a**: add OpenSSF Scorecard to roadmapctl — once the repo is public (T001), Scorecard works and provides supply-chain security signals.

## Alcance

**In**:
- Create `/.github/workflows/scorecard.yml` using `pablontiv/crossbeam/.github/workflows/scorecard.yml@v1`
- Schedule: nightly cron (e.g. `0 4 * * *`)
- Required permissions: `security-events: write`, `id-token: write`, `contents: read`

**Out**:
- No changes to existing ci.yml
- Note: depends on T001 (repo must be public for Scorecard to work)

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/.github/workflows/scorecard.yml` passes
- File calls crossbeam scorecard.yml@v1
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/.github/workflows/scorecard.yml (new)
- /home/shared/rootline/.github/workflows/scorecard.yml (reference)
