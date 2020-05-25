package middeware

import (
	"bytes"
	"fmt"
	"io"

	"yumi/pkg/conf"
	"yumi/pkg/net/ymhttp"
)

func PrintRequest() ymhttp.HandlerFunc {
	return func(c *ymhttp.Context) {
		if conf.Get().Environment == conf.EnvDebug {
			fmt.Printf("url: %s\n", c.Request.URL.String())
			fmt.Printf("header: %s\n", c.Request.Header)
			body := []byte{}
			_, _ = io.Copy(bytes.NewBuffer(body), c.Request.Body)
			fmt.Printf("body: %s\n", string(body))
		}
	}
}
