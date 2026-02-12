package app

import (
	"testing"

	"github.com/motoacs/enque/backend/model"
)

func TestResolveEncoderPath(t *testing.T) {
	tools := model.ToolSnapshot{NVEncC: model.ToolInfo{Found: true, Path: "nvencc", Version: "NVEncC 8.0"}}
	p, err := resolveEncoderPath(model.EncoderTypeNVEncC, tools)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if p != "nvencc" {
		t.Fatalf("unexpected path: %s", p)
	}
}
