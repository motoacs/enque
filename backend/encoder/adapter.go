package encoder

import (
	"context"

	"github.com/motoacs/enque/backend/model"
)

type BuildRequest struct {
	Profile    model.Profile
	AppConfig  model.AppConfig
	InputPath  string
	OutputPath string
}

type BuildResult struct {
	Argv             []string
	DisplayCommand   string
	EffectiveDecoder model.Decoder
}

type Adapter interface {
	Type() model.EncoderType
	BuildArgs(req BuildRequest) (BuildResult, error)
	BuildRetryArgs(req BuildRequest, previous BuildResult) (*BuildResult, bool, error)
	ParseProgress(line string) (model.JobProgress, bool)
	ValidateProfile(profile model.Profile) error
	DetectCapabilities(ctx context.Context, encoderPath string) (map[string]any, error)
}
