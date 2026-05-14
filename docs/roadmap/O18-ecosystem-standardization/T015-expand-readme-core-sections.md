---
estado: Completed
tipo: task
---
# T015: Expand README — Quick Start, Core Idea, AI-Native

**Contribuye a**: align roadmapctl README with rootline narrative structure — T013 added badges and expanded the layer table but omitted the three core narrative sections present in every ecosystem README.

## Alcance

**In**:
- Add Quick Start section (6 numbered+commented commands showing the full workflow)
- Add Core Idea section (mental model: Markdown + .stem + estado semantics + deterministic guards)
- Add AI-Native section (designed for agent invocation, stable JSON contracts, bootstrap API)
- Expand Commands section from basic table to full CLI reference listing all command families

**Out**:
- No changes to Layer Responsibilities, Installation, Skill Source, or Development sections

## Criterios de Aceptación

- `grep "## Quick Start" /home/shared/roadmapctl/README.md` exits 0
- `grep "## Core Idea" /home/shared/roadmapctl/README.md` exits 0
- `grep "## AI-Native" /home/shared/roadmapctl/README.md` exits 0
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/README.md
- /home/shared/rootline/README.md (section structure reference)
