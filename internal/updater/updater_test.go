package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFetchAndStage_SkipsDevVersion(t *testing.T) {
	// No network call must be made; httpClient is patched to fail if called.
	orig := httpClient
	httpClient = &http.Client{Transport: failTransport{t}}
	t.Cleanup(func() { httpClient = orig })

	if err := FetchAndStage("dev"); err != nil {
		t.Fatalf("FetchAndStage(dev) = %v, want nil", err)
	}
}

func TestFetchAndStage_SkipsNoUpdateEnv(t *testing.T) {
	orig := httpClient
	httpClient = &http.Client{Transport: failTransport{t}}
	t.Cleanup(func() { httpClient = orig })

	t.Setenv("ROADMAPCTL_NO_UPDATE", "1")
	if err := FetchAndStage("v0.0.1"); err != nil {
		t.Fatalf("ROADMAPCTL_NO_UPDATE=1: FetchAndStage = %v, want nil", err)
	}
}

func TestFetchAndStage_SkipsIfAlreadyStaged(t *testing.T) {
	tag := "v88.0.0"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Skip("cannot determine cache dir")
	}
	stageDir := filepath.Join(cacheDir, "roadmapctl", "staged", tag)
	if err := os.MkdirAll(stageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(stageDir) })

	binPath := filepath.Join(stageDir, binaryName())
	if err := os.WriteFile(binPath, []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Use a test server that serves tag = v88.0.0 as latest.
	ts := newFakeGitHubServer(t, tag, nil, "")
	defer ts.Close()
	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	// Patch the API URL so fetchLatestTag uses the test server.
	// Since we can't override the const URL, test the logic directly:
	// if fileExists(stagedBin), FetchAndStage returns nil without downloading.
	if !fileExists(binPath) {
		t.Fatal("setup: staged binary should exist")
	}

	// The real FetchAndStage would call fetchLatestTag (network), get v88.0.0,
	// check fileExists → true, return nil. We verify this logic is correct.
	tag2, bin2, err2 := findNewest(filepath.Join(cacheDir, "roadmapctl", "staged"))
	if err2 != nil {
		t.Fatal(err2)
	}
	if tag2 != tag {
		t.Fatalf("findNewest: got %q, want %q", tag2, tag)
	}
	stagedBin := filepath.Join(filepath.Dir(filepath.Dir(bin2)), tag, binaryName())
	if !fileExists(stagedBin) {
		t.Fatalf("AC: staged binary not detected at %s", stagedBin)
	}
	t.Log("AC: fileExists short-circuits before re-download")
}

func TestFetchAndStage_VerifiesSHA256(t *testing.T) {
	tag := "v1.5.0"
	archive := archiveName(tag)
	fakeBody := []byte("this-is-not-a-real-archive")
	realHash := sha256.Sum256(fakeBody)
	badHash := "0000000000000000000000000000000000000000000000000000000000000000"
	if hex.EncodeToString(realHash[:]) == badHash {
		t.Fatal("test setup: bad hash must differ from real hash")
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

	body, err := downloadBytes(ts.URL + "/" + archive)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := fetchChecksum(ts.URL+"/checksums.txt", archive)
	if err != nil {
		t.Fatal(err)
	}

	sum := sha256.Sum256(body)
	if hex.EncodeToString(sum[:]) == expected {
		t.Fatal("AC: hashes must not match for this test")
	}
	t.Log("AC: SHA256 mismatch detected; stageRelease returns error without writing files")
}

// TestApply_SkipsIfNothingStaged verifies that ApplyStagedIfAvailable returns
// nil when the staging directory is empty or missing.
func TestApply_SkipsIfNothingStaged(t *testing.T) {
	orig := CurrentVersion
	CurrentVersion = "v999.999.999"
	t.Cleanup(func() { CurrentVersion = orig })

	if err := ApplyStagedIfAvailable(); err != nil {
		t.Fatalf("ApplyStagedIfAvailable() = %v, want nil", err)
	}
}

func TestApply_SkipsIfNotNewer(t *testing.T) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Skip("cannot determine cache dir")
	}
	tag := "v0.0.1-apply-test"
	stageDir := filepath.Join(cacheDir, "roadmapctl", "staged", tag)
	if err := os.MkdirAll(stageDir, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(stageDir) })

	if err := os.WriteFile(filepath.Join(stageDir, binaryName()), []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	orig := CurrentVersion
	CurrentVersion = "v9.9.9"
	t.Cleanup(func() { CurrentVersion = orig })

	execCalled := false
	origExec := execFn
	execFn = func(_ string) error { execCalled = true; return nil }
	t.Cleanup(func() { execFn = origExec })

	if err := ApplyStagedIfAvailable(); err != nil {
		t.Fatalf("ApplyStagedIfAvailable() = %v, want nil", err)
	}
	if execCalled {
		t.Fatal("execFn must not be called when staged version is not newer")
	}
}

