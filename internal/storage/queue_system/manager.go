package queue

import (
	"context"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type QueueMode int

const (
	ModeMemory QueueMode = iota
	ModeRedis
	ModeBoth
	ModeFallback
)

type QueueManager struct {
	memory Queue
	redis  Queue
	logger utils.Logger
}

func NewQueueManager(memory Queue, redis Queue, logger utils.Logger) *QueueManager {
	return &QueueManager{
		memory: memory,
		redis:  redis,
		logger: logger,
	}
}

func (qm *QueueManager) Start(ctx context.Context, workers int) {
	if qm.memory != nil {
		qm.memory.Start(ctx, workers)
	}
	if qm.redis != nil {
		qm.redis.Start(ctx, workers)
	}
	qm.logger.Infof("[queue_manager] started with %d workers", workers)
}

func (qm *QueueManager) Stop() {
	qm.logger.Infof("[queue_manager] stopping queues...")
	if qm.memory != nil {
		qm.memory.Stop()
	}
	if qm.redis != nil {
		qm.redis.Stop()
	}
	qm.logger.Infof("[queue_manager] all queues stopped")
}

// Enqueue dispatches a job to one or both queue backends based on the given mode.
// ModeBoth enqueues to both memory and redis for redundancy. When a Deduplicator is
// configured on the HandlerRegistry, only the first backend to process the job will
// actually execute the handler; the second will be silently deduplicated.
func (qm *QueueManager) Enqueue(ctx context.Context, job Job, mode QueueMode) error {
	// Normalize once so both backends share the same ID (critical for dedup).
	job.Normalize()
	switch mode {
	case ModeMemory:
		if qm.memory == nil {
			return fmt.Errorf("memory queue is not configured")
		}
		return qm.memory.Enqueue(ctx, job)

	case ModeRedis:
		if qm.redis == nil {
			return fmt.Errorf("redis queue is not configured")
		}
		return qm.redis.Enqueue(ctx, job)

	case ModeBoth:
		if qm.memory == nil && qm.redis == nil {
			return fmt.Errorf("no queue backend is configured")
		}
		var memErr, redisErr error
		if qm.memory != nil {
			memErr = qm.memory.Enqueue(ctx, job)
		}
		if qm.redis != nil {
			redisErr = qm.redis.Enqueue(ctx, job)
		}
		if memErr != nil && redisErr != nil {
			qm.logger.Errorf("[queue_manager] Enqueue: failed to enqueue to both queues %s: memory=%v, redis=%v", job.LogPrefix(), memErr, redisErr)
			return fmt.Errorf("memory: %w, redis: %v", memErr, redisErr)
		}
		if memErr != nil {
			qm.logger.Warnf("[queue_manager] Enqueue: failed to enqueue to memory queue %s: %v", job.LogPrefix(), memErr)
		}
		if redisErr != nil {
			qm.logger.Warnf("[queue_manager] Enqueue: failed to enqueue to redis queue %s: %v", job.LogPrefix(), redisErr)
		}
		return nil

	case ModeFallback:
		if qm.redis != nil {
			if err := qm.redis.Enqueue(ctx, job); err != nil {
				qm.logger.Warnf("[queue_manager] Enqueue: redis failed, falling back to memory %s: %v", job.LogPrefix(), err)
			} else {
				return nil
			}
		}
		if qm.memory == nil {
			return fmt.Errorf("no queue backend available for fallback")
		}
		return qm.memory.Enqueue(ctx, job)

	default:
		return fmt.Errorf("invalid queue mode: %d", mode)
	}
}
