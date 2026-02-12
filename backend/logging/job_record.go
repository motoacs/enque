package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/motoacs/enque/backend/model"
)

type JobLogger struct {
	logDir string
	mu     sync.Mutex
}

func NewJobLogger(baseDir string) *JobLogger {
	return &JobLogger{logDir: filepath.Join(baseDir, "logs")}
}

func (l *JobLogger) Dir() string { return l.logDir }

func (l *JobLogger) WriteJobRecord(record model.JobRecord) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := os.MkdirAll(l.logDir, 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(l.logDir, fmt.Sprintf("%s.json", record.JobID))
	return os.WriteFile(path, b, 0o644)
}

func (l *JobLogger) OpenStderrLog(jobID string) (*os.File, string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := os.MkdirAll(l.logDir, 0o755); err != nil {
		return nil, "", err
	}
	path := filepath.Join(l.logDir, fmt.Sprintf("%s.stderr.log", jobID))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return nil, "", err
	}
	return f, path, nil
}
