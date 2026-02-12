package nvencc

import (
	"strings"
	"testing"

	"github.com/yuta/enque/backend/profile"
)

func defaultProfile() profile.Profile {
	return profile.Profile{
		EncoderType:  "nvencc",
		EncoderOpts:  map[string]any{},
		Codec:        "hevc",
		RateControl:  "qvbr",
		RateValue:    28,
		Preset:       "P4",
		OutputDepth:  10,
		Multipass:    "none",
		OutputRes:    "",
		AQ:           true,
		AQTemporal:   true,
		SplitEnc:     "auto",
		Parallel:     "off",
		Decoder:      "avhw",
		Device:       "auto",
		AudioMode:    "copy",
		AudioBitrate: 256,
		Colormatrix:  "auto",
		Transfer:     "auto",
		Colorprim:    "auto",
		Colorrange:   "auto",
		DHDR10Info:   "off",
		MetadataCopy:      true,
		VideoMetadataCopy: true,
		AudioMetadataCopy: true,
		ChapterCopy:       true,
		SubCopy:           true,
		DataCopy:          true,
		AttachmentCopy:    true,
		RestoreFileTime:   false,
		NVEncCAdvanced:    profile.NVEncCAdvanced{},
		CustomOptions:     "",
	}
}

func argsString(args []string) string {
	return strings.Join(args, " ")
}

func TestBuildArgs_DefaultProfile(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	args, err := a.BuildArgs(p, `C:\input.mp4`, `C:\output.mkv`)
	if err != nil {
		t.Fatal(err)
	}

	s := argsString(args)

	// Check argument order
	avhwIdx := strings.Index(s, "--avhw")
	inputIdx := strings.Index(s, "-i")
	codecIdx := strings.Index(s, "-c hevc")
	outputIdx := strings.LastIndex(s, "-o")

	if avhwIdx >= inputIdx {
		t.Error("--avhw must come before -i")
	}
	if inputIdx >= codecIdx {
		t.Error("-i must come before -c")
	}
	if codecIdx >= outputIdx {
		t.Error("-c must come before -o")
	}

	// Check key arguments exist
	mustContain := []string{
		"--avhw",
		"-i", `C:\input.mp4`,
		"-c", "hevc",
		"--qvbr", "28",
		"--preset", "P4",
		"--output-depth", "10",
		"--aq",
		"--aq-temporal",
		"--split-enc", "auto",
		"--audio-copy",
		"--metadata", "copy",
		"--video-metadata", "copy",
		"--audio-metadata", "copy",
		"--chapter-copy",
		"--sub-copy",
		"--data-copy",
		"--attachment-copy",
		"-o", `C:\output.mkv`,
	}
	for _, arg := range mustContain {
		if !strings.Contains(s, arg) {
			t.Errorf("missing arg: %s in %s", arg, s)
		}
	}
}

func TestBuildArgs_AllCodecs(t *testing.T) {
	a := &NVEncCAdapter{}
	for _, codec := range []string{"h264", "hevc", "av1"} {
		p := defaultProfile()
		p.Codec = codec
		args, err := a.BuildArgs(p, "input.mp4", "output.mkv")
		if err != nil {
			t.Fatalf("codec %s: %v", codec, err)
		}
		s := argsString(args)
		if !strings.Contains(s, "-c "+codec) {
			t.Errorf("codec %s: expected -c %s in %s", codec, codec, s)
		}
	}
}

func TestBuildArgs_AllRateControls(t *testing.T) {
	a := &NVEncCAdapter{}
	tests := []struct {
		rc   string
		flag string
	}{
		{"qvbr", "--qvbr"},
		{"cqp", "--cqp"},
		{"cbr", "--cbr"},
		{"vbr", "--vbr"},
	}
	for _, tt := range tests {
		p := defaultProfile()
		p.RateControl = tt.rc
		p.RateValue = 28
		args, err := a.BuildArgs(p, "in.mp4", "out.mkv")
		if err != nil {
			t.Fatal(err)
		}
		s := argsString(args)
		if !strings.Contains(s, tt.flag+" 28") {
			t.Errorf("rc=%s: expected %s 28 in %s", tt.rc, tt.flag, s)
		}
	}
}

func TestBuildArgs_Decoder_AVSW(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.Decoder = "avsw"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	if !strings.Contains(s, "--avsw") {
		t.Errorf("expected --avsw in %s", s)
	}
	if strings.Contains(s, "--avhw") {
		t.Errorf("should not contain --avhw in %s", s)
	}
}

func TestBuildArgs_Decoder_AVSW_WithDecoder(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.Decoder = "avsw"
	p.NVEncCAdvanced.AVSWDecoder = "h264_cuvid"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	if !strings.Contains(s, "--avsw h264_cuvid") {
		t.Errorf("expected '--avsw h264_cuvid' in %s", s)
	}
}

