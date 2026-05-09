package materialize

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/diff"
	"github.com/pablontiv/roadmapctl/internal/roadmap"
)

const PlanKind = "roadmapctl/materialize-plan"

const baseStemContent = `version: 2
scope:
  match: "*.md"

schema:
  estado:
    type: enum
    required:
      match: ["T*"]
    match: ["O*", "T*"]
    values: [Pending, Specified, In Progress, Completed, Blocked, On Hold, Obsolete]

  tipo:
    type: enum
    required:
      match: ["O*", "T*"]
    match: ["O*", "T*"]
    values: [outcome, task]

  id:
    type: sequence
    match:
      "O*": { prefix: O, digits: 2 }
      "T*": { prefix: T, digits: 3 }

links:
  blocked_by:
    target: '^(\./|\.\./|.*/)T[0-9]{3}-[^/]+\.md$'
  reference:
    target: ".*"

validate:
  - field: tipo
    rule: non_empty
`

const defaultRoadmapctlTOML = `done_statuses = ["Completed", "Obsolete"]
active_statuses = ["Pending", "Specified", "In Progress"]
leaf_filter = "isIndex == false"
outcome_close_verify = []
pr_merge_strategy = "squash"
commit_style = "conventional"
auto_push = true
loop_max_tasks = 0
parallel = true
autonomy = "until_done"
compact_after_task_commit = true
pr_mode = false

[status_values]
pending = "Pending"
specified = "Specified"
in_progress = "In Progress"
completed = "Completed"
blocked = "Blocked"
obsolete = "Obsolete"
`

type Plan struct {
	Version int    `json:"version"`
	Kind    string `json:"kind"`
	Title   string `json:"title,omitempty"`
	Items   []Item `json:"items"`
}

type Item struct {
	Type               string       `json:"type"`
	Slug               string       `json:"slug"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	AcceptanceCriteria []string     `json:"acceptance_criteria"`
	Tasks              []Task       `json:"tasks,omitempty"`
	ContributesTo      []string     `json:"contributes_to,omitempty"`
	Preserves          []string     `json:"preserves,omitempty"`
	Context            string       `json:"context,omitempty"`
	ScopeIn            []string     `json:"scope_in,omitempty"`
	ScopeOut           []string     `json:"scope_out,omitempty"`
	InitialState       string       `json:"initial_state,omitempty"`
	SourceOfTruth      []string     `json:"source_of_truth,omitempty"`
	BlockedBy          []Dependency `json:"blocked_by,omitempty"`
	TechnicalSpec      string       `json:"technical_spec,omitempty"`
}

type Task struct {
	Type               string       `json:"type,omitempty"`
	Slug               string       `json:"slug"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	Preserves          []string     `json:"preserves"`
	Context            string       `json:"context"`
	ScopeIn            []string     `json:"scope_in"`
	ScopeOut           []string     `json:"scope_out"`
	InitialState       string       `json:"initial_state"`
	AcceptanceCriteria []string     `json:"acceptance_criteria"`
	SourceOfTruth      []string     `json:"source_of_truth"`
	BlockedBy          []Dependency `json:"blocked_by,omitempty"`
	TechnicalSpec      string       `json:"technical_spec,omitempty"`
}

type Dependency struct {
	Ref  string `json:"ref,omitempty"`
	Path string `json:"path,omitempty"`
}

type Result struct {
	Changes []Change `json:"changes"`
}

type Change struct {
	Path          string   `json:"path"`
	Operation     string   `json:"operation"`
	Applied       bool     `json:"applied"`
	Content       string   `json:"content,omitempty"`
	Diff          string   `json:"diff,omitempty"`
	Preconditions []string `json:"preconditions,omitempty"`
}

type plannedTask struct {
	Ref  string
	Path string
	Task Task
}

