package materialize

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func TestPlanPathsBasicOutcomeAndTask(t *testing.T) {
	// Create a temporary roadmap directory
	tmpDir, err := os.MkdirTemp("", "roadmap-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create basic bootstrap files
	stemPath := filepath.Join(tmpDir, ".stem")
	if err := os.WriteFile(stemPath, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create .stem: %v", err)
	}

	// Test input
	input := PathPlanInput{
		Version: 1,
		Kind:    PathPlanKind,
		Items: []PathPlanItem{
			{Type: "outcome", Slug: "rebuild-api"},
			{Type: "task", Slug: "add-endpoint"},
		},
	}

	// Plan paths
	result, found, err := PlanPaths(tmpDir, input)

	if err != nil {
		t.Fatalf("PlanPaths failed: %v", err)
	}

	if len(found) > 0 {
		t.Fatalf("unexpected diagnostics: %v", found)
	}

	if len(result.Paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(result.Paths))
	}

	// Check outcome path
	if result.Paths[0].Type != "outcome" {
		t.Fatalf("expected first path to be outcome, got %s", result.Paths[0].Type)
	}
	if result.Paths[0].Operation != "create" {
		t.Fatalf("expected first path operation to be create, got %s", result.Paths[0].Operation)
	}
	if !contains(result.Paths[0].Path, "O") {
		t.Fatalf("expected path to contain O prefix, got %s", result.Paths[0].Path)
	}
	if !contains(result.Paths[0].Path, "rebuild-api") {
		t.Fatalf("expected path to contain slug, got %s", result.Paths[0].Path)
	}

	// Check task path
	if result.Paths[1].Type != "task" {
		t.Fatalf("expected second path to be task, got %s", result.Paths[1].Type)
	}
	if result.Paths[1].Operation != "create" {
		t.Fatalf("expected second path operation to be create, got %s", result.Paths[1].Operation)
	}
}

func TestPlanPathsInvalidVersion(t *testing.T) {
	tmpDir := t.TempDir()

	input := PathPlanInput{
		Version: 2,
		Kind:    PathPlanKind,
		Items:   []PathPlanItem{{Type: "outcome", Slug: "test"}},
	}

	_, found, _ := PlanPaths(tmpDir, input)

	if len(found) == 0 {
		t.Fatalf("expected validation error for invalid version")
	}
	if found[0].ID != diagnostics.DiagnosticMaterializeInputVersionUnsupported {
		t.Fatalf("expected version unsupported diagnostic, got %s", found[0].ID)
	}
}

func TestPlanPathsInvalidKind(t *testing.T) {
	tmpDir := t.TempDir()

	input := PathPlanInput{
		Version: 1,
		Kind:    "invalid-kind",
		Items:   []PathPlanItem{{Type: "outcome", Slug: "test"}},
	}

	_, found, _ := PlanPaths(tmpDir, input)

	if len(found) == 0 {
		t.Fatalf("expected validation error for invalid kind")
	}
	if found[0].ID != diagnostics.DiagnosticMaterializeInputKindInvalid {
		t.Fatalf("expected kind invalid diagnostic, got %s", found[0].ID)
	}
}

func TestPlanPathsEmptyItems(t *testing.T) {
	tmpDir := t.TempDir()

	input := PathPlanInput{
		Version: 1,
		Kind:    PathPlanKind,
		Items:   []PathPlanItem{},
	}

	_, found, _ := PlanPaths(tmpDir, input)

	if len(found) == 0 {
		t.Fatalf("expected validation error for empty items")
	}
	if found[0].ID != diagnostics.DiagnosticMaterializeInputEmpty {
		t.Fatalf("expected empty input diagnostic, got %s", found[0].ID)
	}
}

func TestPlanPathsInvalidSlug(t *testing.T) {
	tmpDir := t.TempDir()

	input := PathPlanInput{
		Version: 1,
		Kind:    PathPlanKind,
		Items: []PathPlanItem{
			{Type: "outcome", Slug: "O-invalid"},
		},
	}

	_, found, _ := PlanPaths(tmpDir, input)

	if len(found) == 0 {
		t.Fatalf("expected validation error for invalid slug")
	}
	if found[0].ID != diagnostics.DiagnosticMaterializeInputSlugInvalid {
		t.Fatalf("expected slug invalid diagnostic, got %s", found[0].ID)
	}
}

func TestValidSlug(t *testing.T) {
	tests := map[string]bool{
		"valid-slug":     true,
		"test":           true,
		"a":              true,
		"123":            true,
		"test-123-slug":  true,
		"O-invalid":      false,
		"T-invalid":      false,
		"test/slash":     false,
		"test..dotdot":   false,
		"-leadinghyphen": false,
		"traininghyphen-": false,
		"":               false,
	}

	for slug, expected := range tests {
		result := validSlug(slug)
		if result != expected {
			t.Fatalf("validSlug(%q): expected %v, got %v", slug, expected, result)
		}
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
