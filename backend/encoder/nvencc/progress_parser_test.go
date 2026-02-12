package nvencc

import "testing"

func TestParseProgress(t *testing.T) {
	line := "42.3% 123.4 fps 5.6 Mbps ETA 00:01:12"
	p, ok := ParseProgress(line)
	if !ok {
		t.Fatalf("expected parse success")
	}
	if p.Percent == nil || *p.Percent != 42.3 {
		t.Fatalf("percent mismatch: %#v", p.Percent)
	}
	if p.FPS == nil || *p.FPS != 123.4 {
		t.Fatalf("fps mismatch")
	}
	if p.BitrateKbps == nil || *p.BitrateKbps != 5600 {
		t.Fatalf("bitrate mismatch: %#v", p.BitrateKbps)
	}
	if p.ETASec == nil || *p.ETASec != 72 {
		t.Fatalf("eta mismatch: %#v", p.ETASec)
	}
}

func TestParseProgressFailure(t *testing.T) {
	p, ok := ParseProgress("this line has no progress")
	if ok {
		t.Fatalf("expected parse failure")
	}
	if p.Percent != nil {
		t.Fatalf("percent should be nil")
	}
}
