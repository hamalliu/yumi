package api

import (
	"time"

	"yumi/api/admin"
	"yumi/api/doc"
	"yumi/api/media"
	"yumi/api/onlyoffice"
	"yumi/conf"
	"yumi/gin"
	"yumi/gin/middleware"
)

//Mount ...
func Mount(r gin.GroupRoutes) {
	ar := r.Group("接口", "api")
	ar.Use(
		middleware.Recovery(),
		middleware.InitCasbin(conf.Get().Casbin.ModelFile, nil),
		middleware.Cors(conf.Get().HttpServer.CORSAllowedOrigins, int(conf.Get().HttpServer.CORSMaxAge.Duration()/time.Second)),
		middleware.CSRF(nil, nil),
		middleware.MaxBytes(0),
		middleware.Debug(conf.IsDebug()),
	)

	// 接口文档接口
	doc.Mount(ar)

	// 业务接口
	admin.Mount(ar)
	media.Mount(ar)
	onlyoffice.Mount(ar)
}
