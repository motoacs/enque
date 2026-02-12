//go:build windows

package queue

import "os/exec"

// platformShutdown initiates Windows shutdown.
func platformShutdown() error {
	cmd := exec.Command("shutdown", "/s", "/t", "60")
	return cmd.Run()
}

// platformSleep puts Windows to sleep.
func platformSleep() error {
	// rundll32 powrprof.dll,SetSuspendState 0,1,0
	cmd := exec.Command("rundll32", "powrprof.dll,SetSuspendState", "0,1,0")
	return cmd.Run()
}
