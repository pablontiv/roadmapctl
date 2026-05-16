# Derived vs Frontmatter Fields in rootline Query Rows

When implementing features that consume `roadmapctl next`, `roadmapctl pending`, or direct calls to `rootline query`, each row returned contains two separate field maps:

- **`frontmatter`**: raw YAML fields from the document header (`estado`, `tipo`, etc.)
- **`derived`**: fields computed by the rootline stem, either via `source:` expressions (e.g., `titulo` from `source: body.h1`) or built-in derivations (`is_done`, `isIndex`, etc.)

## Correct Field Access: `effectiveFields(row)`

The function `effectiveFields(row)` (implemented at `internal/roadmap/model.go:192`) merges both maps with **`derived` fields taking priority**. This is the correct way to access any field from a query row:

```go
func effectiveFields(row map[string]any) map[string]any {
	result := map[string]any{}
	if fm, ok := row["frontmatter"].(map[string]any); ok {
		for k, v := range fm {
			result[k] = v
		}
	}
	if derived, ok := row["derived"].(map[string]any); ok {
		for k, v := range derived {
			result[k] = v
		}
	}
	return result
}
```

Always call `effectiveFields(row)` first, then access fields from the result.

## Common Mistakes

**Bad:** accessing `titulo` directly from `frontmatter`:
```go
titulo := stringField(row["frontmatter"].(map[string]any), "titulo")  // always empty!
```

Since `titulo` is a derived field (computed from `source: body.h1`), it never appears in `frontmatter`.

**Good:** use `effectiveFields` to get the merged map:
```go
fields := effectiveFields(row)
titulo := stringField(fields, "titulo")  // correct
```

## Field Distribution

- **Always in `frontmatter`**: document YAML headers (`estado`, `tipo`, etc.)
- **Always in `derived`**: computed fields like `titulo` (via `source:` rules), `is_done`, `isIndex`
- **Access everything via**: `effectiveFields(row)` to ensure consistency

## Lint Patterns (golangci-lint v2 with gosec)

**G704 SSRF taint analysis** — gosec flags `httpClient.Do(req)` when the URL originates from a package-level variable (even if set to a known endpoint). Add `//nolint:gosec` to the `.Do(req)` call in addition to the `http.NewRequestWithContext` line.

**G602 slice bounds** — gosec flags `b[i]` when `i` comes from `range a` and `b` is a separate array (even if both are `[3]int`). Avoid range loops over one array when indexing another; compare fields directly instead.

**G703 errors on Close** — use `defer func() { _ = f.Close() }()` rather than `defer f.Close()` to satisfy errcheck.

**Cross-platform path tests** — when asserting that output contains a filesystem path, normalize both sides with `filepath.ToSlash` so the test passes on Windows (which uses `/` in Go's temp paths but `filepath.Join` produces `\`).
