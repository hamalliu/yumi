package api

import (
	"yumi/api/admin"
	"yumi/api/media"
	"yumi/pkg/net/gin"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	ar := r.Group("api", DebugLog)

	admin.Mount(ar)
	media.Mount(ar)
}
