package middeware

import (
	"net/http"
	"time"

	"github.com/rs/cors"

	"yumi/pkg/conf"
	"yumi/pkg/gin"
	"yumi/pkg/gin/header"
)

//Cors ...
func Cors(corsConf conf.CORS) gin.HandlerFunc {
	return func(c *gin.Context) {
		opts := cors.Options{
			AllowedOrigins:         corsConf.AllowedOrigins,
			AllowOriginFunc:        nil,
			AllowOriginRequestFunc: nil,
			AllowedMethods:         []string{http.MethodGet, http.MethodPost},
			AllowedHeaders:         header.ReqHeaders(),
			ExposedHeaders:         header.RespHeaders(),
			MaxAge:                 int(corsConf.MaxAge.Duration() / time.Second),
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
