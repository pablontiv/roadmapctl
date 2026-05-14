---
estado: Completed
tipo: task
---
# T008: Add LICENSE (PolyForm Noncommercial 1.0.0)

**Contribuye a**: establish the ecosystem-wide license on roadmapctl (currently no LICENSE file exists).

## Alcance

**In**:
- Create `/LICENSE` with PolyForm Noncommercial 1.0.0 full text, copyright "2026 Pablo Ontiveros"

**Out**:
- No changes to README (badges are added in T013)

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/LICENSE` passes
- LICENSE contains "PolyForm Noncommercial License 1.0.0"
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/LICENSE (new)
- /home/shared/rootline/LICENSE (reference for PolyForm NC text)
