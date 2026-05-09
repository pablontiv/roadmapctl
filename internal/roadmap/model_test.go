package roadmap

import "testing"

func TestRoadmapContextFromTreeSupportsDirectAndOutcomeTasks(t *testing.T) {
	decoded := map[string]any{
		"root": map[string]any{
			"children": []any{
				map[string]any{"name": "T001-direct.md", "path": "T001-direct.md", "is_leaf": true, "estado": "Pending", "completed": float64(0), "total": float64(1)},
				map[string]any{"name": "O01-work", "path": "O01-work", "completed": float64(1), "total": float64(2), "children": []any{
					map[string]any{"name": "T001-first.md", "path": "O01-work/T001-first.md", "is_leaf": true, "estado": "Completed", "completed": float64(1), "total": float64(1)},
					map[string]any{"name": "T002-second.md", "path": "O01-work/T002-second.md", "is_leaf": true, "estado": "On Hold", "completed": float64(0), "total": float64(1)},
				}},
			},
		},
	}

	ctx, err := RoadmapContextFromTree(decoded)
	if err != nil {
		t.Fatalf("RoadmapContextFromTree error = %v", err)
	}
	if len(ctx.Outcomes) != 1 || ctx.Outcomes[0].Path != "O01-work" || ctx.Outcomes[0].Completed != 1 || ctx.Outcomes[0].Total != 2 {
		t.Fatalf("Outcomes = %#v", ctx.Outcomes)
	}
	if len(ctx.Tasks) != 3 {
		t.Fatalf("Tasks = %#v", ctx.Tasks)
	}
	if ctx.Tasks[0].Path != "T001-direct.md" || ctx.Tasks[0].OutcomePath != "" || ctx.Tasks[0].Status != "Pending" {
		t.Fatalf("direct task = %#v", ctx.Tasks[0])
	}
	if ctx.Tasks[2].Path != "O01-work/T002-second.md" || ctx.Tasks[2].OutcomePath != "O01-work" || ctx.Tasks[2].Status != "On Hold" {
		t.Fatalf("outcome task = %#v", ctx.Tasks[2])
	}
}
