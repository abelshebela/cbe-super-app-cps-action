package queue

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Deduplicator checks whether a job has already been processed.
type Deduplicator interface {
	// IsDuplicate returns true if the job ID was already marked as processed.
	IsDuplicate(ctx context.Context, jobID string) (bool, error)
	// MarkProcessed records the job ID so future calls to IsDuplicate return true.
	MarkProcessed(ctx context.Context, jobID string) error
}

// ---------------------------------------------------------------------------
// In-memory implementation
// ---------------------------------------------------------------------------

const defaultDedupTTL = 10 * time.Minute
const defaultCleanupInterval = 1 * time.Minute

type memEntry struct {
	expiresAt time.Time
}

type InMemoryDeduplicator struct {
	mu       sync.Mutex
	seen     map[string]memEntry
	ttl      time.Duration
	cancel   context.CancelFunc
}

type MemDedupOption func(*InMemoryDeduplicator)

func WithDedupTTL(d time.Duration) MemDedupOption {
	return func(dd *InMemoryDeduplicator) {
		if d > 0 {
			dd.ttl = d
		}
	}
}

func NewInMemoryDeduplicator(opts ...MemDedupOption) *InMemoryDeduplicator {
	dd := &InMemoryDeduplicator{
		seen: make(map[string]memEntry),
		ttl:  defaultDedupTTL,
	}
	for _, opt := range opts {
		opt(dd)
	}

	ctx, cancel := context.WithCancel(context.Background())
	dd.cancel = cancel
	go dd.cleanup(ctx)
	return dd
}

func (dd *InMemoryDeduplicator) IsDuplicate(_ context.Context, jobID string) (bool, error) {
	dd.mu.Lock()
	defer dd.mu.Unlock()
	e, ok := dd.seen[jobID]
	if !ok {
		return false, nil
	}
	if time.Now().After(e.expiresAt) {
		delete(dd.seen, jobID)
		return false, nil
	}
	return true, nil
}

func (dd *InMemoryDeduplicator) MarkProcessed(_ context.Context, jobID string) error {
	dd.mu.Lock()
	defer dd.mu.Unlock()
	dd.seen[jobID] = memEntry{expiresAt: time.Now().Add(dd.ttl)}
	return nil
}

// Stop cancels the background cleanup goroutine.
func (dd *InMemoryDeduplicator) Stop() {
	if dd.cancel != nil {
		dd.cancel()
	}
}

func (dd *InMemoryDeduplicator) cleanup(ctx context.Context) {
	ticker := time.NewTicker(defaultCleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			dd.mu.Lock()
			for id, e := range dd.seen {
				if now.After(e.expiresAt) {
					delete(dd.seen, id)
				}
			}
			dd.mu.Unlock()
		}
	}
}

// ---------------------------------------------------------------------------
// Redis implementation
// ---------------------------------------------------------------------------

const defaultRedisDedupTTL = 10 * time.Minute
const redisDedupPrefix = "dedup:"

type RedisDeduplicator struct {
	client *redis.Client
	ttl    time.Duration
	prefix string
}

type RedisDedupOption func(*RedisDeduplicator)

func WithRedisDedupTTL(d time.Duration) RedisDedupOption {
	return func(dd *RedisDeduplicator) {
		if d > 0 {
			dd.ttl = d
		}
	}
}

func WithRedisDedupPrefix(p string) RedisDedupOption {
	return func(dd *RedisDeduplicator) {
		if p != "" {
			dd.prefix = p + ":"
		}
	}
}

func NewRedisDeduplicator(client *redis.Client, opts ...RedisDedupOption) *RedisDeduplicator {
	dd := &RedisDeduplicator{
		client: client,
		ttl:    defaultRedisDedupTTL,
		prefix: redisDedupPrefix,
	}
	for _, opt := range opts {
		opt(dd)
	}
	return dd
}

func (dd *RedisDeduplicator) IsDuplicate(ctx context.Context, jobID string) (bool, error) {
	exists, err := dd.client.Exists(ctx, dd.prefix+jobID).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (dd *RedisDeduplicator) MarkProcessed(ctx context.Context, jobID string) error {
	return dd.client.Set(ctx, dd.prefix+jobID, 1, dd.ttl).Err()
}
