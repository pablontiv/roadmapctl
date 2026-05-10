package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/testutil"
)

func TestMain(m *testing.M) {
	if os.Getenv("ROADMAPCTL_FAKE_ROOTLINE") == "1" {
		os.Exit(fakeRootline(os.Args[1:], os.Stdout, os.Stderr))
	}
	os.Exit(m.Run())
}

func TestCheckGoldenJSONFixtures(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		fixture    string
		wantExit   int
		wantID     string
		goldenName string
	}{
		{name: "invalid single summary file", command: "check", fixture: "invalid-single-summary-file", wantExit: 1, wantID: "RMC_STRUCTURE_SINGLE_FILE_FALLBACK", goldenName: "check-invalid-single-summary-file.json"},
		{name: "valid outcome with tasks", command: "check", fixture: "valid-outcome-with-tasks", wantExit: 0, goldenName: "check-valid-outcome-with-tasks.json"},
		{name: "valid status on hold", command: "check", fixture: "valid-status-on-hold", wantExit: 0, goldenName: "check-valid-status-on-hold.json"},
		{name: "invalid status bogus", command: "check", fixture: "invalid-status-bogus", wantExit: 1, wantID: "RMC_STATUS_UNKNOWN", goldenName: "check-invalid-status-bogus.json"},
		{name: "invalid config role not in schema", command: "check", fixture: "invalid-config-role-not-in-schema", wantExit: 1, wantID: "RMC_CONFIG_STATUS_SCHEMA_MISMATCH", goldenName: "check-invalid-config-role-not-in-schema.json"},
		{name: "invalid stale outcome stem", command: "check", fixture: "invalid-stale-outcome-stem", wantExit: 1, wantID: "RMC_LINT_SCHEMA_OUTCOME_ESTADO_REQUIRED", goldenName: "check-invalid-stale-outcome-stem.json"},
		{name: "missing stem", command: "doctor", fixture: "invalid-missing-stem", wantExit: 2, wantID: "RMC_CONFIG_STEM_MISSING", goldenName: "doctor-invalid-missing-stem.json"},
		{name: "bare blocked_by", command: "check", fixture: "invalid-bare-blocked-by", wantExit: 1, wantID: "RMC_ROOTLINE_VALIDATE_FAILED", goldenName: "check-invalid-bare-blocked-by.json"},
		{name: "root escape", command: "check", fixture: "invalid-root-escape", wantExit: 2, wantID: "RMC_CONFIG_ROADMAP_ROOT_ESCAPE", goldenName: "check-invalid-root-escape.json"},
		{name: "context legacy", command: "context", fixture: "valid-legacy-config-fallback", wantExit: 0, goldenName: "context-valid-legacy-config-fallback.json"},
		{name: "context workspace", command: "context", fixture: "valid-workspace", wantExit: 0, goldenName: "context-valid-workspace.json"},
		{name: "pending direct tasks", command: "pending", fixture: "valid-direct-tasks", wantExit: 0, goldenName: "pending-valid-direct-tasks.json"},
		{name: "pending outcome tasks", command: "pending", fixture: "valid-outcome-with-tasks", wantExit: 0, goldenName: "pending-valid-outcome-with-tasks.json"},
		{name: "pending none", command: "pending", fixture: "valid-no-pending", wantExit: 0, goldenName: "pending-valid-no-pending.json"},
		{name: "next ready blocked", command: "next", fixture: "valid-next-with-blocked", wantExit: 0, goldenName: "next-valid-next-with-blocked.json"},
		{name: "decision reverse dependencies", command: "decision", fixture: "valid-next-with-blocked", wantExit: 0, goldenName: "decision-valid-next-with-blocked.json"},
		{name: "lint valid", command: "lint", fixture: "lint-valid", wantExit: 0, goldenName: "lint-valid.json"},
		{name: "lint missing table row", command: "lint", fixture: "lint-missing-table-row", wantExit: 0, goldenName: "lint-missing-table-row.json"},
		{name: "lint stale table row", command: "lint", fixture: "lint-stale-table-row", wantExit: 0, goldenName: "lint-stale-table-row.json"},
		{name: "lint missing task sections", command: "lint", fixture: "lint-missing-task-sections", wantExit: 0, goldenName: "lint-missing-task-sections.json"},
		{name: "lint case collision", command: "lint", fixture: "lint-case-collision", wantExit: 1, goldenName: "lint-case-collision.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := testutil.FixturePath(t, tt.fixture)
			var stdout, stderr bytes.Buffer
			args := []string{tt.command, "--repo", fixture, "--output", "json"}
			if tt.command == "context" && strings.Contains(tt.fixture, "workspace") {
				args = append(args, "--workspace")
			}
			code := Execute(args, &stdout, &stderr)
			testutil.AssertExit(t, code, tt.wantExit, &stdout, &stderr)
			report := testutil.DecodeJSON(t, stdout.Bytes())
			if tt.wantID != "" {
				testutil.RequireDiagnosticID(t, report, tt.wantID)
			}
			testutil.AssertNoBackslashes(t, report)
			testutil.AssertGoldenJSON(t, testutil.GoldenPath(tt.goldenName), stdout.Bytes(), map[string]string{
				absoluteFixturePath(t, tt.fixture): fmt.Sprintf("<fixture:%s>", tt.fixture),
			})
		})
	}
}

