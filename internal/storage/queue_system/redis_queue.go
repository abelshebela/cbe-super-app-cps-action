package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	redisMaxBackoffDelay    = 5 * time.Minute
	redisErrorBackoffBase   = 500 * time.Millisecond
	redisErrorBackoffMax    = 30 * time.Second
	redisReaperInterval     = 30 * time.Second
	redisProcessingStaleAge = 5 * time.Minute
	redisMaxRetrySetSize    = 10000
	redisStopTimeout        = 30 * time.Second
	redisReaperMaxPerCycle  = 500
)

type RedisQueue struct {
	client  *redis.Client
	logger  utils.Logger
	handler JobHandler
	metrics *QueueMetrics

	mainKey       string
	processingKey string
	retryKey      string
	deadLetterKey string

	processingKeys           []string // cached [processingKey] for Lua scripts
	retryMainKeys            []string // cached [retryKey, mainKey] for Lua scripts
	processingDeadLetterKeys []string // cached [processingKey, deadLetterKey] for Lua scripts
	processingRetryKeys      []string // cached [processingKey, retryKey] for Lua scripts
	processingMainKeys       []string // cached [processingKey, mainKey] for Lua scripts

	mu     sync.Mutex
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

type RedisQueueOption func(*RedisQueue)

func WithPrefix(prefix string) RedisQueueOption {
	return func(rq *RedisQueue) {
		if prefix != "" {
			rq.mainKey = prefix + ":main"
			rq.processingKey = prefix + ":processing"
			rq.retryKey = prefix + ":retry"
			rq.deadLetterKey = prefix + ":dead_letter"
		}
	}
}

func WithJobHandler(h JobHandler) RedisQueueOption {
	return func(rq *RedisQueue) {
		if h != nil {
			rq.handler = h
		}
	}
}

func WithRedisMetrics(m *QueueMetrics) RedisQueueOption {
	return func(rq *RedisQueue) { rq.metrics = m }
}

func NewRedisQueue(client *redis.Client, logger utils.Logger, opts ...RedisQueueOption) *RedisQueue {
	rq := &RedisQueue{
		client:        client,
		logger:        logger,
		handler:       DefaultJobHandler,
		mainKey:       "jobs:main",
		processingKey: "jobs:processing",
		retryKey:      "jobs:retry",
		deadLetterKey: "jobs:dead_letter",
	}
	for _, opt := range opts {
		opt(rq)
	}
	rq.processingKeys = []string{rq.processingKey}
	rq.retryMainKeys = []string{rq.retryKey, rq.mainKey}
	rq.processingDeadLetterKeys = []string{rq.processingKey, rq.deadLetterKey}
	rq.processingRetryKeys = []string{rq.processingKey, rq.retryKey}
	rq.processingMainKeys = []string{rq.processingKey, rq.mainKey}
	return rq
}

func (rq *RedisQueue) Start(ctx context.Context, workers int) {
	rq.mu.Lock()
	if rq.cancel != nil {
		rq.mu.Unlock()
		rq.logger.Warnf("[redis_queue] Start: already running, ignoring duplicate start")
		return
	}
	ctx, cancel := context.WithCancel(ctx)
	rq.cancel = cancel
	rq.mu.Unlock()

	for i := 0; i < workers; i++ {
		rq.wg.Add(1)
		go rq.worker(ctx, i)
	}

	rq.wg.Add(1)
	go rq.retryScheduler(ctx)

	rq.wg.Add(1)
	go rq.processingReaper(ctx)

	rq.logger.Infof("[redis_queue] started with %d workers", workers)
}

func (rq *RedisQueue) Stop() {
	rq.mu.Lock()
	cancel := rq.cancel
	rq.cancel = nil
	rq.mu.Unlock()

	if cancel == nil {
		return
	}
	cancel()
	if !waitWithTimeout(&rq.wg, redisStopTimeout) {
		rq.logger.Errorf("[redis_queue] Stop: timed out waiting for workers to exit after %s", redisStopTimeout)
	}
	rq.logger.Infof("[redis_queue] all workers stopped")
}

func (rq *RedisQueue) Enqueue(ctx context.Context, job Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		rq.logger.Errorf("[redis_queue] Enqueue: failed to marshal job %s: %v", job.LogPrefix(), err)
		return fmt.Errorf("marshal job: %w", err)
	}
	if err := rq.client.LPush(ctx, rq.mainKey, data).Err(); err != nil {
		rq.logger.Errorf("[redis_queue] Enqueue: failed to push job %s: %v", job.LogPrefix(), err)
		return err
	}
	rq.metrics.IncEnqueued("redis", job.Type)
	rq.logger.Debugf("[redis_queue] Enqueue: job enqueued %s", job.LogPrefix())
	return nil
}

