package detector

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/motoacs/enque/backend/model"
)

type Detector struct {
	appDir string
}

func New(appDir string) *Detector {
	return &Detector{appDir: appDir}
}

func (d *Detector) DetectExternalTools(cfg model.AppConfig) model.ToolSnapshot {
	return model.ToolSnapshot{
		NVEncC:  d.detect(cfg.NVEncCPath, []string{"NVEncC64.exe", "NVEncC.exe"}, true),
		QSVEnc:  d.detect(cfg.QSVEncPath, []string{"QSVEncC64.exe", "QSVEncC.exe"}, false),
		FFmpeg:  d.detect(cfg.FFmpegPath, []string{"ffmpeg.exe", "ffmpeg"}, false),
		FFprobe: d.detect(cfg.FFprobePath, []string{"ffprobe.exe", "ffprobe"}, false),
	}
}

func (d *Detector) detect(configuredPath string, candidates []string, checkNVEncMinVersion bool) model.ToolInfo {
	if configuredPath != "" {
		if fi, err := os.Stat(configuredPath); err == nil && !fi.IsDir() {
			version := getVersion(configuredPath)
			info := model.ToolInfo{Found: true, Path: configuredPath, Version: version}
			if checkNVEncMinVersion && !isNVEncVersionSupported(version) {
				info.Warning = "NVEncC 8.x 以降が必要です"
			}
			return info
		}
	}
	for _, candidate := range candidates {
		p := filepath.Join(d.appDir, candidate)
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			version := getVersion(p)
			info := model.ToolInfo{Found: true, Path: p, Version: version}
			if checkNVEncMinVersion && !isNVEncVersionSupported(version) {
				info.Warning = "NVEncC 8.x 以降が必要です"
			}
			return info
		}
	}
	for _, candidate := range candidates {
		if p, err := exec.LookPath(candidate); err == nil {
			version := getVersion(p)
			info := model.ToolInfo{Found: true, Path: p, Version: version}
			if checkNVEncMinVersion && !isNVEncVersionSupported(version) {
				info.Warning = "NVEncC 8.x 以降が必要です"
			}
			return info
		}
	}
	return model.ToolInfo{Found: false}
}

func getVersion(path string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 0 {
		return ""
	}
	return strings.TrimSpace(lines[0])
}

func isNVEncVersionSupported(versionLine string) bool {
	if versionLine == "" {
		return false
	}
	re := regexp.MustCompile(`([0-9]+)\.[0-9]+`)
	m := re.FindStringSubmatch(versionLine)
	if len(m) < 2 {
		return false
	}
	major, err := strconv.Atoi(m[1])
	if err != nil {
		return false
	}
	return major >= 8
}

func EnsureNVEncReady(t model.ToolSnapshot) error {
	if !t.NVEncC.Found {
		return model.NewError(model.ErrToolNotFound, "NVEncC not found")
	}
	if t.NVEncC.Warning != "" {
		return &model.EnqueError{Code: model.ErrToolVersionUnsupported, Message: t.NVEncC.Warning}
	}
	return nil
}

func (d *Detector) GetGPUInfo(nvenccPath string) (model.GPUInfo, error) {
	if nvenccPath == "" {
		return model.GPUInfo{}, errors.New("nvencc path is empty")
	}
	device, err := runCheck(nvenccPath, "--check-device")
	if err != nil {
		return model.GPUInfo{}, fmt.Errorf("check-device: %w", err)
	}
	features, err := runCheck(nvenccPath, "--check-features")
	if err != nil {
		return model.GPUInfo{}, fmt.Errorf("check-features: %w", err)
	}
	return model.GPUInfo{CheckDeviceOutput: device, CheckFeaturesOutput: features}, nil
}

func runCheck(path, arg string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, arg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
