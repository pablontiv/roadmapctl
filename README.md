# roadmapctl

Companion CLI for Rootline-governed roadmaps.

`roadmapctl` owns roadmap-specific guardrails and workflows while `rootline` remains the generic file-based database and constraint engine.

## Skill source

This repository is the canonical home for the `/roadmap` skill:

```text
.claude/skills/roadmap/
```

The git hook in `.githooks/post-merge` installs that skill into the user scope:

```text
~/.claude/skills/roadmap
```

## MVP

First release scope:

- `roadmapctl doctor`
- `roadmapctl check`

The public CLI contract is documented in [docs/cli-contract.md](docs/cli-contract.md). Implemented `/roadmap` commands that write, mutate, execute tasks, or claim roadmap validity must use `roadmapctl` as a blocking guard.

## Development

```bash
go test ./...
go build ./cmd/roadmapctl
```

Non-goals for MVP:

- no roadmap decomposition intelligence;
- no automatic materialization;
- no automatic fixing;
- no roadmap subcommands inside `rootline`.
