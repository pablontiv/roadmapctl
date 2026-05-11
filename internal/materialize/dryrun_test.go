package materialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/templates"
)

func TestDryRunMissingRootIncludesExplicitBootstrapWithoutWriting(t *testing.T) {
	root := filepath.Join(t.TempDir(), "docs", "roadmap")

	result, diagnostics, err := DryRun(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	wantLeading := []string{".", ".stem", ".roadmapctl.toml"}
	for i, want := range wantLeading {
		if result.Changes[i].Path != want || result.Changes[i].Applied {
			t.Fatalf("bootstrap change[%d] = %#v, want %s applied=false", i, result.Changes[i], want)
		}
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("dry-run created missing root: %v", err)
	}
	for _, change := range result.Changes {
		if change.Path == ".roadmapctl.toml" {
			for _, want := range []string{"required_code_coverage = 85.0", "loop_max_tasks = 0", "parallel = true", "autonomy = \"until_done\"", "compact_after_task_commit = true", "pr_mode = false"} {
				if !strings.Contains(change.Content, want) {
					t.Fatalf("materialize bootstrap TOML missing %q:\n%s", want, change.Content)
				}
			}
			return
		}
	}
	t.Fatal("missing .roadmapctl.toml bootstrap change")
}

func TestApplyCreatesCanonicalFiles(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)

	result, diagnostics, err := Apply(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	for _, change := range result.Changes {
		if !change.Applied {
			t.Fatalf("change was not applied: %#v", change)
		}
		if strings.HasSuffix(change.Path, "-tasks.md") {
			t.Fatalf("apply generated forbidden summary file: %s", change.Path)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(change.Path))); err != nil {
			t.Fatalf("applied file %s missing: %v", change.Path, err)
		}
	}
}

