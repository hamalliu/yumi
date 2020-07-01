package admin

import (
	"yumi/pkg/net/gin"
)

func Mount(r *gin.RouterGroup) {
	ar := r.Group("admin")

	ar.POST("", "login", login)
	ar.POST("", "logout", logout)
}
