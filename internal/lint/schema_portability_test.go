package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func TestCheckFilenamePortabilityReportsCaseCollisionAndReservedName(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"T001-task.md", "t001-task.md", "CON.md"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(""), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	found, err := CheckFilenamePortability(root)
	if err != nil {
		t.Fatal(err)
	}
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintFilenameCaseCollision, "t001-task.md", "T001-task.md")
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintFilenameReserved, "CON.md", "CON")
}

func TestCheckSchemaCompatibilityAllowsExtensionsAndRequiresCoreFields(t *testing.T) {
	valid := map[string]any{
		"schema": map[string]any{"estado": map[string]any{}, "tipo": map[string]any{}, "custom": map[string]any{}},
		"links":  map[string]any{"rules": map[string]any{"blocked_by": map[string]any{}, "reference": map[string]any{}}},
	}
	if found := CheckSchemaCompatibility(valid); len(found) != 0 {
		t.Fatalf("valid diagnostics = %#v", found)
	}

	missing := map[string]any{"schema": map[string]any{"tipo": map[string]any{}}, "links": map[string]any{"rules": map[string]any{}}}
	found := CheckSchemaCompatibility(missing)
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintSchemaFieldMissing, ".stem", "estado")
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintSchemaLinkMissing, ".stem", "blocked_by")
}
