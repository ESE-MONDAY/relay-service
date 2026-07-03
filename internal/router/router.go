package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ESE-MONDAY/relay-service/internal/handler"
	"github.com/ESE-MONDAY/relay-service/internal/middleware"
)

func New(log *zap.Logger) *gin.Engine {

	r := gin.New()

	r.Use(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Recovery(log),
	)

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("failed to configure trusted proxies", zap.Error(err))
	}

	return r
}

type Handlers struct {
	Email  *handler.EmailHandler
	Health *handler.HealthHandler
}

func Register(
	r *gin.Engine,
	h *handler.Handler,
) {

	// Health endpoints
	r.GET("/health", h.Health.Health)
	r.GET("/ready", h.Health.Ready)

	// Version 1 API
	v1 := r.Group("/v1")
	{
		v1.POST("/emails", h.Email.CreateEmail)
	}
}
