package api

import (
	"yumi/api/admin"
	"yumi/api/doc"
	"yumi/api/media"
	"yumi/gin"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	ar := r.Group("api", DebugLog)

	admin.Mount(ar)
	media.Mount(ar)
	doc.Mount(ar)
}
