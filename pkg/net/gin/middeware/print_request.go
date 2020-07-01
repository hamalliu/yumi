package middeware

import (
	"bytes"
	"fmt"
	"io"

	"yumi/pkg/conf"
	"yumi/pkg/net/gin"
)

func PrintRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if conf.IsDebug() {
			fmt.Printf("url: %s\n", c.Request.URL.String())
			fmt.Printf("header: %s\n", c.Request.Header)
			body := []byte{}
			_, _ = io.Copy(bytes.NewBuffer(body), c.Request.Body)
			fmt.Printf("body: %s\n", string(body))
		}
	}
}