package queue

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yuta/enque/backend/config"
	"github.com/yuta/enque/backend/encoder"
	"github.com/yuta/enque/backend/events"
	"github.com/yuta/enque/backend/logging"
	"github.com/yuta/enque/backend/metadata"
	"github.com/yuta/enque/backend/profile"
)

// Worker executes encoding jobs from the job channel.
type Worker struct {
	id            int
	session       *Session
	adapter       encoder.Adapter
	runner        *encoder.ProcessRunner
	resolver      *OutputResolver
	tempTracker   *TempTracker
	emitter       *events.Emitter
	manager       *Manager
	prof          profile.Profile
	appCfg        AppConfigSnapshot
	encoderPath   string
	cancelJobFunc context.CancelFunc
	cancelJobMu   chan struct{} // Protects cancelJobFunc access
}

// WorkerConfig holds dependencies for creating a Worker.
type WorkerConfig struct {
	ID          int
	Session     *Session
	Adapter     encoder.Adapter
	Runner      *encoder.ProcessRunner
	Resolver    *OutputResolver
	TempTracker *TempTracker
	Emitter     *events.Emitter
	Manager     *Manager
	Profile     profile.Profile
	AppConfig   AppConfigSnapshot
	EncoderPath string
}

// NewWorker creates a worker with the given config.
func NewWorker(cfg WorkerConfig) *Worker {
	return &Worker{
		id:          cfg.ID,
		session:     cfg.Session,
		adapter:     cfg.Adapter,
		runner:      cfg.Runner,
		resolver:    cfg.Resolver,
		tempTracker: cfg.TempTracker,
		emitter:     cfg.Emitter,
		manager:     cfg.Manager,
		prof:        cfg.Profile,
		appCfg:      cfg.AppConfig,
		encoderPath: cfg.EncoderPath,
		cancelJobMu: make(chan struct{}, 1),
	}
}

// Run processes jobs from the channel until it's closed or the session stops.
func (w *Worker) Run(ctx context.Context, jobs <-chan *QueueJob) {
	for job := range jobs {
		if w.session.IsStopping() || w.session.ShouldSkipJob(job.JobID) {
			w.session.MarkJobStatus(job.JobID, JobSkipped, nil, "skipped by user")
			w.emitJobFinished(job, JobSkipped, nil, "skipped by user")
			continue
		}

		w.executeJob(ctx, job)

		// Check on_error=stop policy
		if job.Status == JobFailed && w.appCfg.OnError == "stop" {
			w.session.RequestStop()
		}
	}
}

// CancelCurrentJob cancels the currently running job on this worker.
func (w *Worker) CancelCurrentJob() {
	w.cancelJobMu <- struct{}{}
	if w.cancelJobFunc != nil {
		w.cancelJobFunc()
	}
	<-w.cancelJobMu
}

