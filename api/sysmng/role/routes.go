package role

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	rr := r.Group("role")

	rr.Handle(http.MethodPost, "", "add", Add)
	rr.Handle(http.MethodPost, "", "delete", Delete)
	rr.Handle(http.MethodPut, "", "update", Update)
	rr.Handle(http.MethodGet, "", "getmodules", GetItem)
	rr.Handle(http.MethodGet, "", "getmodules", Search)
	rr.Handle(http.MethodGet, "", "getmodules", SavePowerUserOfRole)
	rr.Handle(http.MethodGet, "", "getmodules", AddPowerUserOfRole)
	rr.Handle(http.MethodGet, "", "getmodules", GetPowerUserOfRole)
}
