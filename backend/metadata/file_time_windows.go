//go:build windows

package metadata

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGetFileTime  = kernel32.NewProc("GetFileTime")
	procSetFileTime  = kernel32.NewProc("SetFileTime")
)

// FileTimeInfo holds Windows file timestamps.
type FileTimeInfo struct {
	CreationTime   syscall.Filetime
	LastAccessTime syscall.Filetime
	LastWriteTime  syscall.Filetime
}

// GetFileTime reads CreationTime and LastWriteTime from a file using Win32 API.
func GetFileTime(path string) (*FileTimeInfo, error) {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, fmt.Errorf("utf16 path: %w", err)
	}

	handle, err := syscall.CreateFile(
		pathPtr,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer syscall.CloseHandle(handle)

	var info FileTimeInfo
	ret, _, callErr := procGetFileTime.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&info.CreationTime)),
		uintptr(unsafe.Pointer(&info.LastAccessTime)),
		uintptr(unsafe.Pointer(&info.LastWriteTime)),
	)
	if ret == 0 {
		return nil, fmt.Errorf("GetFileTime: %v", callErr)
	}

	return &info, nil
}

// SetFileTime sets CreationTime and LastWriteTime on a file using Win32 API.
func SetFileTime(path string, info *FileTimeInfo) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return fmt.Errorf("utf16 path: %w", err)
	}

	handle, err := syscall.CreateFile(
		pathPtr,
		syscall.FILE_WRITE_ATTRIBUTES,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return fmt.Errorf("open file for write: %w", err)
	}
	defer syscall.CloseHandle(handle)

	ret, _, callErr := procSetFileTime.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&info.CreationTime)),
		uintptr(unsafe.Pointer(&info.LastAccessTime)),
		uintptr(unsafe.Pointer(&info.LastWriteTime)),
	)
	if ret == 0 {
		return fmt.Errorf("SetFileTime: %v", callErr)
	}

	return nil
}

// RestoreFileTime copies CreationTime and LastWriteTime from src to dst.
func RestoreFileTime(srcPath, dstPath string) error {
	info, err := GetFileTime(srcPath)
	if err != nil {
		return fmt.Errorf("get source time: %w", err)
	}
	return SetFileTime(dstPath, info)
}

// IsSupported returns true on Windows.
func IsSupported() bool {
	return true
}

// RestoreFileTimeIfNeeded restores file time if the feature is supported and enabled.
func RestoreFileTimeIfNeeded(srcPath, dstPath string, enabled bool) error {
	if !enabled {
		return nil
	}
	return RestoreFileTime(srcPath, dstPath)
}

// Ensure we reference os to avoid unused import
var _ = os.Stat
