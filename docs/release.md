# Release outline

`roadmapctl` releases are separate from Rootline releases. Rootline remains the generic DBMS/constraint engine; `roadmapctl` is the roadmap-specific guard CLI.

## MVP installation

Before packaged releases are available, install from source with Go:

```bash
go install github.com/pablontiv/roadmapctl/cmd/roadmapctl@latest
```

Verify the binary is on `PATH`:

```bash
roadmapctl --help
```

`roadmapctl doctor` and `roadmapctl check` require a compatible `rootline` executable available via `--rootline`, `ROOTLINE_BIN`, or `PATH`.

## Rootline compatibility policy

Current policy: `roadmapctl` does not hard-fail solely by Rootline version unless a minimum version is explicitly approved. Compatibility is capability-based. `roadmapctl doctor` always reports the detected `rootline --version` for diagnostics and release evidence.

Recommended/tested Rootline: `v0.9.100-33-g40a3fbc` or newer in the `v0.9.100` line. CI currently installs `github.com/pablontiv/rootline/cmd/rootline@latest` on Linux, macOS, and Windows.

Required capabilities by command family:

| Rootline command | Required by | Compatibility expectation |
|------------------|-------------|---------------------------|
| `--version` | `doctor` | Emits version text for environment reports. |
| `validate --all <root> --output json` | `check` | Emits parseable JSON, including when validation exits non-zero. |
| `describe <root>/ --output json` | `check`, future context/lint | Exposes schema enum values via `schema.<field>.values` or supported legacy top-level `values`. |
| `query <root> --where ... --output json` | `check`, future read commands | Emits `rows[]` with `path`, `frontmatter`, and optional `derived`. |
| `graph <root> --where ... --output json` | `check`, future read commands | Emits `cycles[]` and `broken_links[]`. |
| `tree <root> --where ... --output json` | future read commands | Emits recursive tree data with child nodes and completion counts. |
| `set <file> field=value` | future transition commands | May emit text; `roadmapctl` must treat output as raw unless JSON stability is approved. |
| `new <filepath>` | future materialization commands | May emit text; `roadmapctl` must treat output as raw unless JSON stability is approved. |

Rootline compatibility diagnostics should differentiate:

- missing binary: `RMC_ENV_ROOTLINE_MISSING`, `details.kind="missing_binary"`, exit `3`;
- unsupported command or flag: operation diagnostic with `details.kind="incompatible_command"`;
- invalid JSON syntax: operation diagnostic with `details.kind="invalid_json"`;
- valid JSON with unsupported shape: operation diagnostic with `details.kind="invalid_shape"`.

Warnings for versions below a recommended line may be added later, but hard version gates require an explicit release-governance decision.

## CI release gates

Every release candidate should pass:

```bash
go test ./...
go build ./cmd/roadmapctl
```

The initial CI matrix runs those commands on:

- Linux;
- macOS;
- Windows.

## Future GoReleaser plan

When publication is approved, add a GoReleaser configuration that builds archives for:

| OS | Architectures |
|----|---------------|
| linux | amd64, arm64 |
| darwin | amd64, arm64 |
| windows | amd64, arm64 |

Expected artifacts:

- compressed archives per OS/architecture;
- `checksums.txt` for all artifacts;
- generated release notes from git tags;
- no Homebrew/Scoop/Winget publishing in the MVP unless explicitly approved later.

Signing, SBOM generation, package-manager publication, and installer channels are intentionally deferred beyond the MVP.
