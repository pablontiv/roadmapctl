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