func (w *Worker) executeJob(ctx context.Context, job *QueueJob) {
	// Resolve output paths
	outputCfg := OutputConfig{
		FolderMode:    w.appCfg.OutputFolderMode,
		FolderPath:    w.appCfg.OutputFolderPath,
		NameTemplate:  w.appCfg.OutputNameTemplate,
		Container:     w.prof.OutputContainer,
		OverwriteMode: w.appCfg.OverwriteMode,
	}

	resolved, err := w.resolver.Resolve(job.InputPath, outputCfg)
	if err != nil {
		exitCode := -1
		w.session.MarkJobStatus(job.JobID, JobSkipped, &exitCode, err.Error())
		w.emitJobFinished(job, JobSkipped, &exitCode, err.Error())
		return
	}

	job.TempOutputPath = resolved.TempPath
	job.FinalOutputPath = resolved.FinalPath

	// Get input file size
	if info, err := os.Stat(job.InputPath); err == nil {
		job.InputSizeBytes = info.Size()
	}

	// Mark as running
	job.Status = JobRunning
	job.WorkerID = w.id
	job.StartedAt = time.Now()

	w.emitJobStarted(job)

	// Handle overwrite confirmation (ask mode)
	if resolved.NeedsOverwrite {
		// Emit event to frontend for user decision
		w.emitter.JobNeedsOverwrite(map[string]interface{}{
			"session_id":        w.session.ID,
			"job_id":            job.JobID,
			"final_output_path": resolved.FinalPath,
		})

		// Wait for user response (10-minute timeout → skip)
		decision, _ := w.manager.WaitForOverwrite(job.JobID)
		switch decision {
		case "overwrite":
			// Continue with encoding
		case "abort":
			w.session.RequestStop()
			exitCode := -1
			w.session.MarkJobStatus(job.JobID, JobSkipped, &exitCode, "overwrite aborted by user")
			w.emitJobFinished(job, JobSkipped, &exitCode, "overwrite aborted by user")
			w.resolver.Release(resolved.FinalPath)
			return
		default: // "skip" or timeout
			exitCode := -1
			w.session.MarkJobStatus(job.JobID, JobSkipped, &exitCode, "overwrite skipped by user")
			w.emitJobFinished(job, JobSkipped, &exitCode, "overwrite skipped by user")
			w.resolver.Release(resolved.FinalPath)
			return
		}
	}

	// Build args
	args, err := w.adapter.BuildArgs(w.prof, job.InputPath, resolved.TempPath)
	if err != nil {
		exitCode := -1
		w.session.MarkJobStatus(job.JobID, JobFailed, &exitCode, err.Error())
		w.emitJobFinished(job, JobFailed, &exitCode, err.Error())
		w.resolver.Release(resolved.FinalPath)
		return
	}

	// Track temp file
	w.tempTracker.Add(TempEntry{
		TempPath:  resolved.TempPath,
		FinalPath: resolved.FinalPath,
		JobID:     job.JobID,
		SessionID: w.session.ID,
	})

	// Set up job log
	logsDir := filepath.Join(config.LogsDir(), w.session.ID)
	stderrWriter, err := logging.NewStderrWriter(logsDir, job.JobID)
	if err != nil {
		exitCode := -1
		w.session.MarkJobStatus(job.JobID, JobFailed, &exitCode, fmt.Sprintf("create stderr log: %v", err))
		w.emitJobFinished(job, JobFailed, &exitCode, err.Error())
		w.resolver.Release(resolved.FinalPath)
		return
	}
	defer stderrWriter.Close()

	// Create cancellable context for this job
	jobCtx, cancel := context.WithCancel(ctx)
	w.cancelJobMu <- struct{}{}
	w.cancelJobFunc = cancel
	<-w.cancelJobMu
	defer func() {
		w.cancelJobMu <- struct{}{}
		w.cancelJobFunc = nil
		<-w.cancelJobMu
		cancel()
	}()

	// Execute encoder
	result := w.runner.Run(jobCtx, args, stderrWriter,
		func(progress encoder.Progress) {
			w.emitJobProgress(job, progress)
		},
		func(line string) {
			w.emitJobLog(job, line)
		},
	)

	// Handle result
	status := w.determineJobStatus(result, jobCtx)
	w.session.MarkJobStatus(job.JobID, status, &result.ExitCode, result.ErrorMessage)

	// Post-process
	if status == JobCompleted {
		w.postProcessSuccess(job, resolved)
	} else {
		w.postProcessFailure(job, resolved, status)
	}

	// Try decoder fallback if applicable
	if status == JobFailed && w.appCfg.DecoderFallback && w.adapter.SupportsDecoderFallback() {
		if w.prof.Decoder == "avhw" {
			w.retryWithFallback(ctx, job, resolved, logsDir)
			return
		}
	}

	// Save job record
	w.saveJobRecord(job, resolved, args, result, status, false, "")

	w.emitJobFinished(job, status, &result.ExitCode, result.ErrorMessage)
}

func (w *Worker) retryWithFallback(ctx context.Context, job *QueueJob, resolved *ResolveResult, logsDir string) {
	// Build args with avsw decoder
	overriddenProfile := w.prof
	overriddenProfile.Decoder = "avsw"

	args, err := w.adapter.BuildArgs(overriddenProfile, job.InputPath, resolved.TempPath)
	if err != nil {
		return
	}

	stderrWriter, err := logging.NewStderrWriter(logsDir, job.JobID+"_retry")
	if err != nil {
		return
	}
	defer stderrWriter.Close()

	w.emitter.Warning(map[string]interface{}{
		"session_id": w.session.ID,
		"job_id":     job.JobID,
		"message":    "retrying with software decoder (avsw)",
	})

	jobCtx, cancel := context.WithCancel(ctx)
	w.cancelJobMu <- struct{}{}
	w.cancelJobFunc = cancel
	<-w.cancelJobMu
	defer func() {
		w.cancelJobMu <- struct{}{}
		w.cancelJobFunc = nil
		<-w.cancelJobMu
		cancel()
	}()

	result := w.runner.Run(jobCtx, args, stderrWriter,
		func(progress encoder.Progress) {
			w.emitJobProgress(job, progress)
		},
		func(line string) {
			w.emitJobLog(job, line)
		},
	)

	status := w.determineJobStatus(result, jobCtx)
	w.session.MarkJobStatus(job.JobID, status, &result.ExitCode, result.ErrorMessage)

	if status == JobCompleted {
		w.postProcessSuccess(job, resolved)
	} else {
		w.postProcessFailure(job, resolved, status)
	}

	w.saveJobRecord(job, resolved, args, result, status, true, "avhw -> avsw fallback")
	w.emitJobFinished(job, status, &result.ExitCode, result.ErrorMessage)
}

func (w *Worker) determineJobStatus(result encoder.RunResult, ctx context.Context) JobStatus {
	if result.TimedOut {
		return JobTimeout
	}
	if ctx.Err() != nil {
		if w.session.IsAborting() {
			return JobCancelled
		}
		return JobCancelled
	}
	if result.ExitCode == 0 {
		return JobCompleted
	}
	return JobFailed
}

