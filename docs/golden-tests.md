# Golden tests

`roadmapctl` keeps stable JSON expectations under `testdata/golden/`.

To update goldens after an intentional report-schema or diagnostic change:

1. Run the full suite first and inspect the golden diff:

   ```bash
   go test ./...
   ```

2. Rebuild the expected JSON with the same fixture commands covered by
   `internal/cli/golden_test.go`.
3. Normalize machine-local absolute fixture paths to `<fixture:name>` and
   Rootline version strings to `<rootline-version>` before writing files.
4. Re-run:

   ```bash
   go test ./...
   ```

Do not update goldens to hide behavior changes. The fixtures must exercise
`roadmapctl` diagnostics only; they must not require changes to Rootline or
materialize/fix roadmap files during tests.

`roadmapctl lint` goldens live under `testdata/golden/` like `doctor` and
`check` goldens. For the same warning diagnostics, strict and non-strict JSON
reports should remain identical; `--strict` changes only the process exit code.
Lint fixtures under `testdata/fixtures/lint-*` are read-only and cover valid,
warning, and error severities.

## Workspace fixtures

Workspace context tests require `.git` directories inside
`testdata/fixtures/valid-workspace/` and `testdata/fixtures/invalid-workspace-*/`.
Git does not track empty directories and refuses to track files inside a path
component named `.git`, so `TestMain` in `internal/cli/golden_test.go` creates
these directories at test startup via `os.MkdirAll`. Do not attempt to commit
`.gitkeep` files inside those directories.
