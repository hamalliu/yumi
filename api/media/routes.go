package media

import (
	"yumi/gin"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	mr := r.Group("media")

	mr.POST("/uploadmultipart", UploadMultipart)
	mr.POST("/upload", Upload)
}
