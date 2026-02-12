package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// AppLogger provides daily-rotated application logging (30-day retention).
type AppLogger struct {
	dir     string
	current *os.File
	logger  *log.Logger
	day     string
}

// NewAppLogger creates an app logger writing to {dir}/app-{date}.log.
func NewAppLogger(dir string) (*AppLogger, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	al := &AppLogger{dir: dir}
	if err := al.rotate(); err != nil {
		return nil, err
	}
	return al, nil
}

// Info logs an informational message.
func (al *AppLogger) Info(msg string, args ...interface{}) {
	al.checkRotate()
	al.logger.Printf("[INFO] "+msg, args...)
}

// Warn logs a warning message.
func (al *AppLogger) Warn(msg string, args ...interface{}) {
	al.checkRotate()
	al.logger.Printf("[WARN] "+msg, args...)
}

// Error logs an error message.
func (al *AppLogger) Error(msg string, args ...interface{}) {
	al.checkRotate()
	al.logger.Printf("[ERROR] "+msg, args...)
}

// Close closes the current log file.
func (al *AppLogger) Close() error {
	if al.current != nil {
		return al.current.Close()
	}
	return nil
}

func (al *AppLogger) checkRotate() {
	today := time.Now().Format("2006-01-02")
	if today != al.day {
		al.rotate()
	}
}

func (al *AppLogger) rotate() error {
	if al.current != nil {
		al.current.Close()
	}

	today := time.Now().Format("2006-01-02")
	al.day = today

	path := filepath.Join(al.dir, "app-"+today+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}

	al.current = f
	al.logger = log.New(f, "", log.LstdFlags)

	al.cleanup()
	return nil
}

func (al *AppLogger) cleanup() {
	entries, err := os.ReadDir(al.dir)
	if err != nil {
		return
	}

	var logFiles []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "app-") && strings.HasSuffix(e.Name(), ".log") {
			logFiles = append(logFiles, e.Name())
		}
	}

	sort.Strings(logFiles)

	if len(logFiles) > 30 {
		for _, name := range logFiles[:len(logFiles)-30] {
			os.Remove(filepath.Join(al.dir, name))
		}
	}
}
