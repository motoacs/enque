//go:build !windows

package encoder

import "os/exec"

// JobObject is a no-op stub on non-Windows platforms.
type JobObject struct{}

// CreateJobObject returns nil on non-Windows.
func CreateJobObject() (*JobObject, error) {
	return nil, nil
}

// SetupJobObject is a no-op on non-Windows.
func SetupJobObject(cmd *exec.Cmd) (*JobObject, error) {
	return nil, nil
}

// AssignProcess is a no-op on non-Windows.
func (jo *JobObject) AssignProcess(cmd *exec.Cmd) error {
	return nil
}

// Terminate is a no-op on non-Windows.
func (jo *JobObject) Terminate(exitCode uint32) error {
	return nil
}

// Close is a no-op on non-Windows.
func (jo *JobObject) Close() error {
	return nil
}
