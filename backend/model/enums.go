package model

import "time"

type Codec string

type EncoderType string

type RateControl string

type Preset string

type Multipass string

type SplitEnc string

type ParallelMode string

type Decoder string

type AudioMode string

type OnError string

type OverwriteMode string

type PostAction string

type Language string

type JobStatus string

const (
	CodecH264 Codec = "h264"
	CodecHEVC Codec = "hevc"
	CodecAV1  Codec = "av1"
)

const (
	EncoderTypeNVEncC EncoderType = "nvencc"
	EncoderTypeQSVEnc EncoderType = "qsvenc"
	EncoderTypeFFmpeg EncoderType = "ffmpeg"
)

const (
	RateControlQVBR RateControl = "qvbr"
	RateControlCQP  RateControl = "cqp"
	RateControlCBR  RateControl = "cbr"
	RateControlVBR  RateControl = "vbr"
)

const (
	MultipassNone    Multipass = "none"
	MultipassQuarter Multipass = "quarter"
	MultipassFull    Multipass = "full"
)

const (
	SplitEncOff        SplitEnc = "off"
	SplitEncAuto       SplitEnc = "auto"
	SplitEncAutoForced SplitEnc = "auto_forced"
	SplitEncForced2    SplitEnc = "forced_2"
	SplitEncForced3    SplitEnc = "forced_3"
	SplitEncForced4    SplitEnc = "forced_4"
)

const (
	ParallelOff  ParallelMode = "off"
	ParallelAuto ParallelMode = "auto"
	Parallel2    ParallelMode = "2"
	Parallel3    ParallelMode = "3"
)

const (
	DecoderAVHW Decoder = "avhw"
	DecoderAVSW Decoder = "avsw"
)

const (
	AudioModeCopy AudioMode = "copy"
	AudioModeAAC  AudioMode = "aac"
	AudioModeOpus AudioMode = "opus"
)

const (
	OnErrorSkip OnError = "skip"
	OnErrorStop OnError = "stop"
)

const (
	OverwriteModeAsk        OverwriteMode = "ask"
	OverwriteModeAutoRename OverwriteMode = "auto_rename"
)

const (
	PostActionNone     PostAction = "none"
	PostActionShutdown PostAction = "shutdown"
	PostActionSleep    PostAction = "sleep"
	PostActionCustom   PostAction = "custom"
)

const (
	LanguageJA Language = "ja"
	LanguageEN Language = "en"
)

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusTimeout   JobStatus = "timeout"
	JobStatusSkipped   JobStatus = "skipped"
)

type JobProgress struct {
	Percent     *float64 `json:"percent"`
	FPS         *float64 `json:"fps"`
	BitrateKbps *float64 `json:"bitrate_kbps"`
	ETASec      *int64   `json:"eta_sec"`
	RawLine     string   `json:"raw_line,omitempty"`
}

type QueueJob struct {
	JobID          string      `json:"job_id"`
	InputPath      string      `json:"input_path"`
	InputSizeBytes int64       `json:"input_size_bytes"`
	Status         JobStatus   `json:"status"`
	Progress       JobProgress `json:"progress"`
	StartedAt      time.Time   `json:"started_at,omitempty"`
	FinishedAt     time.Time   `json:"finished_at,omitempty"`
	WorkerID       *int        `json:"worker_id"`
	ExitCode       *int        `json:"exit_code"`
	ErrorMessage   string      `json:"error_message"`
}

type EncodeSession struct {
	SessionID      string    `json:"session_id"`
	State          string    `json:"state"`
	StartedAt      time.Time `json:"started_at"`
	FinishedAt     time.Time `json:"finished_at,omitempty"`
	TotalJobs      int       `json:"total_jobs"`
	CompletedJobs  int       `json:"completed_jobs"`
	RunningJobs    int       `json:"running_jobs"`
	FailedJobs     int       `json:"failed_jobs"`
	CancelledJobs  int       `json:"cancelled_jobs"`
	TimeoutJobs    int       `json:"timeout_jobs"`
	SkippedJobs    int       `json:"skipped_jobs"`
	StopRequested  bool      `json:"stop_requested"`
	AbortRequested bool      `json:"abort_requested"`
}
