package fsx

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestResolveInsideAcceptsRelativePath(t *testing.T) {
	repo := t.TempDir()
	resolved, rel, err := ResolveInside(repo, "docs/roadmap")
	if err != nil {
		t.Fatalf("ResolveInside() error = %v", err)
	}
	if resolved != filepath.Join(repo, "docs", "roadmap") {
		t.Fatalf("resolved = %q", resolved)
	}
	if rel != "docs/roadmap" {
		t.Fatalf("rel = %q", rel)
	}
}

func TestResolveInsideRejectsParentEscape(t *testing.T) {
	repo := t.TempDir()
	_, _, err := ResolveInside(repo, "../outside")
	if err == nil {
		t.Fatal("ResolveInside() error = nil, want escape error")
	}
}

func TestResolveInsideNormalizesWindowsSeparators(t *testing.T) {
	repo := t.TempDir()
	resolved, rel, err := ResolveInside(repo, `docs\\roadmap`)
	if err != nil {
		t.Fatalf("ResolveInside() error = %v", err)
	}
	if resolved != filepath.Join(repo, "docs", "roadmap") {
		t.Fatalf("resolved = %q", resolved)
	}
	if rel != "docs/roadmap" {
		t.Fatalf("rel = %q", rel)
	}
}

func TestResolveInsideRejectsWindowsAbsolutePath(t *testing.T) {
	repo := t.TempDir()
	_, _, err := ResolveInside(repo, `C:\\roadmap`)
	if err == nil {
		t.Fatal("ResolveInside() error = nil, want absolute path error")
	}
}

func TestResolveInsideRejectsSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink privileges vary on Windows")
	}
	repo := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(repo, "link")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}

	_, _, err := ResolveInside(repo, "link/file")
	if !errors.Is(err, ErrPathEscape) {
		t.Fatalf("ResolveInside error = %v, want ErrPathEscape", err)
	}
}

func TestResolveInsideRejectsEmptyPath(t *testing.T) {
	_, _, err := ResolveInside(t.TempDir(), "   ")
	if err == nil {
		t.Fatal("ResolveInside empty error = nil")
	}
}

func TestResolveInsideCurrentDir(t *testing.T) {
	repo := t.TempDir()
	resolved, rel, err := ResolveInside(repo, ".")
	if err != nil {
		t.Fatalf("ResolveInside(., ...) error = %v", err)
	}
	if resolved != repo {
		t.Fatalf("resolved = %q, want %q", resolved, repo)
	}
	if rel != "." {
		t.Fatalf("rel = %q, want .", rel)
	}
}

func TestResolveInsideAbsolutePathUnix(t *testing.T) {
	repo := t.TempDir()
	_, _, err := ResolveInside(repo, "/absolute/path")
	if !errors.Is(err, ErrPathEscape) {
		t.Fatalf("ResolveInside /absolute/path error = %v, want ErrPathEscape", err)
	}
}

func TestResolveInsideUNCPath(t *testing.T) {
	repo := t.TempDir()
	_, _, err := ResolveInside(repo, "//server/path")
	if !errors.Is(err, ErrAbsolutePath) {
		t.Fatalf("ResolveInside //server/path error = %v, want ErrAbsolutePath", err)
	}
}

func TestResolveInsideNestedDirectory(t *testing.T) {
	repo := t.TempDir()
	nested := filepath.Join(repo, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	resolved, rel, err := ResolveInside(repo, "a/b/c")
	if err != nil {
		t.Fatalf("ResolveInside error = %v", err)
	}
	if resolved != nested {
		t.Fatalf("resolved = %q, want %q", resolved, nested)
	}
	if rel != "a/b/c" {
		t.Fatalf("rel = %q, want a/b/c", rel)
	}
}

func TestContainedRelDoubleParent(t *testing.T) {
	root := t.TempDir()
	// Try to escape with ../.. pattern (should fail)
	target := filepath.Join(root, "..", "..")
	_, err := containedRel(root, target)
	if err == nil {
		t.Fatal("containedRel with ../ escape should error")
	}
	if !errors.Is(err, ErrPathEscape) {
		t.Fatalf("containedRel error = %v, want ErrPathEscape", err)
	}
}

func TestVerifySymlinkContainmentNonexistentPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink test on Windows")
	}
	repo := t.TempDir()
	// Path that doesn't exist - should still work since evalExistingPrefix handles it
	err := verifySymlinkContainment(repo, filepath.Join(repo, "nonexistent", "file.txt"))
	if err != nil {
		t.Fatalf("verifySymlinkContainment for nonexistent path error = %v", err)
	}
}

func TestEvalExistingPrefixNonexistentPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink test on Windows")
	}
	repo := t.TempDir()
	// Create a real directory
	realDir := filepath.Join(repo, "real")
	if err := os.MkdirAll(realDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Try to evaluate a path where only the parent exists
	target := filepath.Join(realDir, "nonexistent", "deep", "file.txt")
	result, err := evalExistingPrefix(target)
	if err != nil {
		t.Fatalf("evalExistingPrefix error = %v", err)
	}
	// Should return the clean path since the deepest existing prefix is realDir
	if !strings.Contains(result, "real") {
		t.Fatalf("result = %q, should contain real directory", result)
	}
}

func TestResolveInsideLeadingSlash(t *testing.T) {
	repo := t.TempDir()
	// Leading slash should be interpreted as absolute, which escapes the root
	_, _, err := ResolveInside(repo, "/root/path")
	if !errors.Is(err, ErrPathEscape) {
		t.Fatalf("ResolveInside with leading slash error = %v, want ErrPathEscape", err)
	}
}
