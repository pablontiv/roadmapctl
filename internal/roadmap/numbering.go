package roadmap

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/fsx"
)

const (
	DiagnosticMaterializeInputSlugInvalid = diagnostics.DiagnosticMaterializeInputSlugInvalid
	DiagnosticMaterializePlanConflict     = diagnostics.DiagnosticMaterializePlanConflict
)

var (
	outcomeNamePattern = regexp.MustCompile(`^(O[0-9]{2})-.+`)
	taskNamePattern    = regexp.MustCompile(`^(T[0-9]{3})-.+\.md$`)
	slugPattern        = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)
)

type MaterializePathRequest struct {
	Outcomes    []OutcomePathRequest
	DirectTasks []TaskPathRequest
}

type OutcomePathRequest struct {
	Slug  string
	Tasks []TaskPathRequest
}

type TaskPathRequest struct {
	Slug string
}

type MaterializePathPlan struct {
	Outcomes    []OutcomePathPlan
	DirectTasks []TaskPathPlan
}

type OutcomePathPlan struct {
	Slug     string
	Path     string
	Dir      string
	Existing bool
	Tasks    []TaskPathPlan
}

type TaskPathPlan struct {
	Slug string
	Path string
}

func PlanMaterializePaths(roadmapRoot string, request MaterializePathRequest) (MaterializePathPlan, []Diagnostic, error) {
	root := filepath.Clean(roadmapRoot)
	entries, err := os.ReadDir(root)
	if err != nil {
		if !os.IsNotExist(err) {
			return MaterializePathPlan{}, nil, fmt.Errorf("read roadmap root: %w", err)
		}
		entries = nil
	}
	maxOutcome := 0
	maxDirectTask := 0
	for _, entry := range entries {
		if ignoredEntry(entry.Name()) {
			continue
		}
		if entry.IsDir() {
			maxOutcome = max(maxOutcome, numericSuffix(entry.Name(), OutcomeID))
			continue
		}
		maxDirectTask = max(maxDirectTask, numericSuffix(entry.Name(), TaskID))
	}
	plan := MaterializePathPlan{}
	var found []Diagnostic
	planned := map[string]bool{}
	for _, outcome := range request.Outcomes {
		if !validSlug(outcome.Slug) {
			found = append(found, materializePathDiagnostic(DiagnosticMaterializeInputSlugInvalid, "", "outcome slug is not portable", outcome.Slug))
			continue
		}
		dir, readme, existing := existingOutcomeForSlug(root, outcome.Slug)
		nextTaskNumber := 0
		if existing {
			nextTaskNumber = maxTaskInOutcome(root, dir)
		} else {
			maxOutcome++
			dir = fmt.Sprintf("O%02d-%s", maxOutcome, outcome.Slug)
			readme = filepath.ToSlash(filepath.Join(dir, "README.md"))
			if diagnostic, ok := plannedPathDiagnostic(root, readme, planned); ok {
				found = append(found, diagnostic)
				continue
			}
		}
		outcomePlan := OutcomePathPlan{Slug: outcome.Slug, Dir: dir, Path: readme, Existing: existing}
		for _, task := range outcome.Tasks {
			if !validSlug(task.Slug) {
				found = append(found, materializePathDiagnostic(DiagnosticMaterializeInputSlugInvalid, readme, "task slug is not portable", task.Slug))
				continue
			}
			if existingPath := existingOutcomeTaskSlug(root, dir, task.Slug); existingPath != "" {
				found = append(found, materializePathDiagnostic(DiagnosticMaterializePlanConflict, existingPath, "planned task slug collides with an existing outcome task", existingPath))
				continue
			}
			nextTaskNumber++
			taskPath := filepath.ToSlash(filepath.Join(dir, fmt.Sprintf("T%03d-%s.md", nextTaskNumber, task.Slug)))
			if diagnostic, ok := plannedPathDiagnostic(root, taskPath, planned); ok {
				found = append(found, diagnostic)
				continue
			}
			outcomePlan.Tasks = append(outcomePlan.Tasks, TaskPathPlan{Slug: task.Slug, Path: taskPath})
		}
		plan.Outcomes = append(plan.Outcomes, outcomePlan)
	}
	for _, task := range request.DirectTasks {
		if !validSlug(task.Slug) {
			found = append(found, materializePathDiagnostic(DiagnosticMaterializeInputSlugInvalid, "", "task slug is not portable", task.Slug))
			continue
		}
		if existingPath := existingRootTaskSlug(root, task.Slug); existingPath != "" {
			found = append(found, materializePathDiagnostic(DiagnosticMaterializePlanConflict, existingPath, "planned task slug collides with an existing roadmap item", existingPath))
			continue
		}
		maxDirectTask++
		taskPath := fmt.Sprintf("T%03d-%s.md", maxDirectTask, task.Slug)
		if diagnostic, ok := plannedPathDiagnostic(root, taskPath, planned); ok {
			found = append(found, diagnostic)
			continue
		}
		plan.DirectTasks = append(plan.DirectTasks, TaskPathPlan{Slug: task.Slug, Path: taskPath})
	}
	if len(found) > 0 {
		return MaterializePathPlan{}, found, nil
	}
	return plan, nil, nil
}

