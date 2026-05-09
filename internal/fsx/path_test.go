package fsx

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
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
