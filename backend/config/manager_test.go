package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func tempConfigPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "config.json")
}

func TestLoad_CreatesDefaultWhenMissing(t *testing.T) {
	m := NewManager(tempConfigPath(t))
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}

	cfg := m.Get()
	d := Default()

	if cfg.Version != d.Version {
		t.Errorf("version=%d, want %d", cfg.Version, d.Version)
	}
	if cfg.MaxConcurrentJobs != d.MaxConcurrentJobs {
		t.Errorf("max_concurrent_jobs=%d, want %d", cfg.MaxConcurrentJobs, d.MaxConcurrentJobs)
	}
	if cfg.OnError != d.OnError {
		t.Errorf("on_error=%q, want %q", cfg.OnError, d.OnError)
	}
	if cfg.Language != d.Language {
		t.Errorf("language=%q, want %q", cfg.Language, d.Language)
	}
}

func TestLoad_ReadsExistingFile(t *testing.T) {
	path := tempConfigPath(t)
	m := NewManager(path)
	m.Load()

	cfg := m.Get()
	cfg.MaxConcurrentJobs = 4
	m.Save(cfg)

	m2 := NewManager(path)
	m2.Load()
	if m2.Get().MaxConcurrentJobs != 4 {
		t.Errorf("max_concurrent_jobs=%d, want 4", m2.Get().MaxConcurrentJobs)
	}
}

func TestSave_Validation(t *testing.T) {
	m := NewManager(tempConfigPath(t))
	m.Load()

	tests := []struct {
		name    string
		modify  func(*AppConfig)
		wantErr bool
	}{
		{"valid", func(c *AppConfig) {}, false},
		{"max_jobs_0", func(c *AppConfig) { c.MaxConcurrentJobs = 0 }, true},
		{"max_jobs_9", func(c *AppConfig) { c.MaxConcurrentJobs = 9 }, true},
		{"output_timeout_low", func(c *AppConfig) { c.NoOutputTimeoutSec = 10 }, true},
		{"progress_timeout_high", func(c *AppConfig) { c.NoProgressTimeoutSec = 100000 }, true},
		{"empty_template", func(c *AppConfig) { c.OutputNameTemplate = "" }, true},
		{"specified_no_path", func(c *AppConfig) {
			c.OutputFolderMode = "specified"
			c.OutputFolderPath = ""
		}, true},
		{"custom_no_cmd", func(c *AppConfig) {
			c.PostCompleteAction = "custom"
			c.PostCompleteCommand = ""
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Default()
			tt.modify(&cfg)
			err := m.Save(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error=%v, wantErr=%v", err, tt.wantErr)
			}
		})
	}
}

func TestMigration_V0toV1(t *testing.T) {
	cfg := AppConfig{Version: 0}
	result, err := Migrate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if result.Version != CurrentVersion {
		t.Errorf("version=%d, want %d", result.Version, CurrentVersion)
	}
	if result.MaxConcurrentJobs != 1 {
		t.Errorf("max_concurrent_jobs=%d, want 1", result.MaxConcurrentJobs)
	}
	if result.OnError != "skip" {
		t.Errorf("on_error=%q, want skip", result.OnError)
	}
}

func TestAtomicWrite(t *testing.T) {
	path := tempConfigPath(t)
	m := NewManager(path)
	m.Load()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Version != CurrentVersion {
		t.Errorf("version=%d, want %d", cfg.Version, CurrentVersion)
	}
}
