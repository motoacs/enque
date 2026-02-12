package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	appcore "github.com/motoacs/enque/backend/app"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const appDataDirName = "Enque"

type App struct {
	*appcore.App
	emitter *wailsEmitter
}

func NewApp() (*App, error) {
	baseDir, err := resolveBaseDir()
	if err != nil {
		return nil, err
	}
	emitter := &wailsEmitter{}
	core, err := appcore.New(baseDir, resolveAppDir(), emitter)
	if err != nil {
		return nil, err
	}
	return &App{
		App:     core,
		emitter: emitter,
	}, nil
}

func (a *App) startup(ctx context.Context) {
	a.emitter.SetContext(ctx)
}

func resolveBaseDir() (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, appDataDirName), nil
}

func resolveAppDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

type wailsEmitter struct {
	mu  sync.RWMutex
	ctx context.Context
}

func (e *wailsEmitter) SetContext(ctx context.Context) {
	e.mu.Lock()
	e.ctx = ctx
	e.mu.Unlock()
}

func (e *wailsEmitter) Emit(name string, payload any) {
	e.mu.RLock()
	ctx := e.ctx
	e.mu.RUnlock()
	if ctx == nil {
		return
	}
	runtime.EventsEmit(ctx, name, payload)
}
