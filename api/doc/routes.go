package doc

import (
	"yumi/gin"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	r.GET("doc", APIDoc)
}