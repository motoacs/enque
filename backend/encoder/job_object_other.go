//go:build !windows

package encoder

import "os"

type jobHandle struct{}

func attachJobObject(p *os.Process) (jobHandle, bool, error) {
	return jobHandle{}, false, nil
}

func terminateWithJobObject(job jobHandle, pid int) error {
	return terminateProcessTree(pid)
}

func closeJobObject(job jobHandle) {}
