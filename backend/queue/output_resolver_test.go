package queue

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestOutputResolver_BasicResolve(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "video.mp4")
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	r := NewOutputResolver()
	result, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedFinal := filepath.Join(tmpDir, "video_encoded.mkv")
	if result.FinalPath != expectedFinal {
		t.Fatalf("expected %s, got %s", expectedFinal, result.FinalPath)
	}
	if !containsSubstring(result.TempPath, ".tmp.mkv") {
		t.Fatalf("temp path should contain .tmp.mkv: %s", result.TempPath)
	}
	if result.NeedsOverwrite {
		t.Fatal("should not need overwrite")
	}
}

func TestOutputResolver_CustomFolder(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "input", "video.mp4")
	os.MkdirAll(filepath.Dir(inputPath), 0o755)
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	outDir := filepath.Join(tmpDir, "output")

	r := NewOutputResolver()
	result, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "specified",
		FolderPath:    outDir,
		NameTemplate:  "{name}.{ext}",
		Container:     "mp4",
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedFinal := filepath.Join(outDir, "video.mp4")
	if result.FinalPath != expectedFinal {
		t.Fatalf("expected %s, got %s", expectedFinal, result.FinalPath)
	}
}

func TestOutputResolver_AutoRename(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "video.mp4")
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	// Create existing file
	existingPath := filepath.Join(tmpDir, "video_encoded.mkv")
	os.WriteFile(existingPath, []byte("existing"), 0o644)

	r := NewOutputResolver()
	result, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(tmpDir, "video_encoded_001.mkv")
	if result.FinalPath != expected {
		t.Fatalf("expected %s, got %s", expected, result.FinalPath)
	}
}

func TestOutputResolver_AutoRenameSequential(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "video.mp4")
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	existingPath := filepath.Join(tmpDir, "video_encoded.mkv")
	os.WriteFile(existingPath, []byte("existing"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "video_encoded_001.mkv"), []byte("existing"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "video_encoded_002.mkv"), []byte("existing"), 0o644)

	r := NewOutputResolver()
	result, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(tmpDir, "video_encoded_003.mkv")
	if result.FinalPath != expected {
		t.Fatalf("expected %s, got %s", expected, result.FinalPath)
	}
}

func TestOutputResolver_SkipMode(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "video.mp4")
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	existingPath := filepath.Join(tmpDir, "video_encoded.mkv")
	os.WriteFile(existingPath, []byte("existing"), 0o644)

	r := NewOutputResolver()
	_, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "skip",
	})
	if err == nil {
		t.Fatal("expected error for skip mode with existing file")
	}
}

func TestOutputResolver_AskMode(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "video.mp4")
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	existingPath := filepath.Join(tmpDir, "video_encoded.mkv")
	os.WriteFile(existingPath, []byte("existing"), 0o644)

	r := NewOutputResolver()
	result, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "ask",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.NeedsOverwrite {
		t.Fatal("expected NeedsOverwrite to be true")
	}
}

func TestOutputResolver_ConcurrentResolve(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewOutputResolver()
	var wg sync.WaitGroup
	results := make([]string, 10)
	errs := make([]error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			inputPath := filepath.Join(tmpDir, "video.mp4")
			result, err := r.Resolve(inputPath, OutputConfig{
				FolderMode:    "same_as_input",
				NameTemplate:  "{name}_encoded.{ext}",
				Container:     "mkv",
				OverwriteMode: "auto_rename",
			})
			if err != nil {
				errs[idx] = err
				return
			}
			results[idx] = result.FinalPath
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("goroutine %d error: %v", i, err)
		}
	}

	// All final paths should be unique
	seen := make(map[string]bool)
	for _, path := range results {
		if path == "" {
			t.Fatal("empty path returned")
		}
		if seen[path] {
			t.Fatalf("duplicate path: %s", path)
		}
		seen[path] = true
	}
}

func TestOutputResolver_TempPathFormat(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "video.mp4")
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	r := NewOutputResolver()
	result, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Temp path should be: {dir}/{name}_encoded.{shortid}.tmp.mkv
	if !containsSubstring(result.TempPath, ".tmp.mkv") {
		t.Fatalf("expected .tmp.mkv in temp path: %s", result.TempPath)
	}
	if !containsSubstring(result.TempPath, "video_encoded.") {
		t.Fatalf("expected video_encoded. in temp path: %s", result.TempPath)
	}
}

func TestOutputResolver_Release(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "video.mp4")
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	r := NewOutputResolver()

	// First resolve
	result1, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Second resolve should get a different path (auto_rename)
	result2, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result1.FinalPath == result2.FinalPath {
		t.Fatal("second resolve should produce different path")
	}

	// Release the first path
	r.Release(result1.FinalPath)

	// Third resolve should reuse the released path
	result3, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_encoded.{ext}",
		Container:     "mkv",
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result3.FinalPath != result1.FinalPath {
		t.Fatalf("expected released path %s to be reused, got %s", result1.FinalPath, result3.FinalPath)
	}
}

func TestOutputResolver_EmptyContainerFallback(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "video.mp4")
	os.WriteFile(inputPath, []byte("dummy"), 0o644)

	r := NewOutputResolver()
	result, err := r.Resolve(inputPath, OutputConfig{
		FolderMode:    "same_as_input",
		NameTemplate:  "{name}_out.{ext}",
		Container:     "", // Should fallback to input extension
		OverwriteMode: "auto_rename",
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedFinal := filepath.Join(tmpDir, "video_out.mp4")
	if result.FinalPath != expectedFinal {
		t.Fatalf("expected %s, got %s", expectedFinal, result.FinalPath)
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
