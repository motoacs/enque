package encoder

import (
	"strings"
	"testing"

	"github.com/yuta/enque/backend/profile"
)

type mockAdapter struct {
	typ string
}

func (m *mockAdapter) Type() string { return m.typ }
func (m *mockAdapter) BuildArgs(p profile.Profile, input, output string) ([]string, error) {
	return nil, nil
}
func (m *mockAdapter) ParseProgress(line string) Progress  { return Progress{} }
func (m *mockAdapter) SupportsDecoderFallback() bool        { return false }

func TestRegistry_Resolve(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockAdapter{typ: "nvencc"})

	adapter, err := r.Resolve("nvencc")
	if err != nil {
		t.Fatal(err)
	}
	if adapter.Type() != "nvencc" {
		t.Errorf("type=%q, want nvencc", adapter.Type())
	}
}

func TestRegistry_Resolve_NotImplemented(t *testing.T) {
	r := NewRegistry()
	r.Register(&mockAdapter{typ: "nvencc"})

	_, err := r.Resolve("qsvenc")
	if err == nil {
		t.Error("expected error for unregistered adapter")
	}
	if !strings.Contains(err.Error(), ErrEncoderNotImplemented) {
		t.Errorf("error=%q, want to contain %q", err.Error(), ErrEncoderNotImplemented)
	}
}

func TestRegistry_Resolve_FFmpeg_NotImplemented(t *testing.T) {
	r := NewRegistry()
	_, err := r.Resolve("ffmpeg")
	if err == nil {
		t.Error("expected error for unregistered adapter")
	}
	if !strings.Contains(err.Error(), ErrEncoderNotImplemented) {
		t.Errorf("error=%q, want to contain %q", err.Error(), ErrEncoderNotImplemented)
	}
}
