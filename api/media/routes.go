package media

import (
	"yumi/pkg/net/gin"
)

func Mount(r *gin.RouterGroup) {
	mr := r.Group("media")

	mr.POST("", "/uploadmultipart", UploadMultipart)
	mr.POST("", "/upload", Upload)
}
