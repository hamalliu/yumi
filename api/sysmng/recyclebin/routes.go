package recyclebin

import (
	"net/http"

	"yumi/controller"
)

func Mount(r controller.Route) {
	rbr := r.Group("recyclebin")

	rbr.Handle(http.MethodPost, "add", SearchMyUpdate, "0011", nil)
	rbr.Handle(http.MethodPost, "delete", SearchMyDelete, "0011", nil)
	rbr.Handle(http.MethodPut, "update", RbMyUpdate, "0011", nil)
	rbr.Handle(http.MethodGet, "getmodules", RbMydelete, "0011", nil)
}
