package roadmap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

const (
	DiagnosticSingleFileFallback   = "RMC_STRUCTURE_SINGLE_FILE_FALLBACK"
	DiagnosticMissingOutcomeReadme = "RMC_STRUCTURE_MISSING_OUTCOME_README"
	DiagnosticDuplicateID          = "RMC_STRUCTURE_DUPLICATE_ID"
	DiagnosticExtraNesting         = "RMC_STRUCTURE_EXTRA_NESTING"
	DiagnosticInvalidTaskFilename  = "RMC_STRUCTURE_INVALID_TASK_FILENAME"
	DiagnosticInvalidOutcomeDir    = "RMC_STRUCTURE_INVALID_OUTCOME_DIR"
)

type Diagnostic = diagnostics.Diagnostic

func CheckStructure(roadmapRoot string) ([]Diagnostic, error) {
	root := filepath.Clean(roadmapRoot)
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read roadmap root: %w", err)
	}

	var found []Diagnostic
	outcomeIDs := map[string]string{}
	directTaskIDs := map[string]string{}

	for _, entry := range entries {
		name := entry.Name()
		if ignoredEntry(name) {
			continue
		}
		path := filepath.Join(root, name)

		if entry.IsDir() {
			outcomeID, ok := OutcomeID(name)
			if !ok {
				found = append(found, structureDiagnostic(DiagnosticInvalidOutcomeDir, relPath(root, path), "outcome directories must be named OXX-slug"))
				continue
			}
			if first, exists := outcomeIDs[outcomeID]; exists {
				found = append(found, duplicateDiagnostic(root, filepath.Join(path, "README.md"), outcomeID, first))
			} else {
				outcomeIDs[outcomeID] = relPath(root, path)
			}
			outcomeDiagnostics, err := checkOutcomeStructure(root, path)
			if err != nil {
				return nil, err
			}
			found = append(found, outcomeDiagnostics...)
			continue
		}

		if !strings.HasSuffix(name, ".md") {
			continue
		}
		if isSingleFileFallback(name) {
			found = append(found, structureDiagnostic(DiagnosticSingleFileFallback, relPath(root, path), "multiple tasks must be materialized as canonical TXXX files, not a summary file"))
			continue
		}
		taskID, ok := TaskID(name)
		if !ok {
			found = append(found, structureDiagnostic(DiagnosticInvalidTaskFilename, relPath(root, path), "direct tasks must be named TXXX-slug.md"))
			continue
		}
		if first, exists := directTaskIDs[taskID]; exists {
			found = append(found, duplicateDiagnostic(root, path, taskID, first))
		} else {
			directTaskIDs[taskID] = relPath(root, path)
		}
	}

	return found, nil
}

func checkOutcomeStructure(root string, outcomePath string) ([]Diagnostic, error) {
	entries, err := os.ReadDir(outcomePath)
	if err != nil {
		return nil, fmt.Errorf("read outcome %s: %w", relPath(root, outcomePath), err)
	}

	var found []Diagnostic
	readmePath := filepath.Join(outcomePath, "README.md")
	if _, err := os.Stat(readmePath); err != nil {
		if os.IsNotExist(err) {
			found = append(found, structureDiagnostic(DiagnosticMissingOutcomeReadme, relPath(root, readmePath), "outcomes must contain README.md"))
		} else {
			return nil, fmt.Errorf("inspect outcome README %s: %w", relPath(root, readmePath), err)
		}
	}

	taskIDs := map[string]string{}
	for _, entry := range entries {
		name := entry.Name()
		if ignoredEntry(name) {
			continue
		}
		path := filepath.Join(outcomePath, name)
		if entry.IsDir() {
			found = append(found, structureDiagnostic(DiagnosticExtraNesting, relPath(root, path), "outcomes cannot contain nested directories"))
			continue
		}
		if name == "README.md" || !strings.HasSuffix(name, ".md") {
			continue
		}
		if isSingleFileFallback(name) {
			found = append(found, structureDiagnostic(DiagnosticSingleFileFallback, relPath(root, path), "multiple tasks must be materialized as canonical TXXX files, not a summary file"))
			continue
		}
		taskID, ok := TaskID(name)
		if !ok {
			found = append(found, structureDiagnostic(DiagnosticInvalidTaskFilename, relPath(root, path), "outcome tasks must be named TXXX-slug.md"))
			continue
		}
		if first, exists := taskIDs[taskID]; exists {
			found = append(found, duplicateDiagnostic(root, path, taskID, first))
		} else {
			taskIDs[taskID] = relPath(root, path)
		}
	}

	return found, nil
}

func structureDiagnostic(id string, path string, message string) Diagnostic {
	return Diagnostic{ID: id, Severity: diagnostics.SeverityError, Message: message, Path: path}
}

func duplicateDiagnostic(root string, path string, id string, first string) Diagnostic {
	return Diagnostic{
		ID:       DiagnosticDuplicateID,
		Severity: diagnostics.SeverityError,
		Message:  "duplicate roadmap id in the same scope",
		Path:     relPath(root, path),
		Details: map[string]any{
			"id":    id,
			"first": first,
		},
	}
}

func ignoredEntry(name string) bool {
	return strings.HasPrefix(name, ".")
}

func isSingleFileFallback(name string) bool {
	return strings.HasSuffix(name, "-tasks.md")
}

func relPath(root string, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}
	return filepath.ToSlash(rel)
}
