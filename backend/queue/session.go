package queue

import (
	"context"
	"sync"
	"time"

	"github.com/motoacs/enque/backend/model"
)

type EventEmitter interface {
	Emit(name string, payload any)
}

type overwritePending struct {
	jobID string
	ch    chan model.ResolveOverwriteDecision
}

type sessionRuntime struct {
	session model.EncodeSession
	jobs    map[string]*model.QueueJob

	ctx    context.Context
	cancel context.CancelFunc

	mu               sync.Mutex
	overwriteWaiters map[string]overwritePending
	jobCancels       map[string]context.CancelFunc
}

func newSessionRuntime(totalJobs int) *sessionRuntime {
	ctx, cancel := context.WithCancel(context.Background())
	return &sessionRuntime{
		session: model.EncodeSession{
			SessionID:   newSessionID(),
			State:       "running",
			StartedAt:   time.Now().UTC(),
			TotalJobs:   totalJobs,
			RunningJobs: 0,
		},
		jobs:             map[string]*model.QueueJob{},
		ctx:              ctx,
		cancel:           cancel,
		overwriteWaiters: map[string]overwritePending{},
		jobCancels:       map[string]context.CancelFunc{},
	}
}
