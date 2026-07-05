package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ESE-MONDAY/relay-service/internal/dto"
	"github.com/ESE-MONDAY/relay-service/internal/response"
	"github.com/ESE-MONDAY/relay-service/internal/service"
)

type EmailHandler struct {
	service service.EmailService
}

func NewEmailHandler(
	service service.EmailService,
) *EmailHandler {

	return &EmailHandler{
		service: service,
	}
}

func (h *EmailHandler) CreateEmail(
	c *gin.Context,
) {

	var req dto.CreateEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	key := c.GetHeader("Idempotency-Key")

	resp, err := h.service.CreateEmail(
		c.Request.Context(),
		&req,
		key,
	)

	if err != nil {
		log.Printf("CreateEmail error: %+v\n", err)
		response.Error(c, err)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		resp,
	)
}
