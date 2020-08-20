package apidoc

import (
	"yumi/pkg/net/gin"
)

func Mount(r *gin.RouterGroup) {
	r.GET("", "api_doc/*path", ApiDoc)
}
