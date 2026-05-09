package roadmap

import (
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
)

func TestCanStartBlocksOnIncompleteDependencyWithCustomDoneLabels(t *testing.T) {
	model := ReadModel{Tasks: []Task{
		{Path: "O01/T001-prereq.md", Status: "Ready"},
		{Path: "O01/T002-work.md", Status: "Ready", Dependencies: []string{"O01/T001-prereq.md"}},
	}, TaskByPath: map[string]*Task{}}
	for i := range model.Tasks {
		model.TaskByPath[model.Tasks[i].Path] = &model.Tasks[i]
	}
	roles := TransitionRoles{DoneStatuses: []string{"Done"}, ActiveStatuses: []string{"Ready"}, InProgressStatus: "Doing", CompletedStatus: "Done"}

	result := CanStart(model, roles, "O01/T002-work.md")
	if result.Allowed {
		t.Fatalf("Allowed = true, want false")
	}
	if len(result.BlockingDependencies) != 1 || result.BlockingDependencies[0].Path != "O01/T001-prereq.md" || result.BlockingDependencies[0].Status != "Ready" {
		t.Fatalf("BlockingDependencies = %#v", result.BlockingDependencies)
	}
	if len(result.Diagnostics) != 1 || result.Diagnostics[0].ID != diagnostics.DiagnosticTransitionDependencyBlocked {
		t.Fatalf("Diagnostics = %#v", result.Diagnostics)
	}
}

func TestCanStartAllowsWhenDependenciesUseCustomDoneLabels(t *testing.T) {
	model := ReadModel{Tasks: []Task{
		{Path: "O01/T001-prereq.md", Status: "Done", Done: true},
		{Path: "O01/T002-work.md", Status: "Ready", Active: true, Dependencies: []string{"O01/T001-prereq.md"}},
	}, TaskByPath: map[string]*Task{}}
	for i := range model.Tasks {
		model.TaskByPath[model.Tasks[i].Path] = &model.Tasks[i]
	}
	roles := TransitionRoles{DoneStatuses: []string{"Done"}, ActiveStatuses: []string{"Ready"}, InProgressStatus: "Doing", CompletedStatus: "Done"}

	result := CanStart(model, roles, "O01/T002-work.md")
	if !result.Allowed || result.TargetStatus != "Doing" || len(result.Changes) != 1 {
		t.Fatalf("result = %#v", result)
	}
}

func TestCanCompleteRequiresExistingNonDoneTask(t *testing.T) {
	model := ReadModel{Tasks: []Task{{Path: "T001.md", Status: "Ready", Active: true}}, TaskByPath: map[string]*Task{}}
	model.TaskByPath["T001.md"] = &model.Tasks[0]
	roles := TransitionRoles{DoneStatuses: []string{"Done"}, ActiveStatuses: []string{"Ready"}, CompletedStatus: "Done"}

	result := CanComplete(model, roles, "T001.md")
	if !result.Allowed || result.TargetStatus != "Done" || len(result.Changes) != 1 {
		t.Fatalf("result = %#v", result)
	}

	missing := CanComplete(model, roles, "missing.md")
	if missing.Allowed || len(missing.Diagnostics) != 1 || missing.Diagnostics[0].ID != diagnostics.DiagnosticTransitionTaskNotFound {
		t.Fatalf("missing = %#v", missing)
	}
}
