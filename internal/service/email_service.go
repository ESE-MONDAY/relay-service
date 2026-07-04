package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/ESE-MONDAY/relay-service/internal/dto"
	apperrors "github.com/ESE-MONDAY/relay-service/internal/errors"
	"github.com/ESE-MONDAY/relay-service/internal/models"
	"github.com/ESE-MONDAY/relay-service/internal/queue"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
)

const (
	jobTypeSendEmail = "send_email"
)

type EmailService interface {
	CreateEmail(
		ctx context.Context,
		req *dto.CreateEmailRequest,
	) (*dto.EmailResponse, error)
}

type emailService struct {
	repo     repository.EmailStore
	queue    queue.Publisher
	validate *validator.Validate
}

func NewEmailService(
	repo repository.EmailStore,
	queue queue.Publisher,
) EmailService {
	return &emailService{
		repo:     repo,
		queue:    queue,
		validate: validator.New(),
	}
}

func (s *emailService) CreateEmail(
	ctx context.Context,
	req *dto.CreateEmailRequest,
) (*dto.EmailResponse, error) {

	// Validate incoming request.
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	// Handle idempotent retries.
	if req.IdempotencyKey != "" {

		existing, err := s.repo.FindByIdempotencyKey(
			ctx,
			req.IdempotencyKey,
		)

		switch {

		case err == nil:
			// Already processed.
			return s.toResponse(existing), nil

		case errors.Is(err, repository.ErrEmailNotFound):
			// Safe to continue.

		default:
			// Database or unexpected error.
			return nil, apperrors.Internal(err)
		}
	}

	// Create a new email.
	email := s.newEmail(req)

	// Persist it.
	if err := s.saveEmail(ctx, email); err != nil {
		return nil, err
	}

	// Publish background job.
	if err := s.publishJob(ctx, email); err != nil {
		return nil, err
	}

	return s.toResponse(email), nil
}

func (s *emailService) validateRequest(
	req *dto.CreateEmailRequest,
) error {

	if err := s.validate.Struct(req); err != nil {
		return apperrors.Validation(
			"request validation failed",
			err,
		)
	}

	return nil
}

func (s *emailService) publishJob(
	ctx context.Context,
	email *models.Email,
) error {

	job := queue.Job{
		ID:      uuid.New(),
		EmailID: email.ID,
		Type:    jobTypeSendEmail,
	}

	if err := s.queue.Publish(ctx, job); err != nil {
		return apperrors.Internal(
			fmt.Errorf("publish job: %w", err),
		)
	}

	fmt.Printf("Published job: %+v\n", job)

	return nil
}

func (s *emailService) newEmail(
	req *dto.CreateEmailRequest,
) *models.Email {

	return &models.Email{
		ID:             uuid.New(),
		From:           req.From,
		To:             req.To,
		Subject:        req.Subject,
		Text:           req.Text,
		HTML:           req.HTML,
		IdempotencyKey: req.IdempotencyKey,
		Status:         models.EmailQueued,
	}
}

func (s *emailService) saveEmail(
	ctx context.Context,
	email *models.Email,
) error {

	if err := s.repo.Save(ctx, email); err != nil {
		return apperrors.Internal(
			fmt.Errorf("save email: %w", err),
		)
	}

	return nil
}

func (s *emailService) toResponse(
	email *models.Email,
) *dto.EmailResponse {

	return &dto.EmailResponse{
		ID:      email.ID.String(),
		Status:  string(email.Status),
		Message: "Email accepted",
	}
}
