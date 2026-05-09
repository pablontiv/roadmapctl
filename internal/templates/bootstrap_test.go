package templates

import (
	"strings"
	"testing"
)

func TestDefaultRoadmapctlTOMLIncludesRequiredCodeCoverage(t *testing.T) {
	if want := "required_code_coverage = 85.0"; !strings.Contains(DefaultRoadmapctlTOML, want) {
		t.Fatalf("DefaultRoadmapctlTOML missing %q:\n%s", want, DefaultRoadmapctlTOML)
	}
}
