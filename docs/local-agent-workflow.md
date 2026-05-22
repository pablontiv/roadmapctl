# Local Agent Workflow

This document contains roadmapctl-specific agent and verification notes. Shared
skills must stay portable; repo-specific rules like these belong in local docs.

## Core Commands

Use these commands for ordinary local verification:

```bash
go test ./...
go build ./cmd/roadmapctl
golangci-lint run ./...
just coverage-check
```

`just coverage` generates coverage output. `just coverage-check` enforces the
per-package floors declared in `.coverage-floors.toml` and is also used by the
pre-push hook when Go files changed.

## Rootline Dependency

`roadmapctl` shells out to a `rootline` binary. For normal CLI validation, the
binary must be available through `--rootline`, `ROOTLINE_BIN`, or `PATH`.

Tests that do not require real Rootline can use the fake Rootline path in test
helpers. Tests that require real Rootline behavior should skip when Rootline is
not available instead of hiding the requirement in assertions.

When debugging stale tool output, compare:

```bash
rootline --version
git -C /home/shared/rootline log --oneline -1
```

Do not rebuild external toolchains from shared skills. If a local workflow needs
a rebuild, record that requirement here or in the affected repo's local docs.

## Lint Notes

`golangci-lint` v2 is the local lint gate.

- `defer f.Close()` must be wrapped as `defer func() { _ = f.Close() }()` when
  errcheck would otherwise flag it.
- `httpClient.Do(req)` may require `//nolint:gosec` for G704 when the URL comes
  from a variable, even if the variable points at a fixed endpoint.
- Avoid indexing one array with an index from `range` over another array; gosec
  can report G602 even for equal-size arrays.
- For cross-platform path assertions, compare `filepath.ToSlash` output on both
  sides.
- For staticcheck SA5011, add an unreachable `return` after fatal nil-guard
  branches when needed so the analyzer sees that nil cannot fall through.

## Roadmap Validation

Before any roadmap write or task execution:

```bash
roadmapctl doctor --repo . --output json --strict
roadmapctl check --repo . --output json --strict
```

After any roadmap mutation:

```bash
roadmapctl check --repo . --output json --strict
```

For semantic conventions:

```bash
roadmapctl lint --repo testdata/fixtures/lint-valid --output json --strict
```

The repository's full historical roadmap currently contains legacy tasks that
can produce lint warnings under `--strict`; use targeted fixtures for strict
lint smoke tests unless the task explicitly migrates legacy roadmap files.

## Skill Sync

The canonical skill sources are under `.claude/skills/`. User-scope copies live
under `~/.claude/skills/` and are synchronized by:

```bash
scripts/sync-roadmap-skill.sh --check --all
scripts/sync-roadmap-skill.sh --install --all
```

The `.githooks/pre-push` and `.githooks/post-merge` hooks call
`scripts/install-user.sh`, which syncs all skills and rebuilds the installed
`roadmapctl` binary.

## Release Evidence

For release or skill/guard cutovers, collect evidence for:

- `go test ./...`
- `go build ./cmd/roadmapctl`
- `just coverage-check`
- `roadmapctl check --repo testdata/fixtures/valid-outcome-with-tasks --output json --strict`
- `roadmapctl lint --repo testdata/fixtures/lint-valid --output json --strict`
- missing Rootline negative check exits `3` with `RMC_ENV_ROOTLINE_MISSING`
- `scripts/sync-roadmap-skill.sh --check --all`
- `scripts/verify-roadmap-skill-headless.sh --evidence-dir <dir>` for skill or guard changes

Missing headless evidence blocks skill/guard releases.
