package media

import (
	"time"
	
	"yumi/gin"
	"yumi/gin/middleware"
)

//Mount ...
func Mount(r gin.GroupRoutes) {
	mr := r.Group("文件管理", "media")

	mr.Use(
		middleware.AuthToken(""),
		middleware.NoLoginSecurity(nil, time.Second*15),
	)
	mr.POST("上传多个", "/upload_multiple", middleware.RequiresPermissions([]string{"/media:upload_multiple"}), UploadMultiple)
	mr.POST("上传", "/upload", middleware.RequiresPermissions([]string{"/media:upload"}), Upload)
}
