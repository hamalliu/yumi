package apidoc

import (
	"yumi/pkg/gin"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	r.GET("", "api_doc/*path", APIDoc)
}
