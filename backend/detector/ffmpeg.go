package detector

var ffmpegCandidates = []string{"ffmpeg.exe", "ffmpeg"}
var ffprobeCandidates = []string{"ffprobe.exe", "ffprobe"}

// DetectFFmpeg detects ffmpeg (optional).
func DetectFFmpeg(configPath string) ToolInfo {
	info := ToolInfo{Name: "ffmpeg"}
	path := findExecutable(configPath, ffmpegCandidates)
	if path == "" {
		info.Error = "not found (optional)"
		return info
	}
	info.Path = path
	info.Found = true
	info.Supported = false // v1: adapter not implemented
	return info
}

// DetectFFprobe detects ffprobe (optional).
func DetectFFprobe(configPath string) ToolInfo {
	info := ToolInfo{Name: "ffprobe"}
	path := findExecutable(configPath, ffprobeCandidates)
	if path == "" {
		info.Error = "not found (optional)"
		return info
	}
	info.Path = path
	info.Found = true
	info.Supported = false
	return info
}
