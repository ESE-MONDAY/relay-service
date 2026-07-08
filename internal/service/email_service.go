package service

import (
	"context"
	"fmt"
	"log"

	"github.com/ESE-MONDAY/relay-service/internal/event"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/ESE-MONDAY/relay-service/internal/dto"
	apperrors "github.com/ESE-MONDAY/relay-service/internal/errors"
	"github.com/ESE-MONDAY/relay-service/internal/models"
	"github.com/ESE-MONDAY/relay-service/internal/queue"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
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

	// 1. Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, apperrors.Validation("request validation failed", err)
	}

	// 2. FIX: handle *string correctly
	var keyPtr *string
	if idempotencyKey != "" {
		keyPtr = &idempotencyKey
	}

	// 3. Build model
	email := &models.Email{
		ID:             uuid.New(),
		From:           req.From,
		To:             req.To,
		Subject:        req.Subject,
		Text:           req.Text,
		HTML:           req.HTML,
		IdempotencyKey: keyPtr,
		Status:         models.EmailQueued,
	}

	// 4. Save email
	savedEmail, created, err := s.repo.Save(ctx, email)
	log.Println("CreateEmail error:", err)
	if err != nil {
		return nil, apperrors.Internal(fmt.Errorf("save email: %w", err))
	}

	// 5. Publish event only if newly created
	if created {
		ev := event.EmailEvent{
			EventID: uuid.NewString(),
			Type:    "email.created",
			EmailID: savedEmail.ID,
		}

		// FIX: no Time field unless your struct defines it
		// (you previously had mismatch error)

		if err := s.queue.Publish(ctx, ev); err != nil {
			return nil, apperrors.Internal(
				fmt.Errorf("publish email event: %w", err),
			)
		}

		fmt.Printf("Published event: %+v\n", ev)
	}

	// 6. Response
	return &dto.EmailResponse{
		ID:      savedEmail.ID.String(),
		Status:  string(savedEmail.Status),
		Message: "Email accepted",
	}, nil
}
