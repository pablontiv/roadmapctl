package roadmap

import "github.com/pablontiv/roadmapctl/internal/diagnostics"

type TransitionRoles struct {
	DoneStatuses     []string
	ActiveStatuses   []string
	InProgressStatus string
	CompletedStatus  string
}

type TransitionResult struct {
	Allowed              bool
	Reasons              []string
	BlockingDependencies []BlockingDependency
	Changes              []TransitionChange
	Diagnostics          []diagnostics.Diagnostic
	CurrentStatus        string
	TargetStatus         string
	Role                 string
}

type BlockingDependency struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

type TransitionChange struct {
	Path    string `json:"path"`
	Field   string `json:"field"`
	Before  string `json:"before"`
	After   string `json:"after"`
	Applied bool   `json:"applied"`
}

func CanStart(model ReadModel, roles TransitionRoles, path string) TransitionResult {
	task := model.TaskByPath[path]
	if task == nil {
		return transitionTaskNotFound(path)
	}
	result := TransitionResult{CurrentStatus: task.Status, TargetStatus: roles.InProgressStatus, Role: "in_progress"}
	if stringSet(roles.DoneStatuses)[task.Status] {
		result.Diagnostics = append(result.Diagnostics, diagnostics.Diagnostic{ID: diagnostics.DiagnosticTransitionAlreadyDone, Severity: diagnostics.SeverityWarning, Message: "task is already done", Path: path})
		result.Reasons = append(result.Reasons, "task is already done")
		return result
	}
	if !stringSet(roles.ActiveStatuses)[task.Status] {
		result.Diagnostics = append(result.Diagnostics, diagnostics.Diagnostic{ID: diagnostics.DiagnosticTransitionNotActive, Severity: diagnostics.SeverityWarning, Message: "task status is not active", Path: path})
		result.Reasons = append(result.Reasons, "task status is not active")
		return result
	}
	for _, dep := range task.Dependencies {
		dependency := model.TaskByPath[dep]
		status := ""
		if dependency != nil {
			status = dependency.Status
		}
		if dependency == nil || !stringSet(roles.DoneStatuses)[dependency.Status] {
			result.BlockingDependencies = append(result.BlockingDependencies, BlockingDependency{Path: dep, Status: status})
		}
	}
	if len(result.BlockingDependencies) > 0 {
		result.Diagnostics = append(result.Diagnostics, diagnostics.Diagnostic{ID: diagnostics.DiagnosticTransitionDependencyBlocked, Severity: diagnostics.SeverityWarning, Message: "task has dependencies outside done statuses", Path: path})
		result.Reasons = append(result.Reasons, "dependencies are not done")
		return result
	}
	result.Allowed = true
	result.Reasons = append(result.Reasons, "all dependencies are done")
	result.Changes = append(result.Changes, TransitionChange{Path: path, Field: "estado", Before: task.Status, After: roles.InProgressStatus})
	return result
}

func CanComplete(model ReadModel, roles TransitionRoles, path string) TransitionResult {
	task := model.TaskByPath[path]
	if task == nil {
		return transitionTaskNotFound(path)
	}
	result := TransitionResult{CurrentStatus: task.Status, TargetStatus: roles.CompletedStatus, Role: "completed"}
	if stringSet(roles.DoneStatuses)[task.Status] {
		result.Diagnostics = append(result.Diagnostics, diagnostics.Diagnostic{ID: diagnostics.DiagnosticTransitionAlreadyDone, Severity: diagnostics.SeverityWarning, Message: "task is already done", Path: path})
		result.Reasons = append(result.Reasons, "task is already done")
		return result
	}
	result.Allowed = true
	result.Reasons = append(result.Reasons, "task can be completed after caller verification")
	result.Changes = append(result.Changes, TransitionChange{Path: path, Field: "estado", Before: task.Status, After: roles.CompletedStatus})
	return result
}

func SetStatus(model ReadModel, roles TransitionRoles, path string, targetStatus string) TransitionResult {
	if targetStatus == roles.InProgressStatus {
		return CanStart(model, roles, path)
	}
	task := model.TaskByPath[path]
	if task == nil {
		return transitionTaskNotFound(path)
	}
	role := "custom"
	if targetStatus == roles.CompletedStatus {
		role = "completed"
	}
	result := TransitionResult{Allowed: true, CurrentStatus: task.Status, TargetStatus: targetStatus, Role: role, Reasons: []string{"explicit status change planned"}}
	result.Changes = append(result.Changes, TransitionChange{Path: path, Field: "estado", Before: task.Status, After: targetStatus})
	return result
}

func transitionTaskNotFound(path string) TransitionResult {
	return TransitionResult{Reasons: []string{"task not found"}, Diagnostics: []diagnostics.Diagnostic{{ID: diagnostics.DiagnosticTransitionTaskNotFound, Severity: diagnostics.SeverityError, Message: "task not found in roadmap read model", Path: path}}}
}
