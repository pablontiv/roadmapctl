package materialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	if !strings.Contains(result.Changes[1].Content, "# T001: First task") {
		t.Fatalf("task content did not use planned task ID:\n%s", result.Changes[1].Content)
	}
	if !strings.Contains(result.Changes[2].Content, "# T002: Direct task") {
		t.Fatalf("direct task content did not use planned task ID:\n%s", result.Changes[2].Content)
	}
}

func writeBootstrapFiles(t *testing.T, root string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(root, ".stem"), []byte(baseStemContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".roadmapctl.toml"), []byte(defaultRoadmapctlTOML), 0o644); err != nil {
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
