package lint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func TestCheckTaskSectionsValidTaskPasses(t *testing.T) {
	root := t.TempDir()
	writeTask(t, root, "T001-good.md", fullTaskMarkdown())
	found, err := CheckTaskSections(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Fatalf("diagnostics = %#v", found)
	}
}

func TestCheckTaskSectionsReportsMissingSectionAndEmptyLists(t *testing.T) {
	root := t.TempDir()
	writeTask(t, root, "T001-bad.md", `---
estado: Pending
tipo: task
---
# Bad

## Preserva

- invariant

## Contexto

text

## Alcance

**In**:

## Criterios de Aceptación

No bullets here.

## Fuente de verdad

`)
	found, err := CheckTaskSections(root)
	if err != nil {
		t.Fatal(err)
	}
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintTaskSectionMissing, "T001-bad.md", "Estado inicial esperado")
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintAcceptanceCriteriaMissing, "T001-bad.md", "")
	assertLintDiagnostic(t, found, diagnostics.DiagnosticLintSourceOfTruthEmpty, "T001-bad.md", "")
}

func fullTaskMarkdown() string {
	return `---
estado: Pending
tipo: task
---
# Good

## Preserva

- INV1: invariant

## Contexto

text

## Alcance

**In**:
1. work

## Estado inicial esperado

ready

## Criterios de Aceptación

- AC observable

## Fuente de verdad

- internal/lint
`
}

func writeTask(t *testing.T, root string, path string, content string) {
	t.Helper()
	fullPath := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
