---
estado: Completed
tipo: task
---
# T013: Expand README and add badges

**Contribuye a**: bring roadmapctl README to ecosystem standard — currently 99 lines with no badges; rootline/backscroll READMEs are 300+ lines with CI, language, and license badges.

## Estado inicial esperado

- README.md: 99 lines, no badges, minimal command list, no Installation section with install scripts, no Development section

## Alcance

**In**:
- Add badges: CI status, Go version, PolyForm NC license
- Expand to include all standard sections: Installation (with install.sh/install.ps1 snippets), Quick Start, Layer Responsibilities (existing table is good, keep it), Commands (brief summary of all 9+ commands with links to docs/), Documentation (table linking to docs/cli-contract.md etc), Development (go build, go test, golangci-lint), License
- Tone and style: formal, technical, English, consistent with rootline README

**Out**:
- No changes to docs/ subdirectory files (they already exist and are good)
- No changes to install.sh / install.ps1

## Criterios de Aceptación

- README.md has CI, Go, and License badges
- README.md has Installation, Quick Start, Commands, Documentation, Development, License sections
- `wc -l /home/shared/roadmapctl/README.md` shows >= 250 lines
- Style consistent with rootline README (tables, code blocks, formal tone)
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/README.md
- /home/shared/rootline/README.md (style reference)
