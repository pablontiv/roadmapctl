package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pablontiv/roadmapctl/internal/materialize"
)

func TestPlanPathsCommandIntegration(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir := t.TempDir()
	roadmapRoot := filepath.Join(tmpDir, "roadmap")
	if err := os.MkdirAll(roadmapRoot, 0o755); err != nil {
		t.Fatalf("failed to create roadmap root: %v", err)
	}

	// Create bootstrap files
	stemPath := filepath.Join(roadmapRoot, ".stem")
	if err := os.WriteFile(stemPath, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create .stem: %v", err)
	}

	// Create a test input file
	input := materialize.PathPlanInput{
		Version: 1,
		Kind:    materialize.PathPlanKind,
		Items: []materialize.PathPlanItem{
			{Type: "outcome", Slug: "test-outcome"},
			{Type: "task", Slug: "test-task"},
		},
	}

	inputPath := filepath.Join(tmpDir, "input.json")
	inputData, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}
	if err := os.WriteFile(inputPath, inputData, 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	// Run the command
	options := Options{
		Repo:        tmpDir,
		RoadmapRoot: roadmapRoot,
		Output:      "json",
		Timeout:     10 * time.Second,
	}
	report := runPlanPaths(nil, options, inputPath)

	// Verify the report
	if report.Kind != "roadmapctl/plan-paths" {
		t.Fatalf("expected kind roadmapctl/plan-paths, got %s", report.Kind)
	}

	if report.Version != 1 {
		t.Fatalf("expected version 1, got %d", report.Version)
	}

	if len(report.Result.Paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(report.Result.Paths))
	}

	if report.Result.Paths[0].Type != "outcome" {
		t.Fatalf("expected first path type outcome, got %s", report.Result.Paths[0].Type)
	}

	if report.Result.Paths[1].Type != "task" {
		t.Fatalf("expected second path type task, got %s", report.Result.Paths[1].Type)
	}
}

func TestPlanPathsCommandInvalidInput(t *testing.T) {
	tmpDir := t.TempDir()
	roadmapRoot := filepath.Join(tmpDir, "roadmap")
	if err := os.MkdirAll(roadmapRoot, 0o755); err != nil {
		t.Fatalf("failed to create roadmap root: %v", err)
	}

	// Create invalid input file
	inputPath := filepath.Join(tmpDir, "input.json")
	if err := os.WriteFile(inputPath, []byte("invalid json"), 0o644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	// Run the command
	options := Options{
		Repo:        tmpDir,
		RoadmapRoot: roadmapRoot,
		Output:      "json",
		Timeout:     10 * time.Second,
	}
	report := runPlanPaths(nil, options, inputPath)

	// Verify that diagnostics were generated
	if len(report.Diagnostics) == 0 {
		t.Fatalf("expected diagnostics for invalid input")
	}
}

func TestRenderPlanPathsJSON(t *testing.T) {
	report := pathPlanReport{
		Version: 1,
		Kind:    "roadmapctl/plan-paths",
		Root:    "/test/root",
		RoadmapRoot: "/test/root/roadmap",
		Result: materialize.PathPlanResult{
			Version: 1,
			Kind:    materialize.PathPlanResultKind,
			Paths: []materialize.PathPlanEntry{
				{Path: "O01-test/README.md", Operation: "create", Type: "outcome"},
			},
			Collisions:  []materialize.PathPlanCollision{},
			Diagnostics: []materialize.PathPlanDiagnostic{},
		},
	}

	stdout := &bytes.Buffer{}
	exitCode := renderPlanPaths(report, "json", stdout, io.Discard)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	// Verify JSON output is valid
	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if result["kind"] != "roadmapctl/plan-paths" {
		t.Fatalf("expected kind roadmapctl/plan-paths in output")
	}
}

func TestRenderPlanPathsText(t *testing.T) {
	report := pathPlanReport{
		Version: 1,
		Kind:    "roadmapctl/plan-paths",
		Root:    "/test/root",
		RoadmapRoot: "/test/root/roadmap",
		Result: materialize.PathPlanResult{
			Version: 1,
			Kind:    materialize.PathPlanResultKind,
			Paths: []materialize.PathPlanEntry{
				{Path: "O01-test/README.md", Operation: "create", Type: "outcome"},
				{Path: "O01-test/T001-task.md", Operation: "create", Type: "task"},
			},
			Collisions:  []materialize.PathPlanCollision{},
			Diagnostics: []materialize.PathPlanDiagnostic{},
		},
	}

	stdout := &bytes.Buffer{}
	exitCode := renderPlanPaths(report, "text", stdout, io.Discard)

	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	if len(output) == 0 {
		t.Fatalf("expected text output")
	}

	if !strings.Contains(output, "roadmapctl/plan-paths") {
		t.Fatalf("expected kind in output")
	}

	if !strings.Contains(output, "O01-test/README.md") {
		t.Fatalf("expected outcome path in output")
	}

	if !strings.Contains(output, "O01-test/T001-task.md") {
		t.Fatalf("expected task path in output")
	}
}