func Apply(roadmapRoot string, plan Plan) (Result, []diagnostics.Diagnostic, error) {
	result, found, err := DryRun(roadmapRoot, plan)
	if err != nil || len(found) > 0 {
		return result, found, err
	}
	for i := range result.Changes {
		change := &result.Changes[i]
		if change.Operation == "mkdir" {
			continue
		}
		abs := filepath.Join(filepath.Clean(roadmapRoot), filepath.FromSlash(change.Path))
		if _, err := os.Stat(abs); err == nil {
			return result, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializePlanConflict, change.Path, "planned path now exists; dry-run is stale", change.Path)}, nil
		} else if !os.IsNotExist(err) {
			return result, nil, fmt.Errorf("stat planned path: %w", err)
		}
	}
	for i := range result.Changes {
		change := &result.Changes[i]
		abs := filepath.Join(filepath.Clean(roadmapRoot), filepath.FromSlash(change.Path))
		if change.Operation == "mkdir" {
			if err := os.MkdirAll(abs, 0o755); err != nil {
				return result, nil, fmt.Errorf("create materialization directory: %w", err)
			}
			change.Applied = true
			continue
		}
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return result, nil, fmt.Errorf("create parent directory: %w", err)
		}
		file, err := os.OpenFile(abs, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return result, nil, fmt.Errorf("create materialized file: %w", err)
		}
		if _, err := file.WriteString(change.Content); err != nil {
			_ = file.Close()
			return result, nil, fmt.Errorf("write materialized file: %w", err)
		}
		if err := file.Close(); err != nil {
			return result, nil, fmt.Errorf("close materialized file: %w", err)
		}
		change.Applied = true
	}
	return result, nil, nil
}

func ApplyChanges(roadmapRoot string, changes []Change) (Result, []diagnostics.Diagnostic, error) {
	ordered := append([]Change(nil), changes...)
	if len(ordered) == 0 {
		return Result{}, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputEmpty, "", "change set must contain at least one change", "/changes")}, nil
	}
	sort.SliceStable(ordered, func(i int, j int) bool {
		return batchApplyOrder(ordered[i]) < batchApplyOrder(ordered[j])
	})
	for _, change := range ordered {
		if diagnostic, ok := validateBatchChange(change); !ok {
			return Result{Changes: ordered}, []diagnostics.Diagnostic{diagnostic}, nil
		}
		abs := filepath.Join(filepath.Clean(roadmapRoot), filepath.FromSlash(change.Path))
		if change.Operation == "mkdir" {
			if info, err := os.Stat(abs); err == nil && !info.IsDir() {
				return Result{Changes: ordered}, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializePlanConflict, change.Path, "planned directory path exists as a file", change.Path)}, nil
			} else if err != nil && !os.IsNotExist(err) {
				return Result{Changes: ordered}, nil, fmt.Errorf("stat planned directory: %w", err)
			}
			continue
		}
		if _, err := os.Stat(abs); err == nil {
			return Result{Changes: ordered}, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializePlanConflict, change.Path, "planned path now exists; dry-run is stale", change.Path)}, nil
		} else if !os.IsNotExist(err) {
			return Result{Changes: ordered}, nil, fmt.Errorf("stat planned path: %w", err)
		}
	}
	for i := range ordered {
		change := &ordered[i]
		abs := filepath.Join(filepath.Clean(roadmapRoot), filepath.FromSlash(change.Path))
		if change.Operation == "mkdir" {
			if err := os.MkdirAll(abs, 0o755); err != nil {
				return Result{Changes: ordered}, nil, fmt.Errorf("create materialization directory: %w", err)
			}
			change.Applied = true
			continue
		}
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return Result{Changes: ordered}, nil, fmt.Errorf("create parent directory: %w", err)
		}
		file, err := os.OpenFile(abs, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return Result{Changes: ordered}, nil, fmt.Errorf("create materialized file: %w", err)
		}
		if _, err := file.WriteString(change.Content); err != nil {
			_ = file.Close()
			return Result{Changes: ordered}, nil, fmt.Errorf("write materialized file: %w", err)
		}
		if err := file.Close(); err != nil {
			return Result{Changes: ordered}, nil, fmt.Errorf("close materialized file: %w", err)
		}
		change.Applied = true
	}
	return Result{Changes: ordered}, nil, nil
}

