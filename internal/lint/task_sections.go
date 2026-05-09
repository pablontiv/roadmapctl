package lint

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

var requiredTaskSections = []string{
	"Preserva",
	"Contexto",
	"Alcance",
	"Estado inicial esperado",
	"Criterios de Aceptación",
	"Fuente de verdad",
}

var listItemPattern = regexp.MustCompile(`^\s*(?:[-*+]\s+|[0-9]+[.)]\s+)`)

func CheckTaskSections(roadmapRoot string) ([]diagnostics.Diagnostic, error) {
	root := filepath.Clean(roadmapRoot)
	var found []diagnostics.Diagnostic
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !taskFilenamePattern.MatchString(entry.Name()) {
			return nil
		}
		taskDiagnostics, err := checkTaskSections(root, path)
		if err != nil {
			return err
		}
		found = append(found, taskDiagnostics...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sortDiagnostics(found)
	return found, nil
}

func checkTaskSections(root string, path string) ([]diagnostics.Diagnostic, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	doc, err := ParseMarkdown(source)
	if err != nil {
		return nil, err
	}
	var found []diagnostics.Diagnostic
	for _, section := range requiredTaskSections {
		if !doc.HasHeading(section) {
			found = append(found, lintTaskSectionDiagnostic(diagnostics.DiagnosticLintTaskSectionMissing, root, path, "task is missing a required section heading", section))
		}
	}
	if section, ok := markdownSection(source, doc, "Criterios de Aceptación"); ok && !hasListItem(section) {
		found = append(found, lintTaskSectionDiagnostic(diagnostics.DiagnosticLintAcceptanceCriteriaMissing, root, path, "task has no observable acceptance criteria entries", ""))
	}
	if section, ok := markdownSection(source, doc, "Fuente de verdad"); ok && !hasListItem(section) {
		found = append(found, lintTaskSectionDiagnostic(diagnostics.DiagnosticLintSourceOfTruthEmpty, root, path, "task source-of-truth section has no entries", ""))
	}
	return found, nil
}

func markdownSection(source []byte, doc MarkdownDocument, headingText string) (string, bool) {
	lines := strings.Split(string(source), "\n")
	for i, heading := range doc.Headings {
		if heading.Text != headingText {
			continue
		}
		start := heading.StartLine
		end := len(lines)
		for _, next := range doc.Headings[i+1:] {
			if next.Level <= heading.Level {
				end = next.StartLine - 1
				break
			}
		}
		if start < 0 || start >= len(lines) || end < start {
			return "", true
		}
		return strings.Join(lines[start:end], "\n"), true
	}
	return "", false
}

func hasListItem(section string) bool {
	for _, line := range strings.Split(section, "\n") {
		if listItemPattern.MatchString(line) {
			return true
		}
	}
	return false
}

func lintTaskSectionDiagnostic(id string, root string, path string, message string, section string) diagnostics.Diagnostic {
	details := map[string]any{}
	if section != "" {
		details["target"] = section
		details["section"] = section
	}
	return diagnostics.Diagnostic{ID: id, Severity: diagnostics.SeverityWarning, Message: message, Path: relPath(root, path), Details: details}
}
