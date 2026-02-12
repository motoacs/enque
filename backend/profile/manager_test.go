package profile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func tempProfilePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "profiles.json")
}

func TestLoad_CreatesPresetsWhenMissing(t *testing.T) {
	m := NewManager(tempProfilePath(t))
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}

	profiles := m.List()
	if len(profiles) != 4 {
		t.Fatalf("expected 4 presets, got %d", len(profiles))
	}

	names := map[string]bool{}
	for _, p := range profiles {
		names[p.Name] = true
		if !p.IsPreset {
			t.Errorf("profile %q should be preset", p.Name)
		}
		if p.Version != CurrentVersion {
			t.Errorf("profile %q version=%d, want %d", p.Name, p.Version, CurrentVersion)
		}
	}

	for _, expected := range []string{"HEVC Quality", "AV1 Fast", "Camera Archive", "H.264 Compatible"} {
		if !names[expected] {
			t.Errorf("missing preset %q", expected)
		}
	}
}

func TestLoad_ReadsExistingFile(t *testing.T) {
	path := tempProfilePath(t)
	m := NewManager(path)
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}

	m2 := NewManager(path)
	if err := m2.Load(); err != nil {
		t.Fatal(err)
	}

	if len(m2.List()) != 4 {
		t.Fatalf("expected 4 profiles after re-load, got %d", len(m2.List()))
	}
}

func TestUpsert_CreateAndUpdate(t *testing.T) {
	m := NewManager(tempProfilePath(t))
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}

	p := Profile{
		Name:         "Test Profile",
		EncoderType:  "nvencc",
		EncoderOpts:  map[string]any{},
		Codec:        "hevc",
		RateControl:  "qvbr",
		RateValue:    28,
		Preset:       "P4",
		OutputDepth:  10,
		AudioBitrate: 256,
		AudioMode:    "copy",
	}

	if err := m.Upsert(p); err != nil {
		t.Fatal(err)
	}

	profiles := m.List()
	if len(profiles) != 5 {
		t.Fatalf("expected 5 profiles, got %d", len(profiles))
	}

	created := profiles[4]
	if created.Name != "Test Profile" {
		t.Errorf("name=%q, want %q", created.Name, "Test Profile")
	}
	if created.ID == "" {
		t.Error("expected auto-generated ID")
	}

	// Update
	created.RateValue = 30
	if err := m.Upsert(created); err != nil {
		t.Fatal(err)
	}

	updated, ok := m.Get(created.ID)
	if !ok {
		t.Fatal("profile not found after update")
	}
	if updated.RateValue != 30 {
		t.Errorf("rate_value=%f, want 30", updated.RateValue)
	}
}

func TestUpsert_CannotEditPreset(t *testing.T) {
	m := NewManager(tempProfilePath(t))
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}

	preset := m.List()[0]
	preset.RateValue = 99
	if err := m.Upsert(preset); err == nil {
		t.Error("expected error editing preset")
	}
}

func TestDelete(t *testing.T) {
	m := NewManager(tempProfilePath(t))
	if err := m.Load(); err != nil {
		t.Fatal(err)
	}

	p := Profile{
		Name:         "ToDelete",
		EncoderType:  "nvencc",
		EncoderOpts:  map[string]any{},
		RateValue:    28,
		OutputDepth:  10,
		AudioBitrate: 256,
	}
	m.Upsert(p)
	profiles := m.List()
	newID := profiles[len(profiles)-1].ID

	if err := m.Delete(newID); err != nil {
		t.Fatal(err)
	}
	if len(m.List()) != 4 {
		t.Errorf("expected 4 profiles after delete, got %d", len(m.List()))
	}
}

func TestDelete_CannotDeletePreset(t *testing.T) {
	m := NewManager(tempProfilePath(t))
	m.Load()

	preset := m.List()[0]
	if err := m.Delete(preset.ID); err == nil {
		t.Error("expected error deleting preset")
	}
}

