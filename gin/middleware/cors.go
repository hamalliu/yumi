package middleware

import (
	"net/http"

	"github.com/rs/cors"

	"yumi/gin"
	"yumi/gin/header"
)

//Cors ...
func Cors(allowedOrigins []string, maxAge int) gin.HandlerFunc {
	return func(c *gin.Context) {
		opts := cors.Options{
			AllowedOrigins:         allowedOrigins,
			AllowOriginFunc:        nil,
			AllowOriginRequestFunc: nil,
			AllowedMethods:         []string{http.MethodGet, http.MethodPost},
			AllowedHeaders:         header.ReqHeaders(),
			ExposedHeaders:         header.RespHeaders(),
			MaxAge:                 maxAge, 
			AllowCredentials:       true,
			OptionsPassthrough:     false,
			Debug:                  false,
		}
		cors.New(opts).HandlerFunc(c.Writer, c.Request)
		if !opts.OptionsPassthrough &&
			c.Request.Method == http.MethodOptions &&
			c.GetHeader("Access-Control-Request-Method") != "" {
			c.AbortWithStatus(http.StatusOK)
		}

		c.Next()
	}
}