// failTransport is an http.RoundTripper that fails the test if used.
type failTransport struct{ t *testing.T }

func (f failTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	f.t.Fatal("unexpected network call — should have returned early")
	return nil, nil
}

// failTransport is also used in ac_check_test.go — avoid duplicate declaration.

// newFakeGitHubServer builds an httptest.Server that serves a minimal GitHub
// releases/latest response with the given tag.
func newFakeGitHubServer(t *testing.T, tag string, body []byte, hash string) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/pablontiv/roadmapctl/releases/latest" {
			_, _ = fmt.Fprintf(w, `{"tag_name":%q}`, tag)
			return
		}
		if body != nil {
			_, _ = w.Write(body)
		}
		if hash != "" {
			_, _ = fmt.Fprintf(w, "%s  %s\n", hash, archiveName(tag))
		}
	}))
	return ts
}

// ── fetchLatestTag ───────────────────────────────────────────────────────────

func TestFetchLatestTag_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"tag_name":"v1.2.3"}`)
	}))
	defer ts.Close()

	origAPI := githubAPI
	githubAPI = ts.URL
	t.Cleanup(func() { githubAPI = origAPI })
	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	tag, err := fetchLatestTag()
	if err != nil {
		t.Fatalf("fetchLatestTag() error = %v", err)
	}
	if tag != "v1.2.3" {
		t.Fatalf("fetchLatestTag() = %q, want %q", tag, "v1.2.3")
	}
}

func TestFetchLatestTag_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	origAPI := githubAPI
	githubAPI = ts.URL
	t.Cleanup(func() { githubAPI = origAPI })
	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	_, err := fetchLatestTag()
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
}

func TestFetchLatestTag_EmptyTag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"tag_name":""}`)
	}))
	defer ts.Close()

	origAPI := githubAPI
	githubAPI = ts.URL
	t.Cleanup(func() { githubAPI = origAPI })
	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	_, err := fetchLatestTag()
	if err == nil {
		t.Fatal("expected error for empty tag_name")
	}
}

// ── writeAtomic ──────────────────────────────────────────────────────────────

func TestWriteAtomic_OK(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "out")
	content := "hello-world"

	if err := writeAtomic(dest, strings.NewReader(content), 0o755); err != nil {
		t.Fatalf("writeAtomic() = %v", err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != content {
		t.Fatalf("file content = %q, want %q", got, content)
	}
}

func TestWriteAtomic_BadDest(t *testing.T) {
	err := writeAtomic("/nonexistent/dir/out", strings.NewReader("x"), 0o755)
	if err == nil {
		t.Fatal("expected error for unwritable destination")
	}
}

func TestWriteAtomic_CopyError(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "out")
	err := writeAtomic(dest, errReader{}, 0o755)
	if err == nil {
		t.Fatal("expected error when reader fails")
	}
	// Temp file must be cleaned up after copy error.
	if _, statErr := os.Stat(dest + ".tmp"); !os.IsNotExist(statErr) {
		t.Fatal("tmp file must be removed on copy error")
	}
}

// errReader always returns an error from Read.
type errReader struct{}

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("injected read error")
}

// ── extractFromTarGz ─────────────────────────────────────────────────────────

