package model

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const ProfileVersion = 4

type NVEncCAdvanced struct {
	Interlace       string   `json:"interlace"`
	AVSWDecoder     string   `json:"avsw_decoder"`
	InputCSP        string   `json:"input_csp"`
	OutputCSP       string   `json:"output_csp"`
	Tune            string   `json:"tune"`
	MaxBitrate      *int     `json:"max_bitrate"`
	VBRQuality      *float64 `json:"vbr_quality"`
	LookaheadLevel  *int     `json:"lookahead_level"`
	WeightP         bool     `json:"weightp"`
	MVPrecision     string   `json:"mv_precision"`
	RefsForward     *int     `json:"refs_forward"`
	RefsBackward    *int     `json:"refs_backward"`
	Level           string   `json:"level"`
	Profile         string   `json:"profile"`
	Tier            string   `json:"tier"`
	SSIM            bool     `json:"ssim"`
	PSNR            bool     `json:"psnr"`
	Trim            string   `json:"trim"`
	Seek            string   `json:"seek"`
	SeekTo          string   `json:"seekto"`
	VideoMetadata   string   `json:"video_metadata"`
	AudioCopy       string   `json:"audio_copy"`
	AudioCodec      string   `json:"audio_codec"`
	AudioBitrate    string   `json:"audio_bitrate"`
	AudioQuality    string   `json:"audio_quality"`
	AudioSamplerate string   `json:"audio_samplerate"`
	AudioMetadata   string   `json:"audio_metadata"`
	SubCopy         string   `json:"sub_copy"`
	SubMetadata     string   `json:"sub_metadata"`
	DataCopy        string   `json:"data_copy"`
	AttachmentCopy  string   `json:"attachment_copy"`
	Metadata        string   `json:"metadata"`
	OutputThread    *int     `json:"output_thread"`
}

type Profile struct {
	ID             string         `json:"id"`
	Version        int            `json:"version"`
	Name           string         `json:"name"`
	IsPreset       bool           `json:"is_preset"`
	EncoderType    EncoderType    `json:"encoder_type"`
	EncoderOptions map[string]any `json:"encoder_options"`

	Codec        Codec        `json:"codec"`
	RateControl  RateControl  `json:"rate_control"`
	RateValue    float64      `json:"rate_value"`
	Preset       Preset       `json:"preset"`
	OutputDepth  int          `json:"output_depth"`
	Multipass    Multipass    `json:"multipass"`
	OutputRes    string       `json:"output_res"`
	Bframes      *int         `json:"bframes"`
	Ref          *int         `json:"ref"`
	Lookahead    *int         `json:"lookahead"`
	GOPLen       *int         `json:"gop_len"`
	AQ           bool         `json:"aq"`
	AQTemporal   bool         `json:"aq_temporal"`
	SplitEnc     SplitEnc     `json:"split_enc"`
	Parallel     ParallelMode `json:"parallel"`
	Decoder      Decoder      `json:"decoder"`
	Device       string       `json:"device"`
	AudioMode    AudioMode    `json:"audio_mode"`
	AudioBitrate int          `json:"audio_bitrate"`
	ColorMatrix  string       `json:"colormatrix"`
	Transfer     string       `json:"transfer"`
	ColorPrim    string       `json:"colorprim"`
	ColorRange   string       `json:"colorrange"`
	DHDR10Info   string       `json:"dhdr10_info"`

	MetadataCopy      bool `json:"metadata_copy"`
	VideoMetadataCopy bool `json:"video_metadata_copy"`
	AudioMetadataCopy bool `json:"audio_metadata_copy"`
	ChapterCopy       bool `json:"chapter_copy"`
	SubCopy           bool `json:"sub_copy"`
	DataCopy          bool `json:"data_copy"`
	AttachmentCopy    bool `json:"attachment_copy"`
	RestoreFileTime   bool `json:"restore_file_time"`

	NVEncCAdvanced NVEncCAdvanced `json:"nvencc_advanced"`
	CustomOptions  string         `json:"custom_options"`
}

func DefaultProfile() Profile {
	return Profile{
		Version:           ProfileVersion,
		Name:              "HEVC Quality",
		EncoderType:       EncoderTypeNVEncC,
		EncoderOptions:    map[string]any{},
		Codec:             CodecHEVC,
		RateControl:       RateControlQVBR,
		RateValue:         28,
		Preset:            "P4",
		OutputDepth:       10,
		Multipass:         MultipassNone,
		AQ:                true,
		AQTemporal:        true,
		SplitEnc:          SplitEncAuto,
		Parallel:          ParallelOff,
		Decoder:           DecoderAVHW,
		Device:            "auto",
		AudioMode:         AudioModeCopy,
		AudioBitrate:      256,
		ColorMatrix:       "auto",
		Transfer:          "auto",
		ColorPrim:         "auto",
		ColorRange:        "auto",
		DHDR10Info:        "off",
		MetadataCopy:      true,
		VideoMetadataCopy: true,
		AudioMetadataCopy: true,
		ChapterCopy:       true,
		SubCopy:           true,
		DataCopy:          true,
		AttachmentCopy:    true,
	}
}

