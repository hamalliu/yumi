package middleware

import (
	"net/http"

	"yumi/pkg/gin"
)

// MaxBytes ...
func MaxBytes(n int64) gin.HandlerFunc {
	if n <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if c.Request.ContentLength > n {
			c.Writer.WriteHeader(http.StatusRequestEntityTooLarge)
		} else {
			c.Next()
		}
	}
}
