package queue

import (
	"sync"
	"time"

	"github.com/yuta/enque/backend/profile"
)

// SessionState tracks the encoding session lifecycle.
type SessionState string

const (
	StateRunning   SessionState = "running"
	StateStopping  SessionState = "stopping"
	StateAborting  SessionState = "aborting"
	StateCompleted SessionState = "completed"
	StateAborted   SessionState = "aborted"
)

// JobStatus tracks individual job state.
type JobStatus string

const (
	JobPending   JobStatus = "pending"
	JobRunning   JobStatus = "running"
	JobCompleted JobStatus = "completed"
	JobFailed    JobStatus = "failed"
	JobCancelled JobStatus = "cancelled"
	JobTimeout   JobStatus = "timeout"
	JobSkipped   JobStatus = "skipped"
)

// QueueJob represents a single encoding job in the session.
type QueueJob struct {
	JobID          string    `json:"job_id"`
	InputPath      string    `json:"input_path"`
	InputSizeBytes int64     `json:"input_size_bytes"`
	Status         JobStatus `json:"status"`
	WorkerID       int       `json:"worker_id"`
	ExitCode       *int      `json:"exit_code"`
	ErrorMessage   string    `json:"error_message"`
	TempOutputPath string    `json:"temp_output_path"`
	FinalOutputPath string   `json:"final_output_path"`
	StartedAt      time.Time `json:"started_at"`
	FinishedAt     time.Time `json:"finished_at"`
}

// EncodeRequest is the input from StartEncode (design doc 6.2).
type EncodeRequest struct {
	Jobs              []JobInput       `json:"jobs"`
	Profile           profile.Profile  `json:"profile"`
	AppConfigSnapshot AppConfigSnapshot `json:"app_config_snapshot"`
}

// JobInput is a single job in the StartEncode request.
type JobInput struct {
	JobID     string `json:"job_id"`
	InputPath string `json:"input_path"`
}

// AppConfigSnapshot captures config at encode start time.
type AppConfigSnapshot struct {
	MaxConcurrentJobs    int    `json:"max_concurrent_jobs"`
	OnError              string `json:"on_error"`
	DecoderFallback      bool   `json:"decoder_fallback"`
	KeepFailedTemp       bool   `json:"keep_failed_temp"`
	NoOutputTimeoutSec   int    `json:"no_output_timeout_sec"`
	NoProgressTimeoutSec int    `json:"no_progress_timeout_sec"`
	PostCompleteAction   string `json:"post_complete_action"`
	PostCompleteCommand  string `json:"post_complete_command"`
	OutputFolderMode     string `json:"output_folder_mode"`
	OutputFolderPath     string `json:"output_folder_path"`
	OutputNameTemplate   string `json:"output_name_template"`
	OutputContainer      string `json:"output_container"`
	OverwriteMode        string `json:"overwrite_mode"`
	NVEncCPath           string `json:"nvencc_path"`
}

// Session manages state for a single encoding session.
type Session struct {
	mu             sync.RWMutex
	ID             string
	State          SessionState
	Jobs           []*QueueJob
	StartedAt      time.Time
	FinishedAt     time.Time
	StopRequested  bool
	AbortRequested bool
	EncoderType    string
	AppCfg         AppConfigSnapshot

	// Skip set for individual job skipping
	SkipSet map[string]bool

	// Counters
	TotalJobs     int
	CompletedJobs int
	FailedJobs    int
	CancelledJobs int
	TimeoutJobs   int
	SkippedJobs   int
}

// NewSession creates a new encoding session.
func NewSession(id string, jobs []JobInput, encoderType string, appCfg AppConfigSnapshot) *Session {
	queueJobs := make([]*QueueJob, len(jobs))
	for i, j := range jobs {
		queueJobs[i] = &QueueJob{
			JobID:     j.JobID,
			InputPath: j.InputPath,
			Status:    JobPending,
		}
	}
	return &Session{
		ID:          id,
		State:       StateRunning,
		Jobs:        queueJobs,
		StartedAt:   time.Now(),
		TotalJobs:   len(jobs),
		EncoderType: encoderType,
		AppCfg:      appCfg,
		SkipSet:     make(map[string]bool),
	}
}

// RequestStop sets the stop flag (graceful).
func (s *Session) RequestStop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.State == StateRunning {
		s.State = StateStopping
		s.StopRequested = true
	}
}

// RequestAbort sets the abort flag.
func (s *Session) RequestAbort() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.State == StateRunning || s.State == StateStopping {
		s.State = StateAborting
		s.AbortRequested = true
	}
}

// RequestSkipJob marks a single job to be skipped.
func (s *Session) RequestSkipJob(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.SkipSet[jobID] = true
}

// ShouldSkipJob returns whether a specific job should be skipped.
func (s *Session) ShouldSkipJob(jobID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SkipSet[jobID]
}

// IsStopping returns whether stop or abort has been requested.
func (s *Session) IsStopping() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.StopRequested || s.AbortRequested
}

// IsAborting returns whether abort has been requested.
func (s *Session) IsAborting() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.AbortRequested
}

// MarkJobStatus updates a job's status and session counters.
func (s *Session) MarkJobStatus(jobID string, status JobStatus, exitCode *int, errMsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, j := range s.Jobs {
		if j.JobID == jobID {
			j.Status = status
			j.ExitCode = exitCode
			j.ErrorMessage = errMsg
			j.FinishedAt = time.Now()

			switch status {
			case JobCompleted:
				s.CompletedJobs++
			case JobFailed:
				s.FailedJobs++
			case JobCancelled:
				s.CancelledJobs++
			case JobTimeout:
				s.TimeoutJobs++
			case JobSkipped:
				s.SkippedJobs++
			}
			break
		}
	}
}

// Finish marks the session as completed or aborted.
func (s *Session) Finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FinishedAt = time.Now()
	if s.AbortRequested {
		s.State = StateAborted
	} else {
		s.State = StateCompleted
	}
}

// RunningJobs returns count of currently running jobs.
func (s *Session) RunningJobs() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.runningJobsLocked()
}

func (s *Session) runningJobsLocked() int {
	count := 0
	for _, j := range s.Jobs {
		if j.Status == JobRunning {
			count++
		}
	}
	return count
}

// Snapshot returns session state for events.
func (s *Session) Snapshot() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]interface{}{
		"session_id":      s.ID,
		"state":           s.State,
		"encoder_type":    s.EncoderType,
		"started_at":      s.StartedAt.Format(time.RFC3339),
		"total_jobs":      s.TotalJobs,
		"completed_jobs":  s.CompletedJobs,
		"running_jobs":    s.runningJobsLocked(),
		"failed_jobs":     s.FailedJobs,
		"cancelled_jobs":  s.CancelledJobs,
		"timeout_jobs":    s.TimeoutJobs,
		"skipped_jobs":    s.SkippedJobs,
		"stop_requested":  s.StopRequested,
		"abort_requested": s.AbortRequested,
	}
}
