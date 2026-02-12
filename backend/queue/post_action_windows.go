//go:build windows

package queue

import "syscall"

func setSystemSleep() error {
	powrprof := syscall.NewLazyDLL("powrprof.dll")
	setSuspendState := powrprof.NewProc("SetSuspendState")
	ret, _, callErr := setSuspendState.Call(uintptr(0), uintptr(1), uintptr(0))
	if ret == 0 {
		if callErr != syscall.Errno(0) {
			return callErr
		}
	}
	return nil
}
