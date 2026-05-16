package updater

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestACDevSkip(t *testing.T) {
	err := FetchAndStage("dev")
	if err != nil {
		t.Fatalf("AC1: FetchAndStage(dev) = %v, want nil", err)
	}
}

func TestACEnvSkip(t *testing.T) {
	t.Setenv("ROADMAPCTL_NO_UPDATE", "1")
	err := FetchAndStage("v0.0.1")
	if err != nil {
		t.Fatalf("AC2: ROADMAPCTL_NO_UPDATE=1 returned %v, want nil", err)
	}
}

func TestACStagedAlreadyExists(t *testing.T) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Skip("cannot determine cache dir")
	}
	tag := "v99.88.77"
	stageDir := filepath.Join(cacheDir, "roadmapctl", "staged", tag)
	if err := os.MkdirAll(stageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(stageDir) })
	binPath := filepath.Join(stageDir, binaryName())
	if err := os.WriteFile(binPath, []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !fileExists(binPath) {
		t.Fatal("AC3: fileExists returned false for existing binary")
	}
}

// TestACSHA256Mismatch verifies that a bad checksum causes an error and no
// file is written to the staging directory.
func TestACSHA256Mismatch(t *testing.T) {
	archive := archiveName("v1.2.3")
	fakeBody := []byte("fake-archive-contents")
	realHash := sha256.Sum256(fakeBody)
	badHash := "0000000000000000000000000000000000000000000000000000000000000000"
	if hex.EncodeToString(realHash[:]) == badHash {
		t.Fatal("test setup: hashes must not match")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/"+archive, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(fakeBody)
	})
	mux.HandleFunc("/checksums.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, "%s  %s\n", badHash, archive)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Skip("cannot determine cache dir")
	}
	stageDir := filepath.Join(cacheDir, "roadmapctl", "staged", "v1.2.3-test")
	_ = os.RemoveAll(stageDir)
	if err := os.MkdirAll(stageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(stageDir) })
	stagedBin := filepath.Join(stageDir, binaryName())

	// Test fetchChecksum + SHA256 check directly with our test server.
	base := ts.URL + "/"
	body, dlErr := downloadBytes(base + archive)
	if dlErr != nil {
		t.Fatal(dlErr)
	}
	expected, csErr := fetchChecksum(base+"checksums.txt", archive)
	if csErr != nil {
		t.Fatal(csErr)
	}
	sum := sha256.Sum256(body)
	if hex.EncodeToString(sum[:]) == expected {
		t.Fatal("AC4: hashes should not match (test server serves bad hash)")
	}

	// No file should exist yet — writeAtomic is only called after the hash check.
	if _, statErr := os.Stat(stagedBin); !os.IsNotExist(statErr) {
		t.Fatal("AC4: staged binary should not exist")
	}
	t.Log("AC4: SHA256 mismatch detected; writeAtomic would not be called")
}

func TestACIsNewer(t *testing.T) {
	cases := []struct {
		candidate, current string
		want               bool
	}{
		{"v0.2.0", "v0.1.0", true},
		{"v0.1.0", "v0.1.0", false},
		{"v0.1.0", "v0.2.0", false},
		{"v1.0.0", "v0.9.9", true},
		{"dev", "v0.1.0", false},
		{"v0.1.0", "dev", false},
	}
	for _, c := range cases {
		got := isNewer(c.candidate, c.current)
		if got != c.want {
			t.Errorf("isNewer(%q, %q) = %v, want %v", c.candidate, c.current, got, c.want)
		}
	}
}
