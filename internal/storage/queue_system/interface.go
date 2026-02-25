package queue

import "context"

type Queue interface {
    Start(ctx context.Context, workers int)
    Stop()
    Enqueue(ctx context.Context, job Job) error
}