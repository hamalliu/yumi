package user

import (
	"net/http"

	"yumi/pkg/net/ymhttp"
)

func Mount(r *ymhttp.RouterGroup) {
	rr := r.Group("role")

	rr.Handle(http.MethodPost, "", "/add", Add)
	rr.Handle(http.MethodPost, "", "/delete", Delete)
	rr.Handle(http.MethodPut, "", "/update", Update)
	rr.Handle(http.MethodGet, "", "/getmodules", GetItem)
	rr.Handle(http.MethodGet, "", "/getmodules", Search)
	rr.Handle(http.MethodPost, "", "/getmodules", AddPowerRoleOfUser)
	rr.Handle(http.MethodPost, "", "/getmodules", SavePowerRoleOfUser)
	rr.Handle(http.MethodGet, "", "/getmodules", GetPowerRoleOfUser)
	rr.Handle(http.MethodGet, "", "/getmodules", GetAllPowerOfUser)
}
