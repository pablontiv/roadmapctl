package lint

import "testing"

func TestParseMarkdownExtractsHeadingsAndTasksTable(t *testing.T) {
	source := []byte(`# Outcome

Intro.

## Tasks

| Task | Description |
| --- | --- |
| [T001](T001-first.md) | First task |
| [T002](./nested/T002-second.md) | Second task |

## Acceptance Criteria

- Done
`)
	doc, err := ParseMarkdown(source)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Headings) != 3 {
		t.Fatalf("Headings = %#v", doc.Headings)
	}
	if doc.Headings[1].Level != 2 || doc.Headings[1].Text != "Tasks" || doc.Headings[1].StartLine != 5 {
		t.Fatalf("Tasks heading = %#v", doc.Headings[1])
	}
	table := doc.TableBySection("Tasks")
	if table == nil {
		t.Fatalf("missing Tasks table in %#v", doc.Tables)
	}
	if table.StartLine != 7 || len(table.Headers) != 2 || table.Headers[0] != "Task" || table.Headers[1] != "Description" {
		t.Fatalf("table = %#v", table)
	}
	if len(table.Rows) != 2 || table.Rows[0].Cells[0].Text != "T001" || table.Rows[0].Cells[0].Links[0].Destination != "T001-first.md" {
		t.Fatalf("rows = %#v", table.Rows)
	}
	if table.Rows[1].Cells[0].Links[0].Destination != "nested/T002-second.md" {
		t.Fatalf("link not normalized: %#v", table.Rows[1].Cells[0].Links[0])
	}
}

func TestParseMarkdownKeepsTaskSectionsReadOnly(t *testing.T) {
	source := []byte("# Task\n\n## Preserva\n\n- invariant\n\n## Criterios de Aceptación\n\n- AC\n")
	before := string(source)
	doc, err := ParseMarkdown(source)
	if err != nil {
		t.Fatal(err)
	}
	if string(source) != before {
		t.Fatalf("ParseMarkdown mutated input")
	}
	if !doc.HasHeading("Preserva") || !doc.HasHeading("Criterios de Aceptación") {
		t.Fatalf("missing headings: %#v", doc.Headings)
	}
}