func TestBuildArgs_NullFields_Omitted(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	// bframes, ref, lookahead, gop_len are nil by default
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	for _, opt := range []string{"--bframes", "--ref", "--lookahead", "--gop-len"} {
		if strings.Contains(s, opt) {
			t.Errorf("nil field should not produce %s in %s", opt, s)
		}
	}
}

func TestBuildArgs_NullFields_Present(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	bf := 3
	ref := 4
	la := 16
	gop := 300
	p.Bframes = &bf
	p.Ref = &ref
	p.Lookahead = &la
	p.GopLen = &gop
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	for _, expected := range []string{"--bframes 3", "--ref 4", "--lookahead 16", "--gop-len 300"} {
		if !strings.Contains(s, expected) {
			t.Errorf("expected %q in %s", expected, s)
		}
	}
}

func TestBuildArgs_AudioModes(t *testing.T) {
	a := &NVEncCAdapter{}
	tests := []struct {
		mode     string
		expected string
	}{
		{"copy", "--audio-copy"},
		{"aac", "--audio-codec aac --audio-bitrate 256"},
		{"opus", "--audio-codec opus --audio-bitrate 256"},
	}
	for _, tt := range tests {
		p := defaultProfile()
		p.AudioMode = tt.mode
		args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
		s := argsString(args)
		if !strings.Contains(s, tt.expected) {
			t.Errorf("audio_mode=%s: expected %q in %s", tt.mode, tt.expected, s)
		}
	}
}

func TestBuildArgs_Color_NonAuto(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.Colormatrix = "bt709"
	p.Transfer = "smpte2084"
	p.Colorprim = "bt2020"
	p.Colorrange = "full"
	p.DHDR10Info = "copy"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	for _, expected := range []string{
		"--colormatrix bt709",
		"--transfer smpte2084",
		"--colorprim bt2020",
		"--colorrange full",
		"--dhdr10-info copy",
	} {
		if !strings.Contains(s, expected) {
			t.Errorf("expected %q in %s", expected, s)
		}
	}
}

func TestBuildArgs_MetadataOff(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.MetadataCopy = false
	p.VideoMetadataCopy = false
	p.AudioMetadataCopy = false
	p.ChapterCopy = false
	p.SubCopy = false
	p.DataCopy = false
	p.AttachmentCopy = false
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	for _, opt := range []string{"--metadata", "--video-metadata", "--audio-metadata", "--chapter-copy", "--sub-copy", "--data-copy", "--attachment-copy"} {
		if strings.Contains(s, opt) {
			t.Errorf("metadata off: should not contain %s in %s", opt, s)
		}
	}
}

func TestBuildArgs_Multipass(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.Multipass = "full"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	if !strings.Contains(s, "--multipass full") {
		t.Errorf("expected --multipass full in %s", s)
	}
}

func TestBuildArgs_SplitEnc_Off(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.SplitEnc = "off"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	if strings.Contains(s, "--split-enc") {
		t.Errorf("split_enc=off should not produce --split-enc in %s", s)
	}
}

func TestBuildArgs_Device_NonAuto(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.Device = "1"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	if !strings.Contains(s, "--device 1") {
		t.Errorf("expected --device 1 in %s", s)
	}
}

func TestBuildArgs_Advanced_Overrides(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	mb := 50000
	p.NVEncCAdvanced.MaxBitrate = &mb
	p.NVEncCAdvanced.WeightP = true
	p.NVEncCAdvanced.Level = "5.1"
	p.NVEncCAdvanced.Profile = "main10"
	p.NVEncCAdvanced.Tier = "high"
	p.NVEncCAdvanced.SSIM = true
	p.NVEncCAdvanced.PSNR = true
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	for _, expected := range []string{
		"--max-bitrate 50000",
		"--weightp",
		"--level 5.1",
		"--profile main10",
		"--tier high",
		"--ssim",
		"--psnr",
	} {
		if !strings.Contains(s, expected) {
			t.Errorf("expected %q in %s", expected, s)
		}
	}
}

func TestBuildArgs_Advanced_LaterWins(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	// Standard GUI sets metadata copy, advanced overrides with different value
	p.MetadataCopy = true
	p.NVEncCAdvanced.Metadata = "title=test"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	// Standard GUI emits --metadata copy, then advanced emits --metadata title=test
	// NVEncC uses later-wins, so title=test takes effect
	stdIdx := strings.Index(s, "--metadata copy")
	advIdx := strings.LastIndex(s, "--metadata title=test")
	if stdIdx < 0 || advIdx < 0 {
		t.Fatalf("expected both --metadata copy and --metadata title=test in %s", s)
	}
	if advIdx <= stdIdx {
		t.Errorf("advanced --metadata must appear after standard --metadata for later-wins")
	}
}

