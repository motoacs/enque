package queue

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/yuta/enque/backend/config"
	"github.com/yuta/enque/backend/encoder"
	"github.com/yuta/enque/backend/events"
	"github.com/yuta/enque/backend/logging"
)

// Manager orchestrates encoding sessions with a worker pool.
type Manager struct {
	mu                 sync.RWMutex
	session            *Session
	workers            []*Worker
	registry           *encoder.Registry
	emitter            *events.Emitter
	tempTracker        *TempTracker
	logger             *logging.AppLogger
	cancelFunc         context.CancelFunc
	wg                 sync.WaitGroup
	overwriteResponses map[string]chan string
}

// NewManager creates a new queue manager.
func NewManager(registry *encoder.Registry, emitter *events.Emitter, logger *logging.AppLogger) *Manager {
	return &Manager{
		registry:           registry,
		emitter:            emitter,
		logger:             logger,
		tempTracker:        NewTempTracker(config.TempIndexPath()),
		overwriteResponses: make(map[string]chan string),
	}
}

// StartEncode begins a new encoding session.
func (m *Manager) StartEncode(req EncodeRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if a session is already running
	if m.session != nil && (m.session.State == StateRunning || m.session.State == StateStopping || m.session.State == StateAborting) {
		return fmt.Errorf("%s: session already running", encoder.ErrSessionRunning)
	}

	// Resolve adapter
	adapter, err := m.registry.Resolve(req.Profile.EncoderType)
	if err != nil {
		return err
	}

	// Determine encoder path
	encoderPath := m.resolveEncoderPath(req.Profile.EncoderType, req.AppConfigSnapshot)
	if encoderPath == "" {
		return fmt.Errorf("%s: encoder path not configured for %s", encoder.ErrToolNotFound, req.Profile.EncoderType)
	}

	// Create session
	sessionID := generateSessionID()
	session := NewSession(sessionID, req.Jobs, req.Profile.EncoderType, req.AppConfigSnapshot)
	m.session = session

	// Create process runner
	runner := encoder.NewProcessRunner(
		encoderPath,
		adapter,
		req.AppConfigSnapshot.NoOutputTimeoutSec,
		req.AppConfigSnapshot.NoProgressTimeoutSec,
	)

	// Create output resolver
	resolver := NewOutputResolver()

	// Create workers and start
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFunc = cancel

	maxJobs := req.AppConfigSnapshot.MaxConcurrentJobs
	if maxJobs < 1 {
		maxJobs = 1
	}

	// Job channel
	jobCh := make(chan *QueueJob, len(session.Jobs))
	for _, job := range session.Jobs {
		jobCh <- job
	}
	close(jobCh)

	// Emit session started
	if m.logger != nil {
		m.logger.Info("session started: %s (encoder=%s, jobs=%d, workers=%d)", sessionID, req.Profile.EncoderType, len(req.Jobs), maxJobs)
	}
	m.emitter.SessionStarted(session.Snapshot())

	// Launch workers
	m.workers = make([]*Worker, maxJobs)
	for i := 0; i < maxJobs; i++ {
		w := NewWorker(WorkerConfig{
			ID:          i,
			Session:     session,
			Adapter:     adapter,
			Runner:      runner,
			Resolver:    resolver,
			TempTracker: m.tempTracker,
			Emitter:     m.emitter,
			Manager:     m,
			Profile:     req.Profile,
			AppConfig:   req.AppConfigSnapshot,
			EncoderPath: encoderPath,
		})
		m.workers[i] = w
		m.wg.Add(1)
		go func(worker *Worker) {
			defer m.wg.Done()
			worker.Run(ctx, jobCh)
		}(w)
	}

	// Monitor completion in background
	go m.waitForCompletion()

	return nil
}

