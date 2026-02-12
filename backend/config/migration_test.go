package config

import (
	"testing"

	"github.com/motoacs/enque/backend/model"
)

func TestMigrateConfigFromZeroVersion(t *testing.T) {
	cfg := model.AppConfig{}
	m, changed := Migrate(cfg)
	if !changed {
		t.Fatalf("expected changed")
	}
	if m.Version != model.AppConfigVersion {
		t.Fatalf("expected latest version, got %d", m.Version)
	}
	if m.OutputNameTemplate == "" {
		t.Fatalf("output_name_template should be set")
	}
}
