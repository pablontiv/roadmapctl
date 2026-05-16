//go:build windows

package updater

import "fmt"

func unixExec(_ string, _ []string, _ []string) error {
	return fmt.Errorf("syscall.Exec not supported on Windows")
}
