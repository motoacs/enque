package nvencc

import (
	"slices"
	"testing"

	"github.com/motoacs/enque/backend/encoder"
	"github.com/motoacs/enque/backend/model"
)

func TestBuildCommand_OrderAndOverrides(t *testing.T) {
	p := model.DefaultProfile()
	p.ID = "p1"
	p.Name = "test"
	p.Codec = model.CodecHEVC
	p.NVEncCAdvanced.Metadata = "delete"
	p.CustomOptions = `--metadata keep --foo "bar baz"`
	res, err := BuildCommand(encoder.BuildRequest{
		Profile:    p,
		AppConfig:  model.DefaultAppConfig(),
		InputPath:  `C:\in video.mp4`,
		OutputPath: `C:\out.mkv`,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := EnsureOptionOrder(res.Argv); err != nil {
		t.Fatalf("invalid option order: %v", err)
	}
	iMetaDel := slices.Index(res.Argv, "delete")
	iMetaKeep := slices.Index(res.Argv, "keep")
	if iMetaDel == -1 || iMetaKeep == -1 || iMetaKeep <= iMetaDel {
		t.Fatalf("expected custom option to be later than advanced option: %#v", res.Argv)
	}
	if !slices.Contains(res.Argv, "bar baz") {
		t.Fatalf("expected quoted token from custom options")
	}
}

func TestBuildCommand_RetryDecoder(t *testing.T) {
	a := NewAdapter()
	p := model.DefaultProfile()
	p.ID = "p1"
	p.Name = "x"
	p.Decoder = model.DecoderAVHW
	build, err := a.BuildArgs(encoder.BuildRequest{Profile: p, InputPath: "in.mp4", OutputPath: "out.mkv"})
	if err != nil {
		t.Fatal(err)
	}
	retry, ok, err := a.BuildRetryArgs(encoder.BuildRequest{Profile: p, InputPath: "in.mp4", OutputPath: "out.mkv"}, build)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || retry == nil {
		t.Fatalf("expected retry args")
	}
	if retry.EffectiveDecoder != model.DecoderAVSW {
		t.Fatalf("expected avsw retry, got %s", retry.EffectiveDecoder)
	}
}
