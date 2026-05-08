package roadmap

import "github.com/pablontiv/roadmapctl/internal/diagnostics"

func statusDiagnostics(decoded map[string]any, configured []string, schema []string) []Diagnostic {
	allowedStatuses := stringSet(configured)
	if len(schema) > 0 {
		allowedStatuses = intersectStringSets(allowedStatuses, stringSet(schema))
	}
	allowedTypes := map[string]bool{"task": true, "outcome": true}

	var found []Diagnostic
	for _, rowValue := range arrayValue(decoded["rows"]) {
		row, ok := rowValue.(map[string]any)
		if !ok {
			continue
		}
		path := stringField(row, "path")
		frontmatter, _ := row["frontmatter"].(map[string]any)
		status := stringField(frontmatter, "estado")
		if status == "" || !allowedStatuses[status] {
			found = append(found, Diagnostic{
				ID:       DiagnosticStatusUnknown,
				Severity: diagnostics.SeverityError,
				Message:  "task estado is not allowed by roadmap config or Rootline schema",
				Path:     path,
				Details:  map[string]any{"estado": status},
			})
		}
		tipo := stringField(frontmatter, "tipo")
		if tipo == "" || !allowedTypes[tipo] {
			found = append(found, Diagnostic{
				ID:       DiagnosticTypeUnknown,
				Severity: diagnostics.SeverityError,
				Message:  "record tipo is not allowed by Rootline schema",
				Path:     path,
				Details:  map[string]any{"tipo": tipo},
			})
		}
	}
	return found
}

func extractStatusValues(decoded map[string]any) []string {
	if values := stringsFromArray(decoded["values"]); len(values) > 0 {
		return values
	}
	schema, _ := decoded["schema"].(map[string]any)
	estado, _ := schema["estado"].(map[string]any)
	return stringsFromArray(estado["values"])
}

func arrayValue(value any) []any {
	items, _ := value.([]any)
	return items
}

func stringsFromArray(value any) []string {
	items := arrayValue(value)
	result := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

func stringField(fields map[string]any, key string) string {
	if fields == nil {
		return ""
	}
	value, _ := fields[key].(string)
	return value
}

func stringSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, value := range values {
		set[value] = true
	}
	return set
}

func intersectStringSets(left map[string]bool, right map[string]bool) map[string]bool {
	if len(left) == 0 {
		return right
	}
	result := map[string]bool{}
	for value := range left {
		if right[value] {
			result[value] = true
		}
	}
	return result
}
