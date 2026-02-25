package queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	maxBackoffDelay   = 5 * time.Minute
	maxRetryQueueLen  = 10000
	retryTickInterval = 500 * time.Millisecond
)

type RetryJob struct {
	job       Job
	executeAt time.Time
}

type InMemoryQueue struct {
	queue      chan Job
	retryQueue []RetryJob
	logger     utils.Logger
	handler    JobHandler

	mu     sync.Mutex
	wg     sync.WaitGroup
	cancel context.CancelFunc
}

func NewInMemoryQueue(size int, logger utils.Logger, handler ...JobHandler) *InMemoryQueue {
	h := DefaultJobHandler
	if len(handler) > 0 && handler[0] != nil {
		h = handler[0]
	}
	return &InMemoryQueue{
		queue:   make(chan Job, size),
		logger:  logger,
		handler: h,
	}
}

func (q *InMemoryQueue) Start(ctx context.Context, workers int) {
	ctx, cancel := context.WithCancel(ctx)
	q.cancel = cancel

	for i := 0; i < workers; i++ {
		q.wg.Add(1)
		go q.worker(ctx, i)
	}

	q.wg.Add(1)
	go q.retryScheduler(ctx)
	q.logger.Infof("[memory_queue] started with %d workers, buffer_size=%d", workers, cap(q.queue))
}

func (q *InMemoryQueue) Stop() {
	q.cancel()
	q.wg.Wait()
	q.logger.Infof("[memory_queue] all workers stopped")
}

func (q *InMemoryQueue) Enqueue(ctx context.Context, job Job) error {
	select {
	case q.queue <- job:
		q.logger.Debugf("[memory_queue] Enqueue: job enqueued %s", job.LogPrefix())
		return nil
	default:
		q.logger.Errorf("[memory_queue] Enqueue: queue is full, job dropped %s", job.LogPrefix())
		return fmt.Errorf("in-memory queue is full, job %s dropped", job.ID)
	}
}

func (q *InMemoryQueue) worker(ctx context.Context, id int) {
	defer q.wg.Done()
	q.logger.Infof("[memory_queue] worker %d started", id)
	for {
		select {
		case <-ctx.Done():
			q.logger.Infof("[memory_queue] worker %d stopped", id)
			return
		case job := <-q.queue:
			if err := q.handler(ctx, job); err != nil {
				q.logger.Warnf("[memory_queue] worker(%d): job execution failed %s retry=%d: %v", id, job.LogPrefix(), job.Retry, err)
				q.scheduleRetry(job)
			} else {
				q.logger.Debugf("[memory_queue] worker(%d): job completed %s", id, job.LogPrefix())
			}
		}
	}
}

func (q *InMemoryQueue) scheduleRetry(job Job) {
	job.Retry++
	if job.MaxRetry > 0 && job.Retry > job.MaxRetry {
		q.logger.Errorf("[memory_queue] scheduleRetry: job exceeded max retries, discarding %s max_retry=%d", job.LogPrefix(), job.MaxRetry)
		return
	}

	delay := time.Second * time.Duration(1<<job.Retry)
	if delay > maxBackoffDelay {
		delay = maxBackoffDelay
	}

	q.mu.Lock()
	if len(q.retryQueue) >= maxRetryQueueLen {
		q.mu.Unlock()
		q.logger.Errorf("[memory_queue] scheduleRetry: retry queue full (%d), discarding %s", maxRetryQueueLen, job.LogPrefix())
		return
	}
	q.retryQueue = append(q.retryQueue, RetryJob{
		job:       job,
		executeAt: time.Now().Add(delay),
	})
	q.mu.Unlock()

	q.logger.Debugf("[memory_queue] scheduleRetry: job scheduled for retry %s retry=%d delay=%s", job.LogPrefix(), job.Retry, delay)
}

func (q *InMemoryQueue) retryScheduler(ctx context.Context) {
	defer q.wg.Done()
	ticker := time.NewTicker(retryTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()

			q.mu.Lock()
			var remaining []RetryJob
			var ready []Job

			for _, r := range q.retryQueue {
				if r.executeAt.Before(now) {
					ready = append(ready, r.job)
				} else {
					remaining = append(remaining, r)
				}
			}
			q.retryQueue = remaining
			q.mu.Unlock()

			for _, job := range ready {
				select {
				case q.queue <- job:
				default:
					backoff := retryTickInterval * time.Duration(job.Retry+1)
					if backoff > maxBackoffDelay {
						backoff = maxBackoffDelay
					}
					q.logger.Warnf("[memory_queue] retryScheduler: queue full, re-scheduling job %s delay=%s", job.LogPrefix(), backoff)
					q.mu.Lock()
					q.retryQueue = append(q.retryQueue, RetryJob{
						job:       job,
						executeAt: now.Add(backoff),
					})
					q.mu.Unlock()
				}
			}

			if len(ready) > 0 {
				q.logger.Debugf("[memory_queue] retry scheduler: re-enqueued %d jobs", len(ready))
			}
		}
	}
}
