package apidoc

import (
	"yumi/pkg/net/gin"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	r.GET("", "api_doc/*path", APIDoc)
}
