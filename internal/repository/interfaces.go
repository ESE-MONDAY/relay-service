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

	) error

	FindByID(

		ctx context.Context,

		id uuid.UUID,

	) (*models.Email, error)

	Ping(

		ctx context.Context,

	) error
}
