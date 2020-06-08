package sysmng

import (
	"yumi/api/sysmng/menu"
	"yumi/api/sysmng/role"
	"yumi/api/sysmng/user"
	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	ar := r.Group("sysmng")

	menu.Mount(ar)
	role.Mount(ar)
	user.Mount(ar)
}
