package profile

import (
	"github.com/google/uuid"
	"github.com/motoacs/enque/backend/model"
)

func BuiltInPresets() []model.Profile {
	base := model.DefaultProfile()

	hevc := base
	hevc.ID = "preset-hevc-quality"
	hevc.Name = "HEVC Quality"
	hevc.IsPreset = true
	hevc.Codec = model.CodecHEVC
	hevc.RateControl = model.RateControlQVBR
	hevc.RateValue = 28
	hevc.Preset = "P4"
	hevc.OutputDepth = 10
	hevc.SplitEnc = model.SplitEncAuto

	av1 := base
	av1.ID = "preset-av1-fast"
	av1.Name = "AV1 Fast"
	av1.IsPreset = true
	av1.Codec = model.CodecAV1
	av1.RateValue = 32
	av1.Preset = "P1"
	av1.OutputDepth = 10
	av1.AudioMode = model.AudioModeCopy

	camera := base
	camera.ID = "preset-camera-archive"
	camera.Name = "Camera Archive"
	camera.IsPreset = true
	camera.Codec = model.CodecHEVC
	camera.RateValue = 24
	camera.Preset = "P7"
	camera.RestoreFileTime = true
	camera.AudioMode = model.AudioModeCopy
	camera.MetadataCopy = true
	camera.VideoMetadataCopy = true
	camera.AudioMetadataCopy = true
	camera.ChapterCopy = true
	camera.SubCopy = true
	camera.DataCopy = true
	camera.AttachmentCopy = true

	h264 := base
	h264.ID = "preset-h264-compatible"
	h264.Name = "H.264 Compatible"
	h264.IsPreset = true
	h264.Codec = model.CodecH264
	h264.RateValue = 26
	h264.Preset = "P4"
	h264.OutputDepth = 8
	h264.AudioMode = model.AudioModeAAC
	h264.AudioBitrate = 256

	return []model.Profile{hevc, av1, camera, h264}
}

func NewUserProfile(name string) model.Profile {
	p := model.DefaultProfile()
	p.ID = uuid.NewString()
	p.Name = name
	p.IsPreset = false
	return p
}
