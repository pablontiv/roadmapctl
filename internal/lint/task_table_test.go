package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func TestCheckOutcomeTaskTablesValidOutcomePasses(t *testing.T) {
	root := t.TempDir()
	writeOutcome(t, root, "O01-work", `# Work

## Tasks

| Task | Description |
| --- | --- |
| [T001](T001-first.md) | First |
| [T002](T002-second.md) | Second |
`, []string{"T001-first.md", "T002-second.md"})

	found, err := CheckOutcomeTaskTables(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Fatalf("diagnostics = %#v", found)
	}
}

func TestCheckOutcomeTaskTablesFindsMissingAndStaleRows(t *testing.T) {
	root := t.TempDir()
	writeOutcome(t, root, "O01-work", `# Work

## Tasks

| Task | Description |
| --- | --- |
| [T001](T001-first.md) | First |
| [T999](T999-stale.md) | Stale |
`, []string{"T001-first.md", "T002-second.md"})

	found, err := CheckOutcomeTaskTables(root)
	if err != nil {
		t.Fatal(err)
	}
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintTaskTableMissingRow, "O01-work/README.md", "T002-second.md")
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintTaskTableStaleRow, "O01-work/README.md", "T999-stale.md")
}

func TestCheckOutcomeTaskTablesNoMissingDiagnosticWhenTableAbsent(t *testing.T) {
	root := t.TempDir()
	// Outcome without ## Tasks table is valid: tasks are a computed view
	writeOutcome(t, root, "O01-no-table", `# Work
`, []string{"T001-first.md"})

	found, err := CheckOutcomeTaskTables(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range found {
		if d.ID == diagnostics.DiagnosticLintTaskTableMissing {
			t.Fatalf("unexpected missing-table diagnostic: %#v", d)
		}
	}
}

func TestCheckOutcomeTaskTablesFindsinvalidTableLinks(t *testing.T) {
	root := t.TempDir()
	writeOutcome(t, root, "O01-invalid", `# Other

## Tasks

| Task | Description |
| --- | --- |
| [bad](../O02-other/T001-first.md) | Outside |
`, []string{"T001-first.md"})

	found, err := CheckOutcomeTaskTables(root)
	if err != nil {
		t.Fatal(err)
	}
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintTaskTableInvalidLink, "O01-invalid/README.md", "../O02-other/T001-first.md")
}

func writeOutcome(t *testing.T, root string, name string, readme string, tasks []string) {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, task := range tasks {
		if err := os.WriteFile(filepath.Join(dir, task), []byte("---\nestado: Pending\ntipo: task\n---\n# Task\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func assertLintDiagnostic(t *testing.T, found []diagnostics.Diagnostic, id string, path string, target string) {
	t.Helper()
	for _, diagnostic := range found {
		if diagnostic.ID != id || diagnostic.Path != path {
			continue
		}
		if target == "" || diagnostic.Details["target"] == target {
			return
		}
	}
	t.Fatalf("missing diagnostic %s path=%s target=%s in %#v", id, path, target, found)
}
