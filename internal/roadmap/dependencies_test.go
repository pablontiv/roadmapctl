package roadmap

import (
	"context"
	"errors"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/diagnostics"
	"github.com/pablontiv/roadmapctl/internal/rootlinecli"
)

func TestCheckRootlineDetectsCycleFromGraphJSON(t *testing.T) {
	client := &fakeRootlineClient{
		validate: map[string]any{"version": float64(1), "kind": "rootline/validate-batch", "summary": map[string]any{"invalid": float64(0)}},
		describe: map[string]any{"values": []any{"Pending", "Specified", "In Progress", "Completed", "Blocked", "Obsolete"}},
		query:    map[string]any{"rows": []any{}},
		graph: map[string]any{
			"cycles": []any{[]any{"O01-work/T001-a.md", "O01-work/T002-b.md"}},
		},
	}

	diagnostics, err := CheckRootline(context.Background(), client, RootlineCheckOptions{RoadmapRoot: "/repo/docs/roadmap", LeafFilter: "isIndex == false", AllowedStatuses: []string{"Pending", "Specified", "In Progress", "Completed", "Blocked", "Obsolete"}})
	if err != nil {
		t.Fatalf("CheckRootline error = %v", err)
	}
	assertHasDiagnostic(t, diagnostics, DiagnosticGraphCycle, "")
}

func TestCheckRootlineDetectsBrokenBlockedByFromGraphJSON(t *testing.T) {
	client := &fakeRootlineClient{
		validate: map[string]any{"version": float64(1), "kind": "rootline/validate-batch", "summary": map[string]any{"invalid": float64(0)}},
		describe: map[string]any{"values": []any{"Pending", "Specified", "In Progress", "Completed", "Blocked", "Obsolete"}},
		query:    map[string]any{"rows": []any{}},
		graph: map[string]any{
			"broken_links": []any{map[string]any{"source": "O01-work/T001-task.md", "target": "O01-work/T999-missing.md", "type": "blocked_by", "line": float64(6)}},
		},
	}

	diagnostics, err := CheckRootline(context.Background(), client, RootlineCheckOptions{RoadmapRoot: "/repo/docs/roadmap", LeafFilter: "isIndex == false", AllowedStatuses: []string{"Pending", "Specified", "In Progress", "Completed", "Blocked", "Obsolete"}})
	if err != nil {
		t.Fatalf("CheckRootline error = %v", err)
	}
	assertHasDiagnostic(t, diagnostics, diagnosticsPackageInvalidBlockedBy(), "O01-work/T001-task.md")
}

func TestCheckRootlineDetectsStatusOutsideSchemaOrConfig(t *testing.T) {
	client := &fakeRootlineClient{
		validate: map[string]any{"version": float64(1), "kind": "rootline/validate-batch", "summary": map[string]any{"invalid": float64(0)}},
		describe: map[string]any{"values": []any{"Pending", "Completed"}},
		query: map[string]any{"rows": []any{
			map[string]any{"path": "O01-work/T001-task.md", "frontmatter": map[string]any{"estado": "Bogus", "tipo": "task"}},
		}},
		graph: map[string]any{},
	}

	diagnostics, err := CheckRootline(context.Background(), client, RootlineCheckOptions{RoadmapRoot: "/repo/docs/roadmap", LeafFilter: "isIndex == false", AllowedStatuses: []string{"Pending", "Completed"}})
	if err != nil {
		t.Fatalf("CheckRootline error = %v", err)
	}
	assertHasDiagnostic(t, diagnostics, DiagnosticStatusUnknown, "O01-work/T001-task.md")
}

func TestCheckRootlineAllowsSchemaStatusWithoutOperationalRole(t *testing.T) {
	client := &fakeRootlineClient{
		validate: map[string]any{"version": float64(1), "kind": "rootline/validate-batch", "summary": map[string]any{"invalid": float64(0)}},
		describe: map[string]any{"schema": map[string]any{
			"estado": map[string]any{"values": []any{"Pending", "Completed", "On Hold"}},
			"tipo":   map[string]any{"values": []any{"task", "outcome"}},
		}},
		query: map[string]any{"rows": []any{
			map[string]any{"path": "O01-work/T001-task.md", "frontmatter": map[string]any{"estado": "On Hold", "tipo": "task"}},
		}},
		graph: map[string]any{},
	}

	diagnostics, err := CheckRootline(context.Background(), client, RootlineCheckOptions{RoadmapRoot: "/repo/docs/roadmap", LeafFilter: "isIndex == false", AllowedStatuses: []string{"Pending", "Completed"}})
	if err != nil {
		t.Fatalf("CheckRootline error = %v", err)
	}
	assertNoDiagnostic(t, diagnostics, DiagnosticStatusUnknown)
}

func TestCheckRootlineDetectsTypeOutsideSchema(t *testing.T) {
	client := &fakeRootlineClient{
		validate: map[string]any{"version": float64(1), "kind": "rootline/validate-batch", "summary": map[string]any{"invalid": float64(0)}},
		describe: map[string]any{"schema": map[string]any{
			"estado": map[string]any{"values": []any{"Pending", "Completed"}},
			"tipo":   map[string]any{"values": []any{"task", "outcome"}},
		}},
		query: map[string]any{"rows": []any{
			map[string]any{"path": "O01-work/T001-task.md", "frontmatter": map[string]any{"estado": "Pending", "tipo": "story"}},
		}},
		graph: map[string]any{},
	}

	diagnostics, err := CheckRootline(context.Background(), client, RootlineCheckOptions{RoadmapRoot: "/repo/docs/roadmap", LeafFilter: "isIndex == false", AllowedStatuses: []string{"Pending", "Completed"}})
	if err != nil {
		t.Fatalf("CheckRootline error = %v", err)
	}
	assertHasDiagnostic(t, diagnostics, DiagnosticTypeUnknown, "O01-work/T001-task.md")
}

