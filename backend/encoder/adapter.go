package encoder

import "github.com/yuta/enque/backend/profile"

// Progress holds parsed progress information from encoder stderr.
type Progress struct {
	Percent    *float64 `json:"percent"`
	FPS        *float64 `json:"fps"`
	BitrateKbps *float64 `json:"bitrate_kbps"`
	ETASec     *float64 `json:"eta_sec"`
	RawLine    string   `json:"raw_line"`
}

// Adapter defines the interface for encoder-specific implementations.
type Adapter interface {
	// Type returns the encoder type identifier (e.g., "nvencc").
	Type() string

	// BuildArgs generates command-line arguments from a profile.
	BuildArgs(p profile.Profile, inputPath, outputPath string) ([]string, error)

	// ParseProgress parses a single line of stderr output into progress.
	ParseProgress(line string) Progress

	// SupportsDecoderFallback returns whether this adapter supports avhw->avsw fallback.
	SupportsDecoderFallback() bool
}
