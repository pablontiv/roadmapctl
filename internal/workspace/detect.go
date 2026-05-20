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

// MemberRoots returns the absolute paths of immediate subdirectories of
// root that contain a .git entry (file or directory). Root itself is
// excluded even if it has a .git. Results are sorted lexicographically.
// Read errors are silently ignored so that an unreadable subdirectory
// never aborts member discovery.
func MemberRoots(root string) []string {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		absRoot = root
	}
	entries, err := os.ReadDir(absRoot)
	if err != nil {
		return nil
	}
	var repos []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := filepath.Join(absRoot, entry.Name())
		if _, err := os.Stat(filepath.Join(candidate, ".git")); err == nil {
			repos = append(repos, candidate)
		}
	}
	sort.Strings(repos)
	return repos
}

// IsWorkspaceRoot reports whether root should be treated as a roadmap
// workspace root. A directory qualifies when it has neither a .git of
// its own nor a docs/roadmap/ directory at its top level — i.e. it is
// not itself a roadmap-bearing single repo. Downstream callers then
// enumerate MemberRoots to distinguish "populated workspace" from
// "empty workspace".
func IsWorkspaceRoot(root string) bool {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}
	if _, err := os.Stat(filepath.Join(absRoot, ".git")); err == nil {
		return false
	}
	if info, err := os.Stat(filepath.Join(absRoot, "docs", "roadmap")); err == nil && info.IsDir() {
		return false
	}
	return true
}
