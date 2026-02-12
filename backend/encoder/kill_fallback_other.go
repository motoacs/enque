//go:build !windows

package encoder

import (
	"os/exec"
	"syscall"
)

// KillProcessTree kills the process on non-Windows platforms.
// Uses SIGKILL as there's no process tree concept on macOS/Linux
// without platform-specific process group handling.
func KillProcessTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	return cmd.Process.Signal(syscall.SIGKILL)
}
