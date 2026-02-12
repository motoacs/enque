package nvencc

import (
	"fmt"
	"strconv"

	"github.com/yuta/enque/backend/encoder"
	"github.com/yuta/enque/backend/profile"
)

// NVEncCAdapter implements the Adapter interface for NVEncC.
type NVEncCAdapter struct{}

func (a *NVEncCAdapter) Type() string { return "nvencc" }

func (a *NVEncCAdapter) SupportsDecoderFallback() bool { return true }

// BuildArgs generates NVEncC command-line arguments from a profile.
// Argument order is a fixed contract (design doc 9.3.1):
// 1. --avhw/--avsw  2. -i  3. video basic  4. video detail  5. speed
// 6. audio  7. color  8. metadata  9. nvencc_advanced  10. custom_options  11. -o
func (a *NVEncCAdapter) BuildArgs(p profile.Profile, inputPath, outputPath string) ([]string, error) {
	var args []string

	// 1. Decoder (front-positioned)
	args = appendDecoder(args, p)

	// 2. Input
	args = append(args, "-i", inputPath)

	// 3. Video basic
	args = append(args, "-c", p.Codec)
	args = appendRateControl(args, p)
	args = append(args, "--preset", p.Preset)
	args = append(args, "--output-depth", strconv.Itoa(p.OutputDepth))

	// 4. Video detail (standard GUI)
	args = appendVideoDetail(args, p)

	// 5. Speed
	args = appendSpeed(args, p)

	// 6. Audio (standard GUI)
	args = appendAudio(args, p)

	// 7. Color
	args = appendColor(args, p)

	// 8. Metadata (standard GUI)
	args = appendMetadata(args, p)

	// 9. NVEncC Advanced (overwrites standard GUI by later-wins)
	args = appendAdvanced(args, p.NVEncCAdvanced)

	// 10. Custom options (final priority)
	if p.CustomOptions != "" {
		tokens, err := encoder.TokenizeCustomOptions(p.CustomOptions)
		if err != nil {
			return nil, fmt.Errorf("%s: custom_options: %w", encoder.ErrValidation, err)
		}
		args = append(args, tokens...)
	}

	// 11. Output
	args = append(args, "-o", outputPath)

	return args, nil
}

// BuildArgsWithDecoderOverride builds args but forces the specified decoder (for fallback retry).
func (a *NVEncCAdapter) BuildArgsWithDecoderOverride(p profile.Profile, inputPath, outputPath, decoderOverride string) ([]string, error) {
	overridden := p
	overridden.Decoder = decoderOverride
	overridden.NVEncCAdvanced.AVSWDecoder = ""
	return a.BuildArgs(overridden, inputPath, outputPath)
}

func appendDecoder(args []string, p profile.Profile) []string {
	switch p.Decoder {
	case "avsw":
		if p.NVEncCAdvanced.AVSWDecoder != "" {
			args = append(args, "--avsw", p.NVEncCAdvanced.AVSWDecoder)
		} else {
			args = append(args, "--avsw")
		}
	default: // "avhw"
		args = append(args, "--avhw")
	}
	return args
}

func appendRateControl(args []string, p profile.Profile) []string {
	rv := strconv.FormatFloat(p.RateValue, 'f', -1, 64)
	switch p.RateControl {
	case "qvbr":
		args = append(args, "--qvbr", rv)
	case "cqp":
		args = append(args, "--cqp", rv)
	case "cbr":
		args = append(args, "--cbr", rv)
	case "vbr":
		args = append(args, "--vbr", rv)
	}
	return args
}

func appendVideoDetail(args []string, p profile.Profile) []string {
	if p.Multipass != "none" && p.Multipass != "" {
		args = append(args, "--multipass", p.Multipass)
	}
	if p.OutputRes != "" {
		args = append(args, "--output-res", p.OutputRes)
	}
	if p.Bframes != nil {
		args = append(args, "--bframes", strconv.Itoa(*p.Bframes))
	}
	if p.Ref != nil {
		args = append(args, "--ref", strconv.Itoa(*p.Ref))
	}
	if p.Lookahead != nil {
		args = append(args, "--lookahead", strconv.Itoa(*p.Lookahead))
	}
	if p.GopLen != nil {
		args = append(args, "--gop-len", strconv.Itoa(*p.GopLen))
	}
	if p.AQ {
		args = append(args, "--aq")
	}
	if p.AQTemporal {
		args = append(args, "--aq-temporal")
	}
	return args
}

func appendSpeed(args []string, p profile.Profile) []string {
	if p.SplitEnc != "off" && p.SplitEnc != "" {
		args = append(args, "--split-enc", p.SplitEnc)
	}
	if p.Parallel != "off" && p.Parallel != "" {
		args = append(args, "--parallel", p.Parallel)
	}
	if p.Device != "auto" && p.Device != "" {
		args = append(args, "--device", p.Device)
	}
	return args
}

func appendAudio(args []string, p profile.Profile) []string {
	switch p.AudioMode {
	case "copy":
		args = append(args, "--audio-copy")
	case "aac":
		args = append(args, "--audio-codec", "aac", "--audio-bitrate", strconv.Itoa(p.AudioBitrate))
	case "opus":
		args = append(args, "--audio-codec", "opus", "--audio-bitrate", strconv.Itoa(p.AudioBitrate))
	}
	return args
}

