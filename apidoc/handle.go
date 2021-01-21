package apidoc

import (
	"fmt"
	"net/http"

	"yumi/gin"
)

//APIDoc ...
func APIDoc(c *gin.Context) {
	path, _ := c.Params.Get("path")
	if path == "" {
		path = "api_doc/index.html"
	} else {
		path = fmt.Sprintf("api_doc/%s", path)
	}

	http.ServeFile(c.Writer, c.Request, path)
}
