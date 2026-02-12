package qsvenc

import (
	"context"

	"github.com/motoacs/enque/backend/encoder"
	"github.com/motoacs/enque/backend/model"
)

type Adapter struct{}

func NewAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) Type() model.EncoderType { return model.EncoderTypeQSVEnc }

func (a *Adapter) BuildArgs(req encoder.BuildRequest) (encoder.BuildResult, error) {
	return encoder.BuildResult{}, &model.EnqueError{Code: model.ErrEncoderNotImplemented, Message: "qsvenc adapter is not implemented in v1"}
}

func (a *Adapter) BuildRetryArgs(req encoder.BuildRequest, previous encoder.BuildResult) (*encoder.BuildResult, bool, error) {
	return nil, false, nil
}

func (a *Adapter) ParseProgress(line string) (model.JobProgress, bool) {
	return model.JobProgress{RawLine: line}, false
}

func (a *Adapter) ValidateProfile(profile model.Profile) error { return nil }

func (a *Adapter) DetectCapabilities(ctx context.Context, encoderPath string) (map[string]any, error) {
	return nil, &model.EnqueError{Code: model.ErrEncoderNotImplemented, Message: "qsvenc capabilities are not implemented in v1"}
}
