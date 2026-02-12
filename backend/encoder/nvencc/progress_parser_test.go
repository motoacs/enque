package nvencc

import (
	"testing"
)

func TestParseProgress_Normal(t *testing.T) {
	a := &NVEncCAdapter{}

	line := "[53.2%] 1234 frames: 245.67 fps, 12345 kb/s, remain 0:01:23, GPU 45%, VE 78%, VD 12%"
	prog := a.ParseProgress(line)

	if prog.Percent == nil || *prog.Percent != 53.2 {
		t.Errorf("percent=%v, want 53.2", prog.Percent)
	}
	if prog.FPS == nil || *prog.FPS != 245.67 {
		t.Errorf("fps=%v, want 245.67", prog.FPS)
	}
	if prog.BitrateKbps == nil || *prog.BitrateKbps != 12345 {
		t.Errorf("bitrate=%v, want 12345", prog.BitrateKbps)
	}
	if prog.ETASec == nil || *prog.ETASec != 83 {
		t.Errorf("eta=%v, want 83", prog.ETASec)
	}
}

func TestParseProgress_100Percent(t *testing.T) {
	a := &NVEncCAdapter{}

	line := "[100.0%] 5000 frames: 300.00 fps, 8765 kb/s, remain 0:00:00"
	prog := a.ParseProgress(line)

	if prog.Percent == nil || *prog.Percent != 100.0 {
		t.Errorf("percent=%v, want 100.0", prog.Percent)
	}
	if prog.ETASec == nil || *prog.ETASec != 0 {
		t.Errorf("eta=%v, want 0", prog.ETASec)
	}
}

func TestParseProgress_NoRemain(t *testing.T) {
	a := &NVEncCAdapter{}

	line := "[10.5%] 100 frames: 50.00 fps, 5000 kb/s"
	prog := a.ParseProgress(line)

	if prog.Percent == nil || *prog.Percent != 10.5 {
		t.Errorf("percent=%v, want 10.5", prog.Percent)
	}
	if prog.ETASec != nil {
		t.Errorf("eta should be nil when no remain, got %v", prog.ETASec)
	}
}

func TestParseProgress_ParseFail(t *testing.T) {
	a := &NVEncCAdapter{}

	line := "NVEncC (x64) 8.05 (r2994) by rigaya"
	prog := a.ParseProgress(line)

	if prog.Percent != nil {
		t.Errorf("percent should be nil for non-progress line, got %v", prog.Percent)
	}
	if prog.RawLine != line {
		t.Errorf("raw_line=%q, want %q", prog.RawLine, line)
	}
}

func TestParseProgress_Empty(t *testing.T) {
	a := &NVEncCAdapter{}
	prog := a.ParseProgress("")
	if prog.Percent != nil {
		t.Errorf("empty line should have nil percent")
	}
}

func TestParseProgress_HighBitrate(t *testing.T) {
	a := &NVEncCAdapter{}
	line := "[75.3%] 3750 frames: 500.12 fps, 123456 kb/s, remain 0:00:30"
	prog := a.ParseProgress(line)

	if prog.BitrateKbps == nil || *prog.BitrateKbps != 123456 {
		t.Errorf("bitrate=%v, want 123456", prog.BitrateKbps)
	}
}

func TestParseProgress_LongETA(t *testing.T) {
	a := &NVEncCAdapter{}
	line := "[1.0%] 50 frames: 10.00 fps, 5000 kb/s, remain 2:30:45"
	prog := a.ParseProgress(line)

	expected := 2*3600.0 + 30*60.0 + 45.0
	if prog.ETASec == nil || *prog.ETASec != expected {
		t.Errorf("eta=%v, want %v", prog.ETASec, expected)
	}
}