func TestTransitionJSONGoldens(t *testing.T) {
	tests := []struct {
		name       string
		fixture    string
		args       []string
		goldenName string
	}{
		{name: "can start ready", fixture: "valid-next-with-blocked", args: []string{"transition", "can-start", "O01-work/T001-ready.md"}, goldenName: "transition-can-start-ready.json"},
		{name: "can start blocked", fixture: "valid-next-with-blocked", args: []string{"transition", "can-start", "O01-work/T002-blocked.md"}, goldenName: "transition-can-start-blocked.json"},
		{name: "dry run start", fixture: "valid-next-with-blocked", args: []string{"transition", "start", "O01-work/T001-ready.md", "--dry-run"}, goldenName: "transition-start-dry-run.json"},
		{name: "custom labels", fixture: "transition-custom-status", args: []string{"transition", "start", "O01-work/T001-ready.md", "--dry-run"}, goldenName: "transition-custom-status-start-dry-run.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := testutil.FixturePath(t, tt.fixture)
			args := append([]string{}, tt.args...)
			args = append(args, "--repo", fixture, "--output", "json")
			var stdout, stderr bytes.Buffer
			code := Execute(args, &stdout, &stderr)
			testutil.AssertExit(t, code, 0, &stdout, &stderr)
			report := testutil.DecodeJSON(t, stdout.Bytes())
			testutil.AssertNoBackslashes(t, report)
			testutil.AssertGoldenJSON(t, testutil.GoldenPath(tt.goldenName), stdout.Bytes(), map[string]string{
				absoluteFixturePath(t, tt.fixture): fmt.Sprintf("<fixture:%s>", tt.fixture),
			})
		})
	}
}

func TestLintStrictPromotesWarnings(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"lint", "--repo", testutil.FixturePath(t, "lint-missing-table-row"), "--output", "json", "--strict"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 1, &stdout, &stderr)
}

func TestReadOnlyTextGoldens(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		fixture    string
		goldenName string
	}{
		{name: "pending", command: "pending", fixture: "valid-outcome-with-tasks", goldenName: "pending-valid-outcome-with-tasks.txt"},
		{name: "next", command: "next", fixture: "valid-next-with-blocked", goldenName: "next-valid-next-with-blocked.txt"},
		{name: "decision", command: "decision", fixture: "valid-next-with-blocked", goldenName: "decision-valid-next-with-blocked.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Execute([]string{tt.command, "--repo", testutil.FixturePath(t, tt.fixture), "--output", "text"}, &stdout, &stderr)
			testutil.AssertExit(t, code, 0, &stdout, &stderr)
			want, err := os.ReadFile(testutil.GoldenPath(tt.goldenName))
			if err != nil {
				t.Fatalf("read text golden: %v\nactual:\n%s", err, stdout.String())
			}
			if !bytes.Equal(bytes.TrimSpace(want), bytes.TrimSpace(stdout.Bytes())) {
				t.Fatalf("text golden mismatch\nwant:\n%s\ngot:\n%s", bytes.TrimSpace(want), bytes.TrimSpace(stdout.Bytes()))
			}
		})
	}
}

func TestCheckUsesRootlineBinEnvironmentOverride(t *testing.T) {
	t.Setenv("ROADMAPCTL_FAKE_ROOTLINE", "1")
	t.Setenv("ROOTLINE_BIN", os.Args[0])

	var stdout, stderr bytes.Buffer
	code := Execute([]string{"check", "--repo", testutil.FixturePath(t, "valid-outcome-with-tasks"), "--output", "json"}, &stdout, &stderr)
	testutil.AssertExit(t, code, 0, &stdout, &stderr)

	report := testutil.DecodeJSON(t, stdout.Bytes())
	if summary, _ := report["summary"].(map[string]any); summary["status"] != "ok" {
		t.Fatalf("summary = %#v, want ok", summary)
	}
}

func absoluteFixturePath(t *testing.T, fixture string) string {
	t.Helper()
	abs, err := filepath.Abs(testutil.FixturePath(t, fixture))
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(filepath.Clean(abs))
}

func fakeRootline(args []string, stdout *os.File, stderr *os.File) int {
	if len(args) == 1 && args[0] == "--version" {
		fmt.Fprintln(stdout, "rootline version test-fake")
		return 0
	}
	if len(args) == 0 {
		fmt.Fprintln(stderr, "missing command")
		return 2
	}

	switch args[0] {
	case "validate":
		fmt.Fprintln(stdout, `{"version":1,"kind":"rootline/validate-batch","summary":{"total":0,"valid":0,"invalid":0,"errors_count":0,"warnings_count":0}}`)
		return 0
	case "describe":
		fmt.Fprintln(stdout, `{"type":"enum","values":["Pending","Specified","In Progress","Completed","Blocked","On Hold","Obsolete"]}`)
		return 0
	case "query":
		fmt.Fprintln(stdout, `{"version":1,"kind":"rootline/query","meta":{"count":0},"rows":[]}`)
		return 0
	case "graph":
		if strings.Contains(strings.Join(args, " "), "--check") {
			fmt.Fprintln(stderr, "fake rootline expected JSON graph, not --check")
			return 2
		}
		fmt.Fprintln(stdout, `{"version":1,"kind":"rootline/graph","nodes":[],"edges":[],"cycles":[],"broken_links":[]}`)
		return 0
	case "tree":
		fmt.Fprintln(stdout, `{"version":1,"kind":"rootline/tree","root":{"children":[]}}`)
		return 0
	case "set":
		fmt.Fprintln(stdout, `set estado = "Completed"`)
		return 0
	case "new":
		fmt.Fprintln(stdout, `created docs/roadmap/T001-task.md`)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown fake rootline command %q\n", args[0])
		return 2
	}
}
