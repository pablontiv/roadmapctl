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
