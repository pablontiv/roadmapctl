package roadmap

import (
	"fmt"
	"path/filepath"
)

type StatusRole string

const (
	StatusRolePending    StatusRole = "pending"
	StatusRoleSpecified  StatusRole = "specified"
	StatusRoleInProgress StatusRole = "in-progress"
	StatusRoleCompleted  StatusRole = "completed"
	StatusRoleBlocked    StatusRole = "blocked"
	StatusRoleObsolete   StatusRole = "obsolete"
)

type RoadmapContext struct {
	Outcomes     []Outcome
	Tasks        []Task
	Dependencies []Dependency
	StatusRoles  map[StatusRole]string
}

type Outcome struct {
	Name      string
	Path      string
	Completed int
	Total     int
}

type Task struct {
	Name         string
	Path         string
	OutcomePath  string
	Status       string
	Type         string
	Completed    int
	Total        int
	Done         bool
	Active       bool
	Dependencies []string
	Blocks       []string
}

type Dependency struct {
	Source string
	Target string
	Type   string
}

type StatusRoleConfig struct {
	Done   []string
	Active []string
}

type ReadModel struct {
	Outcomes     []Outcome
	Tasks        []Task
	TaskByPath   map[string]*Task
	Dependencies []Dependency
}

func RoadmapContextFromTree(decoded map[string]any) (RoadmapContext, error) {
	root, ok := decoded["root"].(map[string]any)
	if !ok {
		return RoadmapContext{}, fmt.Errorf("rootline tree JSON missing root object")
	}
	ctx := RoadmapContext{StatusRoles: map[StatusRole]string{}}
	for _, childValue := range arrayValue(root["children"]) {
		child, ok := childValue.(map[string]any)
		if !ok {
			continue
		}
		if boolField(child, "is_leaf") {
			ctx.Tasks = append(ctx.Tasks, taskFromTreeNode(child, ""))
			continue
		}
		outcomePath := cleanSlashPath(stringField(child, "path"))
		ctx.Outcomes = append(ctx.Outcomes, Outcome{
			Name:      stringField(child, "name"),
			Path:      outcomePath,
			Completed: numberField(child, "completed"),
			Total:     numberField(child, "total"),
		})
		for _, taskValue := range arrayValue(child["children"]) {
			taskNode, ok := taskValue.(map[string]any)
			if !ok || !boolField(taskNode, "is_leaf") {
				continue
			}
			ctx.Tasks = append(ctx.Tasks, taskFromTreeNode(taskNode, outcomePath))
		}
	}
	return ctx, nil
}

func ReadModelFromRootline(tree map[string]any, query map[string]any, graph map[string]any, roles StatusRoleConfig) (ReadModel, []Diagnostic) {
	ctx, err := RoadmapContextFromTree(tree)
	if err != nil {
		ctx = RoadmapContext{}
	}
	model := ReadModel{Outcomes: ctx.Outcomes, Tasks: ctx.Tasks, TaskByPath: map[string]*Task{}}
	statusByPath := map[string]string{}
	typeByPath := map[string]string{}
	for _, rowValue := range arrayValue(query["rows"]) {
		row, ok := rowValue.(map[string]any)
		if !ok {
			continue
		}
		path := cleanSlashPath(stringField(row, "path"))
		frontmatter, _ := row["frontmatter"].(map[string]any)
		statusByPath[path] = stringField(frontmatter, "estado")
		typeByPath[path] = stringField(frontmatter, "tipo")
	}
	doneSet := stringSet(roles.Done)
	activeSet := stringSet(roles.Active)
	for i := range model.Tasks {
		path := model.Tasks[i].Path
		if status, ok := statusByPath[path]; ok {
			model.Tasks[i].Status = status
		}
		model.Tasks[i].Type = typeByPath[path]
		model.Tasks[i].Done = doneSet[model.Tasks[i].Status]
		model.Tasks[i].Active = activeSet[model.Tasks[i].Status]
		model.TaskByPath[path] = &model.Tasks[i]
	}
	for _, edgeValue := range arrayValue(graph["edges"]) {
		edge, ok := edgeValue.(map[string]any)
		if !ok || stringField(edge, "type") != "blocked_by" {
			continue
		}
		dep := Dependency{Source: cleanSlashPath(stringField(edge, "source")), Target: cleanSlashPath(stringField(edge, "target")), Type: "blocked_by"}
		model.Dependencies = append(model.Dependencies, dep)
		if task := model.TaskByPath[dep.Source]; task != nil {
			task.Dependencies = append(task.Dependencies, dep.Target)
		}
		if task := model.TaskByPath[dep.Target]; task != nil {
			task.Blocks = append(task.Blocks, dep.Source)
		}
	}
	return model, graphDiagnostics(graph)
}

func taskFromTreeNode(node map[string]any, outcomePath string) Task {
	return Task{
		Name:        stringField(node, "name"),
		Path:        cleanSlashPath(stringField(node, "path")),
		OutcomePath: outcomePath,
		Status:      stringField(node, "estado"),
		Completed:   numberField(node, "completed"),
		Total:       numberField(node, "total"),
	}
}

func boolField(fields map[string]any, key string) bool {
	value, _ := fields[key].(bool)
	return value
}

func numberField(fields map[string]any, key string) int {
	return numberValue(fields[key])
}

func numberValue(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}

func cleanSlashPath(path string) string {
	if path == "" {
		return ""
	}
	return filepath.ToSlash(filepath.Clean(path))
}
