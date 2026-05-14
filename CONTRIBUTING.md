# Contributing to roadmapctl

Thank you for your interest in contributing!

roadmapctl is the deterministic governance CLI for Rootline-governed roadmaps. It owns roadmap-specific guardrails, status transitions, and agent workflow coordination while rootline remains the generic file-based database engine.

## Development Setup

```bash
# Clone and enter the repo
git clone https://github.com/pablontiv/roadmapctl.git
cd roadmapctl

# Set up git hooks
git config core.hooksPath .githooks

# Verify environment
go build ./cmd/roadmapctl
go test ./...
golangci-lint run ./...
```

Requires Go 1.21+ and [golangci-lint v2](https://golangci-lint.run/). roadmapctl also expects a compatible `rootline` binary in PATH (or via `--rootline` / `ROOTLINE_BIN`).

## Workflow

1. Fork the repository
2. Create a feature branch from `master`
3. Make your changes
4. Run `go test ./...` and `golangci-lint run ./...`
5. Commit using [Conventional Commits](https://www.conventionalcommits.org/)
6. Open a Pull Request

## Releasing

Releases are fully automated via CI. On push to `master`, CI analyzes conventional commit prefixes, calculates the next semver version, builds multi-platform binaries, and creates a GitHub Release. The skill sync (`scripts/sync-roadmap-skill.sh`) runs automatically via pre-push hook. No manual release steps needed.

## Commit Convention

```
type(scope): description
```

| Type | When to use |
|------|-------------|
| `feat` | New user-facing command or behavior |
| `fix` | Bug fix |
| `refactor` | Internal restructuring, no behavior change |
| `test` | Adding or updating tests |
| `docs` | Documentation only |
| `chore` | Build, CI, dependency updates |

Breaking changes use `!` suffix: `feat!: rename transition subcommand`

## Code Style

- **Formatting**: `gofmt` (enforced in CI and pre-commit hook)
- **Linting**: `golangci-lint v2` (CI lint gate)
- **Testing**: stdlib `testing` package; unit and integration tests
- **CLI framework**: cobra

## Non-Goals

roadmapctl is deliberately constrained. It does **not**:
- Decompose tasks or generate roadmap content (that is the `/roadmap` skill's job)
- Auto-fix roadmap files beyond what `transition` commands permit
- Implement subcommands inside `rootline`

Keep PRs focused on deterministic governance logic. Avoid adding planning intelligence or AI decomposition behavior.

## Reporting Issues

- **Bugs**: Use the bug report template
- **Features**: Use the feature request template
- **Security**: See [SECURITY.md](SECURITY.md) for responsible disclosure
