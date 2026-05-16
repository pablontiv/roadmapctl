package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	repo   = "pablontiv/roadmapctl"
	binary = "roadmapctl"
)

var httpClient = &http.Client{Timeout: 60 * time.Second}

// FetchAndStage downloads the latest roadmapctl release into a staging
// directory without interrupting the current command. It returns nil silently
// on network errors; a SHA256 mismatch is returned as an error because it
// signals a data integrity problem. No files are written on error.
func FetchAndStage(currentVersion string) error {
	if currentVersion == "dev" || os.Getenv("ROADMAPCTL_NO_UPDATE") == "1" {
		return nil
	}

	tag, err := fetchLatestTag()
	if err != nil {
		return nil
	}

	if !isNewer(tag, currentVersion) {
		return nil
	}

	stageDir, err := stagingDir(tag)
	if err != nil {
		return nil
	}

	stagedBin := filepath.Join(stageDir, binaryName())
	if fileExists(stagedBin) {
		return nil
	}

	return stageRelease(tag, stageDir, stagedBin)
}

func fetchLatestTag() (string, error) {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil) //nolint:gosec
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api: status %d", resp.StatusCode)
	}
	var rel struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	if rel.TagName == "" {
		return "", fmt.Errorf("empty tag_name in github response")
	}
	return rel.TagName, nil
}

// isNewer reports whether candidate is strictly newer than current using semver.
func isNewer(candidate, current string) bool {
	a, aOK := parseSemver(candidate)
	b, bOK := parseSemver(current)
	if !aOK || !bOK {
		return false
	}
	for i := range a {
		if a[i] > b[i] {
			return true
		}
		if a[i] < b[i] {
			return false
		}
	}
	return false
}

func parseSemver(v string) ([3]int, bool) {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 3 {
		return [3]int{}, false
	}
	var nums [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return [3]int{}, false
		}
		nums[i] = n
	}
	return nums, true
}

func stagingDir(tag string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(cacheDir, "roadmapctl", "staged", tag)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func binaryName() string {
	if runtime.GOOS == "windows" {
		return binary + ".exe"
	}
	return binary
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func archiveName(tag string) string {
	version := strings.TrimPrefix(tag, "v")
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	if goos == "windows" {
		return fmt.Sprintf("%s_%s_%s_%s.zip", binary, version, goos, goarch)
	}
	return fmt.Sprintf("%s_%s_%s_%s.tar.gz", binary, version, goos, goarch)
}

func stageRelease(tag, stageDir, stagedBin string) error {
	archive := archiveName(tag)
	baseURL := "https://github.com/" + repo + "/releases/download/" + tag + "/"

	body, err := downloadBytes(baseURL + archive)
	if err != nil {
		return nil
	}

	expectedHash, err := fetchChecksum(baseURL+"checksums.txt", archive)
	if err != nil {
		return nil
	}

	sum := sha256.Sum256(body)
	if hex.EncodeToString(sum[:]) != expectedHash {
		return fmt.Errorf("SHA256 mismatch for %s", archive)
	}

	if err := os.MkdirAll(stageDir, 0o755); err != nil {
		return nil
	}

	if runtime.GOOS == "windows" {
		return extractFromZip(body, stagedBin)
	}
	return extractFromTarGz(body, stagedBin)
}

func downloadBytes(url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil) //nolint:gosec
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: status %d", url, resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func fetchChecksum(url, archive string) (string, error) {
	body, err := downloadBytes(url)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == archive {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("checksum not found for %s", archive)
}

func extractFromTarGz(data []byte, dest string) error {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	defer func() { _ = gr.Close() }()

	target := binaryName()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar: %w", err)
		}
		if filepath.Base(hdr.Name) != target {
			continue
		}
		return writeAtomic(dest, tr, 0o755)
	}
	return fmt.Errorf("binary %s not found in archive", target)
}

func extractFromZip(data []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("zip: %w", err)
	}
	target := binaryName()
	for _, f := range r.File {
		if filepath.Base(f.Name) != target {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() { _ = rc.Close() }()
		return writeAtomic(dest, rc, 0o755)
	}
	return fmt.Errorf("binary %s not found in zip", target)
}

// writeAtomic writes r to dest via a temp file to avoid leaving partial files
// on error.
func writeAtomic(dest string, r io.Reader, mode os.FileMode) error {
	tmp := dest + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
