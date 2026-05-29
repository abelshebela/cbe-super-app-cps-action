package queue

import "context"

type Queue interface {
	Start(ctx context.Context, workers int)
	Stop()
	// Enqueue adds a job to the queue. The caller SHOULD call job.Normalize()
	// before enqueuing. QueueManager does this automatically.
	Enqueue(ctx context.Context, job Job) error
}
