package roadmap

import (
	"path/filepath"
	"testing"
)

func TestCheckStructureRejectsSingleSummaryFileFallback(t *testing.T) {
	diagnostics, err := CheckStructure(fixturePath("invalid-single-summary-file"))
	if err != nil {
		t.Fatalf("CheckStructure() error = %v", err)
	}

	assertHasDiagnostic(t, diagnostics, "RMC_STRUCTURE_SINGLE_FILE_FALLBACK", "roadmap-tasks.md")
}

func TestCheckStructureRejectsOutcomeMissingReadme(t *testing.T) {
	diagnostics, err := CheckStructure(fixturePath("invalid-missing-outcome-readme"))
	if err != nil {
		t.Fatalf("CheckStructure() error = %v", err)
	}

	assertHasDiagnostic(t, diagnostics, "RMC_STRUCTURE_MISSING_OUTCOME_README", "O01-no-readme/README.md")
}

func TestCheckStructureAcceptsValidDirectTasks(t *testing.T) {
	diagnostics, err := CheckStructure(fixturePath("valid-direct-tasks"))
	if err != nil {
		t.Fatalf("CheckStructure() error = %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v, want none", diagnostics)
	}
}

func TestCheckStructureAcceptsValidOutcomeWithTasks(t *testing.T) {
	diagnostics, err := CheckStructure(fixturePath("valid-outcome-with-tasks"))
	if err != nil {
		t.Fatalf("CheckStructure() error = %v", err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v, want none", diagnostics)
	}
}

func TestCheckStructureRejectsDuplicateIDsByScope(t *testing.T) {
	diagnostics, err := CheckStructure(fixturePath("invalid-duplicate-ids"))
	if err != nil {
		t.Fatalf("CheckStructure() error = %v", err)
	}

	assertHasDiagnostic(t, diagnostics, "RMC_STRUCTURE_DUPLICATE_ID", "T001-second.md")
	assertHasDiagnostic(t, diagnostics, "RMC_STRUCTURE_DUPLICATE_ID", "O01-second/README.md")
	assertHasDiagnostic(t, diagnostics, "RMC_STRUCTURE_DUPLICATE_ID", "O02-work/T001-second.md")
}

func TestCheckStructureRejectsExtraNestingUnderOutcome(t *testing.T) {
	diagnostics, err := CheckStructure(fixturePath("invalid-extra-nesting"))
	if err != nil {
		t.Fatalf("CheckStructure() error = %v", err)
	}

	assertHasDiagnostic(t, diagnostics, "RMC_STRUCTURE_EXTRA_NESTING", "O01-work/nested")
}

func TestCheckStructureNormalizesDiagnosticPathsToSlashSeparators(t *testing.T) {
	diagnostics, err := CheckStructure(filepath.Join("..", "..", "testdata", "fixtures", "invalid-extra-nesting", "docs", "roadmap"))
	if err != nil {
		t.Fatalf("CheckStructure() error = %v", err)
	}

	assertHasDiagnostic(t, diagnostics, "RMC_STRUCTURE_EXTRA_NESTING", "O01-work/nested")
	for _, diagnostic := range diagnostics {
		if filepath.Separator == '\\' {
			continue
		}
		if containsBackslash(diagnostic.Path) {
			t.Fatalf("diagnostic path = %q, want slash-normalized", diagnostic.Path)
		}
	}
}

func fixturePath(name string) string {
	return filepath.Join("..", "..", "testdata", "fixtures", name, "docs", "roadmap")
}

func assertHasDiagnostic(t *testing.T, diagnostics []Diagnostic, id string, path string) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.ID == id && diagnostic.Path == path {
			return
		}
	}
	t.Fatalf("missing diagnostic id=%q path=%q in %#v", id, path, diagnostics)
}

func containsBackslash(path string) bool {
	for _, r := range path {
		if r == '\\' {
			return true
		}
	}
	return false
}
