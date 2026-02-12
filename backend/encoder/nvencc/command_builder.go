package nvencc

import (
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/motoacs/enque/backend/encoder"
	"github.com/motoacs/enque/backend/model"
)

func BuildCommand(req encoder.BuildRequest) (encoder.BuildResult, error) {
	if err := model.MustValidateProfile(req.Profile); err != nil {
		return encoder.BuildResult{}, &model.EnqueError{Code: model.ErrValidation, Message: err.Error()}
	}
	p := req.Profile
	args := []string{}

	decoder := model.DecoderAVHW
	if p.Decoder == model.DecoderAVSW {
		decoder = model.DecoderAVSW
		if strings.TrimSpace(p.NVEncCAdvanced.AVSWDecoder) != "" {
			args = append(args, "--avsw", strings.TrimSpace(p.NVEncCAdvanced.AVSWDecoder))
		} else {
			args = append(args, "--avsw")
		}
	} else {
		args = append(args, "--avhw")
	}

	args = append(args, "-i", req.InputPath)
	args = append(args, "-c", string(p.Codec))

	// Standard GUI options.
	switch p.RateControl {
	case model.RateControlQVBR:
		args = append(args, "--qvbr", trimFloat(p.RateValue))
	case model.RateControlCQP:
		args = append(args, "--cqp", trimFloat(p.RateValue))
	case model.RateControlCBR:
		args = append(args, "--cbr", trimFloat(p.RateValue))
	case model.RateControlVBR:
		args = append(args, "--vbr", trimFloat(p.RateValue))
	}
	args = append(args, "--preset", string(p.Preset))
	args = append(args, "--output-depth", strconv.Itoa(p.OutputDepth))
	if p.Multipass != model.MultipassNone {
		args = append(args, "--multipass", string(p.Multipass))
	}
	if strings.TrimSpace(p.OutputRes) != "" {
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
	if p.GOPLen != nil {
		args = append(args, "--gop-len", strconv.Itoa(*p.GOPLen))
	}
	if p.AQ {
		args = append(args, "--aq")
	}
	if p.AQTemporal {
		args = append(args, "--aq-temporal")
	}
	if p.SplitEnc != model.SplitEncOff {
		args = append(args, "--split-enc", string(p.SplitEnc))
	}
	if p.Parallel != model.ParallelOff {
		args = append(args, "--parallel", string(p.Parallel))
	}
	if p.Device != "auto" {
		args = append(args, "--device", p.Device)
	}

	switch p.AudioMode {
	case model.AudioModeCopy:
		args = append(args, "--audio-copy")
	case model.AudioModeAAC:
		args = append(args, "--audio-codec", "aac", "--audio-bitrate", strconv.Itoa(p.AudioBitrate))
	case model.AudioModeOpus:
		args = append(args, "--audio-codec", "opus", "--audio-bitrate", strconv.Itoa(p.AudioBitrate))
	}

	if p.ColorMatrix != "auto" && p.ColorMatrix != "" {
		args = append(args, "--colormatrix", p.ColorMatrix)
	}
	if p.Transfer != "auto" && p.Transfer != "" {
		args = append(args, "--transfer", p.Transfer)
	}
	if p.ColorPrim != "auto" && p.ColorPrim != "" {
		args = append(args, "--colorprim", p.ColorPrim)
	}
	if p.ColorRange != "auto" && p.ColorRange != "" {
		args = append(args, "--colorrange", p.ColorRange)
	}
	if p.DHDR10Info == "copy" {
		args = append(args, "--dhdr10-info", "copy")
	}
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

	// Advanced GUI options (later-wins).
	appendAdv := func(flag, value string) {
		if strings.TrimSpace(value) != "" {
			args = append(args, flag, strings.TrimSpace(value))
		}
	}
	appendAdv("--interlace", p.NVEncCAdvanced.Interlace)
	appendAdv("--input-csp", p.NVEncCAdvanced.InputCSP)
	appendAdv("--output-csp", p.NVEncCAdvanced.OutputCSP)
	appendAdv("--tune", p.NVEncCAdvanced.Tune)
	if p.NVEncCAdvanced.MaxBitrate != nil {
		args = append(args, "--max-bitrate", strconv.Itoa(*p.NVEncCAdvanced.MaxBitrate))
	}
	if p.NVEncCAdvanced.VBRQuality != nil {
		args = append(args, "--vbr-quality", trimFloat(*p.NVEncCAdvanced.VBRQuality))
	}
	if p.NVEncCAdvanced.LookaheadLevel != nil {
		args = append(args, "--lookahead-level", strconv.Itoa(*p.NVEncCAdvanced.LookaheadLevel))
	}
	if p.NVEncCAdvanced.WeightP {
		args = append(args, "--weightp")
	}
	appendAdv("--mv-precision", p.NVEncCAdvanced.MVPrecision)
	if p.NVEncCAdvanced.RefsForward != nil {
		args = append(args, "--refs-forward", strconv.Itoa(*p.NVEncCAdvanced.RefsForward))
	}
	if p.NVEncCAdvanced.RefsBackward != nil {
		args = append(args, "--refs-backward", strconv.Itoa(*p.NVEncCAdvanced.RefsBackward))
	}
	appendAdv("--level", p.NVEncCAdvanced.Level)
	appendAdv("--profile", p.NVEncCAdvanced.Profile)
	appendAdv("--tier", p.NVEncCAdvanced.Tier)
	if p.NVEncCAdvanced.SSIM {
		args = append(args, "--ssim")
	}
	if p.NVEncCAdvanced.PSNR {
		args = append(args, "--psnr")
	}
	appendAdv("--trim", p.NVEncCAdvanced.Trim)
	appendAdv("--seek", p.NVEncCAdvanced.Seek)
	appendAdv("--seekto", p.NVEncCAdvanced.SeekTo)
	appendAdv("--video-metadata", p.NVEncCAdvanced.VideoMetadata)
	appendAdv("--audio-copy", p.NVEncCAdvanced.AudioCopy)
	appendAdv("--audio-codec", p.NVEncCAdvanced.AudioCodec)
	appendAdv("--audio-bitrate", p.NVEncCAdvanced.AudioBitrate)
	appendAdv("--audio-quality", p.NVEncCAdvanced.AudioQuality)
	appendAdv("--audio-samplerate", p.NVEncCAdvanced.AudioSamplerate)
	appendAdv("--audio-metadata", p.NVEncCAdvanced.AudioMetadata)
	appendAdv("--sub-copy", p.NVEncCAdvanced.SubCopy)
	appendAdv("--sub-metadata", p.NVEncCAdvanced.SubMetadata)
	appendAdv("--data-copy", p.NVEncCAdvanced.DataCopy)
	appendAdv("--attachment-copy", p.NVEncCAdvanced.AttachmentCopy)
	appendAdv("--metadata", p.NVEncCAdvanced.Metadata)
	if p.NVEncCAdvanced.OutputThread != nil {
		args = append(args, "--output-thread", strconv.Itoa(*p.NVEncCAdvanced.OutputThread))
	}

	if strings.TrimSpace(p.CustomOptions) != "" {
		t, err := encoder.TokenizeCustomOptions(p.CustomOptions)
		if err != nil {
			return encoder.BuildResult{}, err
		}
		args = append(args, t...)
	}

	args = append(args, "-o", req.OutputPath)

	return encoder.BuildResult{
		Argv:             args,
		DisplayCommand:   toDisplayCommand("NVEncC64.exe", args),
		EffectiveDecoder: decoder,
	}, nil
}

func trimFloat(v float64) string {
	return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(v, 'f', 3, 64), "0"), ".")
}

func toDisplayCommand(bin string, args []string) string {
	quoted := make([]string, 0, len(args)+1)
	quoted = append(quoted, quote(bin))
	for _, a := range args {
		quoted = append(quoted, quote(a))
	}
	return strings.Join(quoted, " ")
}

func quote(v string) string {
	if v == "" {
		return `""`
	}
	if strings.ContainsAny(v, " \t\"'") {
		return "\"" + strings.ReplaceAll(v, "\"", `\\"`) + "\""
	}
	return v
}

func DetectCapabilities(ctx context.Context, encoderPath string) (map[string]any, error) {
	cmd := exec.CommandContext(ctx, encoderPath, "--check-features")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"raw":       string(out),
		"checkedAt": time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func EnsureOptionOrder(args []string) error {
	i := slices.Index(args, "-i")
	o := -1
	for idx, v := range args {
		if v == "-o" {
			o = idx
		}
	}
	if i < 0 || o < 0 || o <= i {
		return fmt.Errorf("invalid option order")
	}
	return nil
}
