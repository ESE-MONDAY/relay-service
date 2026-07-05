package service

import (
	"context"
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
		idempotencyKey string,
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
	idempotencyKey string,
) (*dto.EmailResponse, error) {

	// Validate request body.
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	// Build the email model.
	email := s.newEmail(
		req,
		idempotencyKey,
	)

	// Save email.
	savedEmail, created, err := s.saveEmail(
		ctx,
		email,
	)
	if err != nil {
		return nil, err
	}

	// Only enqueue if this is a newly created email.
	if created {
		if err := s.publishJob(
			ctx,
			savedEmail,
		); err != nil {
			return nil, err
		}
	}

	return s.toResponse(savedEmail), nil
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

func (s *emailService) newEmail(
	req *dto.CreateEmailRequest,
	idempotencyKey string,
) *models.Email {

	return &models.Email{
		ID:             uuid.New(),
		From:           req.From,
		To:             req.To,
		Subject:        req.Subject,
		Text:           req.Text,
		HTML:           req.HTML,
		IdempotencyKey: idempotencyKey,
		Status:         models.EmailQueued,
	}
}

func (s *emailService) saveEmail(
	ctx context.Context,
	email *models.Email,
) (
	*models.Email,
	bool,
	error,
) {

	savedEmail, created, err := s.repo.Save(ctx, email)
	if err != nil {

		fmt.Printf("SAVE ERROR: %v\n", err)

		return nil, false, apperrors.Internal(
			fmt.Errorf("save email: %w", err),
		)
	}

	return savedEmail, created, nil
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

	if err := s.queue.Publish(
		ctx,
		job,
	); err != nil {

		return apperrors.Internal(
			fmt.Errorf("publish job: %w", err),
		)
	}

	fmt.Printf("Published job: %+v\n", job)

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
