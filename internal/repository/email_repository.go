package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	return &EmailRepository{pool: pool}
}

// -----------------------------
// SAVE (idempotent insert)
// -----------------------------
func (r *EmailRepository) Save(
	ctx context.Context,
	email *models.Email,
) (*models.Email, bool, error) {

	query := `
INSERT INTO emails (
	id,
	sender,
	recipient,
	subject,
	text_body,
	html_body,
	status,
	idempotency_key
)
VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8
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
		email.IdempotencyKey,
	)

	// success path
	if err == nil {
		return email, true, nil
	}

	// handle duplicate idempotency key safely
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {

		if email.IdempotencyKey == nil {
			return nil, false, fmt.Errorf("duplicate insert but missing idempotency key")
		}

		existing, err := r.FindByIdempotencyKey(ctx, *email.IdempotencyKey)
		if err != nil {
			return nil, false, err
		}

		return existing, false, nil
	}

	return nil, false, fmt.Errorf("save email: %w", err)
}

// -----------------------------
// FIND BY ID
// -----------------------------
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
	idempotency_key,
	created_at
FROM emails
WHERE id = $1
`

	email := &models.Email{}

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&email.ID,
		&email.From,
		&email.To,
		&email.Subject,
		&email.Text,
		&email.HTML,
		&email.Status,
		&email.IdempotencyKey,
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

// -----------------------------
// FIND BY IDEMPOTENCY KEY
// -----------------------------
func (r *EmailRepository) FindByIdempotencyKey(
	ctx context.Context,
	key string,
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
	idempotency_key,
	created_at
FROM emails
WHERE idempotency_key = $1
`

	email := &models.Email{}

	err := r.pool.QueryRow(ctx, query, key).Scan(
		&email.ID,
		&email.From,
		&email.To,
		&email.Subject,
		&email.Text,
		&email.HTML,
		&email.Status,
		&email.IdempotencyKey,
		&email.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEmailNotFound
		}
		return nil, fmt.Errorf("find email by idempotency key: %w", err)
	}

	return email, nil
}

// -----------------------------
// UPDATE STATUS
// -----------------------------
func (r *EmailRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status models.EmailStatus,
) error {

	query := `
UPDATE emails
SET status = $2
WHERE id = $1
`

	cmd, err := r.pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("update email status: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrEmailNotFound
	}

	return nil
}

// -----------------------------
// ATOMIC CLAIM (IMPORTANT FIXED)
// -----------------------------
func (r *EmailRepository) ClaimForProcessing(
	ctx context.Context,
	id uuid.UUID,
) (*models.Email, error) {

	query := `
UPDATE emails
SET status = $2
WHERE id = $1
  AND status = $3
RETURNING
	id,
	sender,
	recipient,
	subject,
	text_body,
	html_body,
	status,
	idempotency_key,
	created_at
`

	email := &models.Email{}

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		models.EmailProcessing,
		models.EmailQueued,
	).Scan(
		&email.ID,
		&email.From,
		&email.To,
		&email.Subject,
		&email.Text,
		&email.HTML,
		&email.Status,
		&email.IdempotencyKey,
		&email.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// IMPORTANT FIX:
			// not found means:
			// - already processed OR
			// - invalid id
			// DO NOT treat as success silently
			return nil, ErrEmailNotFound
		}
		return nil, fmt.Errorf("claim email for processing: %w", err)
	}

	return email, nil
}

// -----------------------------
// PING
// -----------------------------
func (r *EmailRepository) Ping(ctx context.Context) error {
	if err := r.pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}
	return nil
}

func (r *EmailRepository) Close() {
	r.pool.Close()
}
