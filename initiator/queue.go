package initiator

import (
	queue "cbe-super-app-cps-action/internal/storage/queue_system"

	"github.com/redis/go-redis/v9"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	queueMetricsNamespace = "cps_action"
	queueRedisPrefix      = "cps_action:jobs"
	queueBufferSize       = 5000
	queueWorkers          = 4
)

// QueueInfra groups queue-system components for use by the rest of the application.
type QueueInfra struct {
	Manager  *queue.QueueManager
	Registry *queue.HandlerRegistry
	Metrics  *queue.QueueMetrics
}

// InitQueueSystem creates and returns all queue-system components.
// The caller is responsible for calling QueueInfra.Manager.Start and .Stop.
func InitQueueSystem(redisClient *redis.Client, logger utils.Logger) *QueueInfra {
	// Metrics
	metrics := queue.NewQueueMetrics(queueMetricsNamespace)

	// Deduplicator (Redis-backed so it works across restarts and ModeBoth)
	dedup := queue.NewRedisDeduplicator(redisClient,
		queue.WithRedisDedupPrefix(queueRedisPrefix),
	)

	// Handler registry
	registry := queue.NewHandlerRegistry(
		queue.WithDeduplicator(dedup),
		queue.WithRegistryMetrics(metrics),
	)

	handler := registry.JobHandler()

	// Memory queue
	memQueue := queue.NewInMemoryQueue(queueBufferSize, logger,
		queue.WithMemoryHandler(handler),
		queue.WithMemoryMetrics(metrics),
	)

	// Redis queue
	redisQueue := queue.NewRedisQueue(redisClient, logger,
		queue.WithPrefix(queueRedisPrefix),
		queue.WithJobHandler(handler),
		queue.WithRedisMetrics(metrics),
	)

	// Manager (both backends available; callers choose mode per-enqueue)
	manager := queue.NewQueueManager(memQueue, redisQueue, logger)

	logger.Infof("[queue] infrastructure initialized (buffer=%d, workers=%d, prefix=%s)", queueBufferSize, queueWorkers, queueRedisPrefix)
	return &QueueInfra{
		Manager:  manager,
		Registry: registry,
		Metrics:  metrics,
	}
}
