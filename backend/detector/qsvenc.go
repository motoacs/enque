package detector

var qsvencCandidates = []string{"QSVEncC64.exe", "QSVEncC.exe", "QSVEncC64", "QSVEncC"}

// DetectQSVEncC detects QSVEncC (v1: warning only if not found).
func DetectQSVEncC(configPath string) ToolInfo {
	info := ToolInfo{Name: "QSVEncC"}

	path := findExecutable(configPath, qsvencCandidates)
	if path == "" {
		info.Error = "not found (optional)"
		return info
	}

	info.Path = path
	info.Found = true
	info.Supported = false // v1: adapter not implemented
	return info
}