func ValidateProfile(p Profile) map[string]string {
	errMap := map[string]string{}
	name := strings.TrimSpace(p.Name)
	if len(name) < 1 || len(name) > 80 {
		errMap["name"] = "name must be 1..80 chars"
	}
	switch p.EncoderType {
	case EncoderTypeNVEncC, EncoderTypeQSVEnc, EncoderTypeFFmpeg:
	default:
		errMap["encoder_type"] = "unsupported encoder_type"
	}
	if p.RateValue <= 0 {
		errMap["rate_value"] = "rate_value must be >0"
	}
	if p.OutputDepth != 8 && p.OutputDepth != 10 {
		errMap["output_depth"] = "output_depth must be 8 or 10"
	}
	if p.OutputRes != "" {
		r := regexp.MustCompile(`^\d+x\d+(,[^,=]+=[^,=]+)*$`)
		if !r.MatchString(p.OutputRes) {
			errMap["output_res"] = "invalid output_res"
		}
	}
	if p.Bframes != nil && (*p.Bframes < 0 || *p.Bframes > 7) {
		errMap["bframes"] = "bframes must be 0..7"
	}
	if p.Lookahead != nil && (*p.Lookahead < 0 || *p.Lookahead > 32) {
		errMap["lookahead"] = "lookahead must be 0..32"
	}
	if p.AudioBitrate < 32 || p.AudioBitrate > 1024 {
		errMap["audio_bitrate"] = "audio_bitrate must be 32..1024"
	}
	if p.CustomOptions != "" && len([]rune(p.CustomOptions)) > 4096 {
		errMap["custom_options"] = "custom_options too long"
	}
	if p.Device != "auto" {
		if !regexp.MustCompile(`^([0-9]|1[0-5])$`).MatchString(p.Device) {
			errMap["device"] = "device must be auto or 0..15"
		}
	}
	validateAdvanced(p.NVEncCAdvanced, errMap)
	return errMap
}

func validateAdvanced(v NVEncCAdvanced, errMap map[string]string) {
	checkStr := func(field, val string) {
		if len([]rune(val)) > 1024 {
			errMap[field] = "string too long"
		}
	}
	checkStr("nvencc_advanced.interlace", v.Interlace)
	checkStr("nvencc_advanced.avsw_decoder", v.AVSWDecoder)
	checkStr("nvencc_advanced.input_csp", v.InputCSP)
	checkStr("nvencc_advanced.output_csp", v.OutputCSP)
	checkStr("nvencc_advanced.tune", v.Tune)
	checkStr("nvencc_advanced.mv_precision", v.MVPrecision)
	checkStr("nvencc_advanced.level", v.Level)
	checkStr("nvencc_advanced.profile", v.Profile)
	checkStr("nvencc_advanced.tier", v.Tier)
	checkStr("nvencc_advanced.trim", v.Trim)
	checkStr("nvencc_advanced.seek", v.Seek)
	checkStr("nvencc_advanced.seekto", v.SeekTo)
	checkStr("nvencc_advanced.video_metadata", v.VideoMetadata)
	checkStr("nvencc_advanced.audio_copy", v.AudioCopy)
	checkStr("nvencc_advanced.audio_codec", v.AudioCodec)
	checkStr("nvencc_advanced.audio_bitrate", v.AudioBitrate)
	checkStr("nvencc_advanced.audio_quality", v.AudioQuality)
	checkStr("nvencc_advanced.audio_samplerate", v.AudioSamplerate)
	checkStr("nvencc_advanced.audio_metadata", v.AudioMetadata)
	checkStr("nvencc_advanced.sub_copy", v.SubCopy)
	checkStr("nvencc_advanced.sub_metadata", v.SubMetadata)
	checkStr("nvencc_advanced.data_copy", v.DataCopy)
	checkStr("nvencc_advanced.attachment_copy", v.AttachmentCopy)
	checkStr("nvencc_advanced.metadata", v.Metadata)

	if v.MaxBitrate != nil && *v.MaxBitrate <= 0 {
		errMap["nvencc_advanced.max_bitrate"] = "must be >0"
	}
	if v.VBRQuality != nil && *v.VBRQuality <= 0 {
		errMap["nvencc_advanced.vbr_quality"] = "must be >0"
	}
	if v.LookaheadLevel != nil && *v.LookaheadLevel < 0 {
		errMap["nvencc_advanced.lookahead_level"] = "must be >=0"
	}
	if v.RefsForward != nil && *v.RefsForward < 0 {
		errMap["nvencc_advanced.refs_forward"] = "must be >=0"
	}
	if v.RefsBackward != nil && *v.RefsBackward < 0 {
		errMap["nvencc_advanced.refs_backward"] = "must be >=0"
	}
	if v.OutputThread != nil && (*v.OutputThread < 1 || *v.OutputThread > 64) {
		errMap["nvencc_advanced.output_thread"] = "must be 1..64"
	}
}

func MustValidateProfile(p Profile) error {
	errMap := ValidateProfile(p)
	if len(errMap) == 0 {
		return nil
	}
	keys := make([]string, 0, len(errMap))
	for k := range errMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s: %s", k, errMap[k]))
	}
	return fmt.Errorf("validation failed: %s", strings.Join(parts, "; "))
}
