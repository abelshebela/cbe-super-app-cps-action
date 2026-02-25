package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type JobHandler func(ctx context.Context, job Job) error

var (
	handlers   = map[string]func(context.Context, json.RawMessage) error{}
	handlersMu sync.RWMutex
)

func RegisterHandler(name string, fn func(context.Context, json.RawMessage) error) {
	handlersMu.Lock()
	defer handlersMu.Unlock()
	handlers[name] = fn
}

func DefaultJobHandler(ctx context.Context, job Job) error {
	handlersMu.RLock()
	handler, ok := handlers[job.Type]
	handlersMu.RUnlock()

	if !ok {
		return fmt.Errorf("no handler for job type=%s [%s]", job.Type, job.LogPrefix())
	}
	return handler(ctx, job.Payload)
}
