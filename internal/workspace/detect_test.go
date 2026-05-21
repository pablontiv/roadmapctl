package workspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pablontiv/roadmapctl/internal/workspace"
)

func mkdirAll(t *testing.T, paths ...string) {
	t.Helper()
	for _, p := range paths {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", p, err)
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestIsWorkspaceRoot_RootHasGit_ReturnsFalse(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, ".git"))
	mkdirAll(t, filepath.Join(root, "sub", ".git"))
	writeFile(t, filepath.Join(root, "sub", "docs", "roadmap", ".roadmapctl.toml"), "")

	if workspace.IsWorkspaceRoot(root) {
		t.Fatalf("expected false when root has its own .git")
	}
}

func TestIsWorkspaceRoot_NoGitNoMembers_ReturnsTrue(t *testing.T) {
	root := t.TempDir()
	if !workspace.IsWorkspaceRoot(root) {
		t.Fatalf("expected true: no .git at root and no docs/roadmap means workspace mode (downstream handles empty case)")
	}
}

func TestIsWorkspaceRoot_NoGitButHasRoadmapDir_ReturnsFalse(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "docs", "roadmap"))
	if workspace.IsWorkspaceRoot(root) {
		t.Fatalf("expected false: docs/roadmap present means single-repo even without .git")
	}
}

func TestIsWorkspaceRoot_NoGitWithMembers_ReturnsTrue(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha", ".git"))
	mkdirAll(t, filepath.Join(root, "beta", ".git"))
	writeFile(t, filepath.Join(root, "beta", "docs", "roadmap", ".roadmapctl.toml"), "")

	if !workspace.IsWorkspaceRoot(root) {
		t.Fatalf("expected true when root has no .git")
	}
}

func TestMemberRoots_ExcludesRoot(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, ".git"))
	mkdirAll(t, filepath.Join(root, "alpha", ".git"))

	got := workspace.MemberRoots(root)
	for _, m := range got {
		if m == root {
			t.Fatalf("MemberRoots must exclude root itself; got %v", got)
		}
	}
	if len(got) != 1 || filepath.Base(got[0]) != "alpha" {
		t.Fatalf("expected [alpha]; got %v", got)
	}
}

func TestMemberRoots_SortedAndLists(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t,
		filepath.Join(root, "zeta", ".git"),
		filepath.Join(root, "alpha", ".git"),
		filepath.Join(root, "mu", ".git"),
	)

	got := workspace.MemberRoots(root)
	if len(got) != 3 {
		t.Fatalf("expected 3 members; got %d (%v)", len(got), got)
	}
	wantBases := []string{"alpha", "mu", "zeta"}
	for i, m := range got {
		if filepath.Base(m) != wantBases[i] {
			t.Fatalf("order wrong at %d: got %s want %s (full=%v)", i, filepath.Base(m), wantBases[i], got)
		}
	}
}

func TestMemberRoots_EmptyOnRootOnly(t *testing.T) {
	root := t.TempDir()
	got := workspace.MemberRoots(root)
	if len(got) != 0 {
		t.Fatalf("expected empty; got %v", got)
	}
}

func TestMemberRoots_HandlesGitFile(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"))
	writeFile(t, filepath.Join(root, "alpha", ".git"), "gitdir: /elsewhere\n")

	got := workspace.MemberRoots(root)
	if len(got) != 1 || filepath.Base(got[0]) != "alpha" {
		t.Fatalf("expected to detect .git file (worktree/submodule); got %v", got)
	}
}

func TestMemberRoots_IgnoresNonDirectories(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha", ".git"))
	// Create a file (not dir) that would be skipped
	writeFile(t, filepath.Join(root, "regular_file"), "")

	got := workspace.MemberRoots(root)
	if len(got) != 1 || filepath.Base(got[0]) != "alpha" {
		t.Fatalf("expected to ignore non-directories; got %v", got)
	}
}

func TestMemberRoots_IgnoresDirsWithoutGit(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t,
		filepath.Join(root, "alpha", ".git"),
		filepath.Join(root, "beta"),  // no .git
		filepath.Join(root, "gamma", ".git"),
	)

	got := workspace.MemberRoots(root)
	if len(got) != 2 {
		t.Fatalf("expected 2 members (alpha and gamma); got %d (%v)", len(got), got)
	}
	bases := []string{filepath.Base(got[0]), filepath.Base(got[1])}
	if bases[0] != "alpha" || bases[1] != "gamma" {
		t.Fatalf("expected [alpha, gamma]; got %v", bases)
	}
}

func TestMemberRoots_UnreadableSubdir_SilentlyIgnored(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t,
		filepath.Join(root, "alpha", ".git"),
		filepath.Join(root, "unreadable"),
	)
	// Make unreadable dir inaccessible (requires non-Windows)
	if err := os.Chmod(filepath.Join(root, "unreadable"), 0o000); err != nil {
		t.Skipf("cannot chmod to 0o000: %v", err)
		return
	}
	defer func() {
		_ = os.Chmod(filepath.Join(root, "unreadable"), 0o700) //nolint:gosec
	}()

	// MemberRoots should ignore the unreadable dir and still find alpha
	got := workspace.MemberRoots(root)
	if len(got) != 1 || filepath.Base(got[0]) != "alpha" {
		t.Fatalf("expected [alpha] ignoring unreadable dir; got %v", got)
	}
}
