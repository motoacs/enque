package logging

import (
	"io"
	"sync"
)

type SafeWriter struct {
	w  io.Writer
	mu sync.Mutex
}

func NewSafeWriter(w io.Writer) *SafeWriter {
	return &SafeWriter{w: w}
}

func (w *SafeWriter) WriteLine(line string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	_, err := io.WriteString(w.w, line+"\n")
	return err
}
