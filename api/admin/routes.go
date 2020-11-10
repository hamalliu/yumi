package admin

import (
	"yumi/pkg/gin"
)

// Mount ...
func Mount(r *gin.RouterGroup) {
	ar := r.Group("admin")

	ar.POST("", "login", login)
	ar.POST("", "logout", logout)
}
