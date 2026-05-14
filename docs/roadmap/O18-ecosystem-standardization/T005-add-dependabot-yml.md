---
estado: Completed
tipo: task
---
# T005: Add dependabot.yml

**Contribuye a**: keep Go module dependencies and GitHub Actions SHA pins up to date automatically — matching the pattern in rootline and backscroll.

## Alcance

**In**:
- Create `/.github/dependabot.yml` with two update configs:
  - `gomod`: weekly, open-pull-requests-limit 10, labels [dependencies, ci]
  - `github-actions`: weekly, open-pull-requests-limit 10, labels [dependencies, ci]

**Out**:
- No other .github changes

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/.github/dependabot.yml` passes
- File contains both `gomod` and `github-actions` update configs
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/.github/dependabot.yml (new)
- /home/shared/rootline/.github/dependabot.yml (reference)
