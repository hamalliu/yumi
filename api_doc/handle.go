package api_doc

import (
	"fmt"
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func ApiDoc(c *ymhttp.Context) {
	path, _ := c.Params.Get("path")
	if path == "" {
		path = "api_doc/index.html"
	} else {
		path = fmt.Sprintf("api_doc/%s", path)
	}

	http.ServeFile(c.Writer, c.Request, path)
}
