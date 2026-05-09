package diff

import (
	"strings"
	"testing"
)

func TestNewFileFormatsAddedLines(t *testing.T) {
	got := NewFile("docs/file.md", "first\nsecond")
	for _, want := range []string{"--- /dev/null", "+++ b/docs/file.md", "+first\n", "+second\n"} {
		if !strings.Contains(got, want) {
			t.Fatalf("diff missing %q:\n%s", want, got)
		}
	}
}
