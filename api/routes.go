package api

import (
	"yumi/api/admin"
	"yumi/api/media"
	"yumi/api/sysmng"
	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	ar := r.Group("api", Decrypt, PemissionAuth, DebugLog)

	admin.Mount(ar)
	media.Mount(ar)
	sysmng.Mount(ar)
}
