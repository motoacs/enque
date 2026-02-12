package encoder

import (
	"testing"
	"time"
)

func TestTimeoutGuard(t *testing.T) {
	g := NewTimeoutGuard(5*time.Millisecond, 5*time.Millisecond)
	time.Sleep(8 * time.Millisecond)
	if !g.IsOutputTimedOut() {
		t.Fatalf("expected output timeout")
	}
	g.MarkLine()
	if g.IsOutputTimedOut() {
		t.Fatalf("did not expect output timeout immediately after mark")
	}
	if g.IsProgressTimedOut() {
		t.Fatalf("progress timeout should be disabled before MarkProgress")
	}
	g.MarkProgress()
	time.Sleep(8 * time.Millisecond)
	if !g.IsProgressTimedOut() {
		t.Fatalf("expected progress timeout")
	}
}
