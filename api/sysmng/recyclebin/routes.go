package recyclebin

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	rbr := r.Group("recyclebin")

	rbr.Handle(http.MethodPost, "", "add", SearchMyUpdate)
	rbr.Handle(http.MethodPost, "", "delete", SearchMyDelete)
	rbr.Handle(http.MethodPut, "", "update", RbMyUpdate)
	rbr.Handle(http.MethodGet, "", "getmodules", RbMydelete)
}
