package lint

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

var windowsReservedNames = map[string]bool{
	"CON": true, "PRN": true, "AUX": true, "NUL": true,
	"COM1": true, "COM2": true, "COM3": true, "COM4": true, "COM5": true, "COM6": true, "COM7": true, "COM8": true, "COM9": true,
	"LPT1": true, "LPT2": true, "LPT3": true, "LPT4": true, "LPT5": true, "LPT6": true, "LPT7": true, "LPT8": true, "LPT9": true,
}

func CheckFilenamePortability(roadmapRoot string) ([]diagnostics.Diagnostic, error) {
	root := filepath.Clean(roadmapRoot)
	var found []diagnostics.Diagnostic
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() {
			if reserved := reservedWindowsName(entry.Name()); reserved != "" {
				found = append(found, lintNameDiagnostic(diagnostics.DiagnosticLintFilenameReserved, root, path, "roadmap filename is reserved on Windows", reserved))
			}
			return nil
		}
		diagnosticsForDir, err := checkCaseCollisionsInDir(root, path)
		if err != nil {
			return err
		}
		found = append(found, diagnosticsForDir...)
		if path != root && reservedWindowsName(entry.Name()) != "" {
			found = append(found, lintNameDiagnostic(diagnostics.DiagnosticLintFilenameReserved, root, path, "roadmap directory name is reserved on Windows", reservedWindowsName(entry.Name())))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sortDiagnostics(found)
	return found, nil
}

func checkCaseCollisionsInDir(root string, dir string) ([]diagnostics.Diagnostic, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	seen := map[string]string{}
	var found []diagnostics.Diagnostic
	for _, entry := range entries {
		key := strings.ToLower(entry.Name())
		path := filepath.Join(dir, entry.Name())
		if first, ok := seen[key]; ok {
			found = append(found, lintNameDiagnostic(diagnostics.DiagnosticLintFilenameCaseCollision, root, path, "roadmap entries collide on case-insensitive filesystems", first))
			continue
		}
		seen[key] = entry.Name()
	}
	return found, nil
}

func CheckSchemaCompatibility(describe map[string]any) []diagnostics.Diagnostic {
	var found []diagnostics.Diagnostic
	schema, _ := describe["schema"].(map[string]any)
	for _, field := range []string{"estado", "tipo"} {
		if _, ok := schema[field]; !ok {
			found = append(found, lintSchemaDiagnostic(diagnostics.DiagnosticLintSchemaFieldMissing, "effective schema is missing a required field", field))
		}
	}
	links, _ := describe["links"].(map[string]any)
	rules, _ := links["rules"].(map[string]any)
	if _, ok := rules["blocked_by"]; !ok {
		found = append(found, lintSchemaDiagnostic(diagnostics.DiagnosticLintSchemaLinkMissing, "effective schema is missing required blocked_by link rule", "blocked_by"))
	}
	sortDiagnostics(found)
	return found
}

func reservedWindowsName(name string) string {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	base = strings.TrimRight(base, " .")
	upper := strings.ToUpper(base)
	if windowsReservedNames[upper] {
		return upper
	}
	return ""
}

func lintNameDiagnostic(id string, root string, path string, message string, target string) diagnostics.Diagnostic {
	return diagnostics.Diagnostic{ID: id, Severity: diagnostics.SeverityError, Message: message, Path: relPath(root, path), Details: map[string]any{"target": target}}
}

func lintSchemaDiagnostic(id string, message string, target string) diagnostics.Diagnostic {
	return diagnostics.Diagnostic{ID: id, Severity: diagnostics.SeverityError, Message: message, Path: ".stem", Details: map[string]any{"target": target, "schema_key": target}}
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
