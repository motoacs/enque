package logging

import (
	"fmt"
	"os"
	"path/filepath"
)

// StderrWriter writes encoder stderr to a per-job log file.
type StderrWriter struct {
	file *os.File
}

// NewStderrWriter creates a stderr log file at {logsDir}/{jobID}.stderr.log.
func NewStderrWriter(logsDir, jobID string) (*StderrWriter, error) {
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return nil, fmt.Errorf("create logs dir: %w", err)
	}

	path := filepath.Join(logsDir, jobID+".stderr.log")
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create stderr log: %w", err)
	}
	return &StderrWriter{file: f}, nil
}

// Write implements io.Writer.
func (w *StderrWriter) Write(p []byte) (int, error) {
	return w.file.Write(p)
}

// Close flushes and closes the file.
func (w *StderrWriter) Close() error {
	if w.file != nil {
		w.file.Sync()
		return w.file.Close()
	}
	return nil
}
