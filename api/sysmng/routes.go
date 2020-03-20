package sysmng

import (
	"yumi/api/sysmng/menu"
	"yumi/api/sysmng/recyclebin"
	"yumi/api/sysmng/role"
	"yumi/api/sysmng/user"
	"yumi/controller"
)

func Mount(r controller.Route) {
	ar := r.Group("sysmng")

	menu.Mount(ar)
	recyclebin.Mount(ar)
	role.Mount(ar)
	user.Mount(ar)
}
