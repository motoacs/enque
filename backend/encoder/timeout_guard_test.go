package encoder

import (
	"testing"
	"time"
)

func TestTimeoutGuard_NoOutputTimeout(t *testing.T) {
	triggered := make(chan string, 1)
	tg := NewTimeoutGuard(2, 0, func(reason string) {
		triggered <- reason
	})
	tg.Start()
	defer tg.Stop()

	// Wait for timeout
	select {
	case reason := <-triggered:
		if reason != "no_output" {
			t.Fatalf("expected no_output, got %s", reason)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout not triggered within 5s")
	}

	timedOut, reason := tg.TimedOut()
	if !timedOut {
		t.Fatal("expected TimedOut to be true")
	}
	if reason != "no_output" {
		t.Fatalf("expected no_output, got %s", reason)
	}
}

func TestTimeoutGuard_OutputResetsPrevents(t *testing.T) {
	triggered := make(chan string, 1)
	tg := NewTimeoutGuard(2, 0, func(reason string) {
		triggered <- reason
	})
	tg.Start()
	defer tg.Stop()

	// Keep sending output to prevent timeout
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		tg.NotifyOutput()
	}

	select {
	case <-triggered:
		t.Fatal("timeout should not have been triggered")
	case <-time.After(500 * time.Millisecond):
		// Good, no timeout
	}
}

func TestTimeoutGuard_NoProgressTimeout(t *testing.T) {
	triggered := make(chan string, 1)
	tg := NewTimeoutGuard(0, 2, func(reason string) {
		triggered <- reason
	})
	tg.Start()
	defer tg.Stop()

	// Send initial progress, then stop
	tg.NotifyProgress(10.0)
	time.Sleep(500 * time.Millisecond)
	tg.NotifyProgress(10.0) // Same progress, won't reset timer

	select {
	case reason := <-triggered:
		if reason != "no_progress" {
			t.Fatalf("expected no_progress, got %s", reason)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout not triggered within 5s")
	}
}

func TestTimeoutGuard_ProgressAdvancePrevents(t *testing.T) {
	triggered := make(chan string, 1)
	tg := NewTimeoutGuard(0, 2, func(reason string) {
		triggered <- reason
	})
	tg.Start()
	defer tg.Stop()

	// Keep advancing progress
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		tg.NotifyProgress(float64(i * 10))
	}

	select {
	case <-triggered:
		t.Fatal("timeout should not have been triggered")
	case <-time.After(500 * time.Millisecond):
		// Good
	}
}

func TestTimeoutGuard_BothDisabled(t *testing.T) {
	triggered := make(chan string, 1)
	tg := NewTimeoutGuard(0, 0, func(reason string) {
		triggered <- reason
	})
	tg.Start()
	defer tg.Stop()

	time.Sleep(2 * time.Second)

	select {
	case <-triggered:
		t.Fatal("no timeout should trigger when both disabled")
	default:
		// Good
	}
}

func TestTimeoutGuard_StopPreventsCallback(t *testing.T) {
	triggered := make(chan string, 1)
	tg := NewTimeoutGuard(1, 0, func(reason string) {
		triggered <- reason
	})
	tg.Start()
	tg.Stop()

	time.Sleep(2 * time.Second)

	select {
	case <-triggered:
		t.Fatal("timeout should not trigger after Stop")
	default:
		// Good
	}
}
