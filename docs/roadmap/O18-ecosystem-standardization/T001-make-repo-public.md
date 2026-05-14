---
estado: Pending
tipo: task
---
# T001: Make repo public

**Contribuye a**: expose roadmapctl as a public open-source tool — it is the companion CLI for the public rootline ecosystem and must be publicly accessible.

## Alcance

**In**:
- Set repo visibility to PUBLIC via GitHub API (`private: false`)

**Out**:
- No file changes

## Criterios de Aceptación

- `gh repo view pablontiv/roadmapctl --json isPrivate` returns `{"isPrivate":false}`

## Fuente de verdad

- GitHub API: PATCH repos/pablontiv/roadmapctl (`private: false`)
