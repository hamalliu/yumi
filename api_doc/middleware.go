package api_doc

import (
	"net/http"
	"yumi/pkg/net/ymhttp"
)

func ApiDoc(c *ymhttp.Context) {
	http.ServeFile(c.Writer, c.Request, "index.html")
}
