package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/ESE-MONDAY/relay-service/internal/dto"
	apperrors "github.com/ESE-MONDAY/relay-service/internal/errors"
	"github.com/ESE-MONDAY/relay-service/internal/models"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
)

const (
	emailStatusQueued = "queued"
)

type EmailService interface {
	CreateEmail(
		ctx context.Context,
		req *dto.CreateEmailRequest,
	) (*dto.EmailResponse, error)
}

type emailService struct {
	repo     repository.EmailStore
	validate *validator.Validate
}

func NewEmailService(repo repository.EmailStore) EmailService {
	return &emailService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *emailService) CreateEmail(
	ctx context.Context,
	req *dto.CreateEmailRequest,
) (*dto.EmailResponse, error) {

	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	email := s.newEmail(req)

	if err := s.saveEmail(ctx, email); err != nil {
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

func (s *emailService) newEmail(
	req *dto.CreateEmailRequest,
) *models.Email {

	return &models.Email{
		ID:      uuid.New(),
		From:    req.From,
		To:      req.To,
		Subject: req.Subject,
		Text:    req.Text,
		HTML:    req.HTML,
		Status:  emailStatusQueued,
	}
}

func (s *emailService) saveEmail(
	ctx context.Context,
	email *models.Email,
) error {

	if err := s.repo.Save(ctx, email); err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (s *emailService) toResponse(
	email *models.Email,
) *dto.EmailResponse {

	return &dto.EmailResponse{
		ID:      email.ID.String(),
		Status:  email.Status,
		Message: "Email accepted",
	}
}
