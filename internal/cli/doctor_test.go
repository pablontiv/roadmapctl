package cli

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorValidOutcomeFixtureJSONStatusOK(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"doctor", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--output", "json"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("doctor exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}

	report := decodeReport(t, stdout.Bytes())
	if report.Kind != "roadmapctl/doctor" {
		t.Fatalf("Kind = %q, want roadmapctl/doctor", report.Kind)
	}
	if report.Summary.Status != "ok" || report.Summary.Errors != 0 {
		t.Fatalf("Summary = %#v, want ok with zero errors", report.Summary)
	}
	if report.Root == "" || report.RoadmapRoot == "" {
		t.Fatalf("Root/RoadmapRoot must be populated: %#v", report)
	}
}

func TestDoctorMissingRootlineReturnsEnvironmentDiagnostic(t *testing.T) {
	var stdout, stderr bytes.Buffer
	missingRootline := filepath.Join(t.TempDir(), "missing-rootline")
	code := Execute([]string{"doctor", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--rootline", missingRootline, "--output", "json"}, &stdout, &stderr)

	if code != 3 {
		t.Fatalf("doctor exit = %d, want 3; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	report := decodeReport(t, stdout.Bytes())
	assertDiagnostic(t, report, "RMC_ENV_ROOTLINE_MISSING")
}

func TestDoctorMissingConfigReportsConfigMissing(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"doctor", "--repo", t.TempDir(), "--output", "json"}, &stdout, &stderr)

	if code != 2 {
		t.Fatalf("doctor exit = %d, want 2; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	report := decodeReport(t, stdout.Bytes())
	assertDiagnostic(t, report, "RMC_CONFIG_MISSING")
}

func TestDoctorJSONStdoutIsParseableWithoutExtraText(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"doctor", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--output", "json"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("doctor exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if strings.Contains(stdout.String(), "not implemented") || strings.Count(strings.TrimSpace(stdout.String()), "\n") > 0 {
		t.Fatalf("stdout contains extra text: %q", stdout.String())
	}
	_ = decodeReport(t, stdout.Bytes())
}

type doctorReport struct {
	Kind        string `json:"kind"`
	Root        string `json:"root"`
	RoadmapRoot string `json:"roadmap_root"`
	Summary     struct {
		Status string `json:"status"`
		Errors int    `json:"errors"`
	} `json:"summary"`
	Diagnostics []struct {
		ID string `json:"id"`
	} `json:"diagnostics"`
}

func doctorFixturePath(name string) string {
	return filepath.Join("..", "..", "testdata", "fixtures", name)
}

func decodeReport(t *testing.T, data []byte) doctorReport {
	t.Helper()
	var report doctorReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatalf("stdout is not parseable JSON report: %v\n%s", err, string(data))
	}
	return report
}

func assertDiagnostic(t *testing.T, report doctorReport, id string) {
	t.Helper()
	for _, diagnostic := range report.Diagnostics {
		if diagnostic.ID == id {
			return
		}
	}
	t.Fatalf("missing diagnostic %q in %#v", id, report.Diagnostics)
}
