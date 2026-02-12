package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// JobRecord holds the execution record for a single job (design doc 5.5).
type JobRecord struct {
	SchemaVersion  int      `json:"schema_version"`
	JobID          string   `json:"job_id"`
	SessionID      string   `json:"session_id"`
	InputPath      string   `json:"input_path"`
	OutputPath     string   `json:"output_path"`
	TempOutputPath string   `json:"temp_output_path"`
	CommandLine    []string `json:"command_line"`
	EncoderType    string   `json:"encoder_type"`
	EncoderPath    string   `json:"encoder_path"`
	ExitCode       *int     `json:"exit_code"`
	Status         string   `json:"status"`
	ErrorMessage   string   `json:"error_message,omitempty"`
	WorkerID       int      `json:"worker_id"`
	DeviceUsed     string   `json:"device_used,omitempty"`
	AppVersion        string `json:"app_version"`
	ProfileID         string `json:"profile_id"`
	ProfileName       string `json:"profile_name"`
	ProfileVersion    int    `json:"profile_version"`
	Device            string `json:"device"`
	MaxConcurrentJobs int    `json:"max_concurrent_jobs"`
	UsedJobObject     bool   `json:"used_job_object"`
	StartedAt      string   `json:"started_at"`
	FinishedAt     string   `json:"finished_at"`
	DurationSec    float64  `json:"duration_sec"`
	RetryApplied   bool     `json:"retry_applied"`
	RetryDetail    string   `json:"retry_detail,omitempty"`
}

// Save writes the job record to {logsDir}/{jobID}.json atomically.
func (r *JobRecord) Save(logsDir string) error {
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return fmt.Errorf("create logs dir: %w", err)
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal job record: %w", err)
	}

	finalPath := filepath.Join(logsDir, r.JobID+".json")
	tmpPath := finalPath + ".tmp"

	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write temp: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("sync temp: %w", err)
	}
	f.Close()

	return os.Rename(tmpPath, finalPath)
}
