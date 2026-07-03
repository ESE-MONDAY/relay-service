package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ESE-MONDAY/relay-service/internal/models"
)

var (
	ErrEmailNotFound = errors.New("email not found")
)

type EmailRepository struct {
	pool *pgxpool.Pool
}

func NewEmailRepository(pool *pgxpool.Pool) *EmailRepository {
	return &EmailRepository{
		pool: pool,
	}
}

func (r *EmailRepository) Save(
	ctx context.Context,
	email *models.Email,
) error {

	query := `
INSERT INTO emails (
	id,
	sender,
	recipient,
	subject,
	text_body,
	html_body,
	status
)
VALUES (
	$1, $2, $3, $4, $5, $6, $7
)
`

	_, err := r.pool.Exec(
		ctx,
		query,
		email.ID,
		email.From,
		email.To,
		email.Subject,
		email.Text,
		email.HTML,
		email.Status,
	)

	if err != nil {
		return fmt.Errorf("save email: %w", err)
	}

	return nil
}

func (r *EmailRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Email, error) {

	query := `
SELECT
	id,
	sender,
	recipient,
	subject,
	text_body,
	html_body,
	status,
	created_at
FROM emails
WHERE id = $1
`

	email := &models.Email{}

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&email.ID,
		&email.From,
		&email.To,
		&email.Subject,
		&email.Text,
		&email.HTML,
		&email.Status,
		&email.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmailNotFound
		}

		return nil, fmt.Errorf("find email by id: %w", err)
	}

	return email, nil
}

func (r *EmailRepository) Ping(ctx context.Context) error {
	if err := r.pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}

func (r *EmailRepository) Close() {
	r.pool.Close()
}
