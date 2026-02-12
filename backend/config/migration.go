package config

import "fmt"

// Migrate upgrades an AppConfig to the latest version.
func Migrate(cfg AppConfig) (AppConfig, error) {
	for cfg.Version < CurrentVersion {
		switch cfg.Version {
		case 0:
			cfg = migrateV0toV1(cfg)
		default:
			return cfg, fmt.Errorf("unknown config version %d", cfg.Version)
		}
	}
	return cfg, nil
}

func migrateV0toV1(cfg AppConfig) AppConfig {
	d := Default()
	if cfg.MaxConcurrentJobs == 0 {
		cfg.MaxConcurrentJobs = d.MaxConcurrentJobs
	}
	if cfg.OnError == "" {
		cfg.OnError = d.OnError
	}
	if cfg.NoOutputTimeoutSec == 0 {
		cfg.NoOutputTimeoutSec = d.NoOutputTimeoutSec
	}
	if cfg.NoProgressTimeoutSec == 0 {
		cfg.NoProgressTimeoutSec = d.NoProgressTimeoutSec
	}
	if cfg.PostCompleteAction == "" {
		cfg.PostCompleteAction = d.PostCompleteAction
	}
	if cfg.OutputFolderMode == "" {
		cfg.OutputFolderMode = d.OutputFolderMode
	}
	if cfg.OutputNameTemplate == "" {
		cfg.OutputNameTemplate = d.OutputNameTemplate
	}
	if cfg.OutputContainer == "" {
		cfg.OutputContainer = d.OutputContainer
	}
	if cfg.OverwriteMode == "" {
		cfg.OverwriteMode = d.OverwriteMode
	}
	if cfg.Language == "" {
		cfg.Language = d.Language
	}
	cfg.Version = 1
	return cfg
}
