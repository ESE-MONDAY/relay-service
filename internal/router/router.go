package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ESE-MONDAY/relay-service/internal/middleware"
)

func New(log *zap.Logger) *gin.Engine {

	r := gin.New()

	r.Use(

		middleware.RequestID(),

		middleware.Logger(log),

		middleware.Recovery(log),
	)

	return r
}
