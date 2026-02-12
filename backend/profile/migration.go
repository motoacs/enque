package profile

import "github.com/motoacs/enque/backend/model"

func migrateOne(p model.Profile) (model.Profile, bool) {
	changed := false
	if p.Version < 2 {
		if p.EncoderType == "" {
			p.EncoderType = model.EncoderTypeNVEncC
		}
		if p.EncoderOptions == nil {
			p.EncoderOptions = map[string]any{}
		}
		p.Version = 2
		changed = true
	}
	if p.Version < 3 {
		// v3 adds nvencc_advanced with zero defaults
		p.Version = 3
		changed = true
	}
	if p.Version < 4 {
		if p.NVEncCAdvanced.AVSWDecoder == "" {
			p.NVEncCAdvanced.AVSWDecoder = ""
		}
		p.Version = 4
		changed = true
	}
	if p.Version != model.ProfileVersion {
		p.Version = model.ProfileVersion
		changed = true
	}
	return p, changed
}
