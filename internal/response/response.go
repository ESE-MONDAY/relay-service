package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/ESE-MONDAY/relay-service/internal/errors"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {

	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Error: ErrorBody{
				Code:    "INTERNAL_ERROR",
				Message: "internal server error",
			},
		})
		return
	}

	status := http.StatusInternalServerError

	switch appErr.Code {
	case apperrors.ErrValidation:
		status = http.StatusBadRequest
	case apperrors.ErrNotFound:
		status = http.StatusNotFound
	case apperrors.ErrConflict:
		status = http.StatusConflict
	}

	c.JSON(status, ErrorResponse{
		Success: false,
		Error: ErrorBody{
			Code:    string(appErr.Code),
			Message: appErr.Message,
		},
	})
}
