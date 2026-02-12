package profile

import (
	"testing"

	"github.com/motoacs/enque/backend/model"
)

func TestMigrateProfileV1ToV4(t *testing.T) {
	p := model.Profile{
		Version: 1,
		Name:    "legacy",
	}
	m, changed := migrateOne(p)
	if !changed {
		t.Fatalf("expected changed")
	}
	if m.Version != model.ProfileVersion {
		t.Fatalf("expected version %d, got %d", model.ProfileVersion, m.Version)
	}
	if m.EncoderType != model.EncoderTypeNVEncC {
		t.Fatalf("expected encoder_type nvencc")
	}
	if m.EncoderOptions == nil {
		t.Fatalf("expected encoder_options initialized")
	}
}
