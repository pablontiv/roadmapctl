---
estado: Completed
tipo: task
---
# T016: Remove Status callout from README

**Contribuye a**: README parity across the ecosystem — the `> **Status**: ...` callout is being removed from all 4 repos as it doesn't fit the final style.

## Alcance

**In**:
- Remove `> **Status**: Core command families functional...` line from README.md
- Remove adjacent blank line to avoid double spacing

**Out**:
- No other changes

## Criterios de Aceptación

- `grep "> \*\*Status\*\*" /home/shared/roadmapctl/README.md` returns empty
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/README.md
