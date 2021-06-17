package login

import (
	"time"
	
	"yumi/gin"
	"yumi/gin/middleware"
	"yumi/pkg/codec"
)

// Mount ...
func Mount(r *gin.RouterGroup) {
	decrypter := make(map[string]codec.RsaDecrypter)
	r.POST("login", middleware.LoginSecurity(decrypter, time.Second*15), login)
}