func OutcomeID(name string) (string, bool) {
	matches := outcomeNamePattern.FindStringSubmatch(name)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

func TaskID(name string) (string, bool) {
	matches := taskNamePattern.FindStringSubmatch(name)
	if len(matches) != 2 {
		return "", false
	}
	return matches[1], true
}

func numericSuffix(name string, idFunc func(string) (string, bool)) int {
	id, ok := idFunc(name)
	if !ok || len(id) < 2 {
		return 0
	}
	value, err := strconv.Atoi(id[1:])
	if err != nil {
		return 0
	}
	return value
}

func validSlug(slug string) bool {
	if strings.HasPrefix(slug, "O") || strings.HasPrefix(slug, "T") || strings.ContainsAny(slug, `/\\`) || strings.Contains(slug, "..") {
		return false
	}
	return slugPattern.MatchString(slug)
}

func plannedPathDiagnostic(root string, rel string, planned map[string]bool) (Diagnostic, bool) {
	if strings.Contains(rel, "-tasks.md") {
		return materializePathDiagnostic(DiagnosticMaterializeInputSlugInvalid, rel, "planned path must be canonical and not a fallback summary", rel), true
	}
	abs, _, err := fsx.ResolveInside(root, rel)
	if err != nil {
		return materializePathDiagnostic(DiagnosticMaterializeInputSlugInvalid, rel, "planned path escapes roadmap root", rel), true
	}
	key := filepath.ToSlash(filepath.Clean(rel))
	if planned[key] || fileOrDirExists(abs) {
		return materializePathDiagnostic(DiagnosticMaterializePlanConflict, key, "planned path collides with an existing roadmap item", key), true
	}
	planned[key] = true
	return Diagnostic{}, false
}

func existingOutcomeForSlug(root string, slug string) (string, string, bool) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", "", false
	}
	for _, entry := range entries {
		if !entry.IsDir() || ignoredEntry(entry.Name()) {
			continue
		}
		id, ok := OutcomeID(entry.Name())
		if !ok {
			continue
		}
		if strings.TrimPrefix(entry.Name(), id+"-") == slug {
			readme := filepath.ToSlash(filepath.Join(entry.Name(), "README.md"))
			return entry.Name(), readme, true
		}
	}
	return "", "", false
}

func maxTaskInOutcome(root string, dir string) int {
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return 0
	}
	maxTask := 0
	for _, entry := range entries {
		if entry.IsDir() || ignoredEntry(entry.Name()) {
			continue
		}
		maxTask = max(maxTask, numericSuffix(entry.Name(), TaskID))
	}
	return maxTask
}

func existingOutcomeTaskSlug(root string, dir string, slug string) string {
	entries, err := os.ReadDir(filepath.Join(root, filepath.FromSlash(dir)))
	if err != nil {
		return ""
	}
	suffix := "-" + slug + ".md"
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if _, ok := TaskID(entry.Name()); ok && strings.HasSuffix(entry.Name(), suffix) {
			return filepath.ToSlash(filepath.Join(dir, entry.Name()))
		}
	}
	return ""
}

func existingRootTaskSlug(root string, slug string) string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return ""
	}
	suffix := "-" + slug + ".md"
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if _, ok := TaskID(entry.Name()); ok && strings.HasSuffix(entry.Name(), suffix) {
			return entry.Name()
		}
	}
	return ""
}

func fileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func materializePathDiagnostic(id string, path string, message string, target string) Diagnostic {
	return Diagnostic{ID: id, Severity: diagnostics.SeverityError, Message: message, Path: path, Details: map[string]any{"target": target}}
}
