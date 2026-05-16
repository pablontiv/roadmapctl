//go:build !windows

package updater

import "syscall"

func unixExec(path string, args []string, env []string) error {
	return syscall.Exec(path, args, env) //nolint:gosec
}