func TestDuplicate(t *testing.T) {
	m := NewManager(tempProfilePath(t))
	m.Load()

	preset := m.List()[0]
	dup, err := m.Duplicate(preset.ID, "My Copy")
	if err != nil {
		t.Fatal(err)
	}

	if dup.Name != "My Copy" {
		t.Errorf("name=%q, want %q", dup.Name, "My Copy")
	}
	if dup.IsPreset {
		t.Error("duplicate should not be preset")
	}
	if dup.ID == preset.ID {
		t.Error("duplicate should have different ID")
	}
	if len(m.List()) != 5 {
		t.Errorf("expected 5 profiles, got %d", len(m.List()))
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*Profile)
		wantErr bool
	}{
		{"valid", func(p *Profile) {}, false},
		{"empty name", func(p *Profile) { p.Name = "" }, true},
		{"long name", func(p *Profile) { p.Name = string(make([]byte, 81)) }, true},
		{"bad encoder", func(p *Profile) { p.EncoderType = "invalid" }, true},
		{"rate_value zero", func(p *Profile) { p.RateValue = 0 }, true},
		{"bad output_depth", func(p *Profile) { p.OutputDepth = 12 }, true},
		{"bframes too high", func(p *Profile) { v := 8; p.Bframes = &v }, true},
		{"lookahead too high", func(p *Profile) { v := 33; p.Lookahead = &v }, true},
		{"audio_bitrate low", func(p *Profile) { p.AudioBitrate = 10 }, true},
		{"custom_options too long", func(p *Profile) { p.CustomOptions = string(make([]byte, 4097)) }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Profile{
				Name:         "Test",
				EncoderType:  "nvencc",
				RateValue:    28,
				OutputDepth:  10,
				AudioBitrate: 256,
			}
			tt.modify(&p)
			err := Validate(p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error=%v, wantErr=%v", err, tt.wantErr)
			}
		})
	}
}

func TestMigration_V1toV4(t *testing.T) {
	p := Profile{
		Version:     1,
		Name:        "Old Profile",
		RateValue:   28,
		OutputDepth: 10,
	}

	result, err := Migrate(p)
	if err != nil {
		t.Fatal(err)
	}

	if result.Version != CurrentVersion {
		t.Errorf("version=%d, want %d", result.Version, CurrentVersion)
	}
	if result.EncoderType != "nvencc" {
		t.Errorf("encoder_type=%q, want %q", result.EncoderType, "nvencc")
	}
	if result.EncoderOpts == nil {
		t.Error("encoder_options should not be nil")
	}
}

func TestAtomicWrite(t *testing.T) {
	path := tempProfilePath(t)
	m := NewManager(path)
	m.Load()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var pf ProfilesFile
	if err := json.Unmarshal(data, &pf); err != nil {
		t.Fatal(err)
	}
	if len(pf.Profiles) != 4 {
		t.Errorf("expected 4 profiles in file, got %d", len(pf.Profiles))
	}
}

func TestPresetValues(t *testing.T) {
	presets := GeneratePresets()

	checks := map[string]struct {
		codec       string
		rateValue   float64
		preset      string
		outputDepth int
		audioMode   string
		splitEnc    string
		fileTime    bool
	}{
		"HEVC Quality":      {"hevc", 28, "P4", 10, "copy", "auto", false},
		"AV1 Fast":          {"av1", 32, "P1", 10, "copy", "auto", false},
		"Camera Archive":    {"hevc", 24, "P7", 10, "copy", "off", true},
		"H.264 Compatible":  {"h264", 26, "P4", 8, "aac", "off", false},
	}

	for _, p := range presets {
		c, ok := checks[p.Name]
		if !ok {
			continue
		}
		if p.Codec != c.codec {
			t.Errorf("%s: codec=%q, want %q", p.Name, p.Codec, c.codec)
		}
		if p.RateValue != c.rateValue {
			t.Errorf("%s: rate_value=%f, want %f", p.Name, p.RateValue, c.rateValue)
		}
		if p.Preset != c.preset {
			t.Errorf("%s: preset=%q, want %q", p.Name, p.Preset, c.preset)
		}
		if p.OutputDepth != c.outputDepth {
			t.Errorf("%s: output_depth=%d, want %d", p.Name, p.OutputDepth, c.outputDepth)
		}
		if p.AudioMode != c.audioMode {
			t.Errorf("%s: audio_mode=%q, want %q", p.Name, p.AudioMode, c.audioMode)
		}
		if p.SplitEnc != c.splitEnc {
			t.Errorf("%s: split_enc=%q, want %q", p.Name, p.SplitEnc, c.splitEnc)
		}
		if p.RestoreFileTime != c.fileTime {
			t.Errorf("%s: restore_file_time=%v, want %v", p.Name, p.RestoreFileTime, c.fileTime)
		}
	}
}
