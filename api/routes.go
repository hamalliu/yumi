package api

import (
	"time"
	"yumi/api/admin"
	"yumi/api/doc"
	"yumi/api/login"
	"yumi/api/media"
	"yumi/api/onlyoffice"
	"yumi/conf"
	"yumi/gin"
	"yumi/gin/middleware"
)

//Mount ...
func Mount(r *gin.RouterGroup) {
	ar := r.Group("api")
	ar.Use(
		middleware.Recovery(),
		middleware.Cors(conf.Get().HttpServer.CORSAllowedOrigins, int(conf.Get().HttpServer.CORSMaxAge.Duration()/time.Second)),
		middleware.CSRF(nil, nil),
		middleware.MaxBytes(0),
		middleware.Debug(conf.IsDebug()),
	)

	// 接口文档接口
	doc.Mount(ar)

	// 登录接口
	login.Mount(ar)

	//↑↑↑↑ 以上接口不需要这些中间件 ↑↑↑↑
	ar.Use(
		middleware.AuthToken(""),
		middleware.NoLoginSecurity(nil, time.Second*15),
	)
	// 业务接口
	admin.Mount(ar)
	media.Mount(ar)
	onlyoffice.Mount(ar)
}
