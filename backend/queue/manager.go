package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/motoacs/enque/backend/encoder"
	"github.com/motoacs/enque/backend/logging"
	"github.com/motoacs/enque/backend/metadata"
	"github.com/motoacs/enque/backend/model"
)

const (
	eventSessionStarted    = "enque:session_started"
	eventJobStarted        = "enque:job_started"
	eventJobProgress       = "enque:job_progress"
	eventJobLog            = "enque:job_log"
	eventJobNeedsOverwrite = "enque:job_needs_overwrite"
	eventJobFinished       = "enque:job_finished"
	eventSessionState      = "enque:session_state"
	eventSessionFinished   = "enque:session_finished"
	eventWarning           = "enque:warning"
	eventError             = "enque:error"
)

type Manager struct {
	mu         sync.Mutex
	registry   *encoder.Registry
	runner     encoder.ProcessRunner
	resolver   *OutputResolver
	logger     *logging.JobLogger
	emitter    EventEmitter
	appLogger  *slog.Logger
	runtimeDir string

	active *sessionRuntime
}

func NewManager(registry *encoder.Registry, runner encoder.ProcessRunner, resolver *OutputResolver, logger *logging.JobLogger, emitter EventEmitter, appLogger *slog.Logger, runtimeDir string) *Manager {
	if resolver == nil {
		resolver = NewOutputResolver()
	}
	return &Manager{
		registry:   registry,
		runner:     runner,
		resolver:   resolver,
		logger:     logger,
		emitter:    emitter,
		appLogger:  appLogger,
		runtimeDir: runtimeDir,
	}
}

func (m *Manager) StartEncode(ctx context.Context, req model.StartEncodeRequest, encoderPath string) (model.EncodeSession, error) {
	if len(req.Jobs) == 0 {
		return model.EncodeSession{}, &model.EnqueError{Code: model.ErrValidation, Message: "jobs must not be empty"}
	}
	if errs := model.ValidateAppConfig(req.AppConfigSnapshot); len(errs) > 0 {
		return model.EncodeSession{}, &model.EnqueError{Code: model.ErrValidation, Message: "invalid app config", Fields: errs}
	}
	if errs := model.ValidateProfile(req.Profile); len(errs) > 0 {
		return model.EncodeSession{}, &model.EnqueError{Code: model.ErrValidation, Message: "invalid profile", Fields: errs}
	}
	adapter, err := m.registry.Resolve(req.Profile.EncoderType)
	if err != nil {
		return model.EncodeSession{}, err
	}
	if err := adapter.ValidateProfile(req.Profile); err != nil {
		return model.EncodeSession{}, err
	}

	m.mu.Lock()
	if m.active != nil && m.active.session.State == "running" {
		m.mu.Unlock()
		return model.EncodeSession{}, &model.EnqueError{Code: model.ErrSessionRunning, Message: "session already running"}
	}
	s := newSessionRuntime(len(req.Jobs))
	for _, j := range req.Jobs {
		s.jobs[j.JobID] = &model.QueueJob{JobID: j.JobID, InputPath: j.InputPath, Status: model.JobStatusPending}
	}
	m.active = s
	m.mu.Unlock()

	m.emit(eventSessionStarted, map[string]any{
		"session_id":   s.session.SessionID,
		"total_jobs":   s.session.TotalJobs,
		"started_at":   s.session.StartedAt,
		"encoder_type": req.Profile.EncoderType,
	})

	go m.runSession(ctx, s, req, adapter, encoderPath)
	return s.session, nil
}

func (m *Manager) RequestGracefulStop(sessionID string) error {
	s, err := m.requireSession(sessionID)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.session.StopRequested = true
	s.mu.Unlock()
	m.emitSessionState(s)
	return nil
}

func (m *Manager) RequestAbort(sessionID string) error {
	s, err := m.requireSession(sessionID)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.session.AbortRequested = true
	for _, cancel := range s.jobCancels {
		cancel()
	}
	s.mu.Unlock()
	s.cancel()
	m.emitSessionState(s)
	return nil
}

