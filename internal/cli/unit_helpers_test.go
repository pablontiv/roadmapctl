package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/config"
)

// --- sortDecisionItems ---

func TestSortDecisionItemsByScoreDescending(t *testing.T) {
	items := []decisionItem{
		{Path: "z", Score: 1},
		{Path: "a", Score: 10},
		{Path: "m", Score: 5},
	}
	sortDecisionItems(items)
	if items[0].Score != 10 || items[1].Score != 5 || items[2].Score != 1 {
		t.Fatalf("wrong order: %v", items)
	}
}

func TestSortDecisionItemsTieBreaksByPath(t *testing.T) {
	items := []decisionItem{
		{Path: "z/T001.md", Score: 10},
		{Path: "a/T001.md", Score: 10},
	}
	sortDecisionItems(items)
	if items[0].Path != "a/T001.md" {
		t.Fatalf("expected a/ first, got %q", items[0].Path)
	}
}

// --- transitionStatusValue ---

func TestTransitionStatusValueMapsAllRoles(t *testing.T) {
	cfg := &config.Config{
		StatusValues: config.StatusValues{
			Pending:    "Pending",
			Specified:  "Specified",
			InProgress: "In Progress",
			Completed:  "Completed",
			Blocked:    "Blocked",
			Obsolete:   "Obsolete",
		},
	}
	cases := []struct{ input, want string }{
		{"pending", "Pending"},
		{"specified", "Specified"},
		{"in-progress", "In Progress"},
		{"in_progress", "In Progress"},
		{"completed", "Completed"},
		{"blocked", "Blocked"},
		{"obsolete", "Obsolete"},
		{"CustomStatus", "CustomStatus"},
	}
	for _, c := range cases {
		got := transitionStatusValue(cfg, c.input)
		if got != c.want {
			t.Errorf("transitionStatusValue(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

// --- normalizeTransitionPath ---

func TestNormalizeTransitionPathStripsLeadingDotSlash(t *testing.T) {
	got := normalizeTransitionPath("", "./O01/T001.md")
	if got != "O01/T001.md" {
		t.Fatalf("got %q", got)
	}
}

func TestNormalizeTransitionPathAbsoluteRelativizedAgainstRoadmapRoot(t *testing.T) {
	root := t.TempDir()
	abs := filepath.Join(root, "O01", "T001.md")
	got := normalizeTransitionPath(root, abs)
	if got != "O01/T001.md" {
		t.Fatalf("got %q", got)
	}
}

func TestNormalizeTransitionPathAbsoluteOutsideRootReturnedClean(t *testing.T) {
	abs := filepath.Join(os.TempDir(), "other", "T001.md")
	got := normalizeTransitionPath("", abs)
	want := filepath.ToSlash(filepath.Clean(abs))
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// --- contextSchemaValues ---

func TestContextSchemaValuesReturnsTopLevelValuesForEstado(t *testing.T) {
	decoded := map[string]any{
		"values": []any{"Pending", "In Progress", "Completed"},
	}
	got := contextSchemaValues(decoded, "estado")
	if len(got) != 3 || got[0] != "Pending" {
		t.Fatalf("got %v", got)
	}
}

func TestContextSchemaValuesReturnsFieldSchemaForOtherFields(t *testing.T) {
	decoded := map[string]any{
		"values": []any{"X"},
		"schema": map[string]any{
			"tipo": map[string]any{
				"values": []any{"Feature", "Bug"},
			},
		},
	}
	got := contextSchemaValues(decoded, "tipo")
	if len(got) != 2 || got[0] != "Feature" {
		t.Fatalf("got %v", got)
	}
}

func TestContextSchemaValuesEmptyWhenFieldMissing(t *testing.T) {
	got := contextSchemaValues(map[string]any{}, "estado")
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

// --- configSource ---

func TestConfigSourceReturnsTOMLWhenFileExists(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".roadmapctl.toml")
	if err := os.WriteFile(cfgPath, []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{ConfigPath: cfgPath}
	if got := configSource(cfg); got != "toml" {
		t.Fatalf("got %q, want toml", got)
	}
}

func TestConfigSourceReturnsDefaultsWhenFileMissing(t *testing.T) {
	cfg := &config.Config{ConfigPath: filepath.Join(t.TempDir(), "nonexistent.toml")}
	if got := configSource(cfg); got != "defaults" {
		t.Fatalf("got %q, want defaults", got)
	}
}
