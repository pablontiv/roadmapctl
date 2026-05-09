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

func TestCanStartReportsAlreadyDoneAndNotActive(t *testing.T) {
	model := ReadModel{Tasks: []Task{{Path: "done.md", Status: "Done", Done: true}, {Path: "hold.md", Status: "On Hold"}}, TaskByPath: map[string]*Task{}}
	for i := range model.Tasks {
		model.TaskByPath[model.Tasks[i].Path] = &model.Tasks[i]
	}
	roles := TransitionRoles{DoneStatuses: []string{"Done"}, ActiveStatuses: []string{"Ready"}, InProgressStatus: "Doing"}
	alreadyDone := CanStart(model, roles, "done.md")
	if alreadyDone.Allowed || alreadyDone.Diagnostics[0].ID != diagnostics.DiagnosticTransitionAlreadyDone {
		t.Fatalf("alreadyDone = %#v", alreadyDone)
	}
	notActive := CanStart(model, roles, "hold.md")
	if notActive.Allowed || notActive.Diagnostics[0].ID != diagnostics.DiagnosticTransitionNotActive {
		t.Fatalf("notActive = %#v", notActive)
	}
}

func TestSetStatusPlansCustomStatusAndDelegatesInProgressPolicy(t *testing.T) {
	model := ReadModel{Tasks: []Task{{Path: "T001.md", Status: "Ready", Active: true, Dependencies: []string{"missing.md"}}}, TaskByPath: map[string]*Task{}}
	model.TaskByPath["T001.md"] = &model.Tasks[0]
	roles := TransitionRoles{DoneStatuses: []string{"Done"}, ActiveStatuses: []string{"Ready"}, InProgressStatus: "Doing", CompletedStatus: "Done"}
	custom := SetStatus(model, roles, "T001.md", "On Hold")
	if !custom.Allowed || custom.Role != "custom" || custom.Changes[0].After != "On Hold" {
		t.Fatalf("custom = %#v", custom)
	}
	inProgress := SetStatus(model, roles, "T001.md", "Doing")
	if inProgress.Allowed || len(inProgress.BlockingDependencies) != 1 {
		t.Fatalf("inProgress = %#v", inProgress)
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
