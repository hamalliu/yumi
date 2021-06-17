package media

import (
	"yumi/gin"
	"yumi/gin/middleware"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	mr := r.Group("media")

	mr.POST("/upload_multiple", middleware.RequiresPermissions([]string{"/media:upload_multiple"}), UploadMultiple)
	mr.POST("/upload", middleware.RequiresPermissions([]string{"/media:upload"}), Upload)
}
