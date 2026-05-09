package cli

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestContextJSONIncludesEffectiveHelpersAndLegacySource(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Execute([]string{"context", "--repo", doctorFixturePath("valid-legacy-config-fallback"), "--output", "json"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("context exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var report struct {
		Kind         string `json:"kind"`
		ConfigSource string `json:"config_source"`
		Helpers      struct {
			WhereLeaf    string `json:"where_leaf"`
			WhereNotDone string `json:"where_not_done"`
			WhereActive  string `json:"where_active"`
		} `json:"helpers"`
		Schema struct {
			Estado []string `json:"estado"`
			Tipo   []string `json:"tipo"`
		} `json:"schema"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout is not parseable JSON: %v\n%s", err, stdout.String())
	}
	if report.Kind != "roadmapctl/context" {
		t.Fatalf("Kind = %q", report.Kind)
	}
	if report.ConfigSource != "legacy" {
		t.Fatalf("ConfigSource = %q", report.ConfigSource)
	}
	if report.Helpers.WhereLeaf != "isIndex == false" || report.Helpers.WhereNotDone != `not (estado in ["Completed", "Obsolete"])` || report.Helpers.WhereActive != `estado in ["Pending", "Specified", "In Progress"]` {
		t.Fatalf("helpers = %#v", report.Helpers)
	}
	if len(report.Schema.Estado) == 0 || len(report.Schema.Tipo) == 0 {
		t.Fatalf("schema = %#v", report.Schema)
	}
}
