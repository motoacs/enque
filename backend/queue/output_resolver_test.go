package queue

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/motoacs/enque/backend/model"
)

func TestOutputResolverAutoRenameConcurrent(t *testing.T) {
	dir := t.TempDir()
	resolver := NewOutputResolver()
	cfg := model.DefaultAppConfig()
	cfg.OutputFolderMode = "specified"
	cfg.OutputFolderPath = dir
	cfg.OutputNameTemplate = "{name}.{ext}"
	cfg.OutputContainer = "mkv"
	cfg.OverwriteMode = model.OverwriteModeAutoRename

	input := filepath.Join(dir, "sample.mp4")
	if err := os.WriteFile(input, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sample.mkv"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	results := make(chan string, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := resolver.Resolve(input, cfg)
			if err != nil {
				t.Errorf("resolve failed: %v", err)
				return
			}
			results <- res.FinalOutputPath
		}()
	}
	wg.Wait()
	close(results)
	seen := map[string]struct{}{}
	for p := range results {
		if _, ok := seen[p]; ok {
			t.Fatalf("duplicate output path detected: %s", p)
		}
		seen[p] = struct{}{}
	}
}

func TestOutputResolverAsk(t *testing.T) {
	dir := t.TempDir()
	resolver := NewOutputResolver()
	cfg := model.DefaultAppConfig()
	cfg.OutputFolderMode = "specified"
	cfg.OutputFolderPath = dir
	cfg.OutputNameTemplate = "{name}.{ext}"
	cfg.OutputContainer = "mkv"
	cfg.OverwriteMode = model.OverwriteModeAsk

	input := filepath.Join(dir, "sample.mp4")
	_ = os.WriteFile(input, []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "sample.mkv"), []byte("x"), 0o644)

	res, err := resolver.Resolve(input, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !res.NeedsOverwrite {
		t.Fatalf("expected NeedsOverwrite")
	}
}