func ApplyTarget(roadmapRoot string, changes []Change, target string) (Result, []diagnostics.Diagnostic, error) {
	cleanTarget := filepath.ToSlash(filepath.Clean(strings.TrimSpace(target)))
	if strings.TrimSpace(target) == "" || cleanTarget == "." {
		return Result{}, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputFieldMissing, cleanTarget, "target is required", "/target")}, nil
	}
	var matches []Change
	for _, change := range changes {
		if change.Path == cleanTarget {
			matches = append(matches, change)
		}
	}
	if len(matches) == 0 {
		return Result{}, []diagnostics.Diagnostic{materializeDiagnostic("RMC_MATERIALIZE_TARGET_UNKNOWN", cleanTarget, "target is not present in change set", cleanTarget)}, nil
	}
	if len(matches) > 1 {
		return Result{}, []diagnostics.Diagnostic{materializeDiagnostic("RMC_MATERIALIZE_TARGET_DUPLICATE", cleanTarget, "target appears multiple times in change set", cleanTarget)}, nil
	}
	change := matches[0]
	if change.Operation != "create" || change.Content == "" || !isCanonicalMaterializeFileTarget(change.Path) {
		return Result{}, []diagnostics.Diagnostic{materializeDiagnostic("RMC_MATERIALIZE_TARGET_INVALID", cleanTarget, "target must be one canonical roadmap markdown file create change", cleanTarget)}, nil
	}
	abs := filepath.Join(filepath.Clean(roadmapRoot), filepath.FromSlash(change.Path))
	if _, err := os.Stat(abs); err == nil {
		return Result{}, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializePlanConflict, change.Path, "planned path now exists; dry-run is stale", change.Path)}, nil
	} else if !os.IsNotExist(err) {
		return Result{}, nil, fmt.Errorf("stat planned path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return Result{}, nil, fmt.Errorf("create parent directory: %w", err)
	}
	file, err := os.OpenFile(abs, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return Result{}, nil, fmt.Errorf("create materialized file: %w", err)
	}
	if _, err := file.WriteString(change.Content); err != nil {
		_ = file.Close()
		return Result{}, nil, fmt.Errorf("write materialized file: %w", err)
	}
	if err := file.Close(); err != nil {
		return Result{}, nil, fmt.Errorf("close materialized file: %w", err)
	}
	change.Applied = true
	return Result{Changes: []Change{change}}, nil, nil
}

func batchApplyOrder(change Change) int {
	if change.Operation == "mkdir" {
		return 0
	}
	if change.Path == ".stem" || change.Path == ".roadmapctl.toml" {
		return 1
	}
	if strings.HasSuffix(change.Path, "/README.md") {
		return 2
	}
	return 3
}

func validateBatchChange(change Change) (diagnostics.Diagnostic, bool) {
	cleanPath := filepath.ToSlash(filepath.Clean(strings.TrimSpace(change.Path)))
	if strings.TrimSpace(change.Path) == "" || cleanPath != change.Path || strings.HasPrefix(cleanPath, "../") || filepath.IsAbs(change.Path) {
		return materializeDiagnostic("RMC_MATERIALIZE_CHANGE_INVALID", change.Path, "change path must be a clean roadmap-root-relative path", change.Path), false
	}
	if change.Operation == "mkdir" {
		if cleanPath == "." || isOutcomeDir(filepath.Base(cleanPath)) {
			return diagnostics.Diagnostic{}, true
		}
		return materializeDiagnostic("RMC_MATERIALIZE_CHANGE_INVALID", change.Path, "mkdir change must target roadmap root or an outcome directory", change.Path), false
	}
	if change.Operation != "create" {
		return materializeDiagnostic("RMC_MATERIALIZE_CHANGE_INVALID", change.Path, "batch apply supports only mkdir and create changes", change.Path), false
	}
	if change.Content == "" {
		return materializeDiagnostic("RMC_MATERIALIZE_CHANGE_INVALID", change.Path, "create change must include content", change.Path), false
	}
	if change.Path == ".stem" || change.Path == ".roadmapctl.toml" || isCanonicalMaterializeFileTarget(change.Path) {
		return diagnostics.Diagnostic{}, true
	}
	return materializeDiagnostic("RMC_MATERIALIZE_CHANGE_INVALID", change.Path, "create change must target an allowlisted bootstrap or canonical roadmap file", change.Path), false
}

