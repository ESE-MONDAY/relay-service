package sender

import (
	"context"
	"time"

	"github.com/ESE-MONDAY/relay-service/internal/models"
)

type NoopSender struct{}

func NewNoopSender() *NoopSender {
	return &NoopSender{}
}

func (s *NoopSender) Send(
	ctx context.Context,
	email *models.Email,
) error {

	time.Sleep(500 * time.Millisecond)

	return nil
}
