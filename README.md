# roadmapctl

[![CI](https://github.com/pablontiv/roadmapctl/actions/workflows/ci.yml/badge.svg)](https://github.com/pablontiv/roadmapctl/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: PolyForm NC](https://img.shields.io/badge/License-PolyForm%20NC-blue.svg)](LICENSE)

Companion CLI for Rootline-governed roadmaps.

`roadmapctl` owns roadmap-specific guardrails and workflows while `rootline` remains the generic file-based database and constraint engine.

---

## Table of Contents

- [Layer Responsibilities](#layer-responsibilities)
- [Installation](#installation)
- [Commands](#commands)
- [Skill Source](#skill-source)
- [Development](#development)
- [Documentation](#documentation)
- [License](#license)

---

## Layer Responsibilities

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

---

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

`roadmapctl` expects a compatible `rootline` binary via `--rootline`, `ROOTLINE_BIN`, or `PATH`. See [docs/release.md](docs/release.md) for compatibility notes.

---

## Commands

Implemented command families:

| Family | Commands | Purpose |
|--------|----------|---------|
| Guards | `doctor`, `check`, `lint` | Validate environment and roadmap before writes |
| Read-only state | `context`, `pending`, `next`, `decision` | Query roadmap state without mutation |
| Controlled mutation | `transition`, `materialize`, `bootstrap` | Safe, guarded roadmap state transitions |

### Quick Reference

```bash
# Check environment and config
roadmapctl doctor --repo <path>

# Validate roadmap against schema
roadmapctl check --repo <path> --strict

# List pending tasks
roadmapctl pending --repo <path>

# Start a task (In Progress)
roadmapctl transition start <task.md> --repo <path> --apply

# Complete a task
roadmapctl transition complete <task.md> --repo <path> --apply

# Bootstrap config for a new repo
roadmapctl bootstrap --repo <path> --output json
```

The public CLI contract is documented in [docs/cli-contract.md](docs/cli-contract.md). Skill integration details live in [docs/roadmap-skill-integration.md](docs/roadmap-skill-integration.md).

---

## Skill Source

This repository is the canonical home for the `/roadmap` Claude Code skill:

```text
.claude/skills/roadmap/
```

The git hooks in `.githooks/pre-push` and `.githooks/post-merge` keep the user-scope tools current:

```text
~/.claude/skills/roadmap
/usr/local/bin/roadmapctl   # override with ROADMAPCTL_BIN
```

```bash
git config core.hooksPath .githooks
scripts/install-user.sh
scripts/sync-roadmap-skill.sh --check
```

---

## Development

```bash
go test ./...
go build ./cmd/roadmapctl
golangci-lint run ./...   # CI lint gate (golangci-lint v2 required)
```

Non-goals:

- no roadmap decomposition intelligence
- no automatic fixing
- no roadmap subcommands inside `rootline`

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full development workflow.

---

## Documentation

- [docs/cli-contract.md](docs/cli-contract.md) — Public CLI contract (commands, flags, exit codes)
- [docs/roadmap-skill-integration.md](docs/roadmap-skill-integration.md) — Skill integration guide
- [docs/release.md](docs/release.md) — Release outline and compatibility notes
- [docs/roadmap/](docs/roadmap/) — Project roadmap

---

## License

[PolyForm Noncommercial License 1.0.0](LICENSE) — free for personal and noncommercial use.
