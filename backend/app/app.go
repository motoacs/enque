package app

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/motoacs/enque/backend/config"
	"github.com/motoacs/enque/backend/detector"
	"github.com/motoacs/enque/backend/encoder"
	"github.com/motoacs/enque/backend/encoder/ffmpeg"
	"github.com/motoacs/enque/backend/encoder/nvencc"
	"github.com/motoacs/enque/backend/encoder/qsvenc"
	"github.com/motoacs/enque/backend/logging"
	"github.com/motoacs/enque/backend/model"
	"github.com/motoacs/enque/backend/profile"
	"github.com/motoacs/enque/backend/queue"
)

type App struct {
	baseDir    string
	configMgr  *config.Manager
	profileMgr *profile.Manager
	detector   *detector.Detector
	registry   *encoder.Registry
	queueMgr   *queue.Manager
}

func New(baseDir, appDir string, emitter queue.EventEmitter) (*App, error) {
	cfgMgr := config.NewManager(baseDir)
	profileMgr := profile.NewManager(baseDir)
	dt := detector.New(appDir)
	jobLogger := logging.NewJobLogger(baseDir)
	appLogger, err := logging.NewAppLogger(baseDir)
	if err != nil {
		return nil, err
	}
	registry := encoder.NewRegistry(
		nvencc.NewAdapter(),
		qsvenc.NewAdapter(),
		ffmpeg.NewAdapter(),
	)
	queueMgr := queue.NewManager(
		registry,
		encoder.NewOSProcessRunner(),
		queue.NewOutputResolver(),
		jobLogger,
		emitter,
		appLogger,
		filepath.Join(baseDir, "runtime"),
	)

	return &App{
		baseDir:    baseDir,
		configMgr:  cfgMgr,
		profileMgr: profileMgr,
		detector:   dt,
		registry:   registry,
		queueMgr:   queueMgr,
	}, nil
}

func (a *App) Bootstrap() (model.BootstrapResponse, error) {
	cfg, err := a.configMgr.Load()
	if err != nil {
		return model.BootstrapResponse{}, wrapIO(err, "load config")
	}
	profiles, err := a.profileMgr.List()
	if err != nil {
		return model.BootstrapResponse{}, wrapIO(err, "load profiles")
	}
	tools := a.detector.DetectExternalTools(cfg)
	warnings := []string{}
	if tools.NVEncC.Warning != "" {
		warnings = append(warnings, tools.NVEncC.Warning)
	}
	gpuInfo := model.GPUInfo{}
	if tools.NVEncC.Found && tools.NVEncC.Warning == "" {
		if info, err := a.detector.GetGPUInfo(tools.NVEncC.Path); err == nil {
			gpuInfo = info
		}
	}
	return model.BootstrapResponse{Config: cfg, Profiles: profiles, Tools: tools, GPUInfo: gpuInfo, Warnings: warnings}, nil
}

func (a *App) SaveAppConfig(cfg model.AppConfig) error {
	if err := a.configMgr.Save(cfg); err != nil {
		return wrapIO(err, "save config")
	}
	return nil
}

func (a *App) ListProfiles() ([]model.Profile, error) {
	return a.profileMgr.List()
}

func (a *App) UpsertProfile(p model.Profile) (model.Profile, error) {
	return a.profileMgr.Upsert(p)
}

func (a *App) DeleteProfile(profileID string) error {
	return a.profileMgr.Delete(profileID)
}

func (a *App) DuplicateProfile(profileID, newName string) (model.Profile, error) {
	return a.profileMgr.Duplicate(profileID, newName)
}

func (a *App) SetDefaultProfile(profileID string) error {
	cfg, err := a.configMgr.Load()
	if err != nil {
		return err
	}
	cfg.DefaultProfileID = profileID
	return a.configMgr.Save(cfg)
}

func (a *App) GetGPUInfo() (model.GPUInfo, error) {
	cfg, err := a.configMgr.Load()
	if err != nil {
		return model.GPUInfo{}, err
	}
	tools := a.detector.DetectExternalTools(cfg)
	if err := detector.EnsureNVEncReady(tools); err != nil {
		return model.GPUInfo{}, err
	}
	return a.detector.GetGPUInfo(tools.NVEncC.Path)
}

