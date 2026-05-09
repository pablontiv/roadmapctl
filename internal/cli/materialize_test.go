package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/testutil"
)

func TestMaterializeUsageErrors(t *testing.T) {
	for _, args := range [][]string{
		{"materialize", "--dry-run", "--repo", testutil.FixturePath(t, "valid-outcome-with-tasks"), "--output", "json"},
		{"materialize", "--plan", filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json"), "--repo", testutil.FixturePath(t, "valid-outcome-with-tasks"), "--output", "json"},
		{"materialize", "--plan", filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json"), "--dry-run", "--apply", "--repo", testutil.FixturePath(t, "valid-outcome-with-tasks"), "--output", "json"},
		{"materialize", "--plan", filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json"), "--target", "T001-direct-task.md", "--apply", "--repo", copyFixture(t, "valid-outcome-with-tasks"), "--output", "json"},
	} {
		var stdout, stderr bytes.Buffer
		code := Execute(args, &stdout, &stderr)
		testutil.AssertExit(t, code, 2, &stdout, &stderr)
	}
}

func TestMaterializeTextOutputIncludesDiff(t *testing.T) {
	fixture := testutil.FixturePath(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", fixture, "--output", "text"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	if !strings.Contains(stdout.String(), "# create O02-new-outcome/README.md") || !strings.Contains(stdout.String(), "+++ b/O02-new-outcome/README.md") {
		t.Fatalf("text output missing diff:\n%s", stdout.String())
	}
}

func TestMaterializeDryRunGolden(t *testing.T) {
	fixture := testutil.FixturePath(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	testutil.AssertGoldenJSON(t, testutil.GoldenPath("materialize-dry-run-outcome-and-direct.json"), stdout.Bytes(), map[string]string{
		absoluteFixturePath(t, "valid-outcome-with-tasks"): "<fixture:valid-outcome-with-tasks>",
	})
}

func TestMaterializeDryRunMissingRootShowsBootstrap(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", repo, "--roadmap-root", "docs/roadmap", "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	changes, ok := report["changes"].([]any)
	if !ok || len(changes) < 3 {
		t.Fatalf("changes = %#v", report["changes"])
	}
	for i, want := range []string{".", ".stem", ".roadmapctl.toml"} {
		change, _ := changes[i].(map[string]any)
		if change["path"] != want || change["applied"] != false {
			t.Fatalf("bootstrap change[%d] = %#v, want %s applied=false", i, change, want)
		}
	}
	if _, err := os.Stat(filepath.Join(repo, "docs", "roadmap")); !os.IsNotExist(err) {
		t.Fatalf("dry-run created roadmap root: %v", err)
	}
}

func TestMaterializeDependenciesUseExplicitRelativeLinks(t *testing.T) {
	fixture := testutil.FixturePath(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "dependencies-same-cross-outcome.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	if !strings.Contains(stdout.String(), "[[blocked_by:./T001-first.md]]") {
		t.Fatalf("missing same-outcome explicit dependency link:\n%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "[[blocked_by:../O02-alpha/T001-first.md]]") {
		t.Fatalf("missing cross-outcome explicit dependency link:\n%s", stdout.String())
	}
}

func TestMaterializeInvalidInputDoesNotWriteFiles(t *testing.T) {
	fixture := copyFixture(t, "valid-outcome-with-tasks")
	before := listRoadmapFiles(t, fixture)
	plan := filepath.Join("..", "..", "testdata", "plans", "invalid-path-escape.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--apply", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 1, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	testutil.RequireDiagnosticID(t, report, "RMC_MATERIALIZE_INPUT_SLUG_INVALID")
	after := listRoadmapFiles(t, fixture)
	if !bytes.Equal([]byte(before), []byte(after)) {
		t.Fatalf("invalid input wrote files\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

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

func TestMaterializeChangesApplyWritesWholeFrozenChangeSet(t *testing.T) {
	fixture := copyFixture(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	changesPath := filepath.Join(t.TempDir(), "dry-run.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	if err := os.WriteFile(changesPath, stdout.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	code = Execute([]string{"materialize", "--changes", changesPath, "--apply", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	changes, _ := report["changes"].([]any)
	if report["applied"] != true || len(changes) != 3 {
		t.Fatalf("batch apply report = %#v", report)
	}
	for _, rel := range []string{"O02-new-outcome/README.md", "O02-new-outcome/T001-first-task.md", "T001-direct-task.md"} {
		if _, err := os.Stat(filepath.Join(fixture, "docs", "roadmap", filepath.FromSlash(rel))); err != nil {
			t.Fatalf("batch target missing %s: %v", rel, err)
		}
	}
}

func TestMaterializeChangesApplyReportsConflictPathBeforeWriting(t *testing.T) {
	fixture := copyFixture(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	changesPath := filepath.Join(t.TempDir(), "dry-run.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	if err := os.WriteFile(changesPath, stdout.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture, "docs", "roadmap", "T001-direct-task.md"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	code = Execute([]string{"materialize", "--changes", changesPath, "--apply", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 1, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	testutil.RequireDiagnosticID(t, report, "RMC_MATERIALIZE_PLAN_CONFLICT")
	if !strings.Contains(stdout.String(), "T001-direct-task.md") {
		t.Fatalf("conflict report missing concrete path:\n%s", stdout.String())
	}
	if _, err := os.Stat(filepath.Join(fixture, "docs", "roadmap", "O02-new-outcome", "README.md")); !os.IsNotExist(err) {
		t.Fatalf("batch conflict wrote sibling before failing: %v", err)
	}
}

func TestMaterializeTargetApplyWritesOnlySelectedDryRunChange(t *testing.T) {
	fixture := copyFixture(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	changesPath := filepath.Join(t.TempDir(), "dry-run.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	if err := os.WriteFile(changesPath, stdout.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	stdout.Reset()
	stderr.Reset()
	code = Execute([]string{"materialize", "--changes", changesPath, "--target", "T001-direct-task.md", "--apply", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	changes, _ := report["changes"].([]any)
	if report["applied"] != true || len(changes) != 1 {
		t.Fatalf("target apply report = %#v", report)
	}
	if _, err := os.Stat(filepath.Join(fixture, "docs", "roadmap", "T001-direct-task.md")); err != nil {
		t.Fatalf("selected target missing: %v", err)
	}
	for _, sibling := range []string{"O02-new-outcome/README.md", "O02-new-outcome/T001-first-task.md"} {
		if _, err := os.Stat(filepath.Join(fixture, "docs", "roadmap", filepath.FromSlash(sibling))); !os.IsNotExist(err) {
			t.Fatalf("sibling %s was written: %v", sibling, err)
		}
	}
}

func TestMaterializeTargetApplyRejectsUnknownTargetBeforeWriting(t *testing.T) {
	fixture := copyFixture(t, "valid-outcome-with-tasks")
	plan := filepath.Join("..", "..", "testdata", "plans", "outcome-and-direct.json")
	changesPath := filepath.Join(t.TempDir(), "dry-run.json")
	var stdout, stderr bytes.Buffer

	code := Execute([]string{"materialize", "--plan", plan, "--dry-run", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)
	if err := os.WriteFile(changesPath, stdout.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}

	before := listRoadmapFiles(t, fixture)
	stdout.Reset()
	stderr.Reset()
	code = Execute([]string{"materialize", "--changes", changesPath, "--target", "T999-missing.md", "--apply", "--repo", fixture, "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 1, &stdout, &stderr)
	after := listRoadmapFiles(t, fixture)
	if !bytes.Equal([]byte(before), []byte(after)) {
		t.Fatalf("unknown target wrote files\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func listRoadmapFiles(t *testing.T, fixture string) string {
	t.Helper()
	var files []string
	root := filepath.Join(fixture, "docs", "roadmap")
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
	sort.Strings(files)
	return strings.Join(files, "\n")
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
