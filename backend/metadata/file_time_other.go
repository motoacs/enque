//go:build !windows

package metadata

// RestoreFileTime is a no-op on non-Windows platforms.
func RestoreFileTime(srcPath, dstPath string) error {
	return nil
}

// IsSupported returns false on non-Windows platforms.
func IsSupported() bool {
	return false
}

// RestoreFileTimeIfNeeded is a no-op on non-Windows platforms.
func RestoreFileTimeIfNeeded(srcPath, dstPath string, enabled bool) error {
	return nil
}
