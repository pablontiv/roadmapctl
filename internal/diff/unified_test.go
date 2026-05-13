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

func TestUpdateFileFormatsPreviousAndNewLines(t *testing.T) {
	got := UpdateFile("docs/file.md", "old\n", "new line\nsecond")
	for _, want := range []string{"--- a/docs/file.md", "+++ b/docs/file.md", "-old\n", "+new line\n", "+second\n"} {
		if !strings.Contains(got, want) {
			t.Fatalf("diff missing %q:\n%s", want, got)
		}
	}
}

func TestNewFileIgnoresTrailingEmptyFromSplit(t *testing.T) {
	got := NewFile("docs/file.md", "line\n")
	if strings.Count(got, "+line\n") != 1 {
		t.Fatalf("expected exactly one +line, got:\n%s", got)
	}
}
