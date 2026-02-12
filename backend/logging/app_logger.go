package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

func NewAppLogger(baseDir string) (*slog.Logger, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}
	rotate := &lumberjack.Logger{
		Filename:   filepath.Join(baseDir, "app.log"),
		MaxSize:    20,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   false,
	}
	mw := io.MultiWriter(os.Stdout, rotate)
	return slog.New(slog.NewJSONHandler(mw, nil)), nil
}
