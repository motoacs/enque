package queue

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// OutputResolver determines output paths for encoding jobs.
// All path operations are serialized under a mutex to prevent collisions
// when multiple workers resolve paths concurrently.
type OutputResolver struct {
	mu           sync.Mutex
	reservedPaths map[string]bool // Track paths reserved in this session
}

// NewOutputResolver creates a new output resolver.
func NewOutputResolver() *OutputResolver {
	return &OutputResolver{
		reservedPaths: make(map[string]bool),
	}
}

// ResolveResult holds the resolved output paths for a job.
type ResolveResult struct {
	TempPath      string
	FinalPath     string
	NeedsOverwrite bool // True if final path exists and overwrite_mode is "ask"
}

// Resolve determines the temp and final output paths for a job.
// Must be called under the resolver's internal mutex.
func (r *OutputResolver) Resolve(inputPath string, cfg OutputConfig) (*ResolveResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 1. Determine output directory
	outputDir, err := r.resolveOutputDir(inputPath, cfg)
	if err != nil {
		return nil, err
	}

	// 2. Apply template to get output filename
	outputName := r.applyTemplate(inputPath, cfg.NameTemplate, cfg.Container)

	// 3. Build final path
	finalPath := filepath.Join(outputDir, outputName)

	// 4. Handle overwrite mode
	needsOverwrite := false
	if fileExists(finalPath) || r.reservedPaths[finalPath] {
		switch cfg.OverwriteMode {
		case "overwrite":
			// Allow overwrite
		case "auto_rename":
			finalPath = r.autoRename(finalPath)
		case "skip":
			return nil, fmt.Errorf("output file already exists (skip mode): %s", finalPath)
		case "ask":
			if fileExists(finalPath) {
				needsOverwrite = true
			} else if r.reservedPaths[finalPath] {
				finalPath = r.autoRename(finalPath)
			}
		default:
			// Default to ask
			if fileExists(finalPath) {
				needsOverwrite = true
			}
		}
	}

	// 5. Generate temp file path
	shortID := generateShortID()
	tempPath := r.buildTempPath(finalPath, shortID)

	// 6. Reserve the final path
	r.reservedPaths[finalPath] = true

	return &ResolveResult{
		TempPath:       tempPath,
		FinalPath:      finalPath,
		NeedsOverwrite: needsOverwrite,
	}, nil
}

// Release removes a path from the reserved set (e.g., on job failure/skip).
func (r *OutputResolver) Release(finalPath string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.reservedPaths, finalPath)
}

// OutputConfig holds output configuration for path resolution.
type OutputConfig struct {
	FolderMode   string // "same_as_input" or "specified"
	FolderPath   string // Custom output folder path
	NameTemplate string // e.g., "{name}_encoded.{ext}"
	Container    string // e.g., "mkv", "mp4"
	OverwriteMode string // "overwrite", "auto_rename", "skip", "ask"
}

func (r *OutputResolver) resolveOutputDir(inputPath string, cfg OutputConfig) (string, error) {
	var dir string
	switch cfg.FolderMode {
	case "specified":
		if cfg.FolderPath == "" {
			return "", fmt.Errorf("custom output folder path is empty")
		}
		dir = cfg.FolderPath
	default: // "same_as_input"
		dir = filepath.Dir(inputPath)
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	return dir, nil
}

func (r *OutputResolver) applyTemplate(inputPath, template, container string) string {
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	if template == "" {
		template = "{name}_encoded.{ext}"
	}

	if container == "" {
		container = strings.TrimPrefix(ext, ".")
	}

	result := template
	result = strings.ReplaceAll(result, "{name}", name)
	result = strings.ReplaceAll(result, "{ext}", container)

	return result
}

// autoRename appends _001, _002, etc. to avoid collisions.
func (r *OutputResolver) autoRename(path string) string {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)

	for i := 1; i < 10000; i++ {
		candidate := fmt.Sprintf("%s_%03d%s", base, i, ext)
		if !fileExists(candidate) && !r.reservedPaths[candidate] {
			return candidate
		}
	}

	// Fallback with random suffix
	shortID := generateShortID()
	return fmt.Sprintf("%s_%s%s", base, shortID, ext)
}

// buildTempPath creates: {name}.{short_id}.tmp.{ext}
func (r *OutputResolver) buildTempPath(finalPath, shortID string) string {
	ext := filepath.Ext(finalPath)
	base := strings.TrimSuffix(finalPath, ext)
	return fmt.Sprintf("%s.%s.tmp%s", base, shortID, ext)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func generateShortID() string {
	b := make([]byte, 4) // 4 bytes = 8 hex chars
	rand.Read(b)
	return hex.EncodeToString(b)
}
