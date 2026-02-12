//go:build windows

package encoder

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

type jobHandle = windows.Handle

func attachJobObject(p *os.Process) (jobHandle, bool, error) {
	job, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, false, err
	}
	limitInfo := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	limitInfo.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, err := windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&limitInfo)), uint32(unsafe.Sizeof(limitInfo))); err != nil {
		_ = windows.CloseHandle(job)
		return 0, false, err
	}
	proc, err := windows.OpenProcess(windows.PROCESS_TERMINATE|windows.PROCESS_SET_QUOTA|windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(p.Pid))
	if err != nil {
		_ = windows.CloseHandle(job)
		return 0, false, err
	}
	defer windows.CloseHandle(proc)
	if err := windows.AssignProcessToJobObject(job, proc); err != nil {
		_ = windows.CloseHandle(job)
		return 0, false, err
	}
	return job, true, nil
}

func terminateWithJobObject(job jobHandle, pid int) error {
	if job != 0 {
		if err := windows.TerminateJobObject(job, 1); err == nil {
			return nil
		}
	}
	return terminateProcessTree(pid)
}

func closeJobObject(job jobHandle) {
	if job != 0 {
		_ = windows.CloseHandle(job)
	}
}