func isCanonicalMaterializeFileTarget(path string) bool {
	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return isTaskMarkdown(parts[0])
	}
	if len(parts) == 2 && isOutcomeDir(parts[0]) {
		return parts[1] == "README.md" || isTaskMarkdown(parts[1])
	}
	return false
}

func isOutcomeDir(name string) bool {
	return len(name) > 4 && name[0] == 'O' && isDigit(name[1]) && isDigit(name[2]) && name[3] == '-'
}

func isTaskMarkdown(name string) bool {
	return len(name) > 8 && name[0] == 'T' && isDigit(name[1]) && isDigit(name[2]) && isDigit(name[3]) && name[4] == '-' && strings.HasSuffix(name, ".md")
}

func isDigit(value byte) bool {
	return value >= '0' && value <= '9'
}

func DryRun(roadmapRoot string, plan Plan) (Result, []diagnostics.Diagnostic, error) {
	if found := validatePlan(plan); len(found) > 0 {
		return Result{}, found, nil
	}
	request := roadmap.MaterializePathRequest{}
	for _, item := range plan.Items {
		if item.Type == "outcome" {
			outcome := roadmap.OutcomePathRequest{Slug: item.Slug}
			for _, task := range item.Tasks {
				outcome.Tasks = append(outcome.Tasks, roadmap.TaskPathRequest{Slug: task.Slug})
			}
			request.Outcomes = append(request.Outcomes, outcome)
			continue
		}
		request.DirectTasks = append(request.DirectTasks, roadmap.TaskPathRequest{Slug: item.Slug})
	}
	paths, found, err := roadmap.PlanMaterializePaths(roadmapRoot, request)
	if err != nil || len(found) > 0 {
		return Result{}, found, err
	}

	refs := map[string]string{}
	outcomeBySlug := map[string]Item{}
	for _, item := range plan.Items {
		if item.Type == "outcome" {
			outcomeBySlug[item.Slug] = item
		}
	}
	for _, outcome := range paths.Outcomes {
		for _, task := range outcome.Tasks {
			refs[outcome.Slug+"/"+task.Slug] = task.Path
		}
	}
	for _, task := range paths.DirectTasks {
		refs[task.Slug] = task.Path
	}

	result := Result{Changes: bootstrapChanges(roadmapRoot)}
	for _, outcomePlan := range paths.Outcomes {
		item := outcomeBySlug[outcomePlan.Slug]
		content := renderOutcome(item, outcomePlan)
		result.Changes = append(result.Changes, newCreateChange(outcomePlan.Path, content))
		for i, taskPath := range outcomePlan.Tasks {
			if i >= len(item.Tasks) {
				continue
			}
			task := item.Tasks[i]
			links, depDiagnostics := dependencyLinks(roadmapRoot, taskPath.Path, task.BlockedBy, refs)
			if len(depDiagnostics) > 0 {
				return Result{}, depDiagnostics, nil
			}
			content := renderTask(task, item, links, taskID(taskPath.Path))
			result.Changes = append(result.Changes, newCreateChange(taskPath.Path, content))
		}
	}
	for i, taskPath := range paths.DirectTasks {
		item := directItemAt(plan.Items, i)
		task := taskFromItem(item)
		links, depDiagnostics := dependencyLinks(roadmapRoot, taskPath.Path, task.BlockedBy, refs)
		if len(depDiagnostics) > 0 {
			return Result{}, depDiagnostics, nil
		}
		content := renderTask(task, Item{}, links, taskID(taskPath.Path))
		result.Changes = append(result.Changes, newCreateChange(taskPath.Path, content))
	}
	return result, nil, nil
}

