package model

import "strings"

const AppConfigVersion = 1

type AppConfig struct {
	Version              int           `json:"version"`
	NVEncCPath           string        `json:"nvencc_path"`
	QSVEncPath           string        `json:"qsvenc_path"`
	FFmpegPath           string        `json:"ffmpeg_path"`
	FFprobePath          string        `json:"ffprobe_path"`
	MaxConcurrentJobs    int           `json:"max_concurrent_jobs"`
	OnError              OnError       `json:"on_error"`
	DecoderFallback      bool          `json:"decoder_fallback"`
	KeepFailedTemp       bool          `json:"keep_failed_temp"`
	NoOutputTimeoutSec   int           `json:"no_output_timeout_sec"`
	NoProgressTimeoutSec int           `json:"no_progress_timeout_sec"`
	PostCompleteAction   PostAction    `json:"post_complete_action"`
	PostCompleteCommand  string        `json:"post_complete_command"`
	OutputFolderMode     string        `json:"output_folder_mode"`
	OutputFolderPath     string        `json:"output_folder_path"`
	OutputNameTemplate   string        `json:"output_name_template"`
	OutputContainer      string        `json:"output_container"`
	OverwriteMode        OverwriteMode `json:"overwrite_mode"`
	Language             Language      `json:"language"`
	DefaultProfileID     string        `json:"default_profile_id"`
}

func DefaultAppConfig() AppConfig {
	return AppConfig{
		Version:              AppConfigVersion,
		MaxConcurrentJobs:    1,
		OnError:              OnErrorSkip,
		DecoderFallback:      false,
		KeepFailedTemp:       false,
		NoOutputTimeoutSec:   600,
		NoProgressTimeoutSec: 300,
		PostCompleteAction:   PostActionNone,
		OutputFolderMode:     "same_as_input",
		OutputNameTemplate:   "{name}_encoded.{ext}",
		OutputContainer:      "mkv",
		OverwriteMode:        OverwriteModeAsk,
		Language:             LanguageJA,
	}
}

func ValidateAppConfig(c AppConfig) map[string]string {
	errs := map[string]string{}
	if c.MaxConcurrentJobs < 1 || c.MaxConcurrentJobs > 8 {
		errs["max_concurrent_jobs"] = "must be 1..8"
	}
	if c.NoOutputTimeoutSec < 30 || c.NoOutputTimeoutSec > 86400 {
		errs["no_output_timeout_sec"] = "must be 30..86400"
	}
	if c.NoProgressTimeoutSec < 30 || c.NoProgressTimeoutSec > 86400 {
		errs["no_progress_timeout_sec"] = "must be 30..86400"
	}
	if n := strings.TrimSpace(c.OutputNameTemplate); len(n) < 1 || len(n) > 255 {
		errs["output_name_template"] = "must be 1..255"
	}
	if c.OutputFolderMode == "specified" && strings.TrimSpace(c.OutputFolderPath) == "" {
		errs["output_folder_path"] = "required when output_folder_mode=specified"
	}
	if c.PostCompleteAction == PostActionCustom && strings.TrimSpace(c.PostCompleteCommand) == "" {
		errs["post_complete_command"] = "required when post_complete_action=custom"
	}
	switch c.OnError {
	case OnErrorSkip, OnErrorStop:
	default:
		errs["on_error"] = "must be skip or stop"
	}
	switch c.OverwriteMode {
	case OverwriteModeAsk, OverwriteModeAutoRename:
	default:
		errs["overwrite_mode"] = "must be ask or auto_rename"
	}
	switch c.Language {
	case LanguageJA, LanguageEN:
	default:
		errs["language"] = "must be ja or en"
	}
	return errs
}