func (w *Worker) postProcessSuccess(job *QueueJob, resolved *ResolveResult) {
	// Rename temp to final
	if err := os.Rename(resolved.TempPath, resolved.FinalPath); err != nil {
		job.ErrorMessage = fmt.Sprintf("rename temp to final: %v", err)
	}

	// Remove from temp tracker
	w.tempTracker.Remove(resolved.TempPath)

	// Restore file time from input to output (Windows only)
	if w.prof.RestoreFileTime {
		if err := metadata.RestoreFileTimeIfNeeded(job.InputPath, resolved.FinalPath, true); err != nil {
			// Non-fatal, just log warning
			w.emitter.Warning(map[string]interface{}{
				"session_id": w.session.ID,
				"job_id":     job.JobID,
				"message":    fmt.Sprintf("failed to restore file time: %v", err),
			})
		}
	}
}

func (w *Worker) postProcessFailure(job *QueueJob, resolved *ResolveResult, status JobStatus) {
	// Clean up temp file unless keep_failed_temp is set
	if !w.appCfg.KeepFailedTemp {
		os.Remove(resolved.TempPath)
		w.tempTracker.Remove(resolved.TempPath)
	}

	w.resolver.Release(resolved.FinalPath)
}

func (w *Worker) saveJobRecord(job *QueueJob, resolved *ResolveResult, args []string, result encoder.RunResult, status JobStatus, retryApplied bool, retryDetail string) {
	logsDir := filepath.Join(config.LogsDir(), w.session.ID)
	record := &logging.JobRecord{
		SchemaVersion:     1,
		JobID:             job.JobID,
		SessionID:         w.session.ID,
		InputPath:         job.InputPath,
		OutputPath:        resolved.FinalPath,
		TempOutputPath:    resolved.TempPath,
		CommandLine:       append([]string{w.encoderPath}, args...),
		EncoderType:       w.adapter.Type(),
		EncoderPath:       w.encoderPath,
		ExitCode:          &result.ExitCode,
		Status:            string(status),
		ErrorMessage:      result.ErrorMessage,
		WorkerID:          w.id,
		AppVersion:        config.AppVersion,
		ProfileID:         w.prof.ID,
		ProfileName:       w.prof.Name,
		ProfileVersion:    w.prof.Version,
		Device:            w.prof.Device,
		MaxConcurrentJobs: w.appCfg.MaxConcurrentJobs,
		UsedJobObject:     result.UsedJobObject,
		StartedAt:         job.StartedAt.Format(time.RFC3339),
		FinishedAt:        time.Now().Format(time.RFC3339),
		DurationSec:       time.Since(job.StartedAt).Seconds(),
		RetryApplied:      retryApplied,
		RetryDetail:       retryDetail,
	}
	record.Save(logsDir)
}

// Event emission helpers

func (w *Worker) emitJobStarted(job *QueueJob) {
	w.emitter.JobStarted(map[string]interface{}{
		"session_id":        w.session.ID,
		"job_id":            job.JobID,
		"input_path":        job.InputPath,
		"input_size_bytes":  job.InputSizeBytes,
		"worker_id":         w.id,
		"temp_output_path":  job.TempOutputPath,
		"final_output_path": job.FinalOutputPath,
		"encoder_type":      w.adapter.Type(),
	})
}

func (w *Worker) emitJobProgress(job *QueueJob, progress encoder.Progress) {
	data := map[string]interface{}{
		"session_id": w.session.ID,
		"job_id":     job.JobID,
		"worker_id":  w.id,
	}
	if progress.Percent != nil {
		data["percent"] = *progress.Percent
	}
	if progress.FPS != nil {
		data["fps"] = *progress.FPS
	}
	if progress.BitrateKbps != nil {
		data["bitrate_kbps"] = *progress.BitrateKbps
	}
	if progress.ETASec != nil {
		data["eta_sec"] = *progress.ETASec
	}
	w.emitter.JobProgress(data)
}

func (w *Worker) emitJobLog(job *QueueJob, line string) {
	w.emitter.JobLog(map[string]interface{}{
		"session_id": w.session.ID,
		"job_id":     job.JobID,
		"line":       line,
		"ts":         time.Now().Format(time.RFC3339Nano),
	})
}

func (w *Worker) emitJobFinished(job *QueueJob, status JobStatus, exitCode *int, errMsg string) {
	data := map[string]interface{}{
		"session_id":        w.session.ID,
		"job_id":            job.JobID,
		"status":            string(status),
		"error_message":     errMsg,
		"temp_output_path":  job.TempOutputPath,
		"final_output_path": job.FinalOutputPath,
	}
	if exitCode != nil {
		data["exit_code"] = *exitCode
	}
	w.emitter.JobFinished(data)
}
