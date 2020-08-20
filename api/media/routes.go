package media

import (
	"yumi/pkg/net/gin"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	mr := r.Group("media")

	mr.POST("", "/uploadmultipart", UploadMultipart)
	mr.POST("", "/upload", Upload)
}
