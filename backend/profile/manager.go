package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Manager handles profile CRUD and persistence.
type Manager struct {
	mu       sync.RWMutex
	profiles []Profile
	filePath string
}

// NewManager creates a Manager for the given profiles.json path.
func NewManager(filePath string) *Manager {
	return &Manager{filePath: filePath}
}

// Load reads profiles from disk. If file doesn't exist, generates presets.
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			m.profiles = GeneratePresets()
			return m.saveLocked()
		}
		return fmt.Errorf("read profiles: %w", err)
	}

	var pf ProfilesFile
	needsResave := false
	if err := json.Unmarshal(data, &pf); err != nil {
		// Try flat array fallback: profiles.json may be [{ ... }] instead of { "profiles": [{ ... }] }
		var flat []Profile
		if err2 := json.Unmarshal(data, &flat); err2 != nil {
			// Both formats failed: backup broken file and regenerate presets
			backupPath := m.filePath + fmt.Sprintf(".broken.%d", time.Now().Unix())
			os.Rename(m.filePath, backupPath)
			m.profiles = GeneratePresets()
			return m.saveLocked()
		}
		pf.Profiles = flat
		needsResave = true
	}

	migrated := make([]Profile, 0, len(pf.Profiles))
	for _, p := range pf.Profiles {
		mp, err := Migrate(p)
		if err != nil {
			return fmt.Errorf("migrate profile %q: %w", p.Name, err)
		}
		migrated = append(migrated, mp)
	}

	m.profiles = migrated

	// Re-save in correct format after flat array migration
	if needsResave {
		return m.saveLocked()
	}
	return nil
}

// List returns all profiles.
func (m *Manager) List() []Profile {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Profile, len(m.profiles))
	copy(out, m.profiles)
	return out
}

// Get returns a profile by ID.
func (m *Manager) Get(id string) (Profile, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.profiles {
		if p.ID == id {
			return p, true
		}
	}
	return Profile{}, false
}

// Upsert creates or updates a profile.
func (m *Manager) Upsert(p Profile) error {
	if err := Validate(p); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	found := false
	for i, existing := range m.profiles {
		if existing.ID == p.ID {
			if existing.IsPreset {
				return fmt.Errorf("cannot edit preset profile %q", existing.Name)
			}
			m.profiles[i] = p
			found = true
			break
		}
	}
	if !found {
		if p.ID == "" {
			p.ID = uuid.New().String()
		}
		p.Version = CurrentVersion
		m.profiles = append(m.profiles, p)
	}
	return m.saveLocked()
}

// Delete removes a profile by ID. Presets cannot be deleted.
func (m *Manager) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.profiles {
		if p.ID == id {
			if p.IsPreset {
				return fmt.Errorf("cannot delete preset profile %q", p.Name)
			}
			m.profiles = append(m.profiles[:i], m.profiles[i+1:]...)
			return m.saveLocked()
		}
	}
	return fmt.Errorf("profile not found: %s", id)
}

// Duplicate clones a profile with a new name. is_preset is set to false.
func (m *Manager) Duplicate(id string, newName string) (Profile, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	newName = strings.TrimSpace(newName)
	if newName == "" || len(newName) > 80 {
		return Profile{}, fmt.Errorf("name must be 1..80 characters")
	}

	for _, p := range m.profiles {
		if p.ID == id {
			dup := p
			dup.ID = uuid.New().String()
			dup.Name = newName
			dup.IsPreset = false
			m.profiles = append(m.profiles, dup)
			if err := m.saveLocked(); err != nil {
				return Profile{}, err
			}
			return dup, nil
		}
	}
	return Profile{}, fmt.Errorf("profile not found: %s", id)
}

// SetDefault sets the default profile ID (stored in AppConfig, not here).
// This is a no-op on the profile manager; the caller updates AppConfig.
func (m *Manager) SetDefault(id string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.profiles {
		if p.ID == id {
			return nil
		}
	}
	return fmt.Errorf("profile not found: %s", id)
}

// saveLocked writes profiles to disk atomically. Caller must hold mu.
func (m *Manager) saveLocked() error {
	pf := ProfilesFile{Profiles: m.profiles}
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal profiles: %w", err)
	}

	dir := filepath.Dir(m.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	tmpPath := m.filePath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("sync temp file: %w", err)
	}
	f.Close()

	if err := os.Rename(tmpPath, m.filePath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

// Validate checks a profile against schema constraints (design doc 5.2.1).
func Validate(p Profile) error {
	name := strings.TrimSpace(p.Name)
	if len(name) < 1 || len(name) > 80 {
		return fmt.Errorf("E_VALIDATION: name must be 1..80 characters")
	}

	validEncoders := map[string]bool{"nvencc": true, "qsvenc": true, "ffmpeg": true}
	if !validEncoders[p.EncoderType] {
		return fmt.Errorf("E_VALIDATION: invalid encoder_type %q", p.EncoderType)
	}

	if p.RateValue <= 0 {
		return fmt.Errorf("E_VALIDATION: rate_value must be > 0")
	}

	if p.OutputDepth != 8 && p.OutputDepth != 10 {
		return fmt.Errorf("E_VALIDATION: output_depth must be 8 or 10")
	}

	if p.Bframes != nil && (*p.Bframes < 0 || *p.Bframes > 7) {
		return fmt.Errorf("E_VALIDATION: bframes must be 0..7")
	}

	if p.Lookahead != nil && (*p.Lookahead < 0 || *p.Lookahead > 32) {
		return fmt.Errorf("E_VALIDATION: lookahead must be 0..32")
	}

	if p.AudioBitrate < 32 || p.AudioBitrate > 1024 {
		return fmt.Errorf("E_VALIDATION: audio_bitrate must be 32..1024")
	}

	if len(p.CustomOptions) > 4096 {
		return fmt.Errorf("E_VALIDATION: custom_options max 4096 characters")
	}

	adv := p.NVEncCAdvanced
	if adv.MaxBitrate != nil && *adv.MaxBitrate <= 0 {
		return fmt.Errorf("E_VALIDATION: nvencc_advanced.max_bitrate must be > 0")
	}
	if adv.VBRQuality != nil && *adv.VBRQuality <= 0 {
		return fmt.Errorf("E_VALIDATION: nvencc_advanced.vbr_quality must be > 0")
	}
	if adv.LookaheadLevel != nil && *adv.LookaheadLevel < 0 {
		return fmt.Errorf("E_VALIDATION: nvencc_advanced.lookahead_level must be >= 0")
	}
	if adv.RefsForward != nil && *adv.RefsForward < 0 {
		return fmt.Errorf("E_VALIDATION: nvencc_advanced.refs_forward must be >= 0")
	}
	if adv.RefsBackward != nil && *adv.RefsBackward < 0 {
		return fmt.Errorf("E_VALIDATION: nvencc_advanced.refs_backward must be >= 0")
	}
	if adv.OutputThread != nil && (*adv.OutputThread < 1 || *adv.OutputThread > 64) {
		return fmt.Errorf("E_VALIDATION: nvencc_advanced.output_thread must be 1..64")
	}

	return nil
}
