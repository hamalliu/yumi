package middeware

import (
	"net/http"

	"github.com/rs/cors"

	"yumi/pkg/conf"
	"yumi/pkg/net/ymhttp"
	"yumi/pkg/net/ymhttp/header"
)

func Cors() ymhttp.HandlerFunc {
	return func(c *ymhttp.Context) {
		opts := cors.Options{
			AllowedOrigins:         conf.GetCORS().AllowedOrigins,
			AllowOriginFunc:        nil,
			AllowOriginRequestFunc: nil,
			AllowedMethods:         []string{http.MethodGet, http.MethodPost},
			AllowedHeaders:         header.ReqHeaders(),
			ExposedHeaders:         header.RespHeaders(),
			MaxAge:                 conf.GetCORS().MaxAge,
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
