package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ESE-MONDAY/relay-service/internal/models"
)

type EmailStore interface {
	// Persistence
	Save(
		ctx context.Context,
		email *models.Email,
	) (
		*models.Email,
		bool,
		error,
	)

	// Queries
	FindByID(
		ctx context.Context,
		id uuid.UUID,
	) (*models.Email, error)

	FindByIdempotencyKey(
		ctx context.Context,
		key string,
	) (*models.Email, error)

	// Processing
	ClaimForProcessing(
		ctx context.Context,
		id uuid.UUID,
	) (*models.Email, error)

	UpdateStatus(
		ctx context.Context,
		id uuid.UUID,
		status models.EmailStatus,
	) error

	// Health
	Ping(
		ctx context.Context,
	) error
}