func validatePlan(plan Plan) []diagnostics.Diagnostic {
	var found []diagnostics.Diagnostic
	if plan.Version != 1 {
		found = append(found, materializeDiagnostic(diagnostics.DiagnosticMaterializeInputVersionUnsupported, "", "materialize plan version must be 1", "/version"))
	}
	if plan.Kind != PlanKind {
		found = append(found, materializeDiagnostic(diagnostics.DiagnosticMaterializeInputKindInvalid, "", "materialize plan kind is invalid", "/kind"))
	}
	if len(plan.Items) == 0 {
		found = append(found, materializeDiagnostic(diagnostics.DiagnosticMaterializeInputEmpty, "", "materialize plan must contain at least one item", "/items"))
	}
	for i, item := range plan.Items {
		pointer := fmt.Sprintf("/items/%d", i)
		switch item.Type {
		case "outcome":
			found = append(found, validateOutcome(item, pointer)...)
		case "task":
			found = append(found, validateTask(taskFromItem(item), pointer)...)
		default:
			found = append(found, materializeDiagnostic(diagnostics.DiagnosticMaterializeInputFieldMissing, "", "item type must be outcome or task", pointer+"/type"))
		}
	}
	return found
}

func validateOutcome(item Item, pointer string) []diagnostics.Diagnostic {
	var found []diagnostics.Diagnostic
	found = append(found, requireString(item.Slug, pointer+"/slug")...)
	found = append(found, requireString(item.Title, pointer+"/title")...)
	found = append(found, requireString(item.Description, pointer+"/description")...)
	found = append(found, requireStrings(item.AcceptanceCriteria, pointer+"/acceptance_criteria")...)
	if len(item.Tasks) == 0 {
		found = append(found, materializeDiagnostic(diagnostics.DiagnosticMaterializeInputFieldMissing, "", "outcome must contain at least one task", pointer+"/tasks"))
	}
	for i, task := range item.Tasks {
		found = append(found, validateTask(task, fmt.Sprintf("%s/tasks/%d", pointer, i))...)
	}
	return found
}

func validateTask(task Task, pointer string) []diagnostics.Diagnostic {
	var found []diagnostics.Diagnostic
	found = append(found, requireString(task.Slug, pointer+"/slug")...)
	found = append(found, requireString(task.Title, pointer+"/title")...)
	found = append(found, requireString(task.Description, pointer+"/description")...)
	found = append(found, requireStrings(task.Preserves, pointer+"/preserves")...)
	found = append(found, requireString(task.Context, pointer+"/context")...)
	found = append(found, requireStrings(task.ScopeIn, pointer+"/scope_in")...)
	found = append(found, requireStrings(task.ScopeOut, pointer+"/scope_out")...)
	found = append(found, requireString(task.InitialState, pointer+"/initial_state")...)
	found = append(found, requireStrings(task.AcceptanceCriteria, pointer+"/acceptance_criteria")...)
	found = append(found, requireStrings(task.SourceOfTruth, pointer+"/source_of_truth")...)
	for i, dep := range task.BlockedBy {
		if (strings.TrimSpace(dep.Ref) == "") == (strings.TrimSpace(dep.Path) == "") {
			found = append(found, materializeDiagnostic(diagnostics.DiagnosticMaterializeInputDependencyInvalid, "", "dependency must have exactly one of ref or path", fmt.Sprintf("%s/blocked_by/%d", pointer, i)))
		}
	}
	return found
}

func requireString(value string, pointer string) []diagnostics.Diagnostic {
	if strings.TrimSpace(value) == "" {
		return []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputFieldMissing, "", "required field is empty", pointer)}
	}
	return nil
}

func requireStrings(values []string, pointer string) []diagnostics.Diagnostic {
	if len(values) == 0 {
		return []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputFieldMissing, "", "required array is empty", pointer)}
	}
	for i, value := range values {
		if strings.TrimSpace(value) == "" {
			return []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputFieldMissing, "", "required array contains an empty value", fmt.Sprintf("%s/%d", pointer, i))}
		}
	}
	return nil
}

