package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ESE-MONDAY/relay-service/internal/models"
)

type EmailStore interface {
	Save(
		ctx context.Context,
		email *models.Email,
	) (
		*models.Email,
		bool,
		error,
	)

	FindByID(

		ctx context.Context,

		id uuid.UUID,

	) (*models.Email, error)

	Ping(

		ctx context.Context,

	) error
	UpdateStatus(
		ctx context.Context,
		id uuid.UUID,
		status models.EmailStatus,
	) error
	FindByIdempotencyKey(
		ctx context.Context,
		key string,
	) (*models.Email, error)
}
