package roadmap

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

var (
	blockedByLinkPattern = regexp.MustCompile(`\[\[blocked_by:([^\]]+)\]\]`)
	outcomeNamePattern   = regexp.MustCompile(`^(O[0-9]{2})-.+`)
	taskNamePattern      = regexp.MustCompile(`^(T[0-9]{3})-.+\.md$`)
)

func CheckStructure(roadmapRoot string) ([]Diagnostic, error) {
	root := filepath.Clean(roadmapRoot)
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read roadmap root: %w", err)
	}

	var found []Diagnostic
	outcomeIDs := map[string]string{}
	directTaskIDs := map[string]string{}

	rawLinkDiagnostics, err := checkRawBlockedByLinks(root)
	if err != nil {
		return nil, err
	}
	found = append(found, rawLinkDiagnostics...)

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

func checkRawBlockedByLinks(root string) ([]Diagnostic, error) {
	var found []Diagnostic
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if ignoredEntry(entry.Name()) && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		data, err := os.ReadFile(path) //nolint:gosec
		if err != nil {
			return fmt.Errorf("read roadmap record %s: %w", relPath(root, path), err)
		}
		for _, match := range blockedByLinkPattern.FindAllStringSubmatch(string(data), -1) {
			target := strings.TrimSpace(match[1])
			if diagnostic, ok := invalidBlockedByDiagnostic(root, path, target); ok {
				found = append(found, diagnostic)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return found, nil
}

func invalidBlockedByDiagnostic(root string, sourcePath string, target string) (Diagnostic, bool) {
	if !isExplicitBlockedByTarget(target) {
		return blockedByDiagnostic(root, sourcePath, target, "blocked_by target must use explicit relative path to a task file"), true
	}
	resolved := filepath.Clean(filepath.Join(filepath.Dir(sourcePath), filepath.FromSlash(target)))
	if !strings.HasPrefix(resolved, root+string(filepath.Separator)) && resolved != root {
		return blockedByDiagnostic(root, sourcePath, target, "blocked_by target must stay inside roadmap root"), true
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return Diagnostic{}, false
	}
	if info.IsDir() {
		return blockedByDiagnostic(root, sourcePath, target, "blocked_by target must point to a task file"), true
	}
	if _, ok := TaskID(filepath.Base(resolved)); !ok {
		return blockedByDiagnostic(root, sourcePath, target, "blocked_by target must point to a TXXX task file"), true
	}
	return Diagnostic{}, false
}

func isExplicitBlockedByTarget(target string) bool {
	if !strings.HasPrefix(target, "./") && !strings.HasPrefix(target, "../") && !strings.Contains(target, "/") {
		return false
	}
	return strings.HasPrefix(filepath.Base(target), "T") && strings.HasSuffix(target, ".md")
}

func blockedByDiagnostic(root string, sourcePath string, target string, message string) Diagnostic {
	return Diagnostic{
		ID:       diagnostics.DiagnosticInvalidBlockedBy,
		Severity: diagnostics.SeverityError,
		Message:  message,
		Path:     relPath(root, sourcePath),
		Details:  map[string]any{"target": target, "source": "raw-scan"},
	}
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

// OutcomeID extracts the numeric ID from an outcome directory name (e.g., "O01-slug" -> "O01").
func OutcomeID(name string) (string, bool) {
	matches := outcomeNamePattern.FindStringSubmatch(name)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

// TaskID extracts the numeric ID from a task filename (e.g., "T001-slug.md" -> "T001").
func TaskID(name string) (string, bool) {
	matches := taskNamePattern.FindStringSubmatch(name)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}
