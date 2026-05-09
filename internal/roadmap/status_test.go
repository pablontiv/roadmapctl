package roadmap

import "testing"

func TestIntersectStringSetsHandlesEmptyAndMatches(t *testing.T) {
	right := map[string]bool{"Pending": true, "Completed": true}
	if got := intersectStringSets(map[string]bool{}, right); len(got) != 2 || !got["Pending"] || !got["Completed"] {
		t.Fatalf("empty-left intersection = %#v", got)
	}
	got := intersectStringSets(map[string]bool{"Pending": true, "Bogus": true}, right)
	if len(got) != 1 || !got["Pending"] {
		t.Fatalf("intersection = %#v", got)
	}
}

func TestModelHelpersHandleIntAndEmptyPath(t *testing.T) {
	if got := numberValue(3); got != 3 {
		t.Fatalf("numberValue(int) = %d", got)
	}
	if got := cleanSlashPath(""); got != "" {
		t.Fatalf("cleanSlashPath(empty) = %q", got)
	}
}
