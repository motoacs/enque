package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var frontendAssets embed.FS

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	err = wails.Run(&options.App{
		Title:     "Enque",
		Width:     1280,
		Height:    900,
		MinWidth:  1024,
		MinHeight: 720,
		AssetServer: &assetserver.Options{
			Assets: resolveAssetFS(),
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatalf("wails run failed: %v", err)
	}
}

func resolveAssetFS() fs.FS {
	distFS, err := fs.Sub(frontendAssets, "frontend/dist")
	if err == nil {
		return distFS
	}
	sourceFS, err := fs.Sub(frontendAssets, "frontend")
	if err == nil {
		return sourceFS
	}
	return frontendAssets
}
