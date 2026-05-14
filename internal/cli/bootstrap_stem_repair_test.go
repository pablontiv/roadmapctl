package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/templates"
	"github.com/pablontiv/roadmapctl/internal/testutil"
)

// setupRepairRepo creates a temp repo with git init, the docs/roadmap directory,
// and the given .stem content. Returns the repo path.
func setupRepairRepo(t *testing.T, stemContent string) (repo string) {
	t.Helper()
	repo = t.TempDir()
	initGitRepo(t, repo)
	roadmapRoot := filepath.Join(repo, "docs", "roadmap")
	if err := os.MkdirAll(roadmapRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(roadmapRoot, ".stem"), []byte(stemContent), 0o644); err != nil {
		t.Fatal(err)
	}
	return repo
}

func assertRepairExitCode(t *testing.T, code int, want int, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	t.Helper()
	if code != want {
		t.Fatalf("exit = %d, want %d\nstdout:\n%s\nstderr:\n%s", code, want, stdout.String(), stderr.String())
	}
}

// customStemWithExtraField is a stale stem that also has an unrecognized schema field.
const customStemWithExtraField = `version: 2
scope:
  match: "*.md"

schema:
  estado:
    type: enum
    required:
      match: ["O*", "T*"]
    match: ["O*", "T*"]
    values: [Pending, Specified, In Progress, Completed, Blocked, On Hold, Obsolete]

  tipo:
    type: enum
    required:
      match: ["O*", "T*"]
    match: ["O*", "T*"]
    values: [outcome, task]

  priority:
    type: enum
    values: [low, medium, high]

validate:
  - field: estado
    rule: non_empty
  - field: tipo
    rule: non_empty
`

