package encoder

import (
	"sync"
	"time"
)

type TimeoutGuard struct {
	mu              sync.RWMutex
	lastLineAt      time.Time
	lastProgressAt  time.Time
	outputTimeout   time.Duration
	progressTimeout time.Duration
	progressEnabled bool
}

func NewTimeoutGuard(outputTimeout, progressTimeout time.Duration) *TimeoutGuard {
	now := time.Now()
	return &TimeoutGuard{
		lastLineAt:      now,
		lastProgressAt:  now,
		outputTimeout:   outputTimeout,
		progressTimeout: progressTimeout,
	}
}

func (g *TimeoutGuard) MarkLine() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.lastLineAt = time.Now()
}

func (g *TimeoutGuard) MarkProgress() {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := time.Now()
	g.lastLineAt = now
	g.lastProgressAt = now
	g.progressEnabled = true
}

func (g *TimeoutGuard) IsOutputTimedOut() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.outputTimeout <= 0 {
		return false
	}
	return time.Since(g.lastLineAt) > g.outputTimeout
}

func (g *TimeoutGuard) IsProgressTimedOut() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if !g.progressEnabled || g.progressTimeout <= 0 {
		return false
	}
	return time.Since(g.lastProgressAt) > g.progressTimeout
}