func (rq *RedisQueue) worker(ctx context.Context, id int) {
	defer rq.wg.Done()
	rq.logger.Infof("[redis_queue] worker %d started", id)

	errBackoff := redisErrorBackoffBase

	for {
		select {
		case <-ctx.Done():
			rq.logger.Infof("[redis_queue] worker %d stopped", id)
			return
		default:
		}

		res, err := rq.client.BLMove(ctx, rq.mainKey, rq.processingKey, "RIGHT", "LEFT", 2*time.Second).Result()
		if err != nil {
			errBackoff = rq.handlePopError(ctx, id, err, errBackoff)
			continue
		}

		errBackoff = redisErrorBackoffBase
		rq.processJob(ctx, id, res)
	}
}

func (rq *RedisQueue) handlePopError(ctx context.Context, workerID int, err error, backoff time.Duration) time.Duration {
	if err == redis.Nil || ctx.Err() != nil {
		return redisErrorBackoffBase
	}

	rq.logger.Warnf("[redis_queue] worker %d: failed to pop job, backing off %s: %v", workerID, backoff, err)
	timer := time.NewTimer(backoff)
	select {
	case <-timer.C:
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	}

	nextBackoff := backoff * 2
	if nextBackoff > redisErrorBackoffMax {
		nextBackoff = redisErrorBackoffMax
	}
	return nextBackoff
}

func (rq *RedisQueue) processJob(ctx context.Context, workerID int, res string) {
	var job Job
	if err := json.Unmarshal([]byte(res), &job); err != nil {
		rq.logger.Errorf("[redis_queue] worker %d: failed to unmarshal job, discarding: %v", workerID, err)
		rq.client.LRem(ctx, rq.processingKey, 1, res)
		return
	}

	job.ProcessingStartedAt = time.Now().Unix()
	updated, err := json.Marshal(job)
	if err == nil {
		updatedStr := string(updated)
		if pErr := updateProcessingScript.Run(ctx, rq.client, rq.processingKeys, res, updatedStr).Err(); pErr != nil {
			rq.logger.Warnf("[redis_queue] worker(%d): failed to update processing entry %s: %v", workerID, job.LogPrefix(), pErr)
		}
		res = updatedStr
	}

	if err := rq.handler(ctx, job); err != nil {
		rq.metrics.IncFailed("redis", job.Type)
		rq.logger.Warnf("[redis_queue] worker(%d): job execution failed %s retry=%d: %v", workerID, job.LogPrefix(), job.Retry, err)
		rq.handleFailure(ctx, job, res)
		return
	}

	rq.metrics.IncProcessed("redis", job.Type)
	rq.logger.Debugf("[redis_queue] worker(%d): job completed %s", workerID, job.LogPrefix())
	rq.client.LRem(ctx, rq.processingKey, 1, res)
}

