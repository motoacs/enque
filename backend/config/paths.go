package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// DataDir returns the application data directory.
// Windows: %APPDATA%/Enque/
// macOS/Linux: ~/.config/Enque/ (development fallback)
func DataDir() string {
	if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		if appdata != "" {
			return filepath.Join(appdata, "Enque")
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".enque")
	}
	return filepath.Join(home, ".config", "Enque")
}

// ConfigPath returns the path to config.json.
func ConfigPath() string {
	return filepath.Join(DataDir(), "config.json")
}

// ProfilesPath returns the path to profiles.json.
func ProfilesPath() string {
	return filepath.Join(DataDir(), "profiles.json")
}

// LogsDir returns the path to the logs directory.
func LogsDir() string {
	return filepath.Join(DataDir(), "logs")
}

// RuntimeDir returns the path to the runtime directory.
func RuntimeDir() string {
	return filepath.Join(DataDir(), "runtime")
}

// TempIndexPath returns the path to temp_index.json.
func TempIndexPath() string {
	return filepath.Join(RuntimeDir(), "temp_index.json")
}