func (m *Manager) CancelJob(sessionID, jobID string) error {
	s, err := m.requireSession(sessionID)
	if err != nil {
		return err
	}
	s.mu.Lock()
	cancel, ok := s.jobCancels[jobID]
	s.mu.Unlock()
	if ok {
		cancel()
	}
	return nil
}

func (m *Manager) ResolveOverwrite(sessionID, jobID string, decision model.ResolveOverwriteDecision) error {
	s, err := m.requireSession(sessionID)
	if err != nil {
		return err
	}
	s.mu.Lock()
	waiter, ok := s.overwriteWaiters[jobID]
	if ok {
		delete(s.overwriteWaiters, jobID)
	}
	s.mu.Unlock()
	if !ok {
		return nil
	}
	waiter.ch <- decision
	close(waiter.ch)
	return nil
}

func (m *Manager) ListTempArtifacts() ([]string, error) {
	index, err := m.loadTempIndex()
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(index.Artifacts))
	for _, art := range index.Artifacts {
		if _, err := os.Stat(art.Path); err == nil {
			paths = append(paths, art.Path)
		}
	}
	return paths, nil
}

func (m *Manager) CleanupTempArtifacts(paths []string) error {
	for _, p := range paths {
		_ = os.Remove(p)
	}
	index, err := m.loadTempIndex()
	if err != nil {
		return err
	}
	filtered := make([]model.TempArtifact, 0, len(index.Artifacts))
	for _, art := range index.Artifacts {
		keep := true
		for _, p := range paths {
			if art.Path == p {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, art)
		}
	}
	index.Artifacts = filtered
	return m.saveTempIndex(index)
}

func (m *Manager) runSession(ctx context.Context, s *sessionRuntime, req model.StartEncodeRequest, adapter encoder.Adapter, encoderPath string) {
	defer func() {
		s.mu.Lock()
		s.session.FinishedAt = time.Now().UTC()
		if s.session.AbortRequested {
			s.session.State = "aborted"
		} else {
			s.session.State = "completed"
		}
		s.mu.Unlock()
		m.emit(eventSessionFinished, s.session)

		m.mu.Lock()
		if m.active == s {
			m.active = nil
		}
		m.mu.Unlock()
	}()

	jobsCh := make(chan model.StartJob)
	var wg sync.WaitGroup
	workerCount := req.AppConfigSnapshot.MaxConcurrentJobs
	if workerCount < 1 {
		workerCount = 1
	}
	for workerID := 0; workerID < workerCount; workerID++ {
		wg.Add(1)
		go func(wid int) {
			defer wg.Done()
			for j := range jobsCh {
				if s.session.AbortRequested {
					return
				}
				if s.session.StopRequested {
					m.markSkipped(s, j.JobID, "stop requested")
					continue
				}
				m.runJob(ctx, s, req, adapter, encoderPath, wid, j)
			}
		}(workerID)
	}
	for _, j := range req.Jobs {
		jobsCh <- j
	}
	close(jobsCh)
	wg.Wait()

	if !s.session.AbortRequested {
		m.executePostAction(req.AppConfigSnapshot)
	}
}

