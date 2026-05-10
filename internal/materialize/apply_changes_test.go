package materialize

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/templates"
)

func TestApplyChangesAppliesSortedAllowlistedBatch(t *testing.T) {
	root := t.TempDir()
	changes := []Change{
		{Path: "O02-new/T001-task.md", Operation: "create", Content: "---\nestado: Pending\ntipo: task\n---\n# Task\n"},
		{Path: ".roadmapctl.toml", Operation: "create", Content: templates.DefaultRoadmapctlTOML},
		{Path: "O02-new/README.md", Operation: "create", Content: "---\ntipo: outcome\n---\n# Outcome\n"},
		{Path: ".stem", Operation: "create", Content: templates.BaseStemContent},
		{Path: "O02-new", Operation: "mkdir"},
	}

	result, found, err := ApplyChanges(root, changes)
	if err != nil || len(found) != 0 {
		t.Fatalf("ApplyChanges err=%v diagnostics=%#v", err, found)
	}
	wantOrder := []string{"O02-new", ".roadmapctl.toml", ".stem", "O02-new/README.md", "O02-new/T001-task.md"}
	if len(result.Changes) != len(wantOrder) {
		t.Fatalf("changes = %#v", result.Changes)
	}
	for i, want := range wantOrder {
		change := result.Changes[i]
		if change.Path != want || !change.Applied {
			t.Fatalf("change[%d] = %#v, want path %q applied", i, change, want)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(want))); err != nil {
			t.Fatalf("expected %s: %v", want, err)
		}
	}
}

func TestApplyChangesRejectsInvalidOrConflictingBatchBeforeWriting(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(string)
		changes []Change
		wantID  string
	}{
		{name: "empty", changes: nil, wantID: diagnostics.DiagnosticMaterializeInputEmpty},
		{name: "unclean path", changes: []Change{{Path: "O01-new/../T001-task.md", Operation: "create", Content: "x"}}, wantID: "RMC_MATERIALIZE_CHANGE_INVALID"},
		{name: "bad mkdir target", changes: []Change{{Path: "not-outcome", Operation: "mkdir"}}, wantID: "RMC_MATERIALIZE_CHANGE_INVALID"},
		{name: "unsupported operation", changes: []Change{{Path: "T001-task.md", Operation: "delete", Content: "x"}}, wantID: "RMC_MATERIALIZE_CHANGE_INVALID"},
		{name: "missing create content", changes: []Change{{Path: "T001-task.md", Operation: "create"}}, wantID: "RMC_MATERIALIZE_CHANGE_INVALID"},
		{name: "non canonical create", changes: []Change{{Path: "notes.md", Operation: "create", Content: "x"}}, wantID: "RMC_MATERIALIZE_CHANGE_INVALID"},
		{
			name: "directory path exists as file",
			setup: func(root string) {
				if err := os.WriteFile(filepath.Join(root, "O01-new"), []byte("file"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			changes: []Change{{Path: "O01-new", Operation: "mkdir"}},
			wantID:  diagnostics.DiagnosticMaterializePlanConflict,
		},
		{
			name: "planned file exists",
			setup: func(root string) {
				if err := os.WriteFile(filepath.Join(root, "T001-task.md"), []byte("file"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			changes: []Change{{Path: "T001-task.md", Operation: "create", Content: "x"}},
			wantID:  diagnostics.DiagnosticMaterializePlanConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if tt.setup != nil {
				tt.setup(root)
			}
			before := listFiles(t, root)
			_, found, err := ApplyChanges(root, tt.changes)
			if err != nil {
				t.Fatalf("ApplyChanges err=%v", err)
			}
			if !hasDiagnostic(found, tt.wantID) {
				t.Fatalf("diagnostics = %#v, want %s", found, tt.wantID)
			}
			after := listFiles(t, root)
			if before != after {
				t.Fatalf("ApplyChanges wrote files\nbefore:\n%s\nafter:\n%s", before, after)
			}
		})
	}
}

func listFiles(t *testing.T, root string) string {
	t.Helper()
	var files []string
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return strings.Join(files, "\n")
}
