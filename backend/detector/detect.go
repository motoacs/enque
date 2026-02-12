package detector

import "github.com/yuta/enque/backend/config"

// DetectAll runs detection for all external tools.
func DetectAll(cfg config.AppConfig) DetectionResult {
	return DetectionResult{
		NVEncC:  DetectNVEncC(cfg.NVEncCPath),
		QSVEncC: DetectQSVEncC(cfg.QSVEncPath),
		FFmpeg:  DetectFFmpeg(cfg.FFmpegPath),
		FFprobe: DetectFFprobe(cfg.FFprobePath),
	}
}