func TestCheckRootlineMissingRootlineDiagnosticExit3(t *testing.T) {
	client := &fakeRootlineClient{err: &rootlinecli.Error{Kind: rootlinecli.ErrorMissingBinary, Message: "missing rootline", ExitCode: diagnostics.ExitEnvironment}}

	found, err := CheckRootline(context.Background(), client, RootlineCheckOptions{RoadmapRoot: "/repo/docs/roadmap", LeafFilter: "isIndex == false", AllowedStatuses: []string{"Pending"}})
	if err != nil {
		t.Fatalf("CheckRootline error = %v", err)
	}
	assertHasDiagnostic(t, found, diagnostics.DiagnosticRootlineMissing, "")
	if got := diagnostics.ExitCode(diagnostics.NewReport("roadmapctl/check", "/repo", "/repo/docs/roadmap", found), false); got != diagnostics.ExitEnvironment {
		t.Fatalf("ExitCode = %d, want %d", got, diagnostics.ExitEnvironment)
	}
}

func TestCheckRootlineParsesValidateJSONOnNonZeroExit(t *testing.T) {
	client := &fakeRootlineClient{
		validate:    map[string]any{"version": float64(1), "kind": "rootline/validate", "summary": map[string]any{"invalid": float64(1)}},
		validateErr: &rootlinecli.Error{Kind: rootlinecli.ErrorExecution, Message: "rootline command failed", Stderr: "validation failed", ExitCode: 1},
		describe:    map[string]any{"values": []any{"Pending"}},
		query:       map[string]any{"rows": []any{}},
		graph:       map[string]any{},
	}

	found, err := CheckRootline(context.Background(), client, RootlineCheckOptions{RoadmapRoot: "/repo/docs/roadmap", LeafFilter: "isIndex == false", AllowedStatuses: []string{"Pending"}})
	if err != nil {
		t.Fatalf("CheckRootline error = %v", err)
	}
	assertHasDiagnostic(t, found, DiagnosticRootlineValidateFailed, "")
	if !hasDiagnosticDetail(found, DiagnosticRootlineValidateFailed, "invalid", 1) {
		t.Fatalf("missing parsed invalid summary diagnostic in %#v", found)
	}
}

func TestCheckRootlineUsesGenericRootlineJSONCommands(t *testing.T) {
	client := &fakeRootlineClient{
		validate: map[string]any{"summary": map[string]any{"invalid": float64(0)}},
		describe: map[string]any{"values": []any{"Pending"}},
		query:    map[string]any{"rows": []any{}},
		graph:    map[string]any{},
	}

	_, err := CheckRootline(context.Background(), client, RootlineCheckOptions{RoadmapRoot: "docs/roadmap", LeafFilter: "isIndex == false", AllowedStatuses: []string{"Pending"}})
	if err != nil {
		t.Fatalf("CheckRootline error = %v", err)
	}
	want := []string{"validate", "describe", "query", "graph"}
	if len(client.calls) != len(want) {
		t.Fatalf("calls = %#v, want %#v", client.calls, want)
	}
	for i := range want {
		if client.calls[i] != want[i] {
			t.Fatalf("calls = %#v, want %#v", client.calls, want)
		}
	}
	if client.usedGraphCheck {
		t.Fatal("CheckRootline used graph --check, want graph JSON only")
	}
}

type fakeRootlineClient struct {
	validate map[string]any
	describe map[string]any
	query    map[string]any
	graph    map[string]any
	err      error

	validateErr error
	describeErr error
	queryErr    error
	graphErr    error

	calls          []string
	usedGraphCheck bool
}

func (f *fakeRootlineClient) Validate(ctx context.Context, paths ...string) (*rootlinecli.JSONResult, error) {
	f.calls = append(f.calls, "validate")
	return f.result(f.validate, f.validateErr)
}

func (f *fakeRootlineClient) Describe(ctx context.Context, target string, fields ...string) (*rootlinecli.JSONResult, error) {
	f.calls = append(f.calls, "describe")
	return f.result(f.describe, f.describeErr)
}

func (f *fakeRootlineClient) Query(ctx context.Context, root string, wheres ...string) (*rootlinecli.JSONResult, error) {
	f.calls = append(f.calls, "query")
	return f.result(f.query, f.queryErr)
}

func (f *fakeRootlineClient) Graph(ctx context.Context, root string, wheres ...string) (*rootlinecli.JSONResult, error) {
	f.calls = append(f.calls, "graph")
	for _, where := range wheres {
		if where == "--check" || where == "check" {
			f.usedGraphCheck = true
		}
	}
	return f.result(f.graph, f.graphErr)
}

func (f *fakeRootlineClient) result(decoded map[string]any, operationErr error) (*rootlinecli.JSONResult, error) {
	if f.err != nil {
		return nil, f.err
	}
	if decoded == nil {
		return nil, errors.New("missing fake response")
	}
	return &rootlinecli.JSONResult{Decoded: decoded}, operationErr
}

func diagnosticsPackageInvalidBlockedBy() string {
	return diagnostics.DiagnosticInvalidBlockedBy
}
