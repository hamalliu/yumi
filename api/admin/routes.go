package admin

import (
	"time"

	"yumi/gin"
	"yumi/gin/middleware"
	"yumi/pkg/codec"
)

// Mount ...
func Mount(r gin.GroupRoutes) {
	ar := r.Group("行政管理", "admin")

	decrypter := make(map[string]codec.RsaDecrypter)
	ar.POST("登录", "login", middleware.LoginSecurity(decrypter, time.Second*15), login)

	//↑↑↑↑ 以上接口不需要这些中间件 ↑↑↑↑
	ar.Use(
		middleware.AuthToken(""),
		middleware.NoLoginSecurity(nil, time.Second*15),
	)

	ar.POST("注销", "logout", logout)
}
