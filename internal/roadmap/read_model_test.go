package roadmap

import "testing"

func TestReadModelFromRootlineNormalizesDependenciesAndStatusRoles(t *testing.T) {
	tree := map[string]any{"root": map[string]any{"children": []any{
		map[string]any{"name": "T001-direct.md", "path": "T001-direct.md", "is_leaf": true, "estado": "Pending", "completed": float64(0), "total": float64(1)},
		map[string]any{"name": "O01-work", "path": "O01-work", "children": []any{
			map[string]any{"name": "T001-done.md", "path": "O01-work/T001-done.md", "is_leaf": true, "estado": "Done", "completed": float64(1), "total": float64(1)},
			map[string]any{"name": "T002-blocked.md", "path": "O01-work/T002-blocked.md", "is_leaf": true, "estado": "Ready", "completed": float64(0), "total": float64(1)},
		}},
	}}}
	query := map[string]any{"rows": []any{
		map[string]any{"path": "T001-direct.md", "frontmatter": map[string]any{"tipo": "task", "estado": "Pending"}},
		map[string]any{"path": "O01-work/T001-done.md", "frontmatter": map[string]any{"tipo": "task", "estado": "Done"}},
		map[string]any{"path": "O01-work/T002-blocked.md", "frontmatter": map[string]any{"tipo": "task", "estado": "Ready"}},
	}}
	graph := map[string]any{"edges": []any{map[string]any{"source": "O01-work/T002-blocked.md", "target": "O01-work/T001-done.md", "type": "blocked_by"}}, "cycles": []any{}, "broken_links": []any{}}

	model, diagnostics := ReadModelFromRootline(tree, query, graph, StatusRoleConfig{Done: []string{"Done"}, Active: []string{"Pending", "Ready"}})
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	if len(model.Tasks) != 3 {
		t.Fatalf("Tasks = %#v", model.Tasks)
	}
	blocked := model.TaskByPath["O01-work/T002-blocked.md"]
	if !blocked.Active || blocked.Done || len(blocked.Dependencies) != 1 || blocked.Dependencies[0] != "O01-work/T001-done.md" {
		t.Fatalf("blocked task = %#v", blocked)
	}
	done := model.TaskByPath["O01-work/T001-done.md"]
	if !done.Done || len(done.Blocks) != 1 || done.Blocks[0] != "O01-work/T002-blocked.md" {
		t.Fatalf("done task = %#v", done)
	}
}

func TestReadModelFromRootlineFallsBackToQueryRowsWhenTreeOmitsCustomStatuses(t *testing.T) {
	query := map[string]any{"rows": []any{
		map[string]any{"path": "O01-work/T001-ready.md", "frontmatter": map[string]any{"tipo": "task", "estado": "Ready"}},
	}}
	model, diagnostics := ReadModelFromRootline(map[string]any{"root": map[string]any{}}, query, map[string]any{"edges": []any{}, "cycles": []any{}, "broken_links": []any{}}, StatusRoleConfig{Done: []string{"Done"}, Active: []string{"Ready"}})
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	if len(model.Tasks) != 1 || model.Tasks[0].Name != "T001-ready.md" || model.Tasks[0].OutcomePath != "O01-work" || !model.Tasks[0].Active {
		t.Fatalf("model = %#v", model)
	}
}

func TestReadModelFromRootlineReportsGraphDiagnostics(t *testing.T) {
	model, diagnostics := ReadModelFromRootline(map[string]any{"root": map[string]any{}}, map[string]any{"rows": []any{}}, map[string]any{"cycles": []any{[]any{"a", "b"}}, "broken_links": []any{map[string]any{"source": "a", "target": "missing", "type": "blocked_by"}}}, StatusRoleConfig{})
	if len(model.Tasks) != 0 {
		t.Fatalf("Tasks = %#v", model.Tasks)
	}
	if len(diagnostics) != 2 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
}