func appendColor(args []string, p profile.Profile) []string {
	if p.Colormatrix != "auto" && p.Colormatrix != "" {
		args = append(args, "--colormatrix", p.Colormatrix)
	}
	if p.Transfer != "auto" && p.Transfer != "" {
		args = append(args, "--transfer", p.Transfer)
	}
	if p.Colorprim != "auto" && p.Colorprim != "" {
		args = append(args, "--colorprim", p.Colorprim)
	}
	if p.Colorrange != "auto" && p.Colorrange != "" {
		args = append(args, "--colorrange", p.Colorrange)
	}
	if p.DHDR10Info == "copy" {
		args = append(args, "--dhdr10-info", "copy")
	}
	return args
}

func appendMetadata(args []string, p profile.Profile) []string {
	if p.MetadataCopy {
		args = append(args, "--metadata", "copy")
	}
	if p.VideoMetadataCopy {
		args = append(args, "--video-metadata", "copy")
	}
	if p.AudioMetadataCopy {
		args = append(args, "--audio-metadata", "copy")
	}
	if p.ChapterCopy {
		args = append(args, "--chapter-copy")
	}
	if p.SubCopy {
		args = append(args, "--sub-copy")
	}
	if p.DataCopy {
		args = append(args, "--data-copy")
	}
	if p.AttachmentCopy {
		args = append(args, "--attachment-copy")
	}
	return args
}

func appendAdvanced(args []string, adv profile.NVEncCAdvanced) []string {
	if adv.Interlace != "" {
		args = append(args, "--interlace", adv.Interlace)
	}
	if adv.InputCSP != "" {
		args = append(args, "--input-csp", adv.InputCSP)
	}
	if adv.OutputCSP != "" {
		args = append(args, "--output-csp", adv.OutputCSP)
	}
	if adv.Tune != "" {
		args = append(args, "--tune", adv.Tune)
	}
	if adv.MaxBitrate != nil {
		args = append(args, "--max-bitrate", strconv.Itoa(*adv.MaxBitrate))
	}
	if adv.VBRQuality != nil {
		args = append(args, "--vbr-quality", strconv.Itoa(*adv.VBRQuality))
	}
	if adv.LookaheadLevel != nil {
		args = append(args, "--lookahead-level", strconv.Itoa(*adv.LookaheadLevel))
	}
	if adv.WeightP {
		args = append(args, "--weightp")
	}
	if adv.MVPrecision != "" {
		args = append(args, "--mv-precision", adv.MVPrecision)
	}
	if adv.RefsForward != nil {
		args = append(args, "--refs-forward", strconv.Itoa(*adv.RefsForward))
	}
	if adv.RefsBackward != nil {
		args = append(args, "--refs-backward", strconv.Itoa(*adv.RefsBackward))
	}
	if adv.Level != "" {
		args = append(args, "--level", adv.Level)
	}
	if adv.Profile != "" {
		args = append(args, "--profile", adv.Profile)
	}
	if adv.Tier != "" {
		args = append(args, "--tier", adv.Tier)
	}
	if adv.OutputThread != nil {
		args = append(args, "--output-thread", strconv.Itoa(*adv.OutputThread))
	}
	if adv.SSIM {
		args = append(args, "--ssim")
	}
	if adv.PSNR {
		args = append(args, "--psnr")
	}
	if adv.Trim != "" {
		args = append(args, "--trim", adv.Trim)
	}
	if adv.Seek != "" {
		args = append(args, "--seek", adv.Seek)
	}
	if adv.SeekTo != "" {
		args = append(args, "--seekto", adv.SeekTo)
	}
	if adv.VideoMetadata != "" {
		args = append(args, "--video-metadata", adv.VideoMetadata)
	}
	if adv.AudioCopy != "" {
		args = append(args, "--audio-copy", adv.AudioCopy)
	}
	if adv.AudioCodec != "" {
		args = append(args, "--audio-codec", adv.AudioCodec)
	}
	if adv.AudioBitrate != "" {
		args = append(args, "--audio-bitrate", adv.AudioBitrate)
	}
	if adv.AudioQuality != "" {
		args = append(args, "--audio-quality", adv.AudioQuality)
	}
	if adv.AudioSamplerate != "" {
		args = append(args, "--audio-samplerate", adv.AudioSamplerate)
	}
	if adv.AudioMetadata != "" {
		args = append(args, "--audio-metadata", adv.AudioMetadata)
	}
	if adv.SubCopy != "" {
		args = append(args, "--sub-copy", adv.SubCopy)
	}
	if adv.SubMetadata != "" {
		args = append(args, "--sub-metadata", adv.SubMetadata)
	}
	if adv.DataCopy != "" {
		args = append(args, "--data-copy", adv.DataCopy)
	}
	if adv.AttachmentCopy != "" {
		args = append(args, "--attachment-copy", adv.AttachmentCopy)
	}
	if adv.Metadata != "" {
		args = append(args, "--metadata", adv.Metadata)
	}
	return args
}
