package lint

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func isCaseInsensitiveFS(t *testing.T) bool {
	t.Helper()
	dir := t.TempDir()
	f1 := filepath.Join(dir, "CaseTest.txt")
	f2 := filepath.Join(dir, "casetest.txt")
	if err := os.WriteFile(f1, nil, 0600); err != nil {
		return false
	}
	if err := os.WriteFile(f2, nil, 0600); err != nil {
		return false
	}
	entries, _ := os.ReadDir(dir)
	return len(entries) == 1
}

func TestCheckFilenamePortabilityReportsCaseCollisionAndReservedName(t *testing.T) {
	if isCaseInsensitiveFS(t) {
		t.Skip("skipping case collision test on case-insensitive filesystem")
	}
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
		"schema": map[string]any{
			"estado": map[string]any{"required": true, "required_match": map[string]any{"patterns": []any{"T*"}}},
			"tipo":   map[string]any{},
			"custom": map[string]any{},
		},
		"validate": []any{map[string]any{"field": "tipo", "rule": "non_empty"}},
		"links":    map[string]any{"rules": map[string]any{"blocked_by": map[string]any{}, "reference": map[string]any{}}},
	}
	if found := CheckSchemaCompatibility(valid); len(found) != 0 {
		t.Fatalf("valid diagnostics = %#v", found)
	}

	missing := map[string]any{"schema": map[string]any{"tipo": map[string]any{}}, "links": map[string]any{"rules": map[string]any{}}}
	found := CheckSchemaCompatibility(missing)
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintSchemaFieldMissing, ".stem", "estado")
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintSchemaLinkMissing, ".stem", "blocked_by")
}

func TestCheckOutcomeSchemaCompatibilityAllowsTaskScopedEstadoRequirement(t *testing.T) {
	describe := map[string]any{
		"schema": map[string]any{
			"estado": map[string]any{"required": true, "required_match": map[string]any{"patterns": []any{"T*"}}},
		},
		"validate": []any{map[string]any{"field": "tipo", "rule": "non_empty"}},
	}
	if found := CheckOutcomeSchemaCompatibility(describe); len(found) != 0 {
		t.Fatalf("valid outcome schema diagnostics = %#v", found)
	}
}

func TestCheckOutcomeSchemaCompatibilityReportsEstadoRequiredForOutcomes(t *testing.T) {
	describe := map[string]any{
		"schema": map[string]any{
			"estado": map[string]any{"required": true, "required_match": map[string]any{"patterns": []any{"O*", "T*"}}},
			"tipo":   map[string]any{},
		},
		"links": map[string]any{"rules": map[string]any{"blocked_by": map[string]any{}}},
	}
	found := CheckOutcomeSchemaCompatibility(describe)
	assertLintDiagnostic(t, found, "RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED", ".stem", "estado.required_match")
}

func TestCheckOutcomeSchemaCompatibilityReportsGlobalEstadoNonEmptyValidate(t *testing.T) {
	describe := map[string]any{
		"schema": map[string]any{
			"estado": map[string]any{"required": true, "required_match": map[string]any{"patterns": []any{"T*"}}},
			"tipo":   map[string]any{},
		},
		"validate": []any{map[string]any{"field": "estado", "rule": "non_empty"}},
		"links":    map[string]any{"rules": map[string]any{"blocked_by": map[string]any{}}},
	}
	found := CheckOutcomeSchemaCompatibility(describe)
	assertLintDiagnostic(t, found, "RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY", ".stem", "validate.estado.non_empty")
}

func TestCheckFilenamePortabilityNoIssuesOnCleanDir(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"T001-task.md", "T002-feature.md"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	found, err := CheckFilenamePortability(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Fatalf("expected no diagnostics, got %v", found)
	}
}

func TestReservedWindowsNameDetectsAndIgnores(t *testing.T) {
	if got := reservedWindowsName("CON.md"); got != "CON" {
		t.Fatalf("reservedWindowsName(CON.md) = %q, want CON", got)
	}
	if got := reservedWindowsName("T001-task.md"); got != "" {
		t.Fatalf("reservedWindowsName(T001-task.md) = %q, want empty", got)
	}
	if got := reservedWindowsName("NUL"); got != "NUL" {
		t.Fatalf("reservedWindowsName(NUL) = %q, want NUL", got)
	}
}

func TestCheckFilenamePortabilityDetectsReservedName(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cannot create files named CON on Windows")
	}
	root := t.TempDir()
	for _, name := range []string{"CON.md", "T001-task.md"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	found, err := CheckFilenamePortability(root)
	if err != nil {
		t.Fatal(err)
	}
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintFilenameReserved, "CON.md", "CON")
}

func TestLintNameDiagnosticFormat(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "sub", "file.md")
	diag := lintNameDiagnostic(diagnostics.DiagnosticLintFilenameReserved, root, path, "test message", "CON")
	if diag.ID != diagnostics.DiagnosticLintFilenameReserved {
		t.Fatalf("lintNameDiagnostic ID = %q, want %s", diag.ID, diagnostics.DiagnosticLintFilenameReserved)
	}
}

func TestArrayValueStringSliceAndDefault(t *testing.T) {
	got := arrayValue([]string{"a", "b"})
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("arrayValue([]string) = %v, want [a b]", got)
	}
	if got := arrayValue(42); got != nil {
		t.Fatalf("arrayValue(int) = %v, want nil", got)
	}
}
