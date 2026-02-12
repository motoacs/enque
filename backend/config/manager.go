package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Manager handles AppConfig persistence.
type Manager struct {
	mu       sync.RWMutex
	config   AppConfig
	filePath string
}

// NewManager creates a Manager for the given config.json path.
func NewManager(filePath string) *Manager {
	return &Manager{filePath: filePath}
}

// Load reads config from disk. If file doesn't exist, generates defaults.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			m.config = Default()
			return m.saveLocked()
		}
		return fmt.Errorf("read config: %w", err)
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	cfg, err = Migrate(cfg)
	if err != nil {
		return fmt.Errorf("migrate config: %w", err)
	}

	m.config = cfg
	return nil
}

// Get returns the current config.
func (m *Manager) Get() AppConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// Save validates and persists the config.
func (m *Manager) Save(cfg AppConfig) error {
	if err := Validate(cfg); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config = cfg
	return m.saveLocked()
}

// saveLocked writes config atomically. Caller must hold mu.
func (m *Manager) saveLocked() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	tmpPath := m.filePath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write temp: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("sync temp: %w", err)
	}
	f.Close()

	if err := os.Rename(tmpPath, m.filePath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp: %w", err)
	}
	return nil
}

// Validate checks AppConfig constraints (design doc 5.3.1).
func Validate(cfg AppConfig) error {
	if cfg.MaxConcurrentJobs < 1 || cfg.MaxConcurrentJobs > 8 {
		return fmt.Errorf("E_VALIDATION: max_concurrent_jobs must be 1..8")
	}
	if cfg.NoOutputTimeoutSec < 30 || cfg.NoOutputTimeoutSec > 86400 {
		return fmt.Errorf("E_VALIDATION: no_output_timeout_sec must be 30..86400")
	}
	if cfg.NoProgressTimeoutSec < 30 || cfg.NoProgressTimeoutSec > 86400 {
		return fmt.Errorf("E_VALIDATION: no_progress_timeout_sec must be 30..86400")
	}
	tpl := strings.TrimSpace(cfg.OutputNameTemplate)
	if len(tpl) < 1 || len(tpl) > 255 {
		return fmt.Errorf("E_VALIDATION: output_name_template must be 1..255 characters")
	}
	if cfg.OutputFolderMode == "specified" && strings.TrimSpace(cfg.OutputFolderPath) == "" {
		return fmt.Errorf("E_VALIDATION: output_folder_path required when mode is 'specified'")
	}
	if cfg.PostCompleteAction == "custom" && strings.TrimSpace(cfg.PostCompleteCommand) == "" {
		return fmt.Errorf("E_VALIDATION: post_complete_command required when action is 'custom'")
	}
	return nil
}