func (rq *RedisQueue) handleFailure(ctx context.Context, job Job, raw string) {
	job.Retry++

	if job.MaxRetry > 0 && job.Retry > job.MaxRetry {
		rq.metrics.IncDeadLettered("redis", job.Type)
		rq.logger.Errorf("[redis_queue] handleFailure: job exceeded max retries, moving to dead letter queue %s max_retry=%d", job.LogPrefix(), job.MaxRetry)
		dlData, mErr := json.Marshal(job)
		if mErr != nil {
			rq.logger.Errorf("[redis_queue] handleFailure: failed to marshal job for dead letter %s: %v", job.LogPrefix(), mErr)
			rq.client.LRem(ctx, rq.processingKey, 1, raw)
			return
		}
		if err := moveToDeadLetterScript.Run(ctx, rq.client, rq.processingDeadLetterKeys, raw, string(dlData)).Err(); err != nil {
			rq.logger.Errorf("[redis_queue] handleFailure: failed to move job to dead letter %s: %v", job.LogPrefix(), err)
		}
		return
	}

	delay := time.Second * time.Duration(1<<job.Retry)
	if delay > redisMaxBackoffDelay {
		delay = redisMaxBackoffDelay
	}

	data, err := json.Marshal(job)
	if err != nil {
		rq.logger.Errorf("[redis_queue] handleFailure: failed to marshal job for retry %s: %v", job.LogPrefix(), err)
		rq.client.LRem(ctx, rq.processingKey, 1, raw)
		return
	}

	// Cap retry set size to prevent unbounded growth
	if count, cErr := rq.client.ZCard(ctx, rq.retryKey).Result(); cErr == nil && count >= redisMaxRetrySetSize {
		rq.metrics.IncDeadLettered("redis", job.Type)
		rq.logger.Errorf("[redis_queue] handleFailure: retry set full (%d), moving to dead letter %s", count, job.LogPrefix())
		if pErr := moveToDeadLetterScript.Run(ctx, rq.client, rq.processingDeadLetterKeys, raw, string(data)).Err(); pErr != nil {
			rq.logger.Errorf("[redis_queue] handleFailure: failed to move job to dead letter (overflow) %s: %v", job.LogPrefix(), pErr)
		}
		return
	}

	rq.metrics.IncRetried("redis", job.Type)
	score := float64(time.Now().Add(delay).Unix())
	if err := moveToRetryScript.Run(ctx, rq.client, rq.processingRetryKeys, raw, string(data), score).Err(); err != nil {
		rq.logger.Errorf("[redis_queue] handleFailure: failed to move job to retry set %s: %v", job.LogPrefix(), err)
	}

	rq.logger.Debugf("[redis_queue] handleFailure: job scheduled for retry %s retry=%d delay=%s", job.LogPrefix(), job.Retry, delay)
}

var updateProcessingScript = redis.NewScript(`
local key = KEYS[1]
local oldVal = ARGV[1]
local newVal = ARGV[2]
redis.call('LREM', key, 1, oldVal)
redis.call('LPUSH', key, newVal)
return 1
`)

var moveToDeadLetterScript = redis.NewScript(`
local processingKey = KEYS[1]
local deadLetterKey = KEYS[2]
local rawVal = ARGV[1]
local dlData = ARGV[2]
redis.call('LREM', processingKey, 1, rawVal)
redis.call('LPUSH', deadLetterKey, dlData)
return 1
`)

var moveToRetryScript = redis.NewScript(`
local processingKey = KEYS[1]
local retryKey = KEYS[2]
local rawVal = ARGV[1]
local retryData = ARGV[2]
local score = tonumber(ARGV[3])
redis.call('LREM', processingKey, 1, rawVal)
redis.call('ZADD', retryKey, score, retryData)
return 1
`)

var requeueFromProcessingScript = redis.NewScript(`
local processingKey = KEYS[1]
local mainKey = KEYS[2]
local rawVal = ARGV[1]
local updatedData = ARGV[2]
redis.call('LREM', processingKey, 1, rawVal)
redis.call('LPUSH', mainKey, updatedData)
return 1
`)

var retryLuaScript = redis.NewScript(`
local retryKey = KEYS[1]
local mainKey = KEYS[2]
local now = ARGV[1]
local jobs = redis.call('ZRANGEBYSCORE', retryKey, '-inf', now)
local count = 0
for _, job in ipairs(jobs) do
    redis.call('ZREM', retryKey, job)
    redis.call('LPUSH', mainKey, job)
    count = count + 1
end
return count
`)

func (rq *RedisQueue) retryScheduler(ctx context.Context) {
	defer rq.wg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if depth, dErr := rq.client.LLen(ctx, rq.mainKey).Result(); dErr == nil {
				rq.metrics.SetDepth("redis", float64(depth))
			}

			now := time.Now().Unix()

			count, err := retryLuaScript.Run(ctx, rq.client, rq.retryMainKeys, now).Int()
			if err != nil {
				if ctx.Err() == nil {
					rq.logger.Warnf("[redis_queue] retry scheduler: lua script failed: %v", err)
				}
				continue
			}

			if count > 0 {
				rq.logger.Debugf("[redis_queue] retry scheduler: re-enqueued %d jobs", count)
			}
		}
	}
}

