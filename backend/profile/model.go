package profile

// Profile represents an encoding profile (design doc 5.2).
type Profile struct {
	ID           string         `json:"id"`
	Version      int            `json:"version"`
	Name         string         `json:"name"`
	IsPreset     bool           `json:"is_preset"`
	EncoderType  string         `json:"encoder_type"`
	EncoderOpts  map[string]any `json:"encoder_options"`

	// Video basic (nvencc)
	Codec       string  `json:"codec"`
	RateControl string  `json:"rate_control"`
	RateValue   float64 `json:"rate_value"`
	Preset      string  `json:"preset"`
	OutputDepth int     `json:"output_depth"`
	Multipass   string  `json:"multipass"`
	OutputRes   string  `json:"output_res"`

	// Video detail
	Bframes   *int `json:"bframes"`
	Ref       *int `json:"ref"`
	Lookahead *int `json:"lookahead"`
	GopLen    *int `json:"gop_len"`
	AQ        bool `json:"aq"`
	AQTemporal bool `json:"aq_temporal"`

	// Speed
	SplitEnc string `json:"split_enc"`
	Parallel string `json:"parallel"`
	Decoder  string `json:"decoder"`
	Device   string `json:"device"`

	// Audio
	AudioMode    string `json:"audio_mode"`
	AudioBitrate int    `json:"audio_bitrate"`

	// Color
	Colormatrix string `json:"colormatrix"`
	Transfer    string `json:"transfer"`
	Colorprim   string `json:"colorprim"`
	Colorrange  string `json:"colorrange"`
	DHDR10Info  string `json:"dhdr10_info"`

	// Metadata
	MetadataCopy      bool `json:"metadata_copy"`
	VideoMetadataCopy bool `json:"video_metadata_copy"`
	AudioMetadataCopy bool `json:"audio_metadata_copy"`
	ChapterCopy       bool `json:"chapter_copy"`
	SubCopy           bool `json:"sub_copy"`
	DataCopy          bool `json:"data_copy"`
	AttachmentCopy    bool `json:"attachment_copy"`
	RestoreFileTime   bool `json:"restore_file_time"`

	// Advanced NVEncC options
	NVEncCAdvanced NVEncCAdvanced `json:"nvencc_advanced"`

	// Output
	OutputContainer string `json:"output_container"`

	// Custom options (layer 2)
	CustomOptions string `json:"custom_options"`
}

// NVEncCAdvanced holds advanced NVEncC-specific options (design doc 5.2).
type NVEncCAdvanced struct {
	Interlace      string `json:"interlace"`
	AVSWDecoder    string `json:"avsw_decoder"`
	InputCSP       string `json:"input_csp"`
	OutputCSP      string `json:"output_csp"`
	Tune           string `json:"tune"`
	MaxBitrate     *int   `json:"max_bitrate"`
	VBRQuality     *int   `json:"vbr_quality"`
	LookaheadLevel *int   `json:"lookahead_level"`
	WeightP        bool   `json:"weightp"`
	MVPrecision    string `json:"mv_precision"`
	RefsForward    *int   `json:"refs_forward"`
	RefsBackward   *int   `json:"refs_backward"`
	Level          string `json:"level"`
	Profile        string `json:"profile"`
	Tier           string `json:"tier"`
	SSIM           bool   `json:"ssim"`
	PSNR           bool   `json:"psnr"`
	Trim           string `json:"trim"`
	Seek           string `json:"seek"`
	SeekTo         string `json:"seekto"`
	VideoMetadata  string `json:"video_metadata"`
	AudioCopy      string `json:"audio_copy"`
	AudioCodec     string `json:"audio_codec"`
	AudioBitrate   string `json:"audio_bitrate"`
	AudioQuality   string `json:"audio_quality"`
	AudioSamplerate string `json:"audio_samplerate"`
	AudioMetadata  string `json:"audio_metadata"`
	SubCopy        string `json:"sub_copy"`
	SubMetadata    string `json:"sub_metadata"`
	DataCopy       string `json:"data_copy"`
	AttachmentCopy string `json:"attachment_copy"`
	Metadata       string `json:"metadata"`
	OutputThread   *int   `json:"output_thread"`
}

// ProfilesFile is the top-level structure for profiles.json.
type ProfilesFile struct {
	Profiles []Profile `json:"profiles"`
}

// CurrentVersion is the latest profile schema version.
const CurrentVersion = 4
