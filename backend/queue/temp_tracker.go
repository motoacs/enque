package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// TempEntry records a temp file created during encoding.
type TempEntry struct {
	TempPath  string `json:"temp_path"`
	FinalPath string `json:"final_path"`
	JobID     string `json:"job_id"`
	SessionID string `json:"session_id"`
}

// TempTracker manages runtime/temp_index.json for crash recovery.
type TempTracker struct {
	mu      sync.Mutex
	path    string
	entries []TempEntry
}

// NewTempTracker loads or creates the temp index file.
func NewTempTracker(indexPath string) *TempTracker {
	t := &TempTracker{path: indexPath}
	t.load()
	return t
}

// Add registers a temp file entry.
func (t *TempTracker) Add(entry TempEntry) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = append(t.entries, entry)
	return t.save()
}

// Remove deletes a temp file entry by temp path.
func (t *TempTracker) Remove(tempPath string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, e := range t.entries {
		if e.TempPath == tempPath {
			t.entries = append(t.entries[:i], t.entries[i+1:]...)
			return t.save()
		}
	}
	return nil
}

// List returns all tracked temp entries.
func (t *TempTracker) List() []TempEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	result := make([]TempEntry, len(t.entries))
	copy(result, t.entries)
	return result
}

// ListPaths returns all tracked temp file paths.
func (t *TempTracker) ListPaths() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	paths := make([]string, len(t.entries))
	for i, e := range t.entries {
		paths[i] = e.TempPath
	}
	return paths
}

// Clear removes all entries.
func (t *TempTracker) Clear() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = nil
	return t.save()
}

func (t *TempTracker) load() {
	data, err := os.ReadFile(t.path)
	if err != nil {
		t.entries = nil
		return
	}
	var entries []TempEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.entries = nil
		return
	}
	t.entries = entries
}

func (t *TempTracker) save() error {
	dir := dirOf(t.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create runtime dir: %w", err)
	}

	data, err := json.MarshalIndent(t.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal temp index: %w", err)
	}

	tmpPath := t.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("write temp index: %w", err)
	}

	return os.Rename(tmpPath, t.path)
}

func dirOf(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}
