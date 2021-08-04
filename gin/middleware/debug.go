package middleware

import (
	"bytes"
	"fmt"
	"io"

	"yumi/gin"
)

//Debug ...
func Debug(isDebug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isDebug {
			fmt.Printf("url: %s\n", c.Request.URL.String())
			fmt.Printf("header: %s\n", c.Request.Header)
			body := []byte{}
			_, _ = io.Copy(bytes.NewBuffer(body), c.Request.Body)
			fmt.Printf("body: %s\n", string(body))
		}
		c.Next()
	}
}
