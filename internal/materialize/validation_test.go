package materialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func TestDryRunRejectsInvalidPlanMetadata(t *testing.T) {
	_, found, err := DryRun(t.TempDir(), Plan{Version: 2, Kind: "wrong"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		diagnostics.DiagnosticMaterializeInputVersionUnsupported,
		diagnostics.DiagnosticMaterializeInputKindInvalid,
		diagnostics.DiagnosticMaterializeInputEmpty,
	}
	for _, id := range want {
		if !hasDiagnostic(found, id) {
			t.Fatalf("missing %s in %#v", id, found)
		}
	}
}

func TestDryRunRejectsInvalidDependencyShape(t *testing.T) {
	plan := samplePlan()
	plan.Items[1].BlockedBy = []Dependency{{Ref: "x", Path: "y"}}
	_, found, err := DryRun(t.TempDir(), plan)
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(found, diagnostics.DiagnosticMaterializeInputDependencyInvalid) {
		t.Fatalf("diagnostics = %#v", found)
	}
}

func TestDryRunRendersPlanLocalDependencies(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)
	plan := Plan{Version: 1, Kind: PlanKind, Items: []Item{
		{
			Type: "outcome", Slug: "alpha", Title: "Alpha", Description: "Alpha.", AcceptanceCriteria: []string{"Alpha done."},
			Tasks: []Task{
				completeTask("first", "First"),
				withBlockedBy(completeTask("second", "Second"), []Dependency{{Ref: "alpha/first"}}),
			},
		},
		{
			Type: "outcome", Slug: "beta", Title: "Beta", Description: "Beta.", AcceptanceCriteria: []string{"Beta done."},
			Tasks: []Task{withBlockedBy(completeTask("third", "Third"), []Dependency{{Ref: "alpha/first"}})},
		},
	}}

	result, found, err := DryRun(root, plan)
	if err != nil || len(found) != 0 {
		t.Fatalf("err=%v diagnostics=%#v", err, found)
	}
	joined := strings.Join([]string{result.Changes[2].Content, result.Changes[4].Content}, "\n")
	if !strings.Contains(joined, "[[blocked_by:./T001-first.md]]") {
		t.Fatalf("missing same-outcome link:\n%s", joined)
	}
	if !strings.Contains(joined, "[[blocked_by:../O01-alpha/T001-first.md]]") {
		t.Fatalf("missing cross-outcome link:\n%s", joined)
	}
}

func TestDryRunRejectsUnresolvedAndBareDependencies(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)

	plan := samplePlan()
	plan.Items[1].BlockedBy = []Dependency{{Ref: "missing/task"}}
	_, found, err := DryRun(root, plan)
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(found, diagnostics.DiagnosticMaterializeInputDependencyUnresolved) {
		t.Fatalf("unresolved diagnostics = %#v", found)
	}

	plan = samplePlan()
	plan.Items[1].BlockedBy = []Dependency{{Path: "T001-bare.md"}}
	_, found, err = DryRun(root, plan)
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(found, diagnostics.DiagnosticMaterializeInputDependencyInvalid) {
		t.Fatalf("bare diagnostics = %#v", found)
	}

	plan = samplePlan()
	plan.Items[1].BlockedBy = []Dependency{{Path: "./T999-missing.md"}}
	_, found, err = DryRun(root, plan)
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(found, diagnostics.DiagnosticMaterializeInputDependencyUnresolved) {
		t.Fatalf("unresolved explicit path diagnostics = %#v", found)
	}

	plan = samplePlan()
	plan.Items[1].BlockedBy = []Dependency{{Path: "/T001-existing.md"}}
	_, found, err = DryRun(root, plan)
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(found, diagnostics.DiagnosticMaterializeInputDependencyInvalid) {
		t.Fatalf("absolute path diagnostics = %#v", found)
	}
}

func TestDryRunRejectsDirectoryDependencyPath(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)
	if err := os.Mkdir(filepath.Join(root, "T001-directory.md"), 0o755); err != nil {
		t.Fatal(err)
	}

	plan := samplePlan()
	plan.Items[1].BlockedBy = []Dependency{{Path: "./T001-directory.md"}}
	_, found, err := DryRun(root, plan)
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(found, diagnostics.DiagnosticMaterializeInputDependencyInvalid) {
		t.Fatalf("directory dependency diagnostics = %#v", found)
	}
}

func TestDryRunAllowsExistingExplicitDependencyPath(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)
	if err := writeFile(filepath.Join(root, "T001-existing.md"), "---\nestado: Completed\ntipo: task\n---\n# T001: Existing\n"); err != nil {
		t.Fatal(err)
	}

	plan := samplePlan()
	plan.Items[1].BlockedBy = []Dependency{{Path: "./T001-existing.md"}}
	result, found, err := DryRun(root, plan)
	if err != nil || len(found) != 0 {
		t.Fatalf("err=%v diagnostics=%#v", err, found)
	}
	if !strings.Contains(result.Changes[len(result.Changes)-1].Content, "[[blocked_by:./T001-existing.md]]") {
		t.Fatalf("missing existing dependency link:\n%s", result.Changes[len(result.Changes)-1].Content)
	}
}

func TestApplyCreatesMissingRootBootstrap(t *testing.T) {
	root := filepath.Join(t.TempDir(), "docs", "roadmap")
	result, found, err := Apply(root, samplePlan())
	if err != nil || len(found) != 0 {
		t.Fatalf("err=%v diagnostics=%#v", err, found)
	}
	for _, rel := range []string{".stem", ".roadmapctl.toml", "O01-new-outcome/README.md", "O01-new-outcome/T001-first-task.md", "T001-direct-task.md"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("missing %s after apply: %v", rel, err)
		}
	}
	if !result.Changes[0].Applied || result.Changes[0].Operation != "mkdir" {
		t.Fatalf("first change = %#v, want applied mkdir", result.Changes[0])
	}
}

func TestDryRunReportsExistingRootTaskSlugCollision(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)
	if err := writeFile(filepath.Join(root, "T001-direct-task.md"), "existing"); err != nil {
		t.Fatal(err)
	}
	_, found, err := DryRun(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if !hasDiagnostic(found, diagnostics.DiagnosticMaterializePlanConflict) {
		t.Fatalf("diagnostics = %#v", found)
	}
}

func completeTask(slug string, title string) Task {
	return Task{Slug: slug, Title: title, Description: title + ".", Preserves: []string{"Invariant."}, Context: "Context.", ScopeIn: []string{"In."}, ScopeOut: []string{"Out."}, InitialState: "Initial.", AcceptanceCriteria: []string{"AC."}, SourceOfTruth: []string{"docs/materialize-plan-schema.md"}}
}

func withBlockedBy(task Task, deps []Dependency) Task {
	task.BlockedBy = deps
	return task
}

func hasDiagnostic(found []diagnostics.Diagnostic, id string) bool {
	for _, diagnostic := range found {
		if diagnostic.ID == id {
			return true
		}
	}
	return false
}

func writeFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
