package encoder

import (
	"context"
	"testing"

	"github.com/motoacs/enque/backend/model"
)

type fakeAdapter struct{ t model.EncoderType }

func (f fakeAdapter) Type() model.EncoderType { return f.t }
func (f fakeAdapter) BuildArgs(req BuildRequest) (BuildResult, error) {
	return BuildResult{}, nil
}
func (f fakeAdapter) BuildRetryArgs(req BuildRequest, previous BuildResult) (*BuildResult, bool, error) {
	return nil, false, nil
}
func (f fakeAdapter) ParseProgress(line string) (model.JobProgress, bool) {
	return model.JobProgress{}, false
}
func (f fakeAdapter) ValidateProfile(profile model.Profile) error { return nil }
func (f fakeAdapter) DetectCapabilities(ctx context.Context, encoderPath string) (map[string]any, error) {
	return nil, nil
}

func TestRegistryResolve(t *testing.T) {
	r := NewRegistry(fakeAdapter{t: model.EncoderTypeNVEncC})
	if _, err := r.Resolve(model.EncoderTypeNVEncC); err != nil {
		t.Fatalf("expected adapter, got err: %v", err)
	}
	if _, err := r.Resolve(model.EncoderTypeFFmpeg); err == nil {
		t.Fatalf("expected not implemented error")
	}
}
