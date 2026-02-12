package model

type StartEncodeRequest struct {
	Jobs              []StartJob `json:"jobs"`
	Profile           Profile    `json:"profile"`
	AppConfigSnapshot AppConfig  `json:"app_config_snapshot"`
	CommandPreview    string     `json:"command_preview,omitempty"`
}

type StartJob struct {
	JobID     string `json:"job_id"`
	InputPath string `json:"input_path"`
}

type PreviewCommandRequest struct {
	Profile           Profile   `json:"profile"`
	AppConfigSnapshot AppConfig `json:"app_config_snapshot"`
	InputPath         string    `json:"input_path"`
	OutputPath        string    `json:"output_path"`
}

type PreviewCommandResponse struct {
	Argv           []string `json:"argv"`
	DisplayCommand string   `json:"display_command"`
}

type BootstrapResponse struct {
	Config   AppConfig    `json:"config"`
	Profiles []Profile    `json:"profiles"`
	Tools    ToolSnapshot `json:"tools"`
	GPUInfo  GPUInfo      `json:"gpu_info"`
	Warnings []string     `json:"warnings"`
}

type ResolveOverwriteDecision string

const (
	OverwriteDecisionOverwrite ResolveOverwriteDecision = "overwrite"
	OverwriteDecisionSkip      ResolveOverwriteDecision = "skip"
	OverwriteDecisionAbort     ResolveOverwriteDecision = "abort"
)

type ToolInfo struct {
	Found   bool   `json:"found"`
	Path    string `json:"path"`
	Version string `json:"version"`
	Warning string `json:"warning,omitempty"`
}

type ToolSnapshot struct {
	NVEncC  ToolInfo `json:"nvencc"`
	QSVEnc  ToolInfo `json:"qsvenc"`
	FFmpeg  ToolInfo `json:"ffmpeg"`
	FFprobe ToolInfo `json:"ffprobe"`
}

type GPUInfo struct {
	CheckDeviceOutput   string `json:"check_device_output"`
	CheckFeaturesOutput string `json:"check_features_output"`
}