func (m *Manager) runJob(ctx context.Context, s *sessionRuntime, req model.StartEncodeRequest, adapter encoder.Adapter, encoderPath string, workerID int, startJob model.StartJob) {
	job := s.jobs[startJob.JobID]
	resolve, err := m.resolver.Resolve(startJob.InputPath, req.AppConfigSnapshot)
	if err != nil {
		m.finishJob(s, job, model.JobStatusFailed, 1, err.Error(), "", false, "")
		return
	}
	if resolve.NeedsOverwrite {
		decision, err := m.waitOverwriteDecision(s, job.JobID, resolve.FinalOutputPath)
		if err != nil {
			m.finishJob(s, job, model.JobStatusFailed, 1, err.Error(), resolve.FinalOutputPath, false, "")
			return
		}
		switch decision {
		case model.OverwriteDecisionSkip:
			m.finishJob(s, job, model.JobStatusSkipped, 0, "", resolve.FinalOutputPath, false, "")
			return
		case model.OverwriteDecisionAbort:
			_ = m.RequestAbort(s.session.SessionID)
			m.finishJob(s, job, model.JobStatusCancelled, -1, "aborted by overwrite decision", resolve.FinalOutputPath, false, "")
			return
		case model.OverwriteDecisionOverwrite:
		default:
		}
	}

	if err := m.appendTempArtifact(resolve.TempOutputPath); err != nil {
		m.emit(eventWarning, map[string]any{"session_id": s.session.SessionID, "job_id": job.JobID, "message": "failed to append temp index", "error": err.Error()})
	}

	buildReq := encoder.BuildRequest{Profile: req.Profile, AppConfig: req.AppConfigSnapshot, InputPath: startJob.InputPath, OutputPath: resolve.TempOutputPath}
	build, err := adapter.BuildArgs(buildReq)
	if err != nil {
		m.finishJob(s, job, model.JobStatusFailed, 1, err.Error(), resolve.FinalOutputPath, false, "")
		return
	}

	jobCtx, cancel := context.WithCancel(s.ctx)
	s.mu.Lock()
	s.jobCancels[job.JobID] = cancel
	s.session.RunningJobs++
	now := time.Now().UTC()
	job.Status = model.JobStatusRunning
	job.StartedAt = now
	job.WorkerID = &workerID
	s.mu.Unlock()
	m.emit(eventJobStarted, map[string]any{
		"session_id":       s.session.SessionID,
		"job_id":           job.JobID,
		"worker_id":        workerID,
		"input_path":       job.InputPath,
		"temp_output_path": resolve.TempOutputPath,
		"encoder_type":     req.Profile.EncoderType,
	})
	m.emitSessionState(s)

	stderrFile, _, err := m.logger.OpenStderrLog(job.JobID)
	if err != nil {
		cancel()
		m.finishJob(s, job, model.JobStatusFailed, 1, err.Error(), resolve.FinalOutputPath, false, "")
		return
	}
	defer stderrFile.Close()
	stderrWriter := logging.NewSafeWriter(stderrFile)

	lastEmit := time.Time{}
	onProgress := func(line string) {
		progress, ok := adapter.ParseProgress(line)
		if ok {
			now := time.Now()
			if now.Sub(lastEmit) >= 500*time.Millisecond {
				m.emit(eventJobProgress, map[string]any{
					"session_id":   s.session.SessionID,
					"job_id":       job.JobID,
					"percent":      progress.Percent,
					"fps":          progress.FPS,
					"bitrate_kbps": progress.BitrateKbps,
					"eta_sec":      progress.ETASec,
					"raw_line":     progress.RawLine,
				})
				lastEmit = now
			}
		}
	}
	onLog := func(line string) {
		_ = stderrWriter.WriteLine(line)
		m.emit(eventJobLog, map[string]any{"session_id": s.session.SessionID, "job_id": job.JobID, "line": line, "ts": time.Now().UTC()})
	}

	runSpec := encoder.RunSpec{
		Executable:           encoderPath,
		Argv:                 build.Argv,
		NoOutputTimeoutSec:   req.AppConfigSnapshot.NoOutputTimeoutSec,
		NoProgressTimeoutSec: req.AppConfigSnapshot.NoProgressTimeoutSec,
	}
	result := m.runner.Run(jobCtx, runSpec, onProgress, onLog)

	retryApplied := false
	retryDetail := ""
	if result.ExitCode != 0 && req.AppConfigSnapshot.DecoderFallback {
		if retryBuild, ok, retryErr := adapter.BuildRetryArgs(buildReq, build); retryErr == nil && ok {
			retryApplied = true
			retryDetail = "nvencc: avhw->avsw"
			onLog("avhw failed, retrying with avsw")
			retrySpec := runSpec
			retrySpec.Argv = retryBuild.Argv
			result = m.runner.Run(jobCtx, retrySpec, onProgress, onLog)
		}
	}

	cancel()
	s.mu.Lock()
	delete(s.jobCancels, job.JobID)
	s.session.RunningJobs--
	s.mu.Unlock()

	finalStatus := model.JobStatusCompleted
	errMsg := ""
	if result.TimedOut {
		finalStatus = model.JobStatusTimeout
		errMsg = "process timed out"
	} else if result.ExitCode != 0 {
		if errors.Is(jobCtx.Err(), context.Canceled) {
			finalStatus = model.JobStatusCancelled
			errMsg = "job cancelled"
		} else {
			finalStatus = model.JobStatusFailed
			if result.Err != nil {
				errMsg = result.Err.Error()
			} else {
				errMsg = fmt.Sprintf("exit code: %d", result.ExitCode)
			}
		}
	}

	if finalStatus == model.JobStatusCompleted {
		if err := os.Rename(resolve.TempOutputPath, resolve.FinalOutputPath); err != nil {
			finalStatus = model.JobStatusFailed
			errMsg = err.Error()
		} else if req.Profile.RestoreFileTime {
			if err := metadata.RestoreFileTime(startJob.InputPath, resolve.FinalOutputPath); err != nil {
				m.emit(eventWarning, map[string]any{"session_id": s.session.SessionID, "job_id": job.JobID, "message": "failed to restore file time", "error": err.Error()})
			}
		}
	}
	if finalStatus != model.JobStatusCompleted && !req.AppConfigSnapshot.KeepFailedTemp {
		_ = os.Remove(resolve.TempOutputPath)
	}
	_ = m.removeTempArtifact(resolve.TempOutputPath)

	m.finishJob(s, job, finalStatus, result.ExitCode, errMsg, resolve.FinalOutputPath, retryApplied, retryDetail)
	if (finalStatus == model.JobStatusFailed || finalStatus == model.JobStatusTimeout) && req.AppConfigSnapshot.OnError == model.OnErrorStop {
		s.mu.Lock()
		s.session.StopRequested = true
		s.mu.Unlock()
	}

	record := model.JobRecord{
		SchemaVersion:     1,
		AppVersion:        "v1.0.0",
		JobID:             job.JobID,
		ProfileID:         req.Profile.ID,
		ProfileName:       req.Profile.Name,
		ProfileVersion:    req.Profile.Version,
		InputPath:         startJob.InputPath,
		TempOutputPath:    resolve.TempOutputPath,
		FinalOutputPath:   resolve.FinalOutputPath,
		EncoderType:       req.Profile.EncoderType,
		EncoderPath:       encoderPath,
		Argv:              build.Argv,
		Device:            req.Profile.Device,
		MaxConcurrentJobs: req.AppConfigSnapshot.MaxConcurrentJobs,
		WorkerID:          workerID,
		UsedJobObject:     result.UsedJobObject,
		StartedAt:         job.StartedAt,
		FinishedAt:        job.FinishedAt,
		ExitCode:          result.ExitCode,
		Status:            finalStatus,
		ErrorMessage:      errMsg,
		RetryApplied:      retryApplied,
		RetryDetail:       retryDetail,
	}
	if err := m.logger.WriteJobRecord(record); err != nil {
		m.emit(eventWarning, map[string]any{"session_id": s.session.SessionID, "job_id": job.JobID, "message": "failed to write job record", "error": err.Error()})
	}
	m.resolver.Release(resolve.FinalOutputPath)
}