func TestApplyDetectsStaleDryRunCollisionWithoutWritingPlannedFile(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)
	if err := os.WriteFile(filepath.Join(root, "T001-direct-task.md"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, diagnostics, err := Apply(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) == 0 || diagnostics[0].ID != "RMC_MATERIALIZE_PLAN_CONFLICT" {
		t.Fatalf("diagnostics = %#v, want plan conflict", diagnostics)
	}
	data, err := os.ReadFile(filepath.Join(root, "T001-direct-task.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "existing" {
		t.Fatalf("stale apply overwrote existing file: %q", data)
	}
}

func TestApplyDoesNotOverwriteExistingStem(t *testing.T) {
	root := t.TempDir()
	stemPath := filepath.Join(root, ".stem")
	originalStem := []byte("custom stem")
	if err := os.WriteFile(stemPath, originalStem, 0o644); err != nil {
		t.Fatal(err)
	}

	_, diagnostics, err := Apply(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	currentStem, err := os.ReadFile(stemPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(currentStem) != string(originalStem) {
		t.Fatalf("existing .stem overwritten: %q", currentStem)
	}
}

func TestDryRunPlansOutcomeAndDirectTaskWithoutWriting(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "O01-existing"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeBootstrapFiles(t, root)
	if err := os.WriteFile(filepath.Join(root, "T001-existing.md"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}

	result, diagnostics, err := DryRun(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	wantPaths := []string{"O02-new-outcome/README.md", "O02-new-outcome/T001-first-task.md", "T002-direct-task.md"}
	if len(result.Changes) != len(wantPaths) {
		t.Fatalf("changes = %#v", result.Changes)
	}
	for i, want := range wantPaths {
		if result.Changes[i].Path != want || result.Changes[i].Operation != "create" || result.Changes[i].Applied {
			t.Fatalf("change[%d] = %#v, want create %s applied=false", i, result.Changes[i], want)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(want))); !os.IsNotExist(err) {
			t.Fatalf("dry-run wrote %s", want)
		}
	}
	if result.Changes[0].Content == "" || result.Changes[1].Content == "" || result.Changes[2].Content == "" {
		t.Fatalf("changes must include proposed content: %#v", result.Changes)
	}
	if !strings.Contains(result.Changes[0].Content, "tipo: outcome") {
		t.Fatalf("outcome content should include outcome type: %s", result.Changes[0].Content)
	}
	if strings.Contains(result.Changes[0].Content, "estado:") {
		t.Fatalf("outcome content should not include manual estado:\n%s", result.Changes[0].Content)
	}
	if !strings.Contains(result.Changes[1].Content, "estado: Specified") {
		t.Fatalf("task content should start as Specified:\n%s", result.Changes[1].Content)
	}
	if !strings.Contains(result.Changes[1].Content, "# T001: First task") {
		t.Fatalf("task content did not use planned task ID:\n%s", result.Changes[1].Content)
	}
	if !strings.Contains(result.Changes[2].Content, "# T002: Direct task") {
		t.Fatalf("direct task content did not use planned task ID:\n%s", result.Changes[2].Content)
	}
}

func TestDryRunRejectsExistingOutcomeSlugWithoutChanges(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)
	if err := os.Mkdir(filepath.Join(root, "O01-new-outcome"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "O01-new-outcome", "README.md"), []byte("---\ntipo: outcome\n---\n# Existing\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, diagnostics, err := DryRun(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changes) != 0 {
		t.Fatalf("changes = %#v, want none", result.Changes)
	}
	if len(diagnostics) != 1 || diagnostics[0].ID != "RMC_MATERIALIZE_PLAN_CONFLICT" || diagnostics[0].Path != "O01-new-outcome/README.md" {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
}

func TestApplyTargetCreatesOnlySelectedCanonicalFile(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)
	result, diagnostics, err := DryRun(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}

	directPath := result.Changes[len(result.Changes)-1].Path
	applied, diagnostics, err := ApplyTarget(root, result.Changes, directPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	if len(applied.Changes) != 1 || !applied.Changes[0].Applied || applied.Changes[0].Path != directPath {
		t.Fatalf("applied result = %#v", applied)
	}
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(directPath))); err != nil {
		t.Fatalf("selected target missing: %v", err)
	}
	for _, sibling := range []string{"O02-new-outcome/README.md", "O02-new-outcome/T001-first-task.md"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(sibling))); !os.IsNotExist(err) {
			t.Fatalf("sibling %s was written: %v", sibling, err)
		}
	}
}

func TestApplyTargetRejectsInvalidTargetsBeforeWriting(t *testing.T) {
	root := t.TempDir()
	writeBootstrapFiles(t, root)
	result, diagnostics, err := DryRun(root, samplePlan())
	if err != nil {
		t.Fatal(err)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	directPath := result.Changes[len(result.Changes)-1].Path
	withDuplicate := append(append([]Change{}, result.Changes...), result.Changes[len(result.Changes)-1])
	tests := []struct {
		name    string
		target  string
		changes []Change
	}{
		{name: "empty", target: "", changes: result.Changes},
		{name: "unknown", target: "T999-missing.md", changes: result.Changes},
		{name: "duplicate", target: directPath, changes: withDuplicate},
		{name: "non-file", target: ".", changes: []Change{{Path: ".", Operation: "mkdir"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, diagnostics, err := ApplyTarget(root, tt.changes, tt.target)
			if err != nil {
				t.Fatal(err)
			}
			if len(diagnostics) == 0 {
				t.Fatalf("expected diagnostics for target %q", tt.target)
			}
			if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(directPath))); !os.IsNotExist(err) {
				t.Fatalf("invalid target wrote file: %v", err)
			}
		})
	}
}

func writeBootstrapFiles(t *testing.T, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, ".stem"), []byte(templates.BaseStemContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".roadmapctl.toml"), []byte(templates.DefaultRoadmapctlTOML), 0o644); err != nil {
		t.Fatal(err)
	}
}

func samplePlan() Plan {
	return Plan{
		Version: 1,
		Kind:    "roadmapctl/materialize-plan",
		Items: []Item{
			{
				Type:               "outcome",
				Slug:               "new-outcome",
				Title:              "New Outcome",
				Description:        "Outcome description.",
				AcceptanceCriteria: []string{"Outcome works."},
				Tasks: []Task{
					{
						Slug:               "first-task",
						Title:              "First task",
						Description:        "Implement first task.",
						Preserves:          []string{"Dry-run stays read-only."},
						Context:            "Needed for materialization.",
						ScopeIn:            []string{"Create proposed markdown."},
						ScopeOut:           []string{"Do not write files."},
						InitialState:       "Roadmap root exists.",
						AcceptanceCriteria: []string{"Change is listed."},
						SourceOfTruth:      []string{"docs/materialize-plan-schema.md"},
					},
				},
			},
			{
				Type:               "task",
				Slug:               "direct-task",
				Title:              "Direct task",
				Description:        "Implement direct task.",
				Preserves:          []string{"Direct dry-run stays read-only."},
				Context:            "Needed for direct materialization.",
				ScopeIn:            []string{"Create proposed direct task markdown."},
				ScopeOut:           []string{"Do not write files."},
				InitialState:       "Roadmap root exists.",
				AcceptanceCriteria: []string{"Direct path is listed."},
				SourceOfTruth:      []string{"docs/materialize-plan-schema.md"},
			},
		},
	}
}
