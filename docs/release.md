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

Minimum expected Rootline compatibility for the MVP:

- `rootline --version` works;
- JSON output is available for `validate`, `describe`, `query`, and `graph`;
- `rootline set` exists for roadmap loop integrations that mutate frontmatter.

The MVP has been developed against Rootline `v0.9.100-33-g40a3fbc` or newer in the `v0.9.100` line.

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
