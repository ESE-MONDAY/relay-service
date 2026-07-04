package sender

import (
	"context"

	"github.com/ESE-MONDAY/relay-service/internal/models"
)

type Sender interface {
	Send(
		ctx context.Context,
		email *models.Email,
	) error
}
