package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/testutil"
)

func TestCheckAndLintReportStaleOutcomeEstadoStem(t *testing.T) {
	for _, command := range []string{"check", "lint"} {
		t.Run(command, func(t *testing.T) {
			fixture := copyFixture(t, "valid-outcome-with-tasks")
			writeStaleOutcomeEstadoStem(t, fixture)

			var stdout, stderr bytes.Buffer
			code := Execute([]string{command, "--repo", fixture, "--output", "json"}, &stdout, &stderr, "dev")
			testutil.AssertExit(t, code, 1, &stdout, &stderr)
			report := testutil.DecodeJSON(t, stdout.Bytes())
			testutil.RequireDiagnosticID(t, report, "RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED")
			testutil.RequireDiagnosticID(t, report, "RMC_LINT_SCHEMA_OUTCOME_ESTADO_NON_EMPTY")
		})
	}
}

func TestDoctorReportsStaleOutcomeEstadoStem(t *testing.T) {
	fixture := copyFixture(t, "valid-outcome-with-tasks")
	writeStaleOutcomeEstadoStem(t, fixture)

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"doctor", "--repo", fixture, "--output", "json"}, &stdout, &stderr, "dev")
	testutil.AssertExit(t, code, 1, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	testutil.RequireDiagnosticID(t, report, "RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED")
}

func TestBootstrapInitApplyBlocksStaleStemBeforeWritingAdjacentFiles(t *testing.T) {
	repo := t.TempDir()
	initGitRepo(t, repo)
	roadmapRoot := filepath.Join(repo, "docs", "roadmap")
	if err := os.MkdirAll(roadmapRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	writeStaleStemFile(t, filepath.Join(roadmapRoot, ".stem"))

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"bootstrap", "init", "--repo", repo, "--roadmap-root", "docs/roadmap", "--apply", "--output", "json"}, &stdout, &stderr, "dev")
	testutil.AssertExit(t, code, 1, &stdout, &stderr)
	report := testutil.DecodeJSON(t, stdout.Bytes())
	testutil.RequireDiagnosticID(t, report, "RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED")
	if _, err := os.Stat(filepath.Join(roadmapRoot, ".roadmapctl.toml")); !os.IsNotExist(err) {
		t.Fatalf("bootstrap wrote adjacent config despite stale stem: %v", err)
	}
}

func writeStaleOutcomeEstadoStem(t *testing.T, repo string) {
	t.Helper()
	writeStaleStemFile(t, filepath.Join(repo, "docs", "roadmap", ".stem"))
}

func writeStaleStemFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(staleOutcomeEstadoStem), 0o644); err != nil {
		t.Fatal(err)
	}
}

const staleOutcomeEstadoStem = `version: 2
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

links:
  blocked_by:
    target: '^(\./|\.\./|.*/)T[0-9]{3}-[^/]+\.md$'

validate:
  - field: estado
    rule: non_empty
  - field: tipo
    rule: non_empty
`
