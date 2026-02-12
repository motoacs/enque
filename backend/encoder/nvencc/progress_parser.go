package nvencc

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/motoacs/enque/backend/model"
)

var (
	rePercent = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)%`)
	reFPS     = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)\s*fps`)
	reBitrate = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)\s*(k|m)?bps`)
	reETAHMS  = regexp.MustCompile(`(?:eta|remain(?:ing)?)[^0-9]*([0-9]{1,2}):([0-9]{2}):([0-9]{2})`)
)

func ParseProgress(line string) (model.JobProgress, bool) {
	lower := strings.ToLower(line)
	progress := model.JobProgress{RawLine: line}
	matched := false

	if m := rePercent.FindStringSubmatch(lower); len(m) == 2 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			progress.Percent = &v
			matched = true
		}
	}
	if m := reFPS.FindStringSubmatch(lower); len(m) == 2 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			progress.FPS = &v
			matched = true
		}
	}
	if m := reBitrate.FindStringSubmatch(lower); len(m) == 3 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			unit := m[2]
			if unit == "m" {
				v *= 1000
			}
			progress.BitrateKbps = &v
			matched = true
		}
	}
	if m := reETAHMS.FindStringSubmatch(lower); len(m) == 4 {
		h, _ := strconv.ParseInt(m[1], 10, 64)
		mi, _ := strconv.ParseInt(m[2], 10, 64)
		s, _ := strconv.ParseInt(m[3], 10, 64)
		eta := h*3600 + mi*60 + s
		progress.ETASec = &eta
		matched = true
	}
	return progress, matched
}
