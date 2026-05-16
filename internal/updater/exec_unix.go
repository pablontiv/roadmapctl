//go:build !windows

package updater

import (
	"os"
	"syscall"
)

func unixExec(path string, args []string, env []string) error {
	return syscall.Exec(path, args, env) //nolint:gosec
}

// atomicReplace replaces dest with src atomically on Unix.
func atomicReplace(dest, src string) error {
	return os.Rename(src, dest)
}

// platformExec replaces the current process with the binary at path on Unix.
func platformExec(path string) error {
	return unixExec(path, os.Args, os.Environ())
}
