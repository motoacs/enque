//go:build windows

package encoder

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	kernel32Win                    = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObject            = kernel32Win.NewProc("CreateJobObjectW")
	procSetInformationJobObject    = kernel32Win.NewProc("SetInformationJobObject")
	procAssignProcessToJobObject   = kernel32Win.NewProc("AssignProcessToJobObject")
	procTerminateJobObject         = kernel32Win.NewProc("TerminateJobObject")
)

const (
	jobObjectExtendedLimitInformation = 9
	jOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE = 0x00002000
)

type jobObjectExtendedLimitInfo struct {
	BasicLimitInformation struct {
		PerProcessUserTimeLimit int64
		PerJobUserTimeLimit     int64
		LimitFlags              uint32
		MinimumWorkingSetSize   uintptr
		MaximumWorkingSetSize   uintptr
		ActiveProcessLimit      uint32
		Affinity                uintptr
		PriorityClass           uint32
		SchedulingClass         uint32
	}
	IoInfo struct {
		ReadOperationCount  uint64
		WriteOperationCount uint64
		OtherOperationCount uint64
		ReadTransferCount   uint64
		WriteTransferCount  uint64
		OtherTransferCount  uint64
	}
	ProcessMemoryLimit uintptr
	JobMemoryLimit     uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

// JobObject wraps a Windows Job Object handle.
type JobObject struct {
	handle syscall.Handle
}

// CreateJobObject creates a new Windows Job Object configured to kill
// all child processes when the handle is closed (app crash safety).
func CreateJobObject() (*JobObject, error) {
	handle, _, err := procCreateJobObject.Call(0, 0)
	if handle == 0 {
		return nil, fmt.Errorf("CreateJobObject: %v", err)
	}

	// Set KILL_ON_JOB_CLOSE so child processes die if our process dies
	var info jobObjectExtendedLimitInfo
	info.BasicLimitInformation.LimitFlags = jOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE

	ret, _, err := procSetInformationJobObject.Call(
		handle,
		uintptr(jobObjectExtendedLimitInformation),
		uintptr(unsafe.Pointer(&info)),
		unsafe.Sizeof(info),
	)
	if ret == 0 {
		syscall.CloseHandle(syscall.Handle(handle))
		return nil, fmt.Errorf("SetInformationJobObject: %v", err)
	}

	return &JobObject{handle: syscall.Handle(handle)}, nil
}

// AssignProcess adds a process to the job object.
func (jo *JobObject) AssignProcess(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return fmt.Errorf("process not started")
	}

	processHandle, err := syscall.OpenProcess(
		syscall.PROCESS_ALL_ACCESS, false, uint32(cmd.Process.Pid))
	if err != nil {
		return fmt.Errorf("OpenProcess: %w", err)
	}
	defer syscall.CloseHandle(processHandle)

	ret, _, callErr := procAssignProcessToJobObject.Call(
		uintptr(jo.handle),
		uintptr(processHandle),
	)
	if ret == 0 {
		return fmt.Errorf("AssignProcessToJobObject: %v", callErr)
	}

	return nil
}

// Terminate kills all processes in the job object.
func (jo *JobObject) Terminate(exitCode uint32) error {
	ret, _, err := procTerminateJobObject.Call(
		uintptr(jo.handle),
		uintptr(exitCode),
	)
	if ret == 0 {
		return fmt.Errorf("TerminateJobObject: %v", err)
	}
	return nil
}

// Close releases the job object handle.
func (jo *JobObject) Close() error {
	if jo.handle != 0 {
		return syscall.CloseHandle(jo.handle)
	}
	return nil
}

// SetupJobObject creates a Job Object and assigns the process to it.
// Returns the Job Object handle (caller must Close it).
func SetupJobObject(cmd *exec.Cmd) (*JobObject, error) {
	jo, err := CreateJobObject()
	if err != nil {
		return nil, err
	}

	if err := jo.AssignProcess(cmd); err != nil {
		jo.Close()
		return nil, err
	}

	return jo, nil
}