func (a *App) DetectExternalTools() (model.ToolSnapshot, error) {
	cfg, err := a.configMgr.Load()
	if err != nil {
		return model.ToolSnapshot{}, err
	}
	return a.detector.DetectExternalTools(cfg), nil
}

func (a *App) PreviewCommand(req model.PreviewCommandRequest) (model.PreviewCommandResponse, error) {
	adapter, err := a.registry.Resolve(req.Profile.EncoderType)
	if err != nil {
		return model.PreviewCommandResponse{}, err
	}
	build, err := adapter.BuildArgs(encoder.BuildRequest{
		Profile:    req.Profile,
		AppConfig:  req.AppConfigSnapshot,
		InputPath:  req.InputPath,
		OutputPath: req.OutputPath,
	})
	if err != nil {
		return model.PreviewCommandResponse{}, err
	}
	return model.PreviewCommandResponse{Argv: build.Argv, DisplayCommand: build.DisplayCommand}, nil
}

func (a *App) StartEncode(req model.StartEncodeRequest) (model.EncodeSession, error) {
	cfg := req.AppConfigSnapshot
	if cfg.Version == 0 {
		cfg = model.DefaultAppConfig()
	}
	tools := a.detector.DetectExternalTools(cfg)
	if req.Profile.EncoderType == model.EncoderTypeNVEncC {
		if err := detector.EnsureNVEncReady(tools); err != nil {
			return model.EncodeSession{}, err
		}
	}
	encoderPath, err := resolveEncoderPath(req.Profile.EncoderType, tools)
	if err != nil {
		return model.EncodeSession{}, err
	}
	return a.queueMgr.StartEncode(context.Background(), req, encoderPath)
}

func (a *App) RequestGracefulStop(sessionID string) error {
	return a.queueMgr.RequestGracefulStop(sessionID)
}

func (a *App) RequestAbort(sessionID string) error {
	return a.queueMgr.RequestAbort(sessionID)
}

func (a *App) CancelJob(sessionID, jobID string) error {
	return a.queueMgr.CancelJob(sessionID, jobID)
}

func (a *App) ResolveOverwrite(sessionID, jobID string, decision model.ResolveOverwriteDecision) error {
	return a.queueMgr.ResolveOverwrite(sessionID, jobID, decision)
}

func (a *App) ListTempArtifacts() ([]string, error) {
	return a.queueMgr.ListTempArtifacts()
}

func (a *App) CleanupTempArtifacts(paths []string) error {
	return a.queueMgr.CleanupTempArtifacts(paths)
}

func resolveEncoderPath(t model.EncoderType, tools model.ToolSnapshot) (string, error) {
	switch t {
	case model.EncoderTypeNVEncC:
		if tools.NVEncC.Path == "" {
			return "", &model.EnqueError{Code: model.ErrToolNotFound, Message: "NVEncC not found"}
		}
		if tools.NVEncC.Warning != "" {
			return "", &model.EnqueError{Code: model.ErrToolVersionUnsupported, Message: tools.NVEncC.Warning}
		}
		return tools.NVEncC.Path, nil
	case model.EncoderTypeQSVEnc:
		if tools.QSVEnc.Path == "" {
			return "", &model.EnqueError{Code: model.ErrEncoderNotImplemented, Message: "qsvenc adapter is not implemented in v1"}
		}
		return tools.QSVEnc.Path, nil
	case model.EncoderTypeFFmpeg:
		if tools.FFmpeg.Path == "" {
			return "", &model.EnqueError{Code: model.ErrEncoderNotImplemented, Message: "ffmpeg adapter is not implemented in v1"}
		}
		return tools.FFmpeg.Path, nil
	default:
		return "", &model.EnqueError{Code: model.ErrEncoderNotImplemented, Message: fmt.Sprintf("encoder %s is not implemented", t)}
	}
}

func wrapIO(err error, action string) error {
	var enqErr *model.EnqueError
	if errors.As(err, &enqErr) {
		return err
	}
	return &model.EnqueError{Code: model.ErrIO, Message: fmt.Sprintf("%s: %v", action, err)}
}
