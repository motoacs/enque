package profile

import "fmt"

// Migrate upgrades a profile from any version to CurrentVersion.
func Migrate(p Profile) (Profile, error) {
	for p.Version < CurrentVersion {
		var err error
		switch p.Version {
		case 0, 1:
			p, err = migrateV1toV2(p)
		case 2:
			p, err = migrateV2toV3(p)
		case 3:
			p, err = migrateV3toV4(p)
		default:
			return p, fmt.Errorf("unknown profile version %d", p.Version)
		}
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

// migrateV1toV2: add encoder_type and encoder_options.
func migrateV1toV2(p Profile) (Profile, error) {
	if p.EncoderType == "" {
		p.EncoderType = "nvencc"
	}
	if p.EncoderOpts == nil {
		p.EncoderOpts = map[string]any{}
	}
	p.Version = 2
	return p, nil
}

// migrateV2toV3: add nvencc_advanced with defaults.
func migrateV2toV3(p Profile) (Profile, error) {
	// NVEncCAdvanced zero value is already correct defaults (empty strings, nil pointers, false bools)
	p.Version = 3
	return p, nil
}

// migrateV3toV4: remove max_cll/master_display/dolby_vision_rpu, add avsw_decoder.
func migrateV3toV4(p Profile) (Profile, error) {
	// avsw_decoder defaults to empty string (zero value)
	p.Version = 4
	return p, nil
}
