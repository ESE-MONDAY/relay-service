package middleware

import "github.com/gin-gonic/gin"

func GetRequestID(c *gin.Context) string {

	v, ok := c.Get(RequestIDKey)

	if !ok {
		return ""
	}

	id, ok := v.(string)

	if !ok {
		return ""
	}

	return id
}
