package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yuta/enque/backend/config"
	"github.com/yuta/enque/backend/detector"
	"github.com/yuta/enque/backend/encoder"
	"github.com/yuta/enque/backend/encoder/nvencc"
	"github.com/yuta/enque/backend/events"
	"github.com/yuta/enque/backend/logging"
	"github.com/yuta/enque/backend/profile"
	"github.com/yuta/enque/backend/queue"
)

// App is the main application struct exposed to Wails.
type App struct {
	ctx        context.Context
	configMgr  *config.Manager
	profileMgr *profile.Manager
	registry   *encoder.Registry
	queueMgr   *queue.Manager
	logger     *logging.AppLogger
}

// New creates a new App instance.
func New() *App {
	reg := encoder.NewRegistry()
	reg.Register(&nvencc.NVEncCAdapter{})

	return &App{
		configMgr:  config.NewManager(config.ConfigPath()),
		profileMgr: profile.NewManager(config.ProfilesPath()),
		registry:   reg,
	}
}

// Startup is called when the app starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.configMgr.Load()
	a.profileMgr.Load()

	logger, err := logging.NewAppLogger(config.LogsDir())
	if err != nil {
		fmt.Printf("warning: failed to init app logger: %v\n", err)
	}
	a.logger = logger

	emitter := events.NewEmitter(ctx)
	a.queueMgr = queue.NewManager(a.registry, emitter, a.logger)
}

// Shutdown is called when the app is closing.
func (a *App) Shutdown(ctx context.Context) {
	if a.logger != nil {
		a.logger.Close()
	}
}

// --- Bootstrap ---

// BootstrapResult holds the initial data sent to the frontend on startup.
type BootstrapResult struct {
	Config        config.AppConfig         `json:"config"`
	Profiles      []profile.Profile        `json:"profiles"`
	Tools         detector.DetectionResult `json:"tools"`
	TempArtifacts []string                 `json:"temp_artifacts"`
}

// Bootstrap returns initial config, profiles, and tool detection results.
func (a *App) Bootstrap() (*BootstrapResult, error) {
	cfg := a.configMgr.Get()
	profiles := a.profileMgr.List()
	tools := detector.DetectAll(cfg)
	tempArtifacts := a.queueMgr.ListTempArtifacts()

	return &BootstrapResult{
		Config:        cfg,
		Profiles:      profiles,
		Tools:         tools,
		TempArtifacts: tempArtifacts,
	}, nil
}

// --- AppConfig ---

// SaveAppConfig persists the application configuration.
func (a *App) SaveAppConfig(cfgJSON string) error {
	var cfg config.AppConfig
	if err := json.Unmarshal([]byte(cfgJSON), &cfg); err != nil {
		return fmt.Errorf("%s: %w", encoder.ErrValidation, err)
	}
	return a.configMgr.Save(cfg)
}

// --- Profile CRUD ---

// ListProfiles returns all saved profiles.
func (a *App) ListProfiles() ([]profile.Profile, error) {
	return a.profileMgr.List(), nil
}

// UpsertProfile creates or updates a profile.
func (a *App) UpsertProfile(profileJSON string) error {
	var p profile.Profile
	if err := json.Unmarshal([]byte(profileJSON), &p); err != nil {
		return fmt.Errorf("%s: %w", encoder.ErrValidation, err)
	}
	return a.profileMgr.Upsert(p)
}

// DeleteProfile removes a profile by ID.
func (a *App) DeleteProfile(profileID string) error {
	return a.profileMgr.Delete(profileID)
}

// DuplicateProfile duplicates a profile with a new name.
func (a *App) DuplicateProfile(profileID string, newName string) (*profile.Profile, error) {
	dup, err := a.profileMgr.Duplicate(profileID, newName)
	if err != nil {
		return nil, err
	}
	return &dup, nil
}

// SetDefaultProfile sets the default profile ID.
func (a *App) SetDefaultProfile(profileID string) error {
	if err := a.profileMgr.SetDefault(profileID); err != nil {
		return err
	}
	cfg := a.configMgr.Get()
	cfg.DefaultProfileID = profileID
	return a.configMgr.Save(cfg)
}

// --- GPU / Tool Detection ---

// GetGPUInfo returns GPU information from NVEncC --check-device.
func (a *App) GetGPUInfo() (string, error) {
	cfg := a.configMgr.Get()
	tools := detector.DetectAll(cfg)
	if !tools.NVEncC.Found {
		return "", fmt.Errorf("%s: NVEncC not found", encoder.ErrToolNotFound)
	}
	return detector.GetGPUInfo(tools.NVEncC.Path)
}

// DetectExternalTools detects NVEncC, QSVEncC, ffmpeg, ffprobe.
func (a *App) DetectExternalTools() (*detector.DetectionResult, error) {
	cfg := a.configMgr.Get()
	result := detector.DetectAll(cfg)
	return &result, nil
}

// --- Encode Control ---

// StartEncode begins an encoding session.
func (a *App) StartEncode(requestJSON string) error {
	var req queue.EncodeRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return fmt.Errorf("%s: %w", encoder.ErrValidation, err)
	}
	return a.queueMgr.StartEncode(req)
}

// RequestGracefulStop stops the session gracefully.
func (a *App) RequestGracefulStop(sessionID string) error {
	return a.queueMgr.RequestGracefulStop(sessionID)
}

// RequestAbort aborts the session.
func (a *App) RequestAbort(sessionID string) error {
	return a.queueMgr.RequestAbort(sessionID)
}

// CancelJob cancels a single running job.
func (a *App) CancelJob(sessionID string, jobID string) error {
	return a.queueMgr.CancelJob(sessionID, jobID)
}

// ResolveOverwrite responds to an overwrite confirmation prompt.
func (a *App) ResolveOverwrite(sessionID string, jobID string, decision string) error {
	return a.queueMgr.ResolveOverwrite(sessionID, jobID, decision)
}

// --- Temp Cleanup ---

// ListTempArtifacts returns leftover temp files from previous sessions.
func (a *App) ListTempArtifacts() ([]string, error) {
	return a.queueMgr.ListTempArtifacts(), nil
}

// CleanupTempArtifacts deletes specified temp files.
func (a *App) CleanupTempArtifacts(paths []string) error {
	return a.queueMgr.CleanupTempArtifacts(paths)
}

// --- Command Preview ---

// GetCommandPreview returns the command line preview for the given profile.
func (a *App) GetCommandPreview(profileJSON string, inputPath string, outputPath string) (string, error) {
	var p profile.Profile
	if err := json.Unmarshal([]byte(profileJSON), &p); err != nil {
		return "", fmt.Errorf("%s: %w", encoder.ErrValidation, err)
	}

	adapter, err := a.registry.Resolve(p.EncoderType)
	if err != nil {
		return "", err
	}

	args, err := adapter.BuildArgs(p, inputPath, outputPath)
	if err != nil {
		return "", err
	}

	return strings.Join(args, " "), nil
}
