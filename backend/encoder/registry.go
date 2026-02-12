package encoder

import (
	"fmt"
	"sync"
)

// Registry manages encoder adapter resolution.
type Registry struct {
	mu       sync.RWMutex
	adapters map[string]Adapter
}

// NewRegistry creates an empty registry.
func NewRegistry() *Registry {
	return &Registry{
		adapters: make(map[string]Adapter),
	}
}

// Register adds an adapter for an encoder type.
func (r *Registry) Register(adapter Adapter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[adapter.Type()] = adapter
}

// Resolve returns the adapter for the given encoder type.
func (r *Registry) Resolve(encoderType string) (Adapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, ok := r.adapters[encoderType]
	if !ok {
		return nil, fmt.Errorf("%s: %s", ErrEncoderNotImplemented, encoderType)
	}
	return adapter, nil
}
