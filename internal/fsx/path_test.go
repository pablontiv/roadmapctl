package fsx

import (
	"path/filepath"
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
