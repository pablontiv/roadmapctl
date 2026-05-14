---
estado: Completed
tipo: task
---
# T009: Create CONTRIBUTING.md

**Contribuye a**: give contributors guidance on how to work with roadmapctl — matching the completeness of rootline and backscroll CONTRIBUTING.md files.

## Alcance

**In**:
- Create `/CONTRIBUTING.md` with:
  - Dev setup (clone, install rootline binary as dependency, git hooks)
  - Build & test: `go build ./cmd/roadmapctl`, `go test ./...`, `golangci-lint run ./...`, `scripts/check-coverage.sh`
  - Cross-platform smoke tests: `roadmapctl check` and `roadmapctl lint` on testdata fixtures
  - Conventional commits convention (feat/fix/docs/chore/test/refactor)
  - Release: automated via CI (go-release.yml from crossbeam)
  - Quality gates: golangci-lint, go test, coverage script, smoke tests on 3 platforms
  - Reporting issues via SECURITY.md

**Out**:
- No changes to Justfile or scripts

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/CONTRIBUTING.md` passes
- No placeholder strings
- Style consistent with rootline CONTRIBUTING.md (formal, English, tables)
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/CONTRIBUTING.md (new)
- /home/shared/rootline/CONTRIBUTING.md (style reference)
- /home/shared/roadmapctl/.github/workflows/ci.yml (for actual quality gates)