func TestBootstrapRepairDetectsAndPrompts(t *testing.T) {
	requiresRealRootline(t)
	repo := setupRepairRepo(t, staleOutcomeEstadoStem)

	// stdin closed immediately → treated as "N" (no confirmation)
	var stdout, stderr bytes.Buffer
	code := ExecuteWithStdin(
		[]string{"bootstrap", "--repo", repo, "--output", "json"},
		bytes.NewReader(nil),
		&stdout, &stderr, "dev",
	)

	// Should exit non-zero because stem is still stale
	assertRepairExitCode(t, code, 1, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	testutil.RequireDiagnosticID(t, report, "RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED")

	// Diff and prompt must appear in stderr
	stderrStr := stderr.String()
	if !strings.Contains(stderrStr, "--- current .stem") {
		t.Fatalf("stderr missing diff header; got:\n%s", stderrStr)
	}
	if !strings.Contains(stderrStr, "+++ canonical .stem") {
		t.Fatalf("stderr missing canonical header; got:\n%s", stderrStr)
	}
	if !strings.Contains(stderrStr, "Update .stem to canonical schema? [y/N]") {
		t.Fatalf("stderr missing prompt; got:\n%s", stderrStr)
	}

	// .stem must be unchanged
	content, err := os.ReadFile(filepath.Join(repo, "docs", "roadmap", ".stem"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != staleOutcomeEstadoStem {
		t.Fatalf(".stem was modified when user said N")
	}
}

func TestBootstrapRepairAppliesOnYesInput(t *testing.T) {
	requiresRealRootline(t)
	repo := setupRepairRepo(t, staleOutcomeEstadoStem)

	var stdout, stderr bytes.Buffer
	code := ExecuteWithStdin(
		[]string{"bootstrap", "--repo", repo, "--output", "json"},
		strings.NewReader("y\n"),
		&stdout, &stderr, "dev",
	)

	assertRepairExitCode(t, code, 0, &stdout, &stderr)

	content, err := os.ReadFile(filepath.Join(repo, "docs", "roadmap", ".stem"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != templates.BaseStemContent {
		t.Fatalf(".stem not updated to canonical:\n%s", string(content))
	}
}

func TestBootstrapRepairAppliesWithYesFlag(t *testing.T) {
	requiresRealRootline(t)
	repo := setupRepairRepo(t, staleOutcomeEstadoStem)

	var stdout, stderr bytes.Buffer
	code := ExecuteWithStdin(
		[]string{"bootstrap", "--repo", repo, "--yes", "--output", "json"},
		bytes.NewReader(nil),
		&stdout, &stderr, "dev",
	)

	assertRepairExitCode(t, code, 0, &stdout, &stderr)

	content, err := os.ReadFile(filepath.Join(repo, "docs", "roadmap", ".stem"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != templates.BaseStemContent {
		t.Fatalf(".stem not updated to canonical:\n%s", string(content))
	}

	// Prompt must NOT appear in stderr when --yes is set
	if strings.Contains(stderr.String(), "Update .stem to canonical schema?") {
		t.Fatalf("--yes flag did not suppress interactive prompt")
	}
}

func TestBootstrapRepairDoesNotModifyOnNoInput(t *testing.T) {
	requiresRealRootline(t)
	repo := setupRepairRepo(t, staleOutcomeEstadoStem)

	var stdout, stderr bytes.Buffer
	code := ExecuteWithStdin(
		[]string{"bootstrap", "--repo", repo, "--output", "json"},
		strings.NewReader("N\n"),
		&stdout, &stderr, "dev",
	)

	assertRepairExitCode(t, code, 1, &stdout, &stderr)

	content, err := os.ReadFile(filepath.Join(repo, "docs", "roadmap", ".stem"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != staleOutcomeEstadoStem {
		t.Fatalf(".stem was modified despite user answering N")
	}
}

func TestBootstrapRepairUnsupportedCustomStem(t *testing.T) {
	requiresRealRootline(t)
	repo := setupRepairRepo(t, customStemWithExtraField)

	var stdout, stderr bytes.Buffer
	code := ExecuteWithStdin(
		[]string{"bootstrap", "--repo", repo, "--yes", "--output", "json"},
		bytes.NewReader(nil),
		&stdout, &stderr, "dev",
	)

	assertRepairExitCode(t, code, 1, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	testutil.RequireDiagnosticID(t, report, "RMC_BOOTSTRAP_REPAIR_UNSUPPORTED_STEM")

	// .stem must be unchanged
	content, err := os.ReadFile(filepath.Join(repo, "docs", "roadmap", ".stem"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != customStemWithExtraField {
		t.Fatalf(".stem was modified for unsupported custom stem")
	}
}

func TestBootstrapRepairNotTriggeredForCanonicalStem(t *testing.T) {
	requiresRealRootline(t)
	repo := setupRepairRepo(t, templates.BaseStemContent)

	var stdout, stderr bytes.Buffer
	code := ExecuteWithStdin(
		[]string{"bootstrap", "--repo", repo, "--output", "json"},
		bytes.NewReader(nil),
		&stdout, &stderr, "dev",
	)

	assertRepairExitCode(t, code, 0, &stdout, &stderr)

	// No prompt should appear
	if strings.Contains(stderr.String(), "Update .stem to canonical schema?") {
		t.Fatalf("repair prompt appeared for canonical stem")
	}

	// .stem must remain canonical
	content, err := os.ReadFile(filepath.Join(repo, "docs", "roadmap", ".stem"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != templates.BaseStemContent {
		t.Fatalf(".stem was unexpectedly modified")
	}
}

func TestIsStemRecognizedLegacy(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name:    "stale legacy stem",
			content: staleOutcomeEstadoStem,
			want:    true,
		},
		{
			name:    "canonical stem",
			content: templates.BaseStemContent,
			want:    true,
		},
		{
			name:    "custom schema field",
			content: customStemWithExtraField,
			want:    false,
		},
		{
			name: "unknown top-level key",
			content: `version: 2
custom_section:
  foo: bar
schema:
  estado:
    type: enum
`,
			want: false,
		},
		{
			name:    "empty content",
			content: "",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isStemRecognizedLegacy(tt.content)
			if got != tt.want {
				t.Fatalf("isStemRecognizedLegacy = %v, want %v", got, tt.want)
			}
		})
	}
}
