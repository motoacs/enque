package events

import (
	"context"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Emitter wraps Wails event emission (design doc 7).
type Emitter struct {
	ctx context.Context
}

// NewEmitter creates an Emitter with the given Wails context.
func NewEmitter(ctx context.Context) *Emitter {
	return &Emitter{ctx: ctx}
}

func (e *Emitter) emit(name string, data interface{}) {
	wailsRuntime.EventsEmit(e.ctx, name, data)
}

// SessionStarted emits enque:session_started.
func (e *Emitter) SessionStarted(data map[string]interface{}) {
	e.emit("enque:session_started", data)
}

// JobStarted emits enque:job_started.
func (e *Emitter) JobStarted(data map[string]interface{}) {
	e.emit("enque:job_started", data)
}

// JobProgress emits enque:job_progress.
func (e *Emitter) JobProgress(data map[string]interface{}) {
	e.emit("enque:job_progress", data)
}

// JobLog emits enque:job_log.
func (e *Emitter) JobLog(data map[string]interface{}) {
	e.emit("enque:job_log", data)
}

// JobNeedsOverwrite emits enque:job_needs_overwrite.
func (e *Emitter) JobNeedsOverwrite(data map[string]interface{}) {
	e.emit("enque:job_needs_overwrite", data)
}

// JobFinished emits enque:job_finished.
func (e *Emitter) JobFinished(data map[string]interface{}) {
	e.emit("enque:job_finished", data)
}

// SessionState emits enque:session_state.
func (e *Emitter) SessionState(data map[string]interface{}) {
	e.emit("enque:session_state", data)
}

// SessionFinished emits enque:session_finished.
func (e *Emitter) SessionFinished(data map[string]interface{}) {
	e.emit("enque:session_finished", data)
}

// Warning emits enque:warning.
func (e *Emitter) Warning(data map[string]interface{}) {
	e.emit("enque:warning", data)
}

// Error emits enque:error.
func (e *Emitter) Error(data map[string]interface{}) {
	e.emit("enque:error", data)
}
