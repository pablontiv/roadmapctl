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

The hook delegates to an explicit sync script that only replaces the `roadmap`
skill folder and does not touch any other user-scope skills:

```bash
scripts/sync-roadmap-skill.sh --install
scripts/sync-roadmap-skill.sh --check
```

## MVP

First release scope:

- `roadmapctl doctor`
- `roadmapctl check`

The public CLI contract is documented in [docs/cli-contract.md](docs/cli-contract.md). Implemented `/roadmap` commands that write, mutate, execute tasks, or claim roadmap validity must use `roadmapctl` as a blocking guard. Skill integration details live in [docs/roadmap-skill-integration.md](docs/roadmap-skill-integration.md).

## Installation

Until packaged releases are published, install from source:

```bash
go install github.com/pablontiv/roadmapctl/cmd/roadmapctl@latest
```

`roadmapctl` expects a compatible `rootline` binary via `--rootline`, `ROOTLINE_BIN`, or `PATH`. See [docs/release.md](docs/release.md) for the release outline and compatibility notes.

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
