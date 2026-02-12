package detector

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ToolInfo holds detection results for an external tool.
type ToolInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Version   string `json:"version"`
	Found     bool   `json:"found"`
	Error     string `json:"error,omitempty"`
	Supported bool   `json:"supported"`
}

// DetectionResult holds all tool detection results.
type DetectionResult struct {
	NVEncC  ToolInfo `json:"nvencc"`
	QSVEncC ToolInfo `json:"qsvenc"`
	FFmpeg  ToolInfo `json:"ffmpeg"`
	FFprobe ToolInfo `json:"ffprobe"`
}

// findExecutable searches for a tool in: 1) configPath, 2) same dir as app exe, 3) PATH.
func findExecutable(configPath string, candidates []string) string {
	// 1. Config-specified path
	if configPath != "" {
		if isExecutable(configPath) {
			return configPath
		}
	}

	// 2. Same directory as the running executable
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		for _, name := range candidates {
			p := filepath.Join(exeDir, name)
			if isExecutable(p) {
				return p
			}
		}
	}

	// 3. PATH
	for _, name := range candidates {
		p, err := exec.LookPath(name)
		if err == nil {
			return p
		}
	}

	return ""
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return info.Mode()&0o111 != 0
}
