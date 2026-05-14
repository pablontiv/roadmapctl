---
estado: Pending
tipo: task
---
# T011: Create CODE_OF_CONDUCT.md

**Contribuye a**: establish the ecosystem code of conduct for roadmapctl contributors — matching rootline and backscroll.

## Alcance

**In**:
- Create `/CODE_OF_CONDUCT.md` adopting Contributor Covenant v2.1
- Copy from `/home/shared/backscroll/CODE_OF_CONDUCT.md` (same text, no repo-specific placeholders)

**Out**:
- No other changes

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/CODE_OF_CONDUCT.md` passes
- File adopts Contributor Covenant v2.1
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/CODE_OF_CONDUCT.md (new)
- /home/shared/backscroll/CODE_OF_CONDUCT.md (source)
