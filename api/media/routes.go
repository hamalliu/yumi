package media

import (
	"yumi/gin"
	"yumi/gin/middleware"
)

//Mount ...
func Mount(r gin.GroupRoutes) {
	mr := r.Group("文件管理", "media")

	mr.POST("上传多个", "/upload_multiple", middleware.RequiresPermissions([]string{"/media:upload_multiple"}), UploadMultiple)
	mr.POST("上传", "/upload", middleware.RequiresPermissions([]string{"/media:upload"}), Upload)
}
