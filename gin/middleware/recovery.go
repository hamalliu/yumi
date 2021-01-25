package middleware

import (
	"fmt"
	"net/http/httputil"
	"runtime"

	"yumi/pkg/ecode"
	"yumi/gin"
	"yumi/pkg/log"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			var rawReq []byte
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				if c.Request != nil {
					rawReq, _ = httputil.DumpRequest(c.Request, false)
				}
				pl := fmt.Sprintf("http call panic: %s\n%v\n%s\n", string(rawReq), err, buf)
				log.Error(pl)
				c.JSON(nil, ecode.ServerErr(nil))
			}
		}()
		c.Next()
	}
}