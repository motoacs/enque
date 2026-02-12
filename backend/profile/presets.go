package profile

import "github.com/google/uuid"

// GeneratePresets returns the 4 built-in preset profiles (design doc 18).
func GeneratePresets() []Profile {
	return []Profile{
		hevcQuality(),
		av1Fast(),
		cameraArchive(),
		h264Compatible(),
	}
}

func defaultAdvanced() NVEncCAdvanced {
	return NVEncCAdvanced{}
}

func hevcQuality() Profile {
	return Profile{
		ID:          uuid.New().String(),
		Version:     CurrentVersion,
		Name:        "HEVC Quality",
		IsPreset:    true,
		EncoderType: "nvencc",
		EncoderOpts: map[string]any{},

		Codec:       "hevc",
		RateControl: "qvbr",
		RateValue:   28,
		Preset:      "P4",
		OutputDepth: 10,
		Multipass:   "none",
		OutputRes:   "",

		AQ:         true,
		AQTemporal: true,

		SplitEnc: "auto",
		Parallel: "off",
		Decoder:  "avhw",
		Device:   "auto",

		AudioMode:    "copy",
		AudioBitrate: 256,

		Colormatrix: "auto",
		Transfer:    "auto",
		Colorprim:   "auto",
		Colorrange:  "auto",
		DHDR10Info:  "off",

		MetadataCopy:      true,
		VideoMetadataCopy: true,
		AudioMetadataCopy: true,
		ChapterCopy:       true,
		SubCopy:           true,
		DataCopy:          true,
		AttachmentCopy:    true,
		RestoreFileTime:   false,

		OutputContainer: "mp4",

		NVEncCAdvanced: defaultAdvanced(),
		CustomOptions:  "",
	}
}

func av1Fast() Profile {
	return Profile{
		ID:          uuid.New().String(),
		Version:     CurrentVersion,
		Name:        "AV1 Fast",
		IsPreset:    true,
		EncoderType: "nvencc",
		EncoderOpts: map[string]any{},

		Codec:       "av1",
		RateControl: "qvbr",
		RateValue:   32,
		Preset:      "P1",
		OutputDepth: 10,
		Multipass:   "none",
		OutputRes:   "",

		AQ:         true,
		AQTemporal: true,

		SplitEnc: "auto",
		Parallel: "off",
		Decoder:  "avhw",
		Device:   "auto",

		AudioMode:    "copy",
		AudioBitrate: 256,

		Colormatrix: "auto",
		Transfer:    "auto",
		Colorprim:   "auto",
		Colorrange:  "auto",
		DHDR10Info:  "off",

		MetadataCopy:      true,
		VideoMetadataCopy: true,
		AudioMetadataCopy: true,
		ChapterCopy:       true,
		SubCopy:           true,
		DataCopy:          true,
		AttachmentCopy:    true,
		RestoreFileTime:   false,

		OutputContainer: "mp4",

		NVEncCAdvanced: defaultAdvanced(),
		CustomOptions:  "",
	}
}

func cameraArchive() Profile {
	return Profile{
		ID:          uuid.New().String(),
		Version:     CurrentVersion,
		Name:        "Camera Archive",
		IsPreset:    true,
		EncoderType: "nvencc",
		EncoderOpts: map[string]any{},

		Codec:       "hevc",
		RateControl: "qvbr",
		RateValue:   24,
		Preset:      "P7",
		OutputDepth: 10,
		Multipass:   "none",
		OutputRes:   "",

		AQ:         true,
		AQTemporal: true,

		SplitEnc: "off",
		Parallel: "off",
		Decoder:  "avhw",
		Device:   "auto",

		AudioMode:    "copy",
		AudioBitrate: 256,

		Colormatrix: "auto",
		Transfer:    "auto",
		Colorprim:   "auto",
		Colorrange:  "auto",
		DHDR10Info:  "off",

		MetadataCopy:      true,
		VideoMetadataCopy: true,
		AudioMetadataCopy: true,
		ChapterCopy:       true,
		SubCopy:           true,
		DataCopy:          true,
		AttachmentCopy:    true,
		RestoreFileTime:   true,

		OutputContainer: "mp4",

		NVEncCAdvanced: defaultAdvanced(),
		CustomOptions:  "",
	}
}

func h264Compatible() Profile {
	return Profile{
		ID:          uuid.New().String(),
		Version:     CurrentVersion,
		Name:        "H.264 Compatible",
		IsPreset:    true,
		EncoderType: "nvencc",
		EncoderOpts: map[string]any{},

		Codec:       "h264",
		RateControl: "qvbr",
		RateValue:   26,
		Preset:      "P4",
		OutputDepth: 8,
		Multipass:   "none",
		OutputRes:   "",

		AQ:         true,
		AQTemporal: true,

		SplitEnc: "off",
		Parallel: "off",
		Decoder:  "avhw",
		Device:   "auto",

		AudioMode:    "aac",
		AudioBitrate: 256,

		Colormatrix: "auto",
		Transfer:    "auto",
		Colorprim:   "auto",
		Colorrange:  "auto",
		DHDR10Info:  "off",

		MetadataCopy:      true,
		VideoMetadataCopy: true,
		AudioMetadataCopy: true,
		ChapterCopy:       true,
		SubCopy:           true,
		DataCopy:          true,
		AttachmentCopy:    true,
		RestoreFileTime:   false,

		OutputContainer: "mp4",

		NVEncCAdvanced: defaultAdvanced(),
		CustomOptions:  "",
	}
}
