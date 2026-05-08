package testutil

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func FixturePath(t testing.TB, name string) string {
	t.Helper()
	return filepath.Join("..", "..", "testdata", "fixtures", name)
}

func AssertGoldenJSON(t testing.TB, goldenPath string, data []byte, replacements map[string]string) {
	t.Helper()
	normalized := NormalizeJSON(t, data, replacements)
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v\nactual:\n%s", goldenPath, err, normalized)
	}
	want = bytes.TrimSpace(want)
	normalized = bytes.TrimSpace(normalized)
	if !bytes.Equal(want, normalized) {
		t.Fatalf("golden mismatch %s\nwant:\n%s\ngot:\n%s", goldenPath, want, normalized)
	}
}

func NormalizeJSON(t testing.TB, data []byte, replacements map[string]string) []byte {
	t.Helper()
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("decode JSON: %v\n%s", err, string(data))
	}
	value = normalizeValue(value, replacements)
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		t.Fatalf("encode normalized JSON: %v", err)
	}
	return out.Bytes()
}

func normalizeValue(value any, replacements map[string]string) any {
	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, item := range v {
			if key == "rootline_version" {
				out[key] = "<rootline-version>"
				continue
			}
			out[key] = normalizeValue(item, replacements)
		}
		return out
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = normalizeValue(item, replacements)
		}
		return out
	case string:
		return NormalizePathString(v, replacements)
	default:
		return v
	}
}

func NormalizePathString(value string, replacements map[string]string) string {
	value = filepath.ToSlash(value)
	type pair struct{ from, to string }
	pairs := make([]pair, 0, len(replacements))
	for from, to := range replacements {
		from = filepath.ToSlash(from)
		pairs = append(pairs, pair{from: from, to: to})
	}
	// Replace longer prefixes first so nested paths normalize deterministically.
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if len(pairs[j].from) > len(pairs[i].from) {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
	for _, p := range pairs {
		value = strings.ReplaceAll(value, p.from, p.to)
	}
	return value
}

func AssertNoBackslashes(t testing.TB, value any) {
	t.Helper()
	if containsBackslash(value) {
		t.Fatalf("value contains backslash path separator: %#v", value)
	}
}

func containsBackslash(value any) bool {
	switch v := value.(type) {
	case string:
		return strings.Contains(v, `\`)
	case map[string]any:
		for _, item := range v {
			if containsBackslash(item) {
				return true
			}
		}
	case []any:
		for _, item := range v {
			if containsBackslash(item) {
				return true
			}
		}
	}
	return false
}

func DecodeJSON(t testing.TB, data []byte) map[string]any {
	t.Helper()
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("decode JSON: %v\n%s", err, string(data))
	}
	return value
}

func RequireDiagnosticID(t testing.TB, report map[string]any, id string) {
	t.Helper()
	diagnostics, _ := report["diagnostics"].([]any)
	for _, item := range diagnostics {
		diagnostic, _ := item.(map[string]any)
		if diagnostic["id"] == id {
			return
		}
	}
	t.Fatalf("missing diagnostic %s in %#v", id, diagnostics)
}

func AssertExit(t testing.TB, got int, want int, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	t.Helper()
	if got != want {
		t.Fatalf("exit = %d, want %d\nstdout:\n%s\nstderr:\n%s", got, want, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func GoldenPath(parts ...string) string {
	items := append([]string{"..", "..", "testdata", "golden"}, parts...)
	return filepath.Join(items...)
}
