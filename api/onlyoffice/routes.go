package onlyoffice

import (
	"yumi/gin"
)

//Mount ...
func Mount(r gin.GroupRoutes) {
	mr := r.Group("在线office", "media")

	mr.POST("上传", "/upload", Upload)
	mr.POST("上传样板", "/sample", Sample)
	mr.GET("获取编辑器", "/editor", Editor)
	mr.POST("新增编辑记录", "/track", Track)
	mr.GET("转换", "/convert", Convert)
	mr.GET("下载", "/download", Download)
	mr.DELETE("删除", "/file", DeleteFile)
}
