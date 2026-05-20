// Package workspace provides primitives for detecting and enumerating
// members of a multi-repo roadmap workspace. A workspace root is a
// directory that is not itself a git repository but contains one or more
// immediate subdirectories that are git repositories with a roadmap
// config.
package workspace

import (
	"os"
	"path/filepath"
	"sort"
)

// MemberRoots returns the absolute paths of directories under root that
// contain a .git entry (file or directory), excluding root itself.
// Results are sorted lexicographically. Walk errors are silently ignored
// so that an unreadable subtree never aborts member discovery.
func MemberRoots(root string) []string {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		absRoot = root
	}
	var repos []string
	_ = filepath.WalkDir(absRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry == nil {
			return nil
		}
		if entry.Name() != ".git" {
			return nil
		}
		parent := filepath.Dir(path)
		if parent != absRoot {
			repos = append(repos, parent)
		}
		if entry.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	sort.Strings(repos)
	return repos
}

// IsWorkspaceRoot reports whether root looks like a roadmap workspace
// root: root has no .git of its own, and at least one member (returned
// by MemberRoots) contains docs/roadmap/.roadmapctl.toml. The member
// config presence is checked with os.Stat to avoid pulling the config
// loader into this package.
func IsWorkspaceRoot(root string) bool {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(absRoot, ".git")); err == nil {
		return false
	}
	for _, member := range MemberRoots(absRoot) {
		if _, err := os.Stat(filepath.Join(member, "docs", "roadmap", ".roadmapctl.toml")); err == nil {
			return true
		}
	}
	return false
}
