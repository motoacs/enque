package encoder

import (
	"fmt"

	"github.com/motoacs/enque/backend/model"
)

type Registry struct {
	adapters map[model.EncoderType]Adapter
}

func NewRegistry(adapters ...Adapter) *Registry {
	m := map[model.EncoderType]Adapter{}
	for _, a := range adapters {
		m[a.Type()] = a
	}
	return &Registry{adapters: m}
}

func (r *Registry) Resolve(t model.EncoderType) (Adapter, error) {
	a, ok := r.adapters[t]
	if !ok {
		return nil, &model.EnqueError{Code: model.ErrEncoderNotImplemented, Message: fmt.Sprintf("encoder_type %s is not implemented", t)}
	}
	return a, nil
}