func (m *Manager) waitOverwriteDecision(s *sessionRuntime, jobID, finalPath string) (model.ResolveOverwriteDecision, error) {
	ch := make(chan model.ResolveOverwriteDecision, 1)
	s.mu.Lock()
	s.overwriteWaiters[jobID] = overwritePending{jobID: jobID, ch: ch}
	s.mu.Unlock()
	m.emit(eventJobNeedsOverwrite, map[string]any{"session_id": s.session.SessionID, "job_id": jobID, "final_output_path": finalPath})
	select {
	case decision := <-ch:
		if decision == "" {
			return model.OverwriteDecisionSkip, nil
		}
		return decision, nil
	case <-time.After(10 * time.Minute):
		return model.OverwriteDecisionSkip, nil
	case <-s.ctx.Done():
		return model.OverwriteDecisionAbort, s.ctx.Err()
	}
}

func (m *Manager) finishJob(s *sessionRuntime, job *model.QueueJob, status model.JobStatus, exitCode int, errMessage string, finalPath string, retryApplied bool, retryDetail string) {
	now := time.Now().UTC()
	s.mu.Lock()
	job.Status = status
	job.FinishedAt = now
	job.ExitCode = &exitCode
	job.ErrorMessage = errMessage
	switch status {
	case model.JobStatusCompleted:
		s.session.CompletedJobs++
	case model.JobStatusFailed:
		s.session.FailedJobs++
	case model.JobStatusCancelled:
		s.session.CancelledJobs++
	case model.JobStatusTimeout:
		s.session.TimeoutJobs++
	case model.JobStatusSkipped:
		s.session.SkippedJobs++
	}
	s.mu.Unlock()

	m.emit(eventJobFinished, map[string]any{
		"session_id":        s.session.SessionID,
		"job_id":            job.JobID,
		"status":            status,
		"exit_code":         exitCode,
		"error_message":     errMessage,
		"final_output_path": finalPath,
		"retry_applied":     retryApplied,
		"retry_detail":      retryDetail,
	})
	m.emitSessionState(s)
}

