//go:build windows

package encoder

import (
	"fmt"
	"os/exec"
	"strconv"
)

// KillProcessTree kills the process and all its children on Windows.
// Uses taskkill /F /T /PID as a fallback when Job Object is not available.
func KillProcessTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}

	pid := cmd.Process.Pid
	killCmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
	if err := killCmd.Run(); err != nil {
		// Fallback to direct kill
		return fmt.Errorf("taskkill failed for PID %d: %w", pid, err)
	}
	return nil
}
