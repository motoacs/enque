package config

// AppConfig holds application-level settings (design doc 5.3).
type AppConfig struct {
	Version              int    `json:"version"`
	NVEncCPath           string `json:"nvencc_path"`
	QSVEncPath           string `json:"qsvenc_path"`
	FFmpegPath           string `json:"ffmpeg_path"`
	FFprobePath          string `json:"ffprobe_path"`
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
	Language             string `json:"language"`
	DefaultProfileID     string `json:"default_profile_id"`
}

// CurrentVersion is the latest config schema version.
const CurrentVersion = 1

// Default returns the default AppConfig.
func Default() AppConfig {
	return AppConfig{
		Version:              CurrentVersion,
		MaxConcurrentJobs:    1,
		OnError:              "skip",
		DecoderFallback:      false,
		KeepFailedTemp:       false,
		NoOutputTimeoutSec:   600,
		NoProgressTimeoutSec: 300,
		PostCompleteAction:   "none",
		OutputFolderMode:     "same_as_input",
		OutputNameTemplate:   "{name}_encoded.{ext}",
		OutputContainer:      "mkv",
		OverwriteMode:        "ask",
		Language:             "ja",
	}
}
