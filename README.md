# roadmapctl

Companion CLI for Rootline-governed roadmaps.

`roadmapctl` owns roadmap-specific guardrails and workflows while `rootline` remains the generic file-based database and constraint engine.

## Layer responsibilities

`roadmapctl` is intentionally a deterministic governance layer for roadmaps that are stored as Rootline-governed Markdown. Each layer has a separate responsibility:

| Layer | Owns | Does not own |
|-------|------|--------------|
| Rootline | Generic Markdown filesystem database: hierarchy, `.stem` schema, frontmatter, links, validation, graph/query/tree/describe primitives. | Roadmap-specific status policy, prioritization, agent workflow, commits, PRs, or AI decomposition. |
| roadmapctl | Deterministic roadmap semantics on top of Rootline: config, status roles, guards, `check`/`lint`, `pending`/`next`/`decision`, `transition`, `materialize`, stable diagnostics and exit codes. | Free-form planning, product decisions, task implementation, creative prose generation, or replacing Rootline's structural engine. |
| `/roadmap` skill | Conversational/orchestration layer: understand user intent, decompose Outcomes/Tasks conceptually, ask for approval, present results, run implementation loops, coordinate agents. | Recomputing deterministic policy already owned by `roadmapctl`, mutating roadmap state directly with Rootline, or manually numbering/materializing roadmap files. |
| Implementing agent | Read an approved task, modify project code/docs, run acceptance checks, summarize and commit work. | Bypassing roadmap guards or deciding roadmap state transitions directly. |
| Git/CI/release | Reproducible evidence, branch/PR/release mechanics, checksums and distribution. | Serving as the source of truth for roadmap structure or status policy. |

The intended control flow is:

```text
User / agent
  -> /roadmap skill
  -> roadmapctl
  -> Rootline
  -> Markdown + .stem + links + filesystem
```

## Skill source

This repository is the canonical home for the `/roadmap` skill:

```text
.claude/skills/roadmap/
```

The git hooks in `.githooks/pre-push` and `.githooks/post-merge` keep the
user-scope tools current:

```text
~/.claude/skills/roadmap
/usr/local/bin/roadmapctl   # override with ROADMAPCTL_BIN
```

The hooks delegate to explicit install/sync scripts. Skill sync only replaces the
`roadmap` skill folder and does not touch any other user-scope skills:

```bash
git config core.hooksPath .githooks
scripts/install-user.sh
scripts/sync-roadmap-skill.sh --check
```

## Commands

Implemented command families:

- Guards: `doctor`, `check`, `lint`
- Read-only state: `context`, `pending`, `next`, `decision`
- Controlled mutation: `transition`, `materialize`, `bootstrap`

The public CLI contract is documented in [docs/cli-contract.md](docs/cli-contract.md). Implemented `/roadmap` commands that write, mutate, execute tasks, or claim roadmap validity must use `roadmapctl` as a blocking guard. Skill integration details live in [docs/roadmap-skill-integration.md](docs/roadmap-skill-integration.md).

## Installation

**Linux / macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/pablontiv/roadmapctl/master/install.sh | sh
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/pablontiv/roadmapctl/master/install.ps1 | iex
```

**From source (Go 1.21+):**

```bash
go install github.com/pablontiv/roadmapctl/cmd/roadmapctl@latest
```

`roadmapctl` expects a compatible `rootline` binary via `--rootline`, `ROOTLINE_BIN`, or `PATH`. See [docs/release.md](docs/release.md) for the release outline and compatibility notes.

## Development

```bash
go test ./...
go build ./cmd/roadmapctl
```

Non-goals:

- no roadmap decomposition intelligence;
- no automatic fixing;
- no roadmap subcommands inside `rootline`.
