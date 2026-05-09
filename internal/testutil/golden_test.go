package testutil

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeJSONReplacesPathsAndRootlineVersion(t *testing.T) {
	input := []byte(`{"path":"/tmp/repo/docs/roadmap","rootline_version":"rootline version test","items":["/tmp/repo/a"]}`)

	out := string(NormalizeJSON(t, input, map[string]string{"/tmp/repo": "<repo>"}))

	for _, want := range []string{"<repo>/docs/roadmap", "<rootline-version>", "<repo>/a"} {
		if !strings.Contains(out, want) {
			t.Fatalf("normalized JSON missing %q:\n%s", want, out)
		}
	}
}

func TestNormalizePathStringUsesLongestReplacementFirst(t *testing.T) {
	got := NormalizePathString("/tmp/repo/nested/file", map[string]string{
		"/tmp/repo":        "<repo>",
		"/tmp/repo/nested": "<nested>",
	})
	if got != "<nested>/file" {
		t.Fatalf("NormalizePathString = %q, want <nested>/file", got)
	}
}

func TestGoldenHelpersSuccessPaths(t *testing.T) {
	dir := t.TempDir()
	golden := filepath.Join(dir, "report.json")
	data := []byte(`{"root":"/tmp/repo","diagnostics":[{"id":"RMC_TEST"}]}`)
	normalized := NormalizeJSON(t, data, map[string]string{"/tmp/repo": "<repo>"})
	if err := os.WriteFile(golden, normalized, 0o644); err != nil {
		t.Fatal(err)
	}

	AssertGoldenJSON(t, golden, data, map[string]string{"/tmp/repo": "<repo>"})
	report := DecodeJSON(t, data)
	RequireDiagnosticID(t, report, "RMC_TEST")
	AssertNoBackslashes(t, report)
	AssertExit(t, 0, 0, &bytes.Buffer{}, &bytes.Buffer{})

	if got := GoldenPath("example.json"); !strings.HasSuffix(filepath.ToSlash(got), "testdata/golden/example.json") {
		t.Fatalf("GoldenPath = %q", got)
	}
	if got := FixturePath(t, "valid-direct-tasks"); !strings.HasSuffix(filepath.ToSlash(got), "testdata/fixtures/valid-direct-tasks") {
		t.Fatalf("FixturePath = %q", got)
	}
}

func TestContainsBackslashDetectsNestedValues(t *testing.T) {
	value := map[string]any{"items": []any{map[string]any{"path": `docs\\roadmap`}}}
	if !containsBackslash(value) {
		t.Fatal("containsBackslash = false, want true")
	}
	if containsBackslash(map[string]any{"path": "docs/roadmap"}) {
		t.Fatal("containsBackslash = true, want false")
	}
}
