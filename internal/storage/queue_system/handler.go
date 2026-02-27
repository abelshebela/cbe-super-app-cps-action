package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type JobHandler func(ctx context.Context, job Job) error

// HandlerRegistry is a per-instance handler registry with optional dedup.
type HandlerRegistry struct {
	handlers map[string]func(context.Context, json.RawMessage) error
	mu       sync.RWMutex
	dedup    Deduplicator
	metrics  *QueueMetrics
}

type RegistryOption func(*HandlerRegistry)

func WithDeduplicator(d Deduplicator) RegistryOption {
	return func(r *HandlerRegistry) { r.dedup = d }
}

func WithRegistryMetrics(m *QueueMetrics) RegistryOption {
	return func(r *HandlerRegistry) { r.metrics = m }
}

func NewHandlerRegistry(opts ...RegistryOption) *HandlerRegistry {
	r := &HandlerRegistry{
		handlers: make(map[string]func(context.Context, json.RawMessage) error),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *HandlerRegistry) Register(name string, fn func(context.Context, json.RawMessage) error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[name] = fn
}

// Handle dispatches a job to the registered handler, applying dedup when configured.
func (r *HandlerRegistry) Handle(ctx context.Context, job Job) error {
	if r.dedup != nil && job.ID != "" {
		dup, err := r.dedup.IsDuplicate(ctx, job.ID)
		if err == nil && dup {
			if r.metrics != nil {
				r.metrics.IncDeduplicated(job.Type)
			}
			return nil
		}
	}

	r.mu.RLock()
	handler, ok := r.handlers[job.Type]
	r.mu.RUnlock()

	if !ok {
		return fmt.Errorf("no handler for job type=%s [%s]", job.Type, job.LogPrefix())
	}

	if err := handler(ctx, job.Payload); err != nil {
		return err
	}

	if r.dedup != nil && job.ID != "" {
		_ = r.dedup.MarkProcessed(ctx, job.ID)
	}
	return nil
}

// JobHandler returns a JobHandler func suitable for passing to queue constructors.
func (r *HandlerRegistry) JobHandler() JobHandler {
	return r.Handle
}

// ---------------------------------------------------------------------------
// Backward-compatible global registry (delegates to a default instance)
// ---------------------------------------------------------------------------

var defaultRegistry = NewHandlerRegistry()

func RegisterHandler(name string, fn func(context.Context, json.RawMessage) error) {
	defaultRegistry.Register(name, fn)
}

func DefaultJobHandler(ctx context.Context, job Job) error {
	return defaultRegistry.Handle(ctx, job)
}
