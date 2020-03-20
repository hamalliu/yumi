package api

import (
	"yumi/api/admin"
	"yumi/api/media"
	"yumi/api/sysmng"
	"yumi/controller"
)

func Mount(r controller.Route) {
	ar := r.Group("api", Decrypt, PemissionAuth, DebugLog)

	admin.Mount(ar)
	media.Mount(ar)
	sysmng.Mount(ar)
}
