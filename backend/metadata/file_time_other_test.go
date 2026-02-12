//go:build !windows

package metadata

import "testing"

func TestRestoreFileTime_NoOpOnNonWindows(t *testing.T) {
	if err := RestoreFileTime("/tmp/a", "/tmp/b"); err != nil {
		t.Fatalf("expected no-op nil, got %v", err)
	}
}
