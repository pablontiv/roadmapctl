---
estado: Pending
tipo: task
---
# T004: Add CODEOWNERS

**Contribuye a**: ensure all PRs require review from @pablontiv, the sole maintainer — matching the pattern used in crossbeam, rootline, and backscroll.

## Alcance

**In**:
- Create `/.github/CODEOWNERS` with content `* @pablontiv`

**Out**:
- No other .github changes (PR/issue templates are a separate task)

## Criterios de Aceptación

- `cat /home/shared/roadmapctl/.github/CODEOWNERS` returns `* @pablontiv`
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/.github/CODEOWNERS (new)
- /home/shared/backscroll/.github/CODEOWNERS (reference)