func (m *Manager) markSkipped(s *sessionRuntime, jobID, reason string) {
	job := s.jobs[jobID]
	if job == nil {
		return
	}
	m.finishJob(s, job, model.JobStatusSkipped, 0, reason, "", false, "")
}

func (m *Manager) requireSession(sessionID string) (*sessionRuntime, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.active == nil || m.active.session.SessionID != sessionID {
		return nil, model.NewError(model.ErrValidation, "session not found")
	}
	return m.active, nil
}

func (m *Manager) emitSessionState(s *sessionRuntime) {
	s.mu.Lock()
	payload := s.session
	s.mu.Unlock()
	m.emit(eventSessionState, payload)
}

func (m *Manager) emit(name string, payload any) {
	if m.emitter != nil {
		m.emitter.Emit(name, payload)
	}
}

func newSessionID() string {
	return uuid.NewString()
}

func (m *Manager) executePostAction(cfg model.AppConfig) {
	err := runPostAction(cfg)
	if m.appLogger != nil {
		if err != nil {
			m.appLogger.Error("post action failed", "action", cfg.PostCompleteAction, "error", err.Error())
		} else {
			m.appLogger.Info("post action completed", "action", cfg.PostCompleteAction)
		}
	}
}

func (m *Manager) tempIndexPath() string {
	return filepath.Join(m.runtimeDir, "temp_index.json")
}

func (m *Manager) loadTempIndex() (model.TempArtifactIndex, error) {
	path := m.tempIndexPath()
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return model.TempArtifactIndex{}, nil
		}
		return model.TempArtifactIndex{}, err
	}
	var idx model.TempArtifactIndex
	if err := json.Unmarshal(b, &idx); err != nil {
		return model.TempArtifactIndex{}, err
	}
	return idx, nil
}

func (m *Manager) saveTempIndex(idx model.TempArtifactIndex) error {
	path := m.tempIndexPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (m *Manager) appendTempArtifact(path string) error {
	idx, err := m.loadTempIndex()
	if err != nil {
		return err
	}
	for _, art := range idx.Artifacts {
		if art.Path == path {
			return nil
		}
	}
	idx.Artifacts = append(idx.Artifacts, model.TempArtifact{Path: path})
	return m.saveTempIndex(idx)
}

func (m *Manager) removeTempArtifact(path string) error {
	idx, err := m.loadTempIndex()
	if err != nil {
		return err
	}
	filtered := make([]model.TempArtifact, 0, len(idx.Artifacts))
	for _, art := range idx.Artifacts {
		if art.Path != path {
			filtered = append(filtered, art)
		}
	}
	idx.Artifacts = filtered
	return m.saveTempIndex(idx)
}
