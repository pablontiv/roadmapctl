package roadmap

import "github.com/pablontiv/roadmapctl/internal/diagnostics"

func statusDiagnostics(decoded map[string]any, configured []string, schemaStatuses []string, schemaTypes []string) []Diagnostic {
	allowedStatuses := stringSet(configured)
	if len(schemaStatuses) > 0 {
		allowedStatuses = stringSet(schemaStatuses)
	}
	allowedTypes := map[string]bool{"task": true, "outcome": true}
	if len(schemaTypes) > 0 {
		allowedTypes = stringSet(schemaTypes)
	}

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
				Message:  "task estado is not allowed by Rootline schema",
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

func operationalStatusDiagnostics(statuses []OperationalStatus, schemaStatuses []string) []Diagnostic {
	if len(statuses) == 0 || len(schemaStatuses) == 0 {
		return nil
	}
	allowed := stringSet(schemaStatuses)
	seen := map[string]bool{}
	var found []Diagnostic
	for _, status := range statuses {
		if status.Value == "" || allowed[status.Value] {
			continue
		}
		key := status.Source + "\x00" + status.Value
		if seen[key] {
			continue
		}
		seen[key] = true
		found = append(found, Diagnostic{
			ID:       DiagnosticConfigStatusSchemaMismatch,
			Severity: diagnostics.SeverityError,
			Message:  "configured operational status is not allowed by Rootline schema",
			Path:     ".claude/roadmap.local.md",
			Details:  map[string]any{"source": status.Source, "status": status.Value},
		})
	}
	return found
}

func extractStatusValues(decoded map[string]any) []string {
	if values := stringsFromArray(decoded["values"]); len(values) > 0 {
		return values
	}
	return extractSchemaEnumValues(decoded, "estado")
}

func extractTypeValues(decoded map[string]any) []string {
	return extractSchemaEnumValues(decoded, "tipo")
}

func extractSchemaEnumValues(decoded map[string]any, field string) []string {
	schema, _ := decoded["schema"].(map[string]any)
	fieldSchema, _ := schema[field].(map[string]any)
	return stringsFromArray(fieldSchema["values"])
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
