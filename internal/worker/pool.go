package worker

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/ESE-MONDAY/relay-service/internal/processor"
	"github.com/ESE-MONDAY/relay-service/internal/queue"
)

type Pool struct {
	workers []*Worker

	wg sync.WaitGroup

	logger *zap.Logger
}

func NewPool(
	count int,
	queue queue.Consumer,
	processor processor.Processor,
	logger *zap.Logger,
) *Pool {

	workers := make([]*Worker, 0, count)

	for i := 1; i <= count; i++ {

		workers = append(
			workers,
			NewWorker(
				i,
				queue,
				processor,
				logger,
			),
		)
	}

	return &Pool{
		workers: workers,
		logger:  logger,
	}
}

func (p *Pool) Start(
	ctx context.Context,
) {

	for _, worker := range p.workers {

		p.wg.Add(1)

		go func(w *Worker) {
			defer p.wg.Done()

			w.Start(ctx)
		}(worker)
	}
}

func (p *Pool) Wait() {

	p.logger.Info("Waiting for workers...")

	p.wg.Wait()

	p.logger.Info("All workers stopped.")
}
