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
// NOTE: ModeBoth will enqueue to both memory and redis queues independently. Since both
// queues have their own workers, the job WILL be processed twice. Use ModeBoth only when
// handlers are idempotent or when dual-write redundancy is intentional.
func (qm *QueueManager) Enqueue(ctx context.Context, job Job, mode QueueMode) error {
	switch mode {
	case ModeMemory:
		return qm.memory.Enqueue(ctx, job)

	case ModeRedis:
		return qm.redis.Enqueue(ctx, job)

	case ModeBoth:
		memErr := qm.memory.Enqueue(ctx, job)
		redisErr := qm.redis.Enqueue(ctx, job)
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
		if err := qm.redis.Enqueue(ctx, job); err != nil {
			qm.logger.Warnf("[queue_manager] Enqueue: redis failed, falling back to memory %s: %v", job.LogPrefix(), err)
			return qm.memory.Enqueue(ctx, job)
		}
		return nil

	default:
		return fmt.Errorf("invalid queue mode: %d", mode)
	}
}
