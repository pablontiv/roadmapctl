# Release outline

`roadmapctl` releases are separate from Rootline releases. Rootline remains the generic DBMS/constraint engine; `roadmapctl` is the roadmap-specific guard CLI.

## Installation

Before packaged releases are available, install from source with Go:

```bash
go install github.com/pablontiv/roadmapctl/cmd/roadmapctl@latest
```

Verify the binary is on `PATH`:

```bash
roadmapctl --help
```

Commands that inspect or validate roadmap files require a compatible `rootline` executable available via `--rootline`, `ROOTLINE_BIN`, or `PATH`.

## Rootline compatibility policy

Current policy: `roadmapctl` does not hard-fail solely by Rootline version unless a minimum version is explicitly approved. Compatibility is capability-based. `roadmapctl doctor` always reports the detected `rootline --version` for diagnostics and release evidence.

Recommended/tested Rootline: `v0.9.100-33-g40a3fbc` or newer in the `v0.9.100` line. CI currently installs `github.com/pablontiv/rootline/cmd/rootline@latest` on Linux, macOS, and Windows.

Required capabilities by command family:

| Rootline command | Required by | Compatibility expectation |
|------------------|-------------|---------------------------|
| `--version` | `doctor` | Emits version text for environment reports. |
| `validate --all <root> --output json` | `check` | Emits parseable JSON, including when validation exits non-zero. |
| `describe <root>/ --output json` | `check`, `context`, `lint`, `transition` | Exposes schema enum values via `schema.<field>.values` or supported legacy top-level `values`. |
| `query <root> --where ... --output json` | `check`, read commands, transition checks | Emits `rows[]` with `path`, `frontmatter`, and optional `derived`. |
| `graph <root> --where ... --output json` | `check`, read commands, transition checks | Emits `cycles[]` and `broken_links[]`. |
| `tree <root> --where ... --output json` | `context`, `pending`, `next`, `decision` | Emits recursive tree data with child nodes and completion counts. |
| `set <file> field=value` | `transition --apply` | May emit text; `roadmapctl` treats output as raw and validates after mutation. |
| `new <filepath>` | legacy/manual troubleshooting only | Materialization is implemented directly by roadmapctl and does not require Rootline `new`. |

Rootline compatibility diagnostics should differentiate:

- missing binary: `RMC_ENV_ROOTLINE_MISSING`, `details.kind="missing_binary"`, exit `3`;
- unsupported command or flag: operation diagnostic with `details.kind="incompatible_command"`;
- invalid JSON syntax: operation diagnostic with `details.kind="invalid_json"`;
- valid JSON with unsupported shape: operation diagnostic with `details.kind="invalid_shape"`.

Warnings for versions below a recommended line may be added later, but hard version gates require an explicit release-governance decision.

## CI release gates

Every release candidate should pass:

```bash
./scripts/check-coverage.sh
go build ./cmd/roadmapctl
```

`./scripts/check-coverage.sh` runs `go test ./... -coverprofile` and enforces a minimum total statement coverage of 85%. Override only for local experiments with `COVERAGE_THRESHOLD=<percent>`; release and CI runs use the default 85% gate.

The initial CI matrix runs those commands on:

- Linux;
- macOS;
- Windows.

## GoReleaser artifacts

The repository includes `.goreleaser.yml` for local dry-runs and tagged releases when publication is approved. It builds only the `roadmapctl` binary from this module; Rootline is not bundled and remains an external runtime dependency.

Validate config locally with:

```bash
goreleaser check
```

Run a local snapshot without publishing:

```bash
goreleaser release --snapshot --clean
```

The configured build matrix is:

| OS | Architectures |
|----|---------------|
| linux | amd64, arm64 |
| darwin | amd64, arm64 |
| windows | amd64, arm64 |

Expected artifacts:

- compressed archives per OS/architecture (`tar.gz`, Windows `zip`);
- `checksums.txt` for all artifacts;
- generated release notes from git tags;
- included docs: `README.md`, `docs/cli-contract.md`, `docs/release.md`.

Install options:

```bash
# From source
go install github.com/pablontiv/roadmapctl/cmd/roadmapctl@latest

# From a release archive
# 1. Download the archive for your OS/arch.
# 2. Verify it against checksums.txt.
# 3. Put the roadmapctl binary on PATH.
```

Homebrew/Scoop/Winget publishing, signing, SBOM generation, and installer channels are intentionally deferred until explicitly approved.