func (rq *RedisQueue) processingReaper(ctx context.Context) {
	defer rq.wg.Done()
	ticker := time.NewTicker(redisReaperInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rq.reapStaleJobs(ctx)
		}
	}
}

func (rq *RedisQueue) reapStaleJobs(ctx context.Context) {
	const batchSize int64 = 100
	var offset int64
	now := time.Now()
	reaped := 0

	for reaped < redisReaperMaxPerCycle {
		batchReaped, done := rq.reapBatch(ctx, now, batchSize, &offset)
		reaped += batchReaped
		if done {
			break
		}
	}

	if reaped >= redisReaperMaxPerCycle {
		rq.logger.Warnf("[redis_queue] reaper: hit per-cycle cap (%d), deferring remaining to next cycle", redisReaperMaxPerCycle)
	}
	if reaped > 0 {
		rq.logger.Infof("[redis_queue] reaper: recovered %d stale jobs from processing list", reaped)
	}
}

// reapBatch processes one batch of items from the processing list.
// It returns the number of items reaped and whether iteration should stop.
func (rq *RedisQueue) reapBatch(ctx context.Context, now time.Time, batchSize int64, offset *int64) (reaped int, done bool) {
	items, err := rq.client.LRange(ctx, rq.processingKey, *offset, *offset+batchSize-1).Result()
	if err != nil {
		if ctx.Err() == nil {
			rq.logger.Warnf("[redis_queue] reaper: failed to read processing list: %v", err)
		}
		return 0, true
	}
	if len(items) == 0 {
		return 0, true
	}

	for _, raw := range items {
		if rq.reapItem(ctx, raw, now) {
			reaped++
		} else {
			*offset++
		}
	}
	return reaped, false
}

// reapItem checks if a job in the processing list is stale and re-queues it.
// Stale jobs have their retry count incremented to prevent infinite reap loops.
// Jobs exceeding MaxRetry are moved to the dead letter queue.
func (rq *RedisQueue) reapItem(ctx context.Context, raw string, now time.Time) bool {
	var job Job
	if err := json.Unmarshal([]byte(raw), &job); err != nil {
		rq.logger.Warnf("[redis_queue] reaper: unmarshal failed, removing stale entry")
		rq.client.LRem(ctx, rq.processingKey, 1, raw)
		return true
	}

	timestamp := job.ProcessingStartedAt
	if timestamp == 0 {
		timestamp = job.CreatedAt
	}
	jobAge := now.Sub(time.Unix(timestamp, 0))
	if jobAge <= redisProcessingStaleAge {
		return false
	}

	job.Retry++
	job.ProcessingStartedAt = 0

	rq.logger.Warnf("[redis_queue] reaper: stale job found (age=%s), re-queuing %s retry=%d", jobAge, job.LogPrefix(), job.Retry)

	if job.MaxRetry > 0 && job.Retry > job.MaxRetry {
		rq.metrics.IncDeadLettered("redis", job.Type)
		rq.logger.Errorf("[redis_queue] reaper: stale job exceeded max retries, moving to dead letter %s max_retry=%d", job.LogPrefix(), job.MaxRetry)
		dlData, mErr := json.Marshal(job)
		if mErr != nil {
			rq.logger.Errorf("[redis_queue] reaper: failed to marshal dead letter job %s: %v", job.LogPrefix(), mErr)
			rq.client.LRem(ctx, rq.processingKey, 1, raw)
			return true
		}
		if err := moveToDeadLetterScript.Run(ctx, rq.client, rq.processingDeadLetterKeys, raw, string(dlData)).Err(); err != nil {
			rq.logger.Errorf("[redis_queue] reaper: failed to move stale job to dead letter %s: %v", job.LogPrefix(), err)
		}
		return true
	}

	updated, err := json.Marshal(job)
	if err != nil {
		rq.logger.Errorf("[redis_queue] reaper: failed to marshal job for re-queue %s: %v", job.LogPrefix(), err)
		rq.client.LRem(ctx, rq.processingKey, 1, raw)
		return true
	}

	if err := requeueFromProcessingScript.Run(ctx, rq.client, rq.processingMainKeys, raw, string(updated)).Err(); err != nil {
		rq.logger.Errorf("[redis_queue] reaper: failed to requeue stale job %s: %v", job.LogPrefix(), err)
	}
	return true
}
