package worker

import (
	"context"

	"go.uber.org/zap"

	"github.com/ESE-MONDAY/relay-service/internal/processor"
	"github.com/ESE-MONDAY/relay-service/internal/queue"
)

type Worker struct {
	id int

	logger *zap.Logger

	queue queue.Consumer

	processor processor.Processor
}

func NewWorker(
	id int,
	queue queue.Consumer,
	processor processor.Processor,
	logger *zap.Logger,
) *Worker {

	return &Worker{
		id:        id,
		queue:     queue,
		processor: processor,
		logger:    logger,
	}
}

func (w *Worker) Start(
	ctx context.Context,
) {

	w.logger.Info(
		"worker started",
		zap.Int("worker", w.id),
	)

	for {

		select {

		case <-ctx.Done():

			w.logger.Info(
				"worker stopped",
				zap.Int("worker", w.id),
			)

			return

		case job, ok := <-w.queue.Jobs():

			if !ok {

				w.logger.Info(
					"queue closed",
					zap.Int("worker", w.id),
				)

				return
			}

			if err := w.processor.Process(ctx, job); err != nil {

				w.logger.Error(
					"job failed",
					zap.Int("worker", w.id),
					zap.String("job_id", job.ID.String()),
					zap.Error(err),
				)

				continue
			}
		}
	}
}
