package onlyoffice

import (
	"yumi/pkg/net/gin"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	mr := r.Group("media")

	mr.POST("", "/upload", Upload)
	mr.POST("", "/sample", Sample)
	mr.GET("", "/editor", Editor)
	mr.POST("", "/track", Track)
	mr.GET("", "/convert", Convert)
	mr.GET("", "/download", Download)
	mr.DELETE("", "/file", DeleteFile)
}
