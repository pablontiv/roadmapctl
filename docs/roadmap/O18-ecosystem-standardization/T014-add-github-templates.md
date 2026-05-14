---
estado: Pending
tipo: task
---
# T014: Add GitHub PR template and issue templates

**Contribuye a**: provide structured contribution workflows for roadmapctl — matching rootline and backscroll which have PR template + bug_report + feature_request issue templates.

## Alcance

**In**:
- Create `/.github/PULL_REQUEST_TEMPLATE.md`: What/Why/How structure + checklist (go test, golangci-lint, smoke tests, docs updated)
- Create `/.github/ISSUE_TEMPLATE/bug_report.md`: structured bug report adapted to roadmapctl (CLI tool, rootline dependency)
- Create `/.github/ISSUE_TEMPLATE/feature_request.md`: structured feature request

**Out**:
- No changes to existing workflows

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/.github/PULL_REQUEST_TEMPLATE.md` passes
- `test -f /home/shared/roadmapctl/.github/ISSUE_TEMPLATE/bug_report.md` passes
- `test -f /home/shared/roadmapctl/.github/ISSUE_TEMPLATE/feature_request.md` passes
- Style consistent with backscroll templates
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/.github/PULL_REQUEST_TEMPLATE.md (new)
- /home/shared/roadmapctl/.github/ISSUE_TEMPLATE/ (new)
- /home/shared/backscroll/.github/PULL_REQUEST_TEMPLATE.md (style reference)
- /home/shared/backscroll/.github/ISSUE_TEMPLATE/ (style reference)
