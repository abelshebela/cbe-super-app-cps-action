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
)

type RedisQueue struct {
	client  *redis.Client
	logger  utils.Logger
	handler JobHandler

	mainKey       string
	processingKey string
	retryKey      string
	deadLetterKey string

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
	return rq
}

func (rq *RedisQueue) Start(ctx context.Context, workers int) {
	ctx, cancel := context.WithCancel(ctx)
	rq.cancel = cancel

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
	rq.cancel()
	rq.wg.Wait()
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
	select {
	case <-time.After(backoff):
	case <-ctx.Done():
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
		pipe := rq.client.Pipeline()
		pipe.LRem(ctx, rq.processingKey, 1, res)
		pipe.LPush(ctx, rq.processingKey, updated)
		if _, pErr := pipe.Exec(ctx); pErr != nil {
			rq.logger.Warnf("[redis_queue] worker(%d): failed to update processing entry %s: %v", workerID, job.LogPrefix(), pErr)
		}
		res = string(updated)
	}

	if err := rq.handler(ctx, job); err != nil {
		rq.logger.Warnf("[redis_queue] worker(%d): job execution failed %s retry=%d: %v", workerID, job.LogPrefix(), job.Retry, err)
		rq.handleFailure(ctx, job, res)
		return
	}

	rq.logger.Debugf("[redis_queue] worker(%d): job completed %s", workerID, job.LogPrefix())
	rq.client.LRem(ctx, rq.processingKey, 1, res)
}

func (rq *RedisQueue) handleFailure(ctx context.Context, job Job, raw string) {
	job.Retry++

	if job.MaxRetry > 0 && job.Retry > job.MaxRetry {
		rq.logger.Errorf("[redis_queue] handleFailure: job exceeded max retries, moving to dead letter queue %s max_retry=%d", job.LogPrefix(), job.MaxRetry)
		pipe := rq.client.Pipeline()
		pipe.LPush(ctx, rq.deadLetterKey, raw)
		pipe.LRem(ctx, rq.processingKey, 1, raw)
		if _, err := pipe.Exec(ctx); err != nil {
			rq.logger.Errorf("[redis_queue] handleFailure: pipeline exec failed (dead letter) %s: %v", job.LogPrefix(), err)
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

	pipe := rq.client.Pipeline()
	pipe.ZAdd(ctx, rq.retryKey, redis.Z{
		Score:  float64(time.Now().Add(delay).Unix()),
		Member: data,
	})
	pipe.LRem(ctx, rq.processingKey, 1, raw)
	if _, err := pipe.Exec(ctx); err != nil {
		rq.logger.Errorf("[redis_queue] handleFailure: pipeline exec failed (retry) %s: %v", job.LogPrefix(), err)
	}

	rq.logger.Debugf("[redis_queue] handleFailure: job scheduled for retry %s retry=%d delay=%s", job.LogPrefix(), job.Retry, delay)
}

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
			now := time.Now().Unix()

			count, err := retryLuaScript.Run(ctx, rq.client, []string{rq.retryKey, rq.mainKey}, now).Int()
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
	items, err := rq.client.LRange(ctx, rq.processingKey, 0, -1).Result()
	if err != nil {
		if ctx.Err() == nil {
			rq.logger.Warnf("[redis_queue] reaper: failed to read processing list: %v", err)
		}
		return
	}

	now := time.Now()
	reaped := 0
	for _, raw := range items {
		if rq.reapItem(ctx, raw, now) {
			reaped++
		}
	}

	if reaped > 0 {
		rq.logger.Infof("[redis_queue] reaper: recovered %d stale jobs from processing list", reaped)
	}
}

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

	rq.logger.Warnf("[redis_queue] reaper: stale job found (age=%s), re-queuing %s", jobAge, job.LogPrefix())
	pipe := rq.client.Pipeline()
	pipe.LRem(ctx, rq.processingKey, 1, raw)
	pipe.LPush(ctx, rq.mainKey, raw)
	if _, err := pipe.Exec(ctx); err != nil {
		rq.logger.Errorf("[redis_queue] reaper: failed to requeue stale job %s: %v", job.LogPrefix(), err)
	}
	return true
}
