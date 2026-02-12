//go:build !windows

package queue

import "fmt"

// platformShutdown is a no-op on non-Windows platforms.
func platformShutdown() error {
	return fmt.Errorf("shutdown is only supported on Windows")
}

// platformSleep is a no-op on non-Windows platforms.
func platformSleep() error {
	return fmt.Errorf("sleep is only supported on Windows")
}
