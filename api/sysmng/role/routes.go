package role

import (
	"net/http"

	"yumi/controller"
)

func Mount(r controller.Route) {
	rr := r.Group("role")

	rr.Handle(http.MethodPost, "add", Add, "0011", nil)
	rr.Handle(http.MethodPost, "delete", Delete, "0011", nil)
	rr.Handle(http.MethodPut, "update", Update, "0011", nil)
	rr.Handle(http.MethodGet, "getmodules", GetItem, "0011", nil)
	rr.Handle(http.MethodGet, "getmodules", Search, "0011", nil)
	rr.Handle(http.MethodGet, "getmodules", SavePowerUserOfRole, "0011", nil)
	rr.Handle(http.MethodGet, "getmodules", AddPowerUserOfRole, "0011", nil)
	rr.Handle(http.MethodGet, "getmodules", GetPowerUserOfRole, "0011", nil)
}
