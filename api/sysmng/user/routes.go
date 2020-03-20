package user

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
	rr.Handle(http.MethodPost, "getmodules", AddPowerRoleOfUser, "0011", nil)
	rr.Handle(http.MethodPost, "getmodules", SavePowerRoleOfUser, "0011", nil)
	rr.Handle(http.MethodGet, "getmodules", GetPowerRoleOfUser, "0011", nil)
	rr.Handle(http.MethodGet, "getmodules", GetAllPowerOfUser, "0011", nil)
}
