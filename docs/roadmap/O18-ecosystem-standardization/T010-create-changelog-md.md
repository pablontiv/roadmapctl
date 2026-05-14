---
estado: Completed
tipo: task
---
# T010: Create CHANGELOG.md

**Contribuye a**: provide consumers and users a human-readable history of roadmapctl releases, starting from v0.0.1.

## Alcance

**In**:
- Create `/CHANGELOG.md` in the backscroll format ([Unreleased], then versioned entries)
- Seed with v0.0.1 entry (the first release, O17 completion) and key milestones from git log

**Out**:
- CHANGELOG is maintained manually on release

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/CHANGELOG.md` passes
- Contains `[Unreleased]` section and `v0.0.1` entry
- Format consistent with `/home/shared/backscroll/CHANGELOG.md`
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/CHANGELOG.md (new)
- /home/shared/backscroll/CHANGELOG.md (format reference)
