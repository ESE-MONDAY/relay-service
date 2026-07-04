package processor

import (
	"context"

	"github.com/ESE-MONDAY/relay-service/internal/queue"
)

type Processor interface {
	Process(
		ctx context.Context,
		job queue.Job,
	) error
}
