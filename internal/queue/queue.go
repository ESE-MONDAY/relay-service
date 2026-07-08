package queue

import (
	"context"
	"errors"
)

var (
	ErrQueueFull = errors.New("queue is full")
)

type Publisher interface {
	Publish(
		ctx context.Context,
		event any,
	) error
}

type Consumer interface {
	Jobs() <-chan Job
}

type Queue interface {
	Publisher
	Consumer
}

type InMemoryQueue struct {
	jobs chan Job
}

func NewInMemoryQueue(
	bufferSize int,
) *InMemoryQueue {

	return &InMemoryQueue{
		jobs: make(chan Job, bufferSize),
	}
}

func (q *InMemoryQueue) Publish(
	ctx context.Context,
	job Job,
) error {

	select {

	case q.jobs <- job:
		return nil

	case <-ctx.Done():
		return ctx.Err()

	default:
		return ErrQueueFull
	}
}

func (q *InMemoryQueue) Jobs() <-chan Job {
	return q.jobs
}
