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
	Name        string
	Path        string
	OutcomePath string
	Status      string
	Completed   int
	Total       int
}

type Dependency struct {
	Source string
	Target string
	Type   string
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
