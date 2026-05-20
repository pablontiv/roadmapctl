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

func TestIsWorkspaceRoot_NoMembers_ReturnsFalse(t *testing.T) {
	root := t.TempDir()
	if workspace.IsWorkspaceRoot(root) {
		t.Fatalf("expected false on empty root")
	}
}

func TestIsWorkspaceRoot_MembersWithoutConfig_ReturnsFalse(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha", ".git"))
	mkdirAll(t, filepath.Join(root, "beta", ".git"))

	if workspace.IsWorkspaceRoot(root) {
		t.Fatalf("expected false when no member has a roadmap config")
	}
}

func TestIsWorkspaceRoot_AtLeastOneMemberConfigured_ReturnsTrue(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha", ".git"))
	mkdirAll(t, filepath.Join(root, "beta", ".git"))
	writeFile(t, filepath.Join(root, "beta", "docs", "roadmap", ".roadmapctl.toml"), "")

	if !workspace.IsWorkspaceRoot(root) {
		t.Fatalf("expected true when at least one member has docs/roadmap/.roadmapctl.toml")
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
