package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/motoacs/enque/backend/model"
)

type Manager struct {
	baseDir string
	path    string
}

func NewManager(baseDir string) *Manager {
	return &Manager{
		baseDir: baseDir,
		path:    filepath.Join(baseDir, "config.json"),
	}
}

func (m *Manager) Path() string {
	return m.path
}

func (m *Manager) Load() (model.AppConfig, error) {
	if err := os.MkdirAll(m.baseDir, 0o755); err != nil {
		return model.AppConfig{}, fmt.Errorf("mkdir config dir: %w", err)
	}
	b, err := os.ReadFile(m.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg := model.DefaultAppConfig()
			if saveErr := m.Save(cfg); saveErr != nil {
				return model.AppConfig{}, saveErr
			}
			return cfg, nil
		}
		return model.AppConfig{}, fmt.Errorf("read config: %w", err)
	}
	var cfg model.AppConfig
	if err := json.Unmarshal(b, &cfg); err != nil {
		if backupErr := backupBrokenFile(m.path, b); backupErr != nil {
			return model.AppConfig{}, fmt.Errorf("backup invalid config failed: %w", backupErr)
		}
		cfg = model.DefaultAppConfig()
		if saveErr := m.Save(cfg); saveErr != nil {
			return model.AppConfig{}, saveErr
		}
		return cfg, nil
	}
	migrated, changed := Migrate(cfg)
	if changed {
		if err := m.Save(migrated); err != nil {
			return model.AppConfig{}, err
		}
	}
	return migrated, nil
}

func (m *Manager) Save(cfg model.AppConfig) error {
	if errs := model.ValidateAppConfig(cfg); len(errs) > 0 {
		return &model.EnqueError{Code: model.ErrValidation, Message: "app config validation failed", Fields: errs}
	}
	cfg.Version = model.AppConfigVersion
	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	return atomicWriteFile(m.path, b, 0o644)
}

func atomicWriteFile(path string, b []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func backupBrokenFile(path string, b []byte) error {
	backupPath := fmt.Sprintf("%s.broken.%d", path, time.Now().Unix())
	return os.WriteFile(backupPath, b, 0o644)
}
