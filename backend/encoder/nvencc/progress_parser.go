package nvencc

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/yuta/enque/backend/encoder"
)

// Progress parsing regex patterns for NVEncC stderr output.
// Examples:
//   [53.2%] 1234 frames: 245.67 fps, 12345 kb/s, remain 0:01:23, GPU 45%, VE 78%, VD 12%
//   [100.0%] 5000 frames: 300.00 fps, 8765 kb/s, remain 0:00:00
var progressRe = regexp.MustCompile(
	`\[\s*(\d+\.?\d*)\s*%\].*?(\d+\.?\d*)\s*fps.*?(\d+\.?\d*)\s*kb/s(?:.*?remain\s+(\d+):(\d+):(\d+))?`,
)

// ParseProgress extracts progress information from a single line of NVEncC stderr.
func (a *NVEncCAdapter) ParseProgress(line string) encoder.Progress {
	line = strings.TrimSpace(line)
	if line == "" {
		return encoder.Progress{RawLine: line}
	}

	matches := progressRe.FindStringSubmatch(line)
	if matches == nil {
		return encoder.Progress{RawLine: line}
	}

	prog := encoder.Progress{RawLine: line}

	if pct, err := strconv.ParseFloat(matches[1], 64); err == nil {
		prog.Percent = &pct
	}
	if fps, err := strconv.ParseFloat(matches[2], 64); err == nil {
		prog.FPS = &fps
	}
	if br, err := strconv.ParseFloat(matches[3], 64); err == nil {
		prog.BitrateKbps = &br
	}
	if len(matches) > 5 && matches[4] != "" {
		h, _ := strconv.ParseFloat(matches[4], 64)
		m, _ := strconv.ParseFloat(matches[5], 64)
		s, _ := strconv.ParseFloat(matches[6], 64)
		eta := h*3600 + m*60 + s
		prog.ETASec = &eta
	}

	return prog
}
