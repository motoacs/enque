package encoder

import (
	"sync"
	"time"
)

// TimeoutGuard monitors encoder activity with two-stage timeout detection.
// Stage 1: no_output_timeout — no stderr output at all.
// Stage 2: no_progress_timeout — output exists but progress percentage is stuck.
type TimeoutGuard struct {
	mu              sync.Mutex
	noOutputSec     int
	noProgressSec   int
	lastOutputTime  time.Time
	lastProgress    *float64
	lastProgressAt  time.Time
	timedOut        bool
	timeoutReason   string
	stopCh          chan struct{}
	onTimeout       func(reason string)
}

// NewTimeoutGuard creates a timeout guard with the given thresholds.
// If either threshold is 0, that stage is disabled.
func NewTimeoutGuard(noOutputSec, noProgressSec int, onTimeout func(reason string)) *TimeoutGuard {
	now := time.Now()
	return &TimeoutGuard{
		noOutputSec:    noOutputSec,
		noProgressSec:  noProgressSec,
		lastOutputTime: now,
		lastProgressAt: now,
		stopCh:         make(chan struct{}),
		onTimeout:      onTimeout,
	}
}

// Start begins the monitoring ticker (1-second interval).
func (tg *TimeoutGuard) Start() {
	go tg.run()
}

// Stop halts the monitoring ticker.
func (tg *TimeoutGuard) Stop() {
	select {
	case <-tg.stopCh:
	default:
		close(tg.stopCh)
	}
}

// NotifyOutput records that output was received from stderr.
func (tg *TimeoutGuard) NotifyOutput() {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	tg.lastOutputTime = time.Now()
}

// NotifyProgress records a new progress percentage.
func (tg *TimeoutGuard) NotifyProgress(percent float64) {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	tg.lastOutputTime = time.Now()
	if tg.lastProgress == nil || *tg.lastProgress != percent {
		tg.lastProgress = &percent
		tg.lastProgressAt = time.Now()
	}
}

// TimedOut returns whether a timeout has been detected and the reason.
func (tg *TimeoutGuard) TimedOut() (bool, string) {
	tg.mu.Lock()
	defer tg.mu.Unlock()
	return tg.timedOut, tg.timeoutReason
}

func (tg *TimeoutGuard) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-tg.stopCh:
			return
		case <-ticker.C:
			tg.check()
		}
	}
}

func (tg *TimeoutGuard) check() {
	tg.mu.Lock()
	defer tg.mu.Unlock()

	if tg.timedOut {
		return
	}

	now := time.Now()

	// Stage 1: no output at all
	if tg.noOutputSec > 0 {
		elapsed := now.Sub(tg.lastOutputTime)
		if elapsed >= time.Duration(tg.noOutputSec)*time.Second {
			tg.timedOut = true
			tg.timeoutReason = "no_output"
			if tg.onTimeout != nil {
				go tg.onTimeout(tg.timeoutReason)
			}
			return
		}
	}

	// Stage 2: output exists but progress stuck
	if tg.noProgressSec > 0 && tg.lastProgress != nil {
		elapsed := now.Sub(tg.lastProgressAt)
		if elapsed >= time.Duration(tg.noProgressSec)*time.Second {
			tg.timedOut = true
			tg.timeoutReason = "no_progress"
			if tg.onTimeout != nil {
				go tg.onTimeout(tg.timeoutReason)
			}
			return
		}
	}
}
