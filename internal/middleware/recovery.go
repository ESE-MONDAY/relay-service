package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {

	return gin.CustomRecovery(

		func(c *gin.Context, recovered interface{}) {

			log.Error(
				"panic",

				zap.Any(
					"error",
					recovered,
				),

				zap.String(
					"request_id",
					GetRequestID(c),
				),
			)

			c.AbortWithStatusJSON(

				http.StatusInternalServerError,

				gin.H{

					"success": false,

					"error": gin.H{

						"code": "INTERNAL_SERVER_ERROR",

						"message": "internal server error",
					},

					"request_id": GetRequestID(c),
				},
			)
		},
	)
}
