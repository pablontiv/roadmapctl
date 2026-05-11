package diff

import "strings"

func NewFile(path string, content string) string {
	var b strings.Builder
	b.WriteString("--- /dev/null\n")
	b.WriteString("+++ b/")
	b.WriteString(path)
	b.WriteString("\n")
	writePrefixedLines(&b, '+', content)
	return b.String()
}

func UpdateFile(path string, previous string, content string) string {
	var b strings.Builder
	b.WriteString("--- a/")
	b.WriteString(path)
	b.WriteString("\n")
	b.WriteString("+++ b/")
	b.WriteString(path)
	b.WriteString("\n")
	writePrefixedLines(&b, '-', previous)
	writePrefixedLines(&b, '+', content)
	return b.String()
}

func writePrefixedLines(b *strings.Builder, prefix byte, content string) {
	for _, line := range strings.SplitAfter(content, "\n") {
		if line == "" {
			continue
		}
		b.WriteByte(prefix)
		b.WriteString(line)
		if !strings.HasSuffix(line, "\n") {
			b.WriteByte('\n')
		}
	}
}
