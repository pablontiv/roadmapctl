package materialize

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/roadmap"
)

const PathPlanKind = "roadmapctl/path-plan"
const PathPlanResultKind = "roadmapctl/path-plan-result"

// PathPlanInput is the compact input format for planning paths.
// It requires only slugs and relationships, not prose or full semantic details.
type PathPlanInput struct {
	Version int               `json:"version"`
	Kind    string            `json:"kind"`
	Items   []PathPlanItem    `json:"items"`
}

// PathPlanItem represents a single outcome or task in a path plan input.
type PathPlanItem struct {
	Type         string `json:"type"` // "outcome" or "task"
	Slug         string `json:"slug"`
	OutcomeSlug  string `json:"outcome_slug,omitempty"` // For tasks nested in outcomes
}

// PathPlanResult is the output format with proposed paths and diagnostics.
type PathPlanResult struct {
	Version     int                  `json:"version"`
	Kind        string               `json:"kind"`
	Paths       []PathPlanEntry      `json:"paths"`
	Collisions  []PathPlanCollision  `json:"collisions"`
	Diagnostics []PathPlanDiagnostic `json:"diagnostics"`
}

// PathPlanEntry describes a single proposed path in the result.
type PathPlanEntry struct {
	Path      string `json:"path"`      // e.g., "O14-rebuild-api/README.md"
	Operation string `json:"operation"` // "create" or "update"
	Type      string `json:"type"`      // "outcome" or "task"
}

// PathPlanCollision describes a path conflict detected during planning.
type PathPlanCollision struct {
	Path    string `json:"path"`
	Reason  string `json:"reason"`
	Planned string `json:"planned,omitempty"` // What we were trying to create
}

// PathPlanDiagnostic provides diagnostic information (warnings/errors).
type PathPlanDiagnostic struct {
	ID       string `json:"id"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Target   string `json:"target,omitempty"`
}

// PlanPaths computes canonical paths for Outcomes and Tasks without rendering content.
// It detects collisions, validates slugs, and determines whether each path is a create or update.
func PlanPaths(roadmapRoot string, input PathPlanInput) (PathPlanResult, []diagnostics.Diagnostic, error) {
	// Validate input
	if found := validatePathPlanInput(input); len(found) > 0 {
		return PathPlanResult{}, found, nil
	}

	// Build a request for roadmap.PlanMaterializePaths
	req := roadmap.MaterializePathRequest{}
	for _, item := range input.Items {
		if item.Type == "outcome" {
			req.Outcomes = append(req.Outcomes, roadmap.OutcomePathRequest{Slug: item.Slug})
			continue
		}
		req.DirectTasks = append(req.DirectTasks, roadmap.TaskPathRequest{Slug: item.Slug})
	}

	// Use the existing path planning logic
	paths, found, err := roadmap.PlanMaterializePaths(roadmapRoot, req)
	if err != nil || len(found) > 0 {
		return PathPlanResult{}, found, err
	}

	// Convert the result to PathPlanResult format
	result := PathPlanResult{
		Version:     1,
		Kind:        PathPlanResultKind,
		Paths:       []PathPlanEntry{},
		Collisions:  []PathPlanCollision{},
		Diagnostics: []PathPlanDiagnostic{},
	}

	// Note: collision detection is handled by roadmap.PlanMaterializePaths
	// which checks for existing outcomes and task slug conflicts

	// Add outcome paths
	for _, outcomePlan := range paths.Outcomes {
		operation := "create"
		if outcomePlan.Existing {
			operation = "update"
		}
		result.Paths = append(result.Paths, PathPlanEntry{
			Path:      filepath.ToSlash(outcomePlan.Path),
			Operation: operation,
			Type:      "outcome",
		})

		// Add task paths
		for _, taskPlan := range outcomePlan.Tasks {
			result.Paths = append(result.Paths, PathPlanEntry{
				Path:      filepath.ToSlash(taskPlan.Path),
				Operation: "create",
				Type:      "task",
			})
		}
	}

	// Add direct task paths
	for _, taskPlan := range paths.DirectTasks {
		result.Paths = append(result.Paths, PathPlanEntry{
			Path:      filepath.ToSlash(taskPlan.Path),
			Operation: "create",
			Type:      "task",
		})
	}

	return result, nil, nil
}

func validatePathPlanInput(input PathPlanInput) []diagnostics.Diagnostic {
	var found []diagnostics.Diagnostic

	if input.Version != 1 {
		found = append(found, pathPlanDiagnostic(
			diagnostics.DiagnosticMaterializeInputVersionUnsupported,
			"path plan version must be 1",
			"/version",
		))
	}

	if input.Kind != PathPlanKind {
		found = append(found, pathPlanDiagnostic(
			diagnostics.DiagnosticMaterializeInputKindInvalid,
			"path plan kind is invalid",
			"/kind",
		))
	}

	if len(input.Items) == 0 {
		found = append(found, pathPlanDiagnostic(
			diagnostics.DiagnosticMaterializeInputEmpty,
			"path plan must contain at least one item",
			"/items",
		))
	}

	for i, item := range input.Items {
		pointer := fmt.Sprintf("/items/%d", i)

		if item.Type != "outcome" && item.Type != "task" {
			found = append(found, pathPlanDiagnostic(
				diagnostics.DiagnosticMaterializeInputFieldMissing,
				"item type must be outcome or task",
				pointer+"/type",
			))
			continue
		}

		if strings.TrimSpace(item.Slug) == "" {
			found = append(found, pathPlanDiagnostic(
				diagnostics.DiagnosticMaterializeInputFieldMissing,
				"slug is required",
				pointer+"/slug",
			))
			continue
		}

		if !validSlug(item.Slug) {
			found = append(found, pathPlanDiagnostic(
				diagnostics.DiagnosticMaterializeInputSlugInvalid,
				"slug is not portable",
				pointer+"/slug",
			))
		}

		if item.Type == "task" && strings.TrimSpace(item.OutcomeSlug) == "" {
			// Task must either be in an outcome or standalone; if OutcomeSlug is empty, treat as direct task
		}
	}

	return found
}

func validSlug(slug string) bool {
	if strings.HasPrefix(slug, "O") || strings.HasPrefix(slug, "T") || strings.ContainsAny(slug, `/\\`) || strings.Contains(slug, "..") {
		return false
	}
	// Must be lowercase letters, digits, hyphens; start and end with letter or digit
	if slug == "" || (len(slug) > 1 && (slug[0] == '-' || slug[len(slug)-1] == '-')) {
		return false
	}
	for _, r := range slug {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}

func pathPlanDiagnostic(id string, message string, pointer string) diagnostics.Diagnostic {
	return diagnostics.Diagnostic{
		ID:       id,
		Severity: diagnostics.SeverityError,
		Message:  message,
		Details:  map[string]any{"pointer": pointer},
	}
}
