package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/motoacs/enque/backend/model"
)

type Manager struct {
	baseDir string
	path    string
}

func NewManager(baseDir string) *Manager {
	return &Manager{baseDir: baseDir, path: filepath.Join(baseDir, "profiles.json")}
}

func (m *Manager) Path() string {
	return m.path
}

func (m *Manager) List() ([]model.Profile, error) {
	if err := os.MkdirAll(m.baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("mkdir profile dir: %w", err)
	}
	b, err := os.ReadFile(m.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			profiles := BuiltInPresets()
			if err := m.saveAll(profiles); err != nil {
				return nil, err
			}
			return profiles, nil
		}
		return nil, fmt.Errorf("read profiles: %w", err)
	}
	var profiles []model.Profile
	if err := json.Unmarshal(b, &profiles); err != nil {
		if backupErr := os.WriteFile(fmt.Sprintf("%s.broken.%d", m.path, time.Now().Unix()), b, 0o644); backupErr != nil {
			return nil, backupErr
		}
		profiles = BuiltInPresets()
		if err := m.saveAll(profiles); err != nil {
			return nil, err
		}
		return profiles, nil
	}
	changed := false
	for i := range profiles {
		migrated, c := migrateOne(profiles[i])
		if c {
			profiles[i] = migrated
			changed = true
		}
	}
	if changed {
		if err := m.saveAll(profiles); err != nil {
			return nil, err
		}
	}
	return profiles, nil
}

func (m *Manager) Upsert(p model.Profile) (model.Profile, error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	p.Version = model.ProfileVersion
	if p.EncoderOptions == nil {
		p.EncoderOptions = map[string]any{}
	}
	if errs := model.ValidateProfile(p); len(errs) > 0 {
		return model.Profile{}, &model.EnqueError{Code: model.ErrValidation, Message: "profile validation failed", Fields: errs}
	}
	profiles, err := m.List()
	if err != nil {
		return model.Profile{}, err
	}
	updated := false
	for i := range profiles {
		if profiles[i].ID == p.ID {
			if profiles[i].IsPreset {
				return model.Profile{}, model.NewError(model.ErrValidation, "preset profile is immutable")
			}
			profiles[i] = p
			updated = true
			break
		}
	}
	if !updated {
		profiles = append(profiles, p)
	}
	if err := m.saveAll(profiles); err != nil {
		return model.Profile{}, err
	}
	return p, nil
}

func (m *Manager) Delete(profileID string) error {
	profiles, err := m.List()
	if err != nil {
		return err
	}
	idx := slices.IndexFunc(profiles, func(p model.Profile) bool { return p.ID == profileID })
	if idx == -1 {
		return nil
	}
	if profiles[idx].IsPreset {
		return model.NewError(model.ErrValidation, "preset profile is immutable")
	}
	profiles = append(profiles[:idx], profiles[idx+1:]...)
	return m.saveAll(profiles)
}

func (m *Manager) Duplicate(profileID, newName string) (model.Profile, error) {
	profiles, err := m.List()
	if err != nil {
		return model.Profile{}, err
	}
	for _, p := range profiles {
		if p.ID == profileID {
			dup := p
			dup.ID = uuid.NewString()
			dup.IsPreset = false
			dup.Name = strings.TrimSpace(newName)
			if dup.Name == "" {
				dup.Name = p.Name + " Copy"
			}
			return m.Upsert(dup)
		}
	}
	return model.Profile{}, model.NewError(model.ErrValidation, "profile not found")
}

func (m *Manager) saveAll(profiles []model.Profile) error {
	b, err := json.MarshalIndent(profiles, "", "  ")
	if err != nil {
		return err
	}
	tmp := m.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, m.path)
}
