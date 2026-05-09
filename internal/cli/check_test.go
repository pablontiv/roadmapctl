package cli

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestCheckInvalidCycleExitsValidation(t *testing.T) {
	code, report, stderr := runCheckJSON(t, "invalid-cycle")
	if code != 1 {
		t.Fatalf("check exit = %d, want 1; stderr=%q report=%#v", code, stderr, report)
	}
	assertDiagnostic(t, report, "RMC_GRAPH_CYCLE")
}

func TestCheckBrokenBlockedByExitsValidation(t *testing.T) {
	code, report, stderr := runCheckJSON(t, "invalid-broken-blocked-by")
	if code != 1 {
		t.Fatalf("check exit = %d, want 1; stderr=%q report=%#v", code, stderr, report)
	}
	assertDiagnostic(t, report, "RMC_GRAPH_INVALID_BLOCKED_BY")
}

func TestCheckStatusMismatchExitsValidation(t *testing.T) {
	code, report, stderr := runCheckJSON(t, "invalid-status-mismatch")
	if code != 1 {
		t.Fatalf("check exit = %d, want 1; stderr=%q report=%#v", code, stderr, report)
	}
	assertDiagnostic(t, report, "RMC_STATUS_UNKNOWN")
}

func TestCheckAllowsStatusDeclaredOnlyInStem(t *testing.T) {
	code, report, stderr := runCheckJSON(t, "valid-status-on-hold")
	if code != 0 {
		t.Fatalf("check exit = %d, want 0; stderr=%q report=%#v", code, stderr, report)
	}
}

func TestCheckConfigRoleOutsideSchemaExitsValidation(t *testing.T) {
	code, report, stderr := runCheckJSON(t, "invalid-config-role-not-in-schema")
	if code != 1 {
		t.Fatalf("check exit = %d, want 1; stderr=%q report=%#v", code, stderr, report)
	}
	assertDiagnostic(t, report, "RMC_CONFIG_STATUS_SCHEMA_MISMATCH")
}

func TestCheckMissingRootlineExitsEnvironment(t *testing.T) {
	var stdout, stderr bytes.Buffer
	missingRootline := filepath.Join(t.TempDir(), "missing-rootline")
	code := Execute([]string{"check", "--repo", doctorFixturePath("valid-outcome-with-tasks"), "--rootline", missingRootline, "--output", "json"}, &stdout, &stderr)
	if code != 3 {
		t.Fatalf("check exit = %d, want 3; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	report := decodeReport(t, stdout.Bytes())
	assertDiagnostic(t, report, "RMC_ENV_ROOTLINE_MISSING")
}

func runCheckJSON(t *testing.T, fixture string) (int, doctorReport, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"check", "--repo", doctorFixturePath(fixture), "--output", "json"}, &stdout, &stderr)
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	return code, decodeReport(t, stdout.Bytes()), stderr.String()
}
