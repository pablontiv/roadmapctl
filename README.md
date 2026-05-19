# roadmapctl

[![CI](https://github.com/pablontiv/roadmapctl/actions/workflows/ci.yml/badge.svg)](https://github.com/pablontiv/roadmapctl/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: PolyForm NC](https://img.shields.io/badge/License-PolyForm%20NC-blue.svg)](LICENSE)

Companion CLI for Rootline-governed roadmaps.

`roadmapctl` owns roadmap-specific guardrails and workflows while `rootline` remains the generic file-based database and constraint engine.

| Concern | Owner |
|---------|-------|
| Filesystem database, `.stem` schema, frontmatter, validation | `rootline` |
| Status policy, transition guards, pending/next/decision, materialize | `roadmapctl` |
| Intent decomposition, agent orchestration, user dialogue | `/roadmap` skill |
| Code/docs changes, acceptance checks | Implementing agent |

---

## Table of Contents

- [Quick Start](#quick-start)
- [Workspace mode](#workspace-mode)
- [Core Idea](#core-idea)
- [Layer Responsibilities](#layer-responsibilities)
- [Installation](#installation)
- [Auto-update](#auto-update)
- [Commands](#commands)
- [AI-Native](#ai-native)
- [Skill Source](#skill-source)
- [Development](#development)
- [Documentation](#documentation)
- [License](#license)

---

## Quick Start

```bash
# 1. Check environment — rootline binary found, config valid
roadmapctl doctor
# Output:
# roadmapctl/doctor
# status: ok
# errors: 0
# warnings: 0
# infos: 1
# [info] RMC_ENV_PATH docs/roadmap: roadmapctl doctor paths resolved

# 2. Validate roadmap against schema (strict mode required before any writes)
roadmapctl check --strict
# Output:
# roadmapctl/check
# status: ok
# errors: 0
# warnings: 0
# infos: 0

# 3. List all pending tasks
roadmapctl pending
# Output:
# roadmapctl/pending
# status: ok
# pending: 2

# 4. See what to work on next
roadmapctl next --output json
# Output (trimmed):
# {
#   "version": 1, "kind": "roadmapctl/next",
#   "summary": {"status": "ok", "errors": 0, "warnings": 0, "infos": 0},
#   "ready": [{"path": "OXX-slug/TXXX-slug.md", "status": "Specified"}],
#   "blocked": []
# }

# 5. Prioritize across ready/blocked tasks
roadmapctl decision
# Output:
# roadmapctl/decision
# status: ok
# recommendations: 2
# quick_wins: 2
# blocked: 0

# 6. Resolve effective config for agents (skill bootstrap)
roadmapctl bootstrap --output json
# Output (trimmed):
# {
#   "version": 1, "kind": "roadmapctl/bootstrap",
#   "summary": {"status": "ok", "errors": 0, ...},
#   "root": "/abs/path/to/repo",
#   "roadmap_root": "/abs/path/to/repo/docs/roadmap",
#   "config_path": "docs/roadmap/.roadmapctl.toml",
#   "helpers": {"where_leaf": "isIndex == false", ...}
# }

# 7. Start a task — transitions estado Pending → In Progress
roadmapctl transition start docs/roadmap/T001-my-task.md --apply --output json
# Output (trimmed):
# {
#   "kind": "roadmapctl/transition", "action": "start",
#   "allowed": true, "current_status": "Specified", "target_status": "In Progress",
#   "changes": [{"field": "estado", "before": "Specified", "after": "In Progress", "applied": true}]
# }

# 8. Complete a task — transitions estado In Progress → Completed
roadmapctl transition complete docs/roadmap/T001-my-task.md --apply
```

---

## Workspace mode

When `roadmapctl` runs from a directory without a `.git` directory but containing sibling repos with their own `.git`, it operates in **workspace mode**.

**Convention**: each participating repo maintains its own complete roadmap under `<repo>/docs/roadmap/` with `.stem`, `.roadmapctl.toml`, outcomes, and tasks. Each repo is autonomous — its roadmap, its tasks, its commits.

**Layout example**:

```text
my-workspace/                    # parent dir without .git
├── docs/                        # repo 1: own .git + docs/roadmap/
│   ├── .git/
│   └── docs/roadmap/
│       ├── .stem
│       └── .roadmapctl.toml
├── tsg-valuecreation-core/      # repo 2: same layout
│   ├── .git/
│   └── docs/roadmap/
│       └── ...
└── tsg-valuecreation-frontend/  # repo 3: same layout
    ├── .git/
    └── docs/roadmap/
        └── ...
```

**Invocation**: most commands operate on a single repo at a time. Pass `--repo <path>` to target a specific repo. `roadmapctl pending --workspace` iterates the discovered repos and aggregates results.

**Anti-pattern**: do not create "code repos" without their own roadmap, expecting a central roadmap repo to commit on their behalf. Each repo is autonomous; cross-repo commit routing is not supported. If a repo needs to participate in the workspace, it needs its own `<repo>/docs/roadmap/`.

---

## Core Idea

Roadmaps are Markdown files. Rootline governs their structure via `.stem` schemas. roadmapctl adds the **governance layer** that makes those files an operational system.

- Status transitions are explicit and guarded — `transition` commands validate preconditions before writing
- `doctor` and `check` are blocking guards — agents cannot write without a clean preflight
- `pending` and `next` give agents a deterministic queue without requiring product judgment
- All output is stable JSON with versioned contracts — safe for automated pipelines

roadmapctl does not plan, decompose, or generate content. It **enforces the invariants**.

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

## Auto-update

roadmapctl updates itself automatically using a **staged async** pattern. A new release is downloaded in the background during run N and applied transparently at the start of run N+1 — no prompts, no interruptions.

```bash
# Disable auto-update entirely
ROADMAPCTL_NO_UPDATE=1 roadmapctl <command>
```

- Development builds (`version == "dev"`) never auto-update.
- Permission errors applying an update are suppressed silently — the command always runs.
- See [docs/auto-update.md](docs/auto-update.md) for the full pattern, OS-specific behavior, and troubleshooting.

---

## Commands

```bash
# Guards (blocking — run before any write or mutation)
roadmapctl doctor --repo <path>                       # Verify environment: rootline binary, config, schema
roadmapctl check --repo <path> [--strict] [--output json]  # Validate roadmap against .stem schema
roadmapctl lint --repo <path>                         # Check format conventions

# Read-only state (safe to call at any time)
roadmapctl bootstrap --repo <path> --output json      # Effective config for agents — helpers, thresholds, flags
roadmapctl pending --repo <path>                      # All tasks not in a done status
roadmapctl next --repo <path>                         # Suggested next task based on priority/order
roadmapctl decision <query> --repo <path>             # Query indexed decisions

# Controlled mutation (require --apply; blocked if preflight fails)
roadmapctl transition start <task.md> --repo <path> --apply
roadmapctl transition complete <task.md> --repo <path> --apply
roadmapctl materialize <spec> --repo <path> --apply
```

The public CLI contract is documented in [docs/cli-contract.md](docs/cli-contract.md). Skill integration details live in [docs/roadmap-skill-integration.md](docs/roadmap-skill-integration.md).

---

## AI-Native

roadmapctl is designed to be invoked by AI agents without human supervision.

- `--output json` on all guards produces stable versioned contracts (`"version": 1, "kind": "roadmapctl/..."`)
- `roadmapctl bootstrap` is the configuration API for agents — resolves helpers, thresholds, and flags from a single call
- `roadmapctl check --strict` returns a non-zero exit code when the roadmap is invalid — agents can gate on this
- The `/roadmap` skill delegates every status decision to roadmapctl — no policy is reimplemented in the skill

```bash
# Agent bootstrap pattern
config=$(roadmapctl bootstrap --output json)
pending=$(roadmapctl pending --output json)
next=$(roadmapctl next --output json)
```

---

## Skill Source

This repository is the canonical home for the `/roadmap` and `/retrospective` Claude Code skills:

```text
.claude/skills/roadmap/
.claude/skills/retrospective/
```

The git hooks in `.githooks/pre-push` and `.githooks/post-merge` keep the user-scope tools current:

```text
~/.claude/skills/roadmap
~/.claude/skills/retrospective
/usr/local/bin/roadmapctl   # override with ROADMAPCTL_BIN
```

```bash
git config core.hooksPath .githooks
scripts/install-user.sh
scripts/sync-roadmap-skill.sh --check
scripts/sync-roadmap-skill.sh --check --skill retrospective
```

`scripts/sync-roadmap-skill.sh` accepts `--skill NAME` to sync any skill under `.claude/skills/`
(default: `roadmap`). `install-user.sh` syncs all registered skills automatically.

---

## Development

```bash
go test ./...
go build ./cmd/roadmapctl
golangci-lint run ./...   # CI lint gate (golangci-lint v2 required)
./scripts/check-coverage.sh  # coverage gate (≥85.0%)
```

Common lint constraints: `defer f.Close()` must be wrapped as `defer func() { _ = f.Close() }()` (errcheck); `httpClient.Do(req)` requires `//nolint:gosec` when the URL comes from a variable (G704); avoid indexing array `b[i]` inside `for i := range a` across two arrays (G602 false positive). Cross-platform path assertions: use `filepath.ToSlash` on both sides when comparing output against `filepath.Join` paths.

Non-goals:

- no roadmap decomposition intelligence
- no automatic fixing
- no roadmap subcommands inside `rootline`

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full development workflow.

---

## Documentation

| Topic | Description |
|-------|-------------|
| [CLI Contract](docs/cli-contract.md) | Commands, flags, exit codes, JSON output shapes |
| [Auto-update](docs/auto-update.md) | Staged async update pattern, OS behavior, escape hatches |
| [Skill Integration](docs/roadmap-skill-integration.md) | How the `/roadmap` skill delegates to roadmapctl |
| [Release](docs/release.md) | Release outline and rootline compatibility notes |
| [Roadmap](docs/roadmap/) | Project roadmap (governed by rootline + roadmapctl) |

---

## License

[PolyForm Noncommercial 1.0.0](LICENSE) — free for non-commercial use.
