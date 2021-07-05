package doc

import (
	"yumi/gin"
)

//Mount ...
func Mount(r gin.GroupRoutes) {
	r.GET("接口文档", "doc", APIDoc)
}