func dependencyLinks(roadmapRoot string, currentPath string, dependencies []Dependency, refs map[string]string) ([]string, []diagnostics.Diagnostic) {
	var links []string
	planned := map[string]bool{}
	for _, path := range refs {
		planned[filepath.ToSlash(filepath.Clean(path))] = true
	}
	for _, dep := range dependencies {
		if dep.Ref != "" {
			target, ok := refs[dep.Ref]
			if !ok {
				return nil, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputDependencyUnresolved, currentPath, "dependency ref cannot be resolved", dep.Ref)}
			}
			links = append(links, explicitRelative(currentPath, target))
			continue
		}
		if filepath.Base(dep.Path) == dep.Path {
			return nil, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputDependencyInvalid, currentPath, "dependency path must be explicit relative path or roadmap-root relative path", dep.Path)}
		}
		target, ok := dependencyTargetPath(currentPath, dep.Path)
		if !ok || !isTaskMarkdown(filepath.Base(target)) {
			return nil, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputDependencyInvalid, currentPath, "dependency path must resolve inside the roadmap to a task markdown file", dep.Path)}
		}
		if !planned[target] {
			abs := filepath.Join(filepath.Clean(roadmapRoot), filepath.FromSlash(target))
			info, err := os.Stat(abs)
			if os.IsNotExist(err) {
				return nil, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputDependencyUnresolved, currentPath, "dependency path cannot be resolved to an existing or planned task", dep.Path)}
			} else if err != nil {
				return nil, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputDependencyInvalid, currentPath, "dependency path cannot be checked", dep.Path)}
			} else if !info.Mode().IsRegular() {
				return nil, []diagnostics.Diagnostic{materializeDiagnostic(diagnostics.DiagnosticMaterializeInputDependencyInvalid, currentPath, "dependency path must resolve to a regular task markdown file", dep.Path)}
			}
		}
		links = append(links, explicitRelative(currentPath, target))
	}
	return links, nil
}

func dependencyTargetPath(currentPath string, dependencyPath string) (string, bool) {
	if strings.HasPrefix(dependencyPath, "/") {
		return "", false
	}
	if strings.HasPrefix(dependencyPath, "./") || strings.HasPrefix(dependencyPath, "../") {
		fromDir := filepath.Dir(filepath.FromSlash(currentPath))
		if fromDir == "." {
			fromDir = ""
		}
		target := filepath.ToSlash(filepath.Clean(filepath.Join(fromDir, filepath.FromSlash(dependencyPath))))
		return target, !strings.HasPrefix(target, "../") && target != "."
	}
	target := filepath.ToSlash(filepath.Clean(filepath.FromSlash(dependencyPath)))
	return target, !strings.HasPrefix(target, "../") && target != "."
}

func explicitRelative(currentPath string, targetPath string) string {
	fromDir := filepath.Dir(filepath.FromSlash(currentPath))
	rel, err := filepath.Rel(fromDir, filepath.FromSlash(targetPath))
	if err != nil {
		return filepath.ToSlash(targetPath)
	}
	rel = filepath.ToSlash(rel)
	if !strings.HasPrefix(rel, "../") && !strings.HasPrefix(rel, "./") {
		rel = "./" + rel
	}
	return rel
}

func renderOutcome(item Item, plan roadmap.OutcomePathPlan) string {
	var b strings.Builder
	b.WriteString("---\ntipo: outcome\n---\n")
	fmt.Fprintf(&b, "# %s\n\n", item.Title)
	b.WriteString(item.Description)
	b.WriteString("\n\n## Criterios de Aceptación\n\n")
	writeBullets(&b, item.AcceptanceCriteria)
	b.WriteString("\n## Tasks\n\n| Task | Descripción |\n|------|-------------|\n")
	for i, taskPlan := range plan.Tasks {
		description := ""
		if i < len(item.Tasks) {
			description = item.Tasks[i].Description
		}
		fmt.Fprintf(&b, "| [%s](%s) | %s |\n", taskID(taskPlan.Path), filepath.Base(taskPlan.Path), description)
	}
	return b.String()
}

