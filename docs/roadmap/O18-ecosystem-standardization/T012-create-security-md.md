---
estado: Completed
tipo: task
---
# T012: Create SECURITY.md

**Contribuye a**: establish a security policy for roadmapctl — it mutates roadmap files and invokes the rootline binary, so path traversal and subprocess injection are relevant concerns.

## Alcance

**In**:
- Create `/SECURITY.md` with:
  - Reporting: do NOT open public issues; use GitHub private advisory or email
  - SLA: 48h ACK, 7d detailed response
  - Scope: roadmapctl is a CLI that validates/mutates Markdown roadmaps and invokes rootline subprocess; security concerns include path traversal via --repo/--roadmap-root, command injection via rootline subprocess invocation, TOML/frontmatter parsing edge cases

**Out**:
- No code changes

## Criterios de Aceptación

- `test -f /home/shared/roadmapctl/SECURITY.md` passes
- File describes scope and reporting procedure
- Style consistent with rootline/backscroll SECURITY.md
- `git -C /home/shared/roadmapctl log --oneline -1` shows a conventional commit

## Fuente de verdad

- /home/shared/roadmapctl/SECURITY.md (new)
- /home/shared/rootline/SECURITY.md (style reference)
