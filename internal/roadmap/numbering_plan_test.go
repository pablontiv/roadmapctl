package roadmap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlanMaterializePathsHandlesMixedScopes(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, "O01-existing"))
	writeFile(t, filepath.Join(root, "T001-direct.md"))
	mkdir(t, filepath.Join(root, "O02-work"))
	writeFile(t, filepath.Join(root, "O02-work", "T003-old.md"))

	plan, diagnostics, err := PlanMaterializePaths(root, MaterializePathRequest{
		Outcomes:    []OutcomePathRequest{{Slug: "new-outcome", Tasks: []TaskPathRequest{{Slug: "first"}, {Slug: "second"}}}},
		DirectTasks: []TaskPathRequest{{Slug: "new-direct"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	if plan.Outcomes[0].Path != "O03-new-outcome/README.md" || plan.Outcomes[0].Tasks[0].Path != "O03-new-outcome/T001-first.md" || plan.Outcomes[0].Tasks[1].Path != "O03-new-outcome/T002-second.md" {
		t.Fatalf("outcome plan = %#v", plan.Outcomes[0])
	}
	if plan.DirectTasks[0].Path != "T002-new-direct.md" {
		t.Fatalf("direct plan = %#v", plan.DirectTasks)
	}
}

func TestPlanMaterializePathsRejectsExistingOutcomeSlug(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "O08-soporte-scip-todos-los-repos", "README.md"))

	plan, diagnostics, err := PlanMaterializePaths(root, MaterializePathRequest{
		Outcomes: []OutcomePathRequest{{Slug: "soporte-scip-todos-los-repos", Tasks: []TaskPathRequest{{Slug: "new-task"}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Outcomes) != 0 {
		t.Fatalf("plan = %#v, want no duplicate outcome plan", plan)
	}
	assertHasDiagnostic(t, diagnostics, DiagnosticMaterializePlanConflict, "O08-soporte-scip-todos-los-repos/README.md")
}

func TestPlanMaterializePathsDetectsCollisionsInvalidSlugsAndEscapes(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "T001-existing.md"))
	plan, diagnostics, err := PlanMaterializePaths(root, MaterializePathRequest{DirectTasks: []TaskPathRequest{{Slug: "bad/tasks"}, {Slug: "existing"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.DirectTasks) != 0 {
		t.Fatalf("plan = %#v", plan)
	}
	assertHasDiagnostic(t, diagnostics, DiagnosticMaterializeInputSlugInvalid, "")
	assertHasDiagnostic(t, diagnostics, DiagnosticMaterializePlanConflict, "T001-existing.md")
}

func TestPlanMaterializePathsUsesSlashSeparators(t *testing.T) {
	root := t.TempDir()
	plan, diagnostics, err := PlanMaterializePaths(root, MaterializePathRequest{Outcomes: []OutcomePathRequest{{Slug: "portable", Tasks: []TaskPathRequest{{Slug: "task"}}}}})
	if err != nil || len(diagnostics) != 0 {
		t.Fatalf("err=%v diagnostics=%#v", err, diagnostics)
	}
	if filepath.Separator == '\\' && plan.Outcomes[0].Tasks[0].Path != "O01-portable/T001-task.md" {
		t.Fatalf("path = %q", plan.Outcomes[0].Tasks[0].Path)
	}
	if plan.Outcomes[0].Tasks[0].Path != "O01-portable/T001-task.md" {
		t.Fatalf("path = %q", plan.Outcomes[0].Tasks[0].Path)
	}
}

func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	mkdir(t, filepath.Dir(path))
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
}
