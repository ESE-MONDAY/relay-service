package handler

import (
	"github.com/ESE-MONDAY/relay-service/internal/service"
)

type Handler struct {
	Email  *EmailHandler
	Health *HealthHandler
}

func New(
	emailService service.EmailService,
) *Handler {

	return &Handler{
		Email:  NewEmailHandler(emailService),
		Health: NewHealthHandler(),
	}
}
