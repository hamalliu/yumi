package admin

import (
	"yumi/gin"
)

// Mount ...
func Mount(r *gin.RouterGroup) {
	ar := r.Group("admin")

	ar.POST("logout", logout)
}
