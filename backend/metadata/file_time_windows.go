//go:build windows

package metadata

import (
	"os"
	"syscall"
	"time"
)

func RestoreFileTime(srcPath, dstPath string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	t := srcInfo.ModTime()
	atime := t
	mtime := t
	if err := os.Chtimes(dstPath, atime, mtime); err != nil {
		return err
	}
	// CreationTime restore requires syscall; best-effort with SetFileTime.
	p, err := syscall.UTF16PtrFromString(dstPath)
	if err != nil {
		return nil
	}
	h, err := syscall.CreateFile(p, syscall.FILE_WRITE_ATTRIBUTES, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil
	}
	defer syscall.CloseHandle(h)
	ft := syscall.NsecToFiletime(timeToNsec(t))
	_ = syscall.SetFileTime(h, &ft, nil, &ft)
	return nil
}

func timeToNsec(t time.Time) int64 {
	return t.UnixNano()
}