func TestBuildArgs_CustomOptions_OverrideAll(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.CustomOptions = "--gop-len 300 --vpp-nlmeans sigma=0.005"
	args, err := a.BuildArgs(p, "in.mp4", "out.mkv")
	if err != nil {
		t.Fatal(err)
	}
	s := argsString(args)
	customIdx := strings.Index(s, "--gop-len 300")
	outputIdx := strings.LastIndex(s, "-o")
	if customIdx >= outputIdx {
		t.Errorf("custom options should appear before -o")
	}
	if !strings.Contains(s, "--vpp-nlmeans") {
		t.Errorf("expected --vpp-nlmeans in %s", s)
	}
}

func TestBuildArgs_CustomOptions_InvalidQuote(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.CustomOptions = `--opt "unclosed`
	_, err := a.BuildArgs(p, "in.mp4", "out.mkv")
	if err == nil {
		t.Error("expected error for unclosed quote")
	}
}

func TestBuildArgs_OutputRes(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.OutputRes = "1920x1080,preserve_aspect_ratio=decrease"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	if !strings.Contains(s, "--output-res 1920x1080,preserve_aspect_ratio=decrease") {
		t.Errorf("expected --output-res in %s", s)
	}
}

func TestBuildArgs_EmptyOutputRes_Omitted(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.OutputRes = ""
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	if strings.Contains(s, "--output-res") {
		t.Errorf("empty output_res should not produce --output-res in %s", s)
	}
}

func TestBuildArgs_ArgumentOrder_FullPipeline(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	bf := 3
	p.Bframes = &bf
	p.Multipass = "quarter"
	p.SplitEnc = "forced_3"
	p.Colormatrix = "bt709"
	mb := 50000
	p.NVEncCAdvanced.MaxBitrate = &mb
	p.CustomOptions = "--vpp-nlmeans sigma=0.005"

	args, err := a.BuildArgs(p, "in.mp4", "out.mkv")
	if err != nil {
		t.Fatal(err)
	}

	// Find position of each argument in the args slice
	indexOf := func(target string) int {
		for i, a := range args {
			if a == target {
				return i
			}
		}
		return -1
	}

	order := []struct {
		arg  string
		name string
	}{
		{"--avhw", "decoder"},
		{"-i", "input"},
		{"-c", "codec"},
		{"--multipass", "multipass"},
		{"--bframes", "bframes"},
		{"--split-enc", "split-enc"},
		{"--audio-copy", "audio"},
		{"--max-bitrate", "advanced"},
		{"--vpp-nlmeans", "custom"},
		{"-o", "output"},
	}

	for i := 1; i < len(order); i++ {
		prevIdx := indexOf(order[i-1].arg)
		currIdx := indexOf(order[i].arg)
		if prevIdx < 0 {
			t.Fatalf("missing %s in args", order[i-1].arg)
		}
		if currIdx < 0 {
			t.Fatalf("missing %s in args", order[i].arg)
		}
		if currIdx <= prevIdx {
			t.Errorf("%s (%s, idx %d) should come after %s (%s, idx %d)",
				order[i].arg, order[i].name, currIdx,
				order[i-1].arg, order[i-1].name, prevIdx)
		}
	}
}

func TestBuildArgs_Parallel(t *testing.T) {
	a := &NVEncCAdapter{}
	for _, mode := range []string{"auto", "2", "3"} {
		p := defaultProfile()
		p.Parallel = mode
		args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
		s := argsString(args)
		if !strings.Contains(s, "--parallel "+mode) {
			t.Errorf("parallel=%s: expected --parallel %s in %s", mode, mode, s)
		}
	}
}

func TestBuildArgs_AdvancedTrimSeek(t *testing.T) {
	a := &NVEncCAdapter{}
	p := defaultProfile()
	p.NVEncCAdvanced.Trim = "0:100,200:300"
	p.NVEncCAdvanced.Seek = "10.0"
	p.NVEncCAdvanced.SeekTo = "30.0"
	args, _ := a.BuildArgs(p, "in.mp4", "out.mkv")
	s := argsString(args)
	if !strings.Contains(s, "--trim 0:100,200:300") {
		t.Errorf("expected --trim in %s", s)
	}
	if !strings.Contains(s, "--seek 10.0") {
		t.Errorf("expected --seek in %s", s)
	}
	if !strings.Contains(s, "--seekto 30.0") {
		t.Errorf("expected --seekto in %s", s)
	}
}
