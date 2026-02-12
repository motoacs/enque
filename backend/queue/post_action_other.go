//go:build !windows

package queue

import "fmt"

func setSystemSleep() error {
	return fmt.Errorf("sleep action is only supported on windows")
}