func renderTask(task Task, outcome Item, links []string, id string) string {
	var b strings.Builder
	b.WriteString("---\nestado: Specified\ntipo: task\n---\n")
	fmt.Fprintf(&b, "# %s: %s\n\n", id, task.Title)
	if outcome.Title != "" {
		fmt.Fprintf(&b, "**Outcome**: [%s](README.md)\n\n", outcome.Title)
	}
	for _, link := range links {
		fmt.Fprintf(&b, "[[blocked_by:%s]]\n", link)
	}
	if len(links) > 0 {
		b.WriteByte('\n')
	}
	b.WriteString("## Preserva\n\n")
	writeBullets(&b, task.Preserves)
	b.WriteString("\n## Contexto\n\n")
	b.WriteString(task.Context)
	b.WriteString("\n\n## Alcance\n\n**In**:\n")
	writeNumbered(&b, task.ScopeIn)
	b.WriteString("\n**Out**:\n")
	writeNumbered(&b, task.ScopeOut)
	b.WriteString("\n## Estado inicial esperado\n\n")
	b.WriteString(task.InitialState)
	b.WriteString("\n\n")
	if strings.TrimSpace(task.TechnicalSpec) != "" {
		b.WriteString("## Especificación Técnica\n\n")
		b.WriteString(strings.TrimSpace(task.TechnicalSpec))
		b.WriteString("\n\n")
	}
	b.WriteString("## Criterios de Aceptación\n\n")
	writeBullets(&b, task.AcceptanceCriteria)
	b.WriteString("\n## Fuente de verdad\n\n")
	writeBullets(&b, task.SourceOfTruth)
	return b.String()
}

func writeBullets(b *strings.Builder, items []string) {
	for _, item := range items {
		fmt.Fprintf(b, "- %s\n", item)
	}
}

func writeNumbered(b *strings.Builder, items []string) {
	for i, item := range items {
		fmt.Fprintf(b, "%d. %s\n", i+1, item)
	}
}

func taskID(path string) string {
	base := filepath.Base(path)
	if len(base) >= 4 {
		return base[:4]
	}
	return "T000"
}

func directItemAt(items []Item, index int) Item {
	seen := 0
	for _, item := range items {
		if item.Type != "task" {
			continue
		}
		if seen == index {
			return item
		}
		seen++
	}
	return Item{}
}

func taskFromItem(item Item) Task {
	return Task{Type: item.Type, Slug: item.Slug, Title: item.Title, Description: item.Description, Preserves: item.Preserves, Context: item.Context, ScopeIn: item.ScopeIn, ScopeOut: item.ScopeOut, InitialState: item.InitialState, AcceptanceCriteria: item.AcceptanceCriteria, SourceOfTruth: item.SourceOfTruth, BlockedBy: item.BlockedBy, TechnicalSpec: item.TechnicalSpec}
}

func bootstrapChanges(roadmapRoot string) []Change {
	var changes []Change
	rootMissing := false
	if _, err := os.Stat(roadmapRoot); os.IsNotExist(err) {
		rootMissing = true
		changes = append(changes, Change{Path: ".", Operation: "mkdir", Applied: false, Preconditions: []string{"roadmap root must not exist or must be a directory"}})
	}
	stemPath := filepath.Join(roadmapRoot, ".stem")
	if _, err := os.Stat(stemPath); os.IsNotExist(err) {
		changes = append(changes, newCreateChange(".stem", baseStemContent))
	}
	configPath := filepath.Join(roadmapRoot, ".roadmapctl.toml")
	if rootMissing {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			changes = append(changes, newCreateChange(".roadmapctl.toml", defaultRoadmapctlTOML))
		}
	}
	return changes
}

func newCreateChange(path string, content string) Change {
	return Change{Path: path, Operation: "create", Applied: false, Content: content, Diff: diff.NewFile(path, content), Preconditions: []string{"path must not exist at apply time"}}
}

func materializeDiagnostic(id string, path string, message string, pointer string) diagnostics.Diagnostic {
	return diagnostics.Diagnostic{ID: id, Severity: diagnostics.SeverityError, Message: message, Path: path, Details: map[string]any{"pointer": pointer}}
}
