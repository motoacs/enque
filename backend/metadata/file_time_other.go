//go:build !windows

package metadata

func RestoreFileTime(srcPath, dstPath string) error {
	return nil
}