// RequestGracefulStop stops the current session gracefully.
func (m *Manager) RequestGracefulStop(sessionID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || m.session.ID != sessionID {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	m.session.RequestStop()
	m.emitter.SessionState(m.session.Snapshot())
	return nil
}

// RequestAbort aborts the current session.
func (m *Manager) RequestAbort(sessionID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || m.session.ID != sessionID {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	m.session.RequestAbort()

	// Cancel context to kill running processes
	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	m.emitter.SessionState(m.session.Snapshot())
	return nil
}

// SkipJob marks a pending job to be skipped when the worker picks it up.
func (m *Manager) SkipJob(sessionID, jobID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || m.session.ID != sessionID {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	m.session.RequestSkipJob(jobID)
	return nil
}

// CancelJob cancels a single running job.
func (m *Manager) CancelJob(sessionID, jobID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || m.session.ID != sessionID {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Find the worker running this job and cancel it
	for _, w := range m.workers {
		// Workers track their current job through the cancel context
		// For now, we mark the job for cancellation via session
		w.CancelCurrentJob()
	}

	return nil
}

// GetSession returns the current session (thread-safe snapshot).
func (m *Manager) GetSession() *Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.session
}

// GetSessionID returns the current session ID or empty string.
func (m *Manager) GetSessionID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.session == nil {
		return ""
	}
	return m.session.ID
}

// ListTempArtifacts returns leftover temp files from previous sessions.
func (m *Manager) ListTempArtifacts() []string {
	return m.tempTracker.ListPaths()
}

// CleanupTempArtifacts deletes specified temp files.
func (m *Manager) CleanupTempArtifacts(paths []string) error {
	for _, p := range paths {
		if err := removeFileIfExists(p); err != nil {
			m.emitter.Warning(map[string]interface{}{
				"message": fmt.Sprintf("failed to cleanup temp file: %s: %v", p, err),
			})
		}
		m.tempTracker.Remove(p)
	}
	return nil
}

// ResolveOverwrite responds to an overwrite confirmation for a job.
func (m *Manager) ResolveOverwrite(sessionID, jobID, decision string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.session == nil || m.session.ID != sessionID {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	ch, ok := m.overwriteResponses[jobID]
	if !ok {
		return fmt.Errorf("no pending overwrite for job: %s", jobID)
	}

	select {
	case ch <- decision:
	default:
		return fmt.Errorf("overwrite already resolved for job: %s", jobID)
	}
	return nil
}

// WaitForOverwrite waits for an overwrite decision with a 10-minute timeout.
// Returns the decision ("overwrite", "skip", or "abort"), or "skip" on timeout.
func (m *Manager) WaitForOverwrite(jobID string) (string, error) {
	ch := make(chan string, 1)

	m.mu.Lock()
	m.overwriteResponses[jobID] = ch
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.overwriteResponses, jobID)
		m.mu.Unlock()
	}()

	select {
	case decision := <-ch:
		return decision, nil
	case <-time.After(10 * time.Minute):
		return "skip", nil
	}
}

func (m *Manager) waitForCompletion() {
	m.wg.Wait()

	m.mu.Lock()
	if m.session != nil {
		m.session.Finish()
		snapshot := m.session.Snapshot()
		m.mu.Unlock()

		if m.logger != nil {
			m.logger.Info("session finished: %s (state=%s, completed=%d, failed=%d)", m.session.ID, string(m.session.State), m.session.CompletedJobs, m.session.FailedJobs)
		}
		m.emitter.SessionFinished(snapshot)

		// Post-complete action
		m.handlePostAction()
	} else {
		m.mu.Unlock()
	}
}

func (m *Manager) handlePostAction() {
	m.mu.RLock()
	session := m.session
	m.mu.RUnlock()

	if session == nil {
		return
	}

	// Do not execute post-action if session was aborted
	if session.State == StateAborted {
		return
	}

	action := session.AppCfg.PostCompleteAction
	command := session.AppCfg.PostCompleteCommand
	if action == "" || action == "none" {
		return
	}

	if err := ExecutePostAction(action, command, m.logger); err != nil {
		if m.logger != nil {
			m.logger.Error("post-complete action failed: %v", err)
		}
	}
}

func (m *Manager) resolveEncoderPath(encoderType string, cfg AppConfigSnapshot) string {
	switch encoderType {
	case "nvencc":
		return cfg.NVEncCPath
	default:
		return ""
	}
}

func generateSessionID() string {
	return fmt.Sprintf("s_%d_%s", time.Now().UnixMilli(), generateShortID())
}

func removeFileIfExists(path string) error {
	if !fileExists(path) {
		return nil
	}
	return os.Remove(path)
}
