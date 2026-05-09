package roadmap

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func TestInvalidBlockedByDiagnosticCoversRawTargetCases(t *testing.T) {
	root := t.TempDir()
	sourceDir := filepath.Join(root, "O01-work")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(sourceDir, "T001-task.md")
	taskTarget := filepath.Join(sourceDir, "T002-task.md")
	nonTaskTarget := filepath.Join(sourceDir, "README.md")
	for _, path := range []string{source, taskTarget, nonTaskTarget} {
		if err := os.WriteFile(path, []byte("---\nestado: Pending\ntipo: task\n---\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name        string
		target      string
		wantInvalid bool
	}{
		{name: "bare", target: "T002-task.md", wantInvalid: true},
		{name: "escape", target: "../../outside/T002-task.md", wantInvalid: true},
		{name: "non task", target: "./README.md", wantInvalid: true},
		{name: "valid explicit task", target: "./T002-task.md", wantInvalid: false},
		{name: "missing explicit target deferred to graph", target: "./T999-missing.md", wantInvalid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagnostic, invalid := invalidBlockedByDiagnostic(root, source, tt.target)
			if invalid != tt.wantInvalid {
				t.Fatalf("invalid = %v, want %v; diagnostic=%#v", invalid, tt.wantInvalid, diagnostic)
			}
			if invalid && diagnostic.ID != diagnostics.DiagnosticInvalidBlockedBy {
				t.Fatalf("diagnostic = %#v", diagnostic)
			}
		})
	}
}
