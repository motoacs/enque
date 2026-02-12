package model

import "time"

type JobRecord struct {
	SchemaVersion int    `json:"schema_version"`
	AppVersion    string `json:"app_version"`
	JobID         string `json:"job_id"`

	ProfileID      string `json:"profile_id"`
	ProfileName    string `json:"profile_name"`
	ProfileVersion int    `json:"profile_version"`

	InputPath       string `json:"input_path"`
	TempOutputPath  string `json:"temp_output_path"`
	FinalOutputPath string `json:"final_output_path"`

	EncoderType EncoderType `json:"encoder_type"`
	EncoderPath string      `json:"encoder_path"`
	Argv        []string    `json:"argv"`
	Device      string      `json:"device"`

	MaxConcurrentJobs int  `json:"max_concurrent_jobs"`
	WorkerID          int  `json:"worker_id"`
	UsedJobObject     bool `json:"used_job_object"`

	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
	ExitCode     int       `json:"exit_code"`
	Status       JobStatus `json:"status"`
	ErrorMessage string    `json:"error_message"`
	RetryApplied bool      `json:"retry_applied"`
	RetryDetail  string    `json:"retry_detail"`
}
