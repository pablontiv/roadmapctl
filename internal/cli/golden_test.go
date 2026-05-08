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
		{name: "missing stem", command: "doctor", fixture: "invalid-missing-stem", wantExit: 2, wantID: "RMC_CONFIG_STEM_MISSING", goldenName: "doctor-invalid-missing-stem.json"},
		{name: "bare blocked_by", command: "check", fixture: "invalid-bare-blocked-by", wantExit: 1, wantID: "RMC_ROOTLINE_VALIDATE_FAILED", goldenName: "check-invalid-bare-blocked-by.json"},
		{name: "root escape", command: "check", fixture: "invalid-root-escape", wantExit: 2, wantID: "RMC_CONFIG_ROADMAP_ROOT_ESCAPE", goldenName: "check-invalid-root-escape.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := testutil.FixturePath(t, tt.fixture)
			var stdout, stderr bytes.Buffer
			code := Execute([]string{tt.command, "--repo", fixture, "--output", "json"}, &stdout, &stderr)
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
	default:
		fmt.Fprintf(stderr, "unknown fake rootline command %q\n", args[0])
		return 2
	}
}
