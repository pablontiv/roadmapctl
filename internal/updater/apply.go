package updater

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
)

// CurrentVersion must be set by main before calling ApplyStagedIfAvailable.
// It should match the version string injected via ldflags (e.g. "v0.2.0" or "dev").
var CurrentVersion = "dev"

// execFn replaces the running process with the binary at path. Overriding this
// in tests prevents actual re-exec during unit testing.
var execFn = platformExec

// ApplyStagedIfAvailable detects the newest staged binary, and if it is newer
// than CurrentVersion, atomically replaces the running binary and re-execs.
// All errors (permissions, filesystem) are swallowed silently — the function
// never interrupts the current command.
func ApplyStagedIfAvailable() error {
	if CurrentVersion == "dev" {
		return nil
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil
	}
	stagedBase := filepath.Join(cacheDir, "roadmapctl", "staged")

	newestTag, newestBin, err := findNewest(stagedBase)
	if err != nil || newestTag == "" {
		return nil
	}

	if !isNewer(newestTag, CurrentVersion) {
		return nil
	}

	currentBin, err := os.Executable()
	if err != nil {
		return nil
	}

	if err := atomicReplace(currentBin, newestBin); err != nil {
		return nil
	}

	// Re-exec — on Unix this never returns on success; on Windows we exit after
	// launching the new process.
	_ = execFn(currentBin)
	return nil
}

// findNewest scans stagedBase for version directories and returns the newest
// tag and its binary path.
func findNewest(stagedBase string) (tag string, binPath string, err error) {
	entries, err := os.ReadDir(stagedBase)
	if err != nil {
		return "", "", err
	}

	type candidate struct {
		tag string
		bin string
	}
	var candidates []candidate
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		bin := filepath.Join(stagedBase, e.Name(), binaryName())
		if fileExists(bin) {
			candidates = append(candidates, candidate{e.Name(), bin})
		}
	}
	if len(candidates) == 0 {
		return "", "", nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return isNewer(candidates[i].tag, candidates[j].tag)
	})
	best := candidates[0]
	return best.tag, best.bin, nil
}

// atomicReplace replaces dest with src using a platform-appropriate strategy.
func atomicReplace(dest, src string) error {
	if runtime.GOOS == "windows" {
		// On Windows, a binary in use cannot be renamed over; copy to dest.tmp
		// and rename, which works when the old file is not open for write.
		tmp := dest + ".old"
		if err := os.Rename(dest, tmp); err != nil {
			return err
		}
		if err := copyFile(src, dest); err != nil {
			// Best-effort restore.
			_ = os.Rename(tmp, dest)
			return err
		}
		_ = os.Remove(tmp)
		return nil
	}
	return os.Rename(src, dest)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src) //nolint:gosec
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	if _, err := copyIO(in, out); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func copyIO(src io.Reader, dst io.Writer) (int64, error) {
	return io.Copy(dst, src)
}

// platformExec re-executes the binary at path with the current process args and
// environment. On Unix it replaces the process in-place (never returns on
// success). On Windows it spawns a new process and exits the current one.
func platformExec(path string) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command(path, os.Args[1:]...) //nolint:gosec
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		os.Exit(0)
	}
	return unixExec(path, os.Args, os.Environ())
}
