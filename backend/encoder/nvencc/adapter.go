package nvencc

import (
	"context"

	"github.com/motoacs/enque/backend/encoder"
	"github.com/motoacs/enque/backend/model"
)

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Type() model.EncoderType {
	return model.EncoderTypeNVEncC
}

func (a *Adapter) BuildArgs(req encoder.BuildRequest) (encoder.BuildResult, error) {
	return BuildCommand(req)
}

func (a *Adapter) BuildRetryArgs(req encoder.BuildRequest, previous encoder.BuildResult) (*encoder.BuildResult, bool, error) {
	if previous.EffectiveDecoder != model.DecoderAVHW {
		return nil, false, nil
	}
	retryReq := req
	retryReq.Profile.Decoder = model.DecoderAVSW
	res, err := BuildCommand(retryReq)
	if err != nil {
		return nil, false, err
	}
	return &res, true, nil
}

func (a *Adapter) ParseProgress(line string) (model.JobProgress, bool) {
	return ParseProgress(line)
}

func (a *Adapter) ValidateProfile(profile model.Profile) error {
	return model.MustValidateProfile(profile)
}

func (a *Adapter) DetectCapabilities(ctx context.Context, encoderPath string) (map[string]any, error) {
	return DetectCapabilities(ctx, encoderPath)
}
