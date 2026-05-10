package config

import "testing"

func TestConfigDiffersComparesOperationalFields(t *testing.T) {
	left := defaultConfig("/repo")
	right := defaultConfig("/repo")
	if configDiffers(left, right) {
		t.Fatal("identical default configs differ")
	}

	right.RequiredCodeCoverage = 90
	if !configDiffers(left, right) {
		t.Fatal("RequiredCodeCoverage change was not detected")
	}

	right = defaultConfig("/repo")
	right.DoneStatuses = []string{"Completed"}
	if !configDiffers(left, right) {
		t.Fatal("DoneStatuses change was not detected")
	}
}

func TestScalarHelpersParseNumbersAndRejectUnsupportedTypes(t *testing.T) {
	for _, value := range []any{12, int64(12), 12.9} {
		if got, ok := intValue(value); !ok || got != 12 {
			t.Fatalf("intValue(%T) = %v %v", value, got, ok)
		}
	}
	if _, ok := intValue("12"); ok {
		t.Fatal("intValue(string) ok = true, want false")
	}
	if got, ok := floatValue(12); !ok || got != 12 {
		t.Fatalf("floatValue(int) = %v %v", got, ok)
	}
	if got, ok := floatValue(int64(13)); !ok || got != 13 {
		t.Fatalf("floatValue(int64) = %v %v", got, ok)
	}
	if _, ok := floatValue("13"); ok {
		t.Fatal("floatValue(string) ok = true, want false")
	}
	if got := parseScalar("91.5"); got != 91.5 {
		t.Fatalf("parseScalar float = %#v", got)
	}
}
