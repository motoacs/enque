package config

import "github.com/motoacs/enque/backend/model"

func Migrate(cfg model.AppConfig) (model.AppConfig, bool) {
	changed := false
	if cfg.Version == 0 {
		defaults := model.DefaultAppConfig()
		cfg.Version = defaults.Version
		if cfg.OutputNameTemplate == "" {
			cfg.OutputNameTemplate = defaults.OutputNameTemplate
		}
		if cfg.OutputContainer == "" {
			cfg.OutputContainer = defaults.OutputContainer
		}
		if cfg.OutputFolderMode == "" {
			cfg.OutputFolderMode = defaults.OutputFolderMode
		}
		if cfg.OnError == "" {
			cfg.OnError = defaults.OnError
		}
		if cfg.OverwriteMode == "" {
			cfg.OverwriteMode = defaults.OverwriteMode
		}
		if cfg.Language == "" {
			cfg.Language = defaults.Language
		}
		if cfg.MaxConcurrentJobs == 0 {
			cfg.MaxConcurrentJobs = defaults.MaxConcurrentJobs
		}
		if cfg.NoOutputTimeoutSec == 0 {
			cfg.NoOutputTimeoutSec = defaults.NoOutputTimeoutSec
		}
		if cfg.NoProgressTimeoutSec == 0 {
			cfg.NoProgressTimeoutSec = defaults.NoProgressTimeoutSec
		}
		changed = true
	}
	if cfg.Version != model.AppConfigVersion {
		cfg.Version = model.AppConfigVersion
		changed = true
	}
	return cfg, changed
}
