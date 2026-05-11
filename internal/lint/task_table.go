package lint

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

var taskFilenamePattern = regexp.MustCompile(`^T[0-9][0-9][0-9]-.+\.md$`)

func CheckOutcomeTaskTables(roadmapRoot string) ([]diagnostics.Diagnostic, error) {
	root := filepath.Clean(roadmapRoot)
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var found []diagnostics.Diagnostic
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "O") {
			continue
		}
		outcomePath := filepath.Join(root, entry.Name())
		diagnosticsForOutcome, err := checkOutcomeTaskTable(root, outcomePath)
		if err != nil {
			return nil, err
		}
		found = append(found, diagnosticsForOutcome...)
	}
	return found, nil
}

func checkOutcomeTaskTable(root string, outcomePath string) ([]diagnostics.Diagnostic, error) {
	readmePath := filepath.Join(outcomePath, "README.md")
	tasks, err := outcomeTaskFiles(outcomePath)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}
	source, err := os.ReadFile(readmePath)
	if err != nil {
		return nil, err
	}
	doc, err := ParseMarkdown(source)
	if err != nil {
		return nil, err
	}
	table := doc.TableBySection("Tasks")
	if table == nil {
		// ## Tasks is a computed view; its absence is expected and not an error
		return nil, nil
	}
	linked := map[string]bool{}
	var found []diagnostics.Diagnostic
	for _, row := range table.Rows {
		if len(row.Cells) == 0 || len(row.Cells[0].Links) == 0 {
			found = append(found, lintTaskTableDiagnostic(diagnostics.DiagnosticLintTaskTableInvalidLink, root, readmePath, "task table row must link to a child task file", ""))
			continue
		}
		for _, link := range row.Cells[0].Links {
			target := normalizeTableTaskTarget(link.Destination)
			if !validChildTaskLink(target) {
				found = append(found, lintTaskTableDiagnostic(diagnostics.DiagnosticLintTaskTableInvalidLink, root, readmePath, "task table link must target a child TXXX markdown task", link.Destination))
				continue
			}
			linked[target] = true
			if !tasks[target] {
				found = append(found, lintTaskTableDiagnostic(diagnostics.DiagnosticLintTaskTableStaleRow, root, readmePath, "task table row links to a missing child task", target))
			}
		}
	}
	for task := range tasks {
		if !linked[task] {
			found = append(found, lintTaskTableDiagnostic(diagnostics.DiagnosticLintTaskTableMissingRow, root, readmePath, "child task is missing from ## Tasks table", task))
		}
	}
	sortDiagnostics(found)
	return found, nil
}

func outcomeTaskFiles(outcomePath string) (map[string]bool, error) {
	entries, err := os.ReadDir(outcomePath)
	if err != nil {
		return nil, err
	}
	tasks := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() || !taskFilenamePattern.MatchString(entry.Name()) {
			continue
		}
		tasks[entry.Name()] = true
	}
	return tasks, nil
}

func normalizeTableTaskTarget(target string) string {
	target = strings.Split(target, "#")[0]
	target = strings.Split(target, "?")[0]
	return filepath.ToSlash(strings.TrimPrefix(target, "./"))
}

func validChildTaskLink(target string) bool {
	return !strings.Contains(target, "/") && taskFilenamePattern.MatchString(target)
}

func lintTaskTableDiagnostic(id string, root string, readmePath string, message string, target string) diagnostics.Diagnostic {
	details := map[string]any{}
	if target != "" {
		details["target"] = target
	}
	return diagnostics.Diagnostic{ID: id, Severity: diagnostics.SeverityWarning, Message: message, Path: relPath(root, readmePath), Details: details}
}

func sortDiagnostics(found []diagnostics.Diagnostic) {
	sort.Slice(found, func(i int, j int) bool {
		if found[i].Path != found[j].Path {
			return found[i].Path < found[j].Path
		}
		if found[i].ID != found[j].ID {
			return found[i].ID < found[j].ID
		}
		return detailTarget(found[i]) < detailTarget(found[j])
	})
}

func detailTarget(diagnostic diagnostics.Diagnostic) string {
	if diagnostic.Details == nil {
		return ""
	}
	value, _ := diagnostic.Details["target"].(string)
	return value
}

func relPath(root string, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}
	return filepath.ToSlash(rel)
}
