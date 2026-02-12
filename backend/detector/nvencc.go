package detector

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var nvenccCandidates = []string{"NVEncC64.exe", "NVEncC.exe", "NVEncC64", "NVEncC"}

// DetectNVEncC detects NVEncC and checks version.
func DetectNVEncC(configPath string) ToolInfo {
	info := ToolInfo{Name: "NVEncC"}

	path := findExecutable(configPath, nvenccCandidates)
	if path == "" {
		info.Error = "E_TOOL_NOT_FOUND"
		return info
	}

	info.Path = path
	info.Found = true

	version, err := getNVEncCVersion(path)
	if err != nil {
		info.Error = fmt.Sprintf("version detection failed: %v", err)
		info.Supported = false
		return info
	}

	info.Version = version
	major, err := parseMajorVersion(version)
	if err != nil {
		info.Error = fmt.Sprintf("version parse failed: %v", err)
		info.Supported = false
		return info
	}

	if major < 8 {
		info.Error = "E_TOOL_VERSION_UNSUPPORTED"
		info.Supported = false
	} else {
		info.Supported = true
	}

	return info
}

// GetGPUInfo runs NVEncC --check-device and returns the output.
func GetGPUInfo(nvenccPath string) (string, error) {
	if nvenccPath == "" {
		return "", fmt.Errorf("NVEncC path not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, nvenccPath, "--check-device")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("check-device failed: %w", err)
	}
	return string(out), nil
}

// GetGPUFeatures runs NVEncC --check-features and returns the output.
func GetGPUFeatures(nvenccPath string) (string, error) {
	if nvenccPath == "" {
		return "", fmt.Errorf("NVEncC path not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, nvenccPath, "--check-features")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("check-features failed: %w", err)
	}
	return string(out), nil
}

func getNVEncCVersion(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, path, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Some versions output version info even on error exit
		if len(out) > 0 {
			return parseVersionString(string(out))
		}
		return "", fmt.Errorf("run --version: %w", err)
	}
	return parseVersionString(string(out))
}

var versionRe = regexp.MustCompile(`(\d+\.\d+[\.\d]*)`)

func parseVersionString(output string) (string, error) {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(strings.ToLower(line), "nvencc") || strings.Contains(line, "version") {
			match := versionRe.FindString(line)
			if match != "" {
				return match, nil
			}
		}
	}
	match := versionRe.FindString(output)
	if match != "" {
		return match, nil
	}
	return "", fmt.Errorf("no version found in output")
}

func parseMajorVersion(version string) (int, error) {
	parts := strings.SplitN(version, ".", 2)
	if len(parts) == 0 {
		return 0, fmt.Errorf("no version parts")
	}
	return strconv.Atoi(parts[0])
}
