package diff

import "strings"

func NewFile(path string, content string) string {
	var b strings.Builder
	b.WriteString("--- /dev/null\n")
	b.WriteString("+++ b/")
	b.WriteString(path)
	b.WriteString("\n")
	for _, line := range strings.SplitAfter(content, "\n") {
		if line == "" {
			continue
		}
		b.WriteByte('+')
		b.WriteString(line)
		if !strings.HasSuffix(line, "\n") {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
