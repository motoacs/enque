package queue

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/motoacs/enque/backend/model"
)

type ResolveResult struct {
	FinalOutputPath string
	TempOutputPath  string
	NeedsOverwrite  bool
}

type OutputResolver struct {
	mu       sync.Mutex
	reserved map[string]struct{}
}

func NewOutputResolver() *OutputResolver {
	return &OutputResolver{reserved: map[string]struct{}{}}
}

func (r *OutputResolver) Resolve(inputPath string, cfg model.AppConfig) (ResolveResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	dir := filepath.Dir(inputPath)
	if cfg.OutputFolderMode == "specified" && strings.TrimSpace(cfg.OutputFolderPath) != "" {
		dir = cfg.OutputFolderPath
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ResolveResult{}, err
	}
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	ext := strings.TrimPrefix(cfg.OutputContainer, ".")
	if ext == "" {
		ext = "mkv"
	}
	name := applyTemplate(cfg.OutputNameTemplate, baseName, ext)
	finalPath := filepath.Join(dir, name)

	if cfg.OverwriteMode == model.OverwriteModeAutoRename {
		finalPath = r.uniqueFinalPath(finalPath)
	} else {
		if _, ok := r.reserved[finalPath]; ok {
			return ResolveResult{FinalOutputPath: finalPath, NeedsOverwrite: true}, nil
		}
		if _, err := os.Stat(finalPath); err == nil {
			return ResolveResult{FinalOutputPath: finalPath, NeedsOverwrite: true}, nil
		}
	}

	r.reserved[finalPath] = struct{}{}
	tempPath := filepath.Join(dir, fmt.Sprintf("%s.%s.tmp.%s", strings.TrimSuffix(filepath.Base(finalPath), filepath.Ext(finalPath)), shortID(8), ext))
	return ResolveResult{FinalOutputPath: finalPath, TempOutputPath: tempPath}, nil
}

func (r *OutputResolver) Release(finalPath string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.reserved, finalPath)
}

func (r *OutputResolver) uniqueFinalPath(path string) string {
	if !r.exists(path) {
		return path
	}
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(filepath.Base(path), ext)
	for i := 1; ; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s_%03d%s", base, i, ext))
		if !r.exists(candidate) {
			return candidate
		}
	}
}

func (r *OutputResolver) exists(path string) bool {
	if _, ok := r.reserved[path]; ok {
		return true
	}
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func applyTemplate(tpl, name, ext string) string {
	if tpl == "" {
		tpl = "{name}_encoded.{ext}"
	}
	repl := strings.ReplaceAll(tpl, "{name}", name)
	repl = strings.ReplaceAll(repl, "{ext}", ext)
	return repl
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func shortID(n int) string {
	buf := make([]byte, n)
	_, _ = rand.Read(buf)
	out := make([]byte, n)
	for i := range buf {
		out[i] = charset[int(buf[i])%len(charset)]
	}
	return string(out)
}