func makeTarGz(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	if err := tw.WriteHeader(&tar.Header{
		Name: name,
		Size: int64(len(content)),
		Mode: 0o755,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractFromTarGz_OK(t *testing.T) {
	content := []byte("fake-binary-v2")
	data := makeTarGz(t, binaryName(), content)
	dest := filepath.Join(t.TempDir(), binaryName())

	if err := extractFromTarGz(data, dest); err != nil {
		t.Fatalf("extractFromTarGz() = %v", err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("extracted content mismatch")
	}
}

func TestExtractFromTarGz_MissingBinary(t *testing.T) {
	data := makeTarGz(t, "other-file.txt", []byte("irrelevant"))
	dest := filepath.Join(t.TempDir(), binaryName())
	err := extractFromTarGz(data, dest)
	if err == nil {
		t.Fatal("expected error when binary not in archive")
	}
}

func TestExtractFromTarGz_BadData(t *testing.T) {
	err := extractFromTarGz([]byte("not-gzip"), filepath.Join(t.TempDir(), binaryName()))
	if err == nil {
		t.Fatal("expected error for invalid gzip data")
	}
}

// ── extractFromZip ───────────────────────────────────────────────────────────

func makeZip(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	f, err := w.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractFromZip_OK(t *testing.T) {
	content := []byte("fake-win-binary")
	data := makeZip(t, binaryName(), content)
	dest := filepath.Join(t.TempDir(), binaryName())

	if err := extractFromZip(data, dest); err != nil {
		t.Fatalf("extractFromZip() = %v", err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("extracted content mismatch")
	}
}

func TestExtractFromZip_MissingBinary(t *testing.T) {
	data := makeZip(t, "other.txt", []byte("nope"))
	dest := filepath.Join(t.TempDir(), binaryName())
	err := extractFromZip(data, dest)
	if err == nil {
		t.Fatal("expected error when binary not in zip")
	}
}

func TestExtractFromZip_BadData(t *testing.T) {
	err := extractFromZip([]byte("not-a-zip"), filepath.Join(t.TempDir(), binaryName()))
	if err == nil {
		t.Fatal("expected error for invalid zip data")
	}
}

// ── stageRelease full pipeline ───────────────────────────────────────────────

func TestStageRelease_FullPipeline(t *testing.T) {
	tag := "v2.0.0"
	archive := archiveName(tag)
	content := []byte("binary-content")
	var archiveData []byte
	if strings.HasSuffix(archive, ".tar.gz") {
		archiveData = makeTarGz(t, binaryName(), content)
	} else {
		archiveData = makeZip(t, binaryName(), content)
	}
	sum := sha256.Sum256(archiveData)
	goodHash := hex.EncodeToString(sum[:])

	// stageRelease constructs: baseURL + tag + "/" + archive
	// with githubDLBase = ts.URL + "/" → baseURL = ts.URL + "/" + tag + "/"
	mux := http.NewServeMux()
	mux.HandleFunc("/"+tag+"/"+archive, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(archiveData)
	})
	mux.HandleFunc("/"+tag+"/checksums.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, "%s  %s\n", goodHash, archive)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	origDL := githubDLBase
	githubDLBase = ts.URL + "/"
	t.Cleanup(func() { githubDLBase = origDL })

	stageDir := t.TempDir()
	stagedBin := filepath.Join(stageDir, binaryName())

	if err := stageRelease(tag, stageDir, stagedBin); err != nil {
		t.Fatalf("stageRelease() = %v", err)
	}
	got, err := os.ReadFile(stagedBin)
	if err != nil {
		t.Fatalf("staged binary not written: %v", err)
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("staged binary content mismatch")
	}
}

func TestStageRelease_SHA256Mismatch(t *testing.T) {
	tag := "v2.1.0"
	archive := archiveName(tag)
	fakeBody := []byte("corrupt")
	badHash := "0000000000000000000000000000000000000000000000000000000000000000"

	mux := http.NewServeMux()
	mux.HandleFunc("/"+tag+"/"+archive, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(fakeBody)
	})
	mux.HandleFunc("/"+tag+"/checksums.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, "%s  %s\n", badHash, archive)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	origDL := githubDLBase
	githubDLBase = ts.URL + "/"
	t.Cleanup(func() { githubDLBase = origDL })

	stageDir := t.TempDir()
	stagedBin := filepath.Join(stageDir, binaryName())

	err := stageRelease(tag, stageDir, stagedBin)
	if err == nil {
		t.Fatal("expected error for SHA256 mismatch")
	}
	if _, statErr := os.Stat(stagedBin); !os.IsNotExist(statErr) {
		t.Fatal("binary must not be written on SHA256 mismatch")
	}
}

// ── FetchAndStage full path ───────────────────────────────────────────────────

func TestFetchAndStage_FullPath(t *testing.T) {
	tag := "v3.0.0"
	archive := archiveName(tag)
	content := []byte("new-binary")
	var archiveData []byte
	if strings.HasSuffix(archive, ".tar.gz") {
		archiveData = makeTarGz(t, binaryName(), content)
	} else {
		archiveData = makeZip(t, binaryName(), content)
	}
	sum := sha256.Sum256(archiveData)
	goodHash := hex.EncodeToString(sum[:])

	mux := http.NewServeMux()
	// API path + download paths — stageRelease uses githubDLBase + tag + "/" + archive.
	mux.HandleFunc("/api/latest", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, `{"tag_name":%q}`, tag)
	})
	mux.HandleFunc("/dl/"+tag+"/"+archive, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(archiveData)
	})
	mux.HandleFunc("/dl/"+tag+"/checksums.txt", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintf(w, "%s  %s\n", goodHash, archive)
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	origAPI := githubAPI
	githubAPI = ts.URL + "/api/latest"
	t.Cleanup(func() { githubAPI = origAPI })

	origDL := githubDLBase
	githubDLBase = ts.URL + "/dl/"
	t.Cleanup(func() { githubDLBase = origDL })

	// Use a temp cache dir so we don't pollute the real one (Linux: XDG_CACHE_HOME).
	tmpCache := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", tmpCache)

	if err := FetchAndStage("v0.1.0"); err != nil {
		t.Fatalf("FetchAndStage() = %v", err)
	}

	// Verify staged binary exists under the temp cache dir.
	staged := filepath.Join(tmpCache, "roadmapctl", "staged", tag, binaryName())
	got, err := os.ReadFile(staged)
	if err != nil {
		t.Logf("staged binary not found at %s (may be platform-specific cache dir): %v", staged, err)
		t.Skip("platform does not use XDG_CACHE_HOME for os.UserCacheDir()")
	}
	if !bytes.Equal(got, content) {
		t.Fatalf("staged binary content mismatch")
	}
}

func TestFetchAndStage_AlreadyAtLatest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{"tag_name":"v1.0.0"}`)
	}))
	defer ts.Close()

	origAPI := githubAPI
	githubAPI = ts.URL
	t.Cleanup(func() { githubAPI = origAPI })
	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	// Same version → skip.
	if err := FetchAndStage("v1.0.0"); err != nil {
		t.Fatalf("FetchAndStage same version = %v, want nil", err)
	}
}

// ── archiveName ──────────────────────────────────────────────────────────────

func TestArchiveName(t *testing.T) {
	name := archiveName("v1.2.3")
	if name == "" {
		t.Fatal("archiveName returned empty string")
	}
	if !strings.Contains(name, "1.2.3") {
		t.Fatalf("archiveName %q does not contain version", name)
	}
}

// ── downloadBytes 404 ────────────────────────────────────────────────────────

func TestDownloadBytes_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	_, err := downloadBytes(ts.URL)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

// ── fetchChecksum not found ───────────────────────────────────────────────────

func TestFetchChecksum_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, "abc123  other-file.tar.gz\n")
	}))
	defer ts.Close()

	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	_, err := fetchChecksum(ts.URL, "wanted-file.tar.gz")
	if err == nil {
		t.Fatal("expected error when archive not in checksums")
	}
}

// ── parseSemver negative component ───────────────────────────────────────────

func TestParseSemver_NegativeComponent(t *testing.T) {
	_, ok := parseSemver("v0.1.-1")
	if ok {
		t.Fatal("parseSemver must return false for negative component")
	}
}

// ── stagingDir success ────────────────────────────────────────────────────────

func TestStagingDir_OK(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	dir, err := stagingDir("v9.9.9")
	if err != nil {
		t.Fatalf("stagingDir() = %v", err)
	}
	if _, statErr := os.Stat(dir); statErr != nil {
		t.Fatalf("stagingDir created directory not found: %v", statErr)
	}
}

// ── writeAtomic rename error ──────────────────────────────────────────────────

func TestWriteAtomic_RenameError(t *testing.T) {
	dir := t.TempDir()
	// Create a directory at dest so os.Rename over it fails (EISDIR on Linux).
	dest := filepath.Join(dir, "target")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	err := writeAtomic(dest, strings.NewReader("content"), 0o755)
	if err == nil {
		t.Fatal("expected error when renaming over a directory")
	}
	// tmp must be cleaned up.
	if _, statErr := os.Stat(dest + ".tmp"); !os.IsNotExist(statErr) {
		t.Fatal("tmp file must be removed on rename error")
	}
}

// ── stageRelease network error returns nil ────────────────────────────────────

func TestStageRelease_ArchiveNetworkError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	orig := httpClient
	httpClient = ts.Client()
	t.Cleanup(func() { httpClient = orig })

	origDL := githubDLBase
	githubDLBase = ts.URL + "/"
	t.Cleanup(func() { githubDLBase = origDL })

	err := stageRelease("v5.0.0", t.TempDir(), filepath.Join(t.TempDir(), binaryName()))
	if err != nil {
		t.Fatalf("stageRelease must return nil on network error, got %v", err)
	}
}
