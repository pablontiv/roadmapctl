package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/testutil"
)

func TestMaterializeApplyWritesFilesAndRunsPostcheck(t *testing.T) {
	fixture := copyFixture(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--apply", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	if report["applied"] != true {
		t.Fatalf("applied = %v; report = %#v", report["applied"], report)
	}
	for _, rel := range []string{"O02-new-outcome/README.md", "O02-new-outcome/T001-first-task.md", "T001-direct-task.md"} {
		if strings.HasSuffix(rel, "-tasks.md") {
			t.Fatalf("forbidden summary file path: %s", rel)
		}
		if _, err := os.Stat(filepath.Join(fixture, "docs", "roadmap", filepath.FromSlash(rel))); err != nil {
			t.Fatalf("expected applied file %s: %v", rel, err)
		}
	}
}

func TestMaterializeDryRunJSONDoesNotWrite(t *testing.T) {
	fixture := testutil.FixturePath(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	if report["kind"] != "roadmapctl/materialize" {
		t.Fatalf("kind = %v", report["kind"])
	}
	changes, ok := report["changes"].([]any)
	if !ok || len(changes) != 3 {
		t.Fatalf("changes = %#v", report["changes"])
	}
	wantPaths := []string{"O02-new-outcome/README.md", "O02-new-outcome/T001-first-task.md", "T001-direct-task.md"}
	for i, want := range wantPaths {
		change, _ := changes[i].(map[string]any)
		if change["path"] != want || change["operation"] != "create" || change["applied"] != false {
			t.Fatalf("change[%d] = %#v, want create %s applied=false", i, change, want)
		}
		if _, err := os.Stat(filepath.Join(fixture, "docs", "roadmap", filepath.FromSlash(want))); !os.IsNotExist(err) {
			t.Fatalf("dry-run wrote %s", want)
		}
	}
}
